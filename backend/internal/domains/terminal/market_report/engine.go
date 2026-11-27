package market_report

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/domains/terminal/shared"
)

const defaultStockbitVerifyURL = "https://exodus.stockbit.com/marketdetectors/BBCA?from=2023-01-01&to=2023-01-02"

// Engine manages report generation, token updates, log ring-buffer, and report archives.
type Engine struct {
	mu                sync.Mutex
	logs              []string
	maxLogs           int
	state             EngineState
	reportDir         string
	logFile           string
	smEnvFile         string
	pythonBin         string
	scriptPath        string
	smtrackerSrcDir   string
	stockbitVerifyURL string
	httpClient        *http.Client
	runCmd            CommandRunner
}

// NewEngine initializes the Market Report Engine without mutating disk on read-only boot.
func NewEngine(deps Deps) *Engine {
	reportDir := deps.ReportDir
	if reportDir == "" {
		reportDir = "/root/lazytrader/reports"
	}
	if _, err := os.Stat(reportDir); os.IsNotExist(err) {
		_ = os.MkdirAll(reportDir, 0o755)
	}

	smDir := deps.SMtrackerDir
	if smDir == "" {
		smDir = "/root/SMtracker"
	}
	scriptPath := deps.ScriptPath
	if scriptPath == "" {
		scriptPath = "/root/SMtracker/tools/terminal/generate_universe_report.py"
	}
	verifyURL := deps.StockbitVerifyURL
	if verifyURL == "" {
		verifyURL = defaultStockbitVerifyURL
	}
	client := deps.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 8 * time.Second}
	}
	runner := deps.Runner
	if runner == nil {
		runner = defaultCommandRunner
	}

	e := &Engine{
		logs:    make([]string, 0, 512),
		maxLogs: 1000,
		state: EngineState{
			Status:    "idle",
			IsRunning: false,
			Session:   1,
			Progress:  0,
			Total:     962,
			Pct:       0.0,
			Message:   "Siap untuk membuat laporan.",
		},
		reportDir:         reportDir,
		logFile:           filepath.Join(reportDir, "execution.log"),
		smEnvFile:         filepath.Join(smDir, ".env"),
		pythonBin:         filepath.Join(smDir, "venv", "bin", "python"),
		scriptPath:        scriptPath,
		smtrackerSrcDir:   filepath.Join(smDir, "src"),
		stockbitVerifyURL: verifyURL,
		httpClient:        client,
		runCmd:            runner,
	}

	e.loadLogsFromDisk()
	// Record startup line in memory only (avoid mutating /root/rev/reports/execution.log during read-only boot/parity).
	e.appendLogLocked("Sistem Market Report Engine siap.", false)
	e.loadLatestFromDiskLocked()
	return e
}

func getWIBNow() string {
	return time.Now().In(shared.WIBLocation()).Format("2006-01-02 15:04:05 WIB")
}

func getWIBTime() string {
	return time.Now().In(shared.WIBLocation()).Format("15:04:05")
}

func (e *Engine) loadLogsFromDisk() {
	f, err := os.Open(e.logFile)
	if err != nil {
		return
	}
	defer f.Close()

	var all []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			all = append(all, line)
		}
	}
	if len(all) > 400 {
		all = all[len(all)-400:]
	}
	e.logs = append(e.logs, all...)
	if len(e.logs) > e.maxLogs {
		e.logs = e.logs[len(e.logs)-e.maxLogs:]
	}
}

func (e *Engine) appendLogLocked(text string, persist bool) {
	line := fmt.Sprintf("[%s] %s", getWIBTime(), text)
	e.logs = append(e.logs, line)
	if len(e.logs) > e.maxLogs {
		excess := len(e.logs) - e.maxLogs
		copy(e.logs, e.logs[excess:])
		e.logs = e.logs[:e.maxLogs]
	}
	if persist {
		if f, err := os.OpenFile(e.logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err == nil {
			_, _ = f.WriteString(line + "\n")
			_ = f.Close()
		}
	}
}

func (e *Engine) addLog(text string, persist bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.appendLogLocked(text, persist)
}

func (e *Engine) loadLatestFromDiskLocked() {
	matches, err := filepath.Glob(filepath.Join(e.reportDir, "market_report_sesi*.txt"))
	if err != nil || len(matches) == 0 {
		return
	}
	type fileEntry struct {
		path    string
		modTime time.Time
	}
	entries := make([]fileEntry, 0, len(matches))
	for _, p := range matches {
		if strings.HasSuffix(p, "_unmasked.txt") {
			continue
		}
		if info, err := os.Stat(p); err == nil {
			entries = append(entries, fileEntry{path: p, modTime: info.ModTime()})
		}
	}
	if len(entries) == 0 {
		return
	}
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].modTime.After(entries[j].modTime)
	})

	maskedPath := entries[0].path
	base := filepath.Base(maskedPath)
	unmaskedPath := strings.TrimSuffix(maskedPath, ".txt") + "_unmasked.txt"

	maskedBytes, err := os.ReadFile(maskedPath)
	if err != nil {
		return
	}
	maskedText := string(maskedBytes)
	unmaskedText := maskedText
	if ub, err := os.ReadFile(unmaskedPath); err == nil {
		unmaskedText = string(ub)
	}

	session := 2
	if strings.Contains(base, "sesi1") {
		session = 1
	}
	modWIB := entries[0].modTime.In(shared.WIBLocation())
	e.state.LatestReport = &LatestReport{
		Session:      session,
		Filename:     base,
		Date:         modWIB.Format("2006-01-02"),
		GeneratedAt:  modWIB.Format("2006-01-02 15:04:05 WIB"),
		MaskedText:   maskedText,
		UnmaskedText: &unmaskedText,
	}
}

// GetStatus returns a deep copy of current status, token validity, time, and recent 150 logs.
func (e *Engine) GetStatus(includeUnmasked bool) StatusData {
	e.mu.Lock()
	st := e.state
	var latestCopy *LatestReport
	if st.LatestReport != nil {
		cp := *st.LatestReport
		if !includeUnmasked {
			cp.UnmaskedText = nil
		} else if cp.UnmaskedText != nil {
			u := *cp.UnmaskedText
			cp.UnmaskedText = &u
		}
		latestCopy = &cp
	}
	recentStart := len(e.logs) - 150
	if recentStart < 0 {
		recentStart = 0
	}
	recent := make([]string, len(e.logs)-recentStart)
	copy(recent, e.logs[recentStart:])
	e.mu.Unlock()

	return StatusData{
		Status:         st.Status,
		IsRunning:      st.IsRunning,
		Session:        st.Session,
		Progress:       st.Progress,
		Total:          st.Total,
		Pct:            st.Pct,
		Message:        st.Message,
		StartedAt:      st.StartedAt,
		CompletedAt:    st.CompletedAt,
		Error:          st.Error,
		LatestReport:   latestCopy,
		TokenValid:     e.CheckTokenAvailable(),
		CurrentTimeWIB: getWIBNow(),
		RecentLogs:     recent,
	}
}

// GetLogs returns all buffered execution logs.
func (e *Engine) GetLogs() []string {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make([]string, len(e.logs))
	copy(out, e.logs)
	return out
}

// ResetError clears error state and sets engine status back to idle if not running.
func (e *Engine) ResetError() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.state.IsRunning {
		e.state.Error = nil
		e.state.Status = "idle"
		e.state.Message = "Siap untuk membuat laporan."
	}
}

// GetHistory returns archived masked report metadata ordered newest-first by mtime.
func (e *Engine) GetHistory() []HistoryItem {
	matches, err := filepath.Glob(filepath.Join(e.reportDir, "market_report_sesi*.txt"))
	if err != nil {
		return []HistoryItem{}
	}
	type fileInfo struct {
		name    string
		modTime time.Time
		size    int64
	}
	list := make([]fileInfo, 0, len(matches))
	for _, p := range matches {
		if strings.HasSuffix(p, "_unmasked.txt") {
			continue
		}
		if info, err := os.Stat(p); err == nil {
			list = append(list, fileInfo{
				name:    filepath.Base(p),
				modTime: info.ModTime(),
				size:    info.Size(),
			})
		}
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].modTime.After(list[j].modTime)
	})

	out := make([]HistoryItem, 0, len(list))
	for _, item := range list {
		session := 2
		if strings.Contains(item.name, "sesi1") {
			session = 1
		}
		modWIB := item.modTime.In(shared.WIBLocation())
		out = append(out, HistoryItem{
			Filename:   item.name,
			Session:    session,
			Date:       modWIB.Format("2006-01-02"),
			ModifiedAt: modWIB.Format("2006-01-02 15:04:05 WIB"),
			SizeBytes:  item.size,
		})
	}
	return out
}

// ReadReport reads masked and unmasked text of a validated historical report file.
func (e *Engine) ReadReport(filename string, includeUnmasked bool) *ReportDetail {
	if strings.Contains(filename, "..") || strings.ContainsAny(filename, "/\\") {
		return nil
	}
	safeFn := filepath.Base(filename)
	if !strings.HasPrefix(safeFn, "market_report_sesi") || !strings.HasSuffix(safeFn, ".txt") {
		return nil
	}

	baseMasked := strings.TrimSuffix(safeFn, "_unmasked.txt")
	if !strings.HasSuffix(baseMasked, ".txt") {
		baseMasked += ".txt"
	}
	baseUnmasked := strings.TrimSuffix(baseMasked, ".txt") + "_unmasked.txt"

	maskedPath := filepath.Join(e.reportDir, baseMasked)
	unmaskedPath := filepath.Join(e.reportDir, baseUnmasked)

	maskedBytes, err := os.ReadFile(maskedPath)
	if err != nil {
		return nil
	}
	maskedText := string(maskedBytes)

	var unmaskedPtr *string
	if includeUnmasked {
		unmaskedText := maskedText
		if ub, err := os.ReadFile(unmaskedPath); err == nil {
			unmaskedText = string(ub)
		}
		unmaskedPtr = &unmaskedText
	}

	session := 2
	if strings.Contains(baseMasked, "sesi1") {
		session = 1
	}
	return &ReportDetail{
		Session:      session,
		Filename:     baseMasked,
		MaskedText:   maskedText,
		UnmaskedText: unmaskedPtr,
	}
}
