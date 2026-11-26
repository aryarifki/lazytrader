package system

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// JobConfigs mirrors routers/system.py JOB_CONFIGS.
var JobConfigs = map[string]JobConfig{
	"scanner": {
		Key: "scanner", Label: "Scanner", Script: "scanner.py", LogPrefix: "scanner",
		Description: "AI & Smart Money Stock Scanner (Wyckoff + Flow)",
		Category:    "daily", Params: []string{},
	},
	"auditor": {
		Key: "auditor", Label: "Auditor", Script: "auditor.py", LogPrefix: "auditor",
		Description: "Evaluasi sinyal OPEN & Update TP/SL ke Telegram",
		Category:    "daily", Params: []string{},
	},
	"backfill": {
		Key: "backfill", Label: "Daily Backfill", Script: "backfill_daily.py", LogPrefix: "backfill",
		Description: "Backfill data transaksi & broker harian (IDX + Stockbit)",
		Category:    "daily", SupportsUniverse: true, Params: []string{"universe", "days", "force", "rate_limit"},
	},
	"backfill_monthly": {
		Key: "backfill_monthly", Label: "Monthly Backfill", Script: "backfill_monthly.py", LogPrefix: "backfill_monthly",
		Description: "Backfill data broker bulanan historis (Hybrid IDX Sequential + Stockbit)",
		Category:    "maintenance", SupportsUniverse: true, Params: []string{"universe", "months", "force", "rate_limit", "status_only"},
	},
	"etl_flow": {
		Key: "etl_flow", Label: "ETL Foreign Flow", Script: "etl_foreign_flow.py", LogPrefix: "etl_flow",
		Description: "Komputasi analitik Foreign Flow (HMM, VAR, Network, Heatmap)",
		Category:    "daily", SupportsUniverse: true, SupportsWindow: true, Params: []string{"universe", "windows", "workers"},
	},
	"backfill_fundamentals": {
		Key: "backfill_fundamentals", Label: "Fundamentals & Governance", Script: "backfill_fundamentals.py", LogPrefix: "cron_fundamentals",
		Description: "Tarik data fundamental, rasio finansial, corporate actions, dan governance dari IDX",
		Category:    "maintenance", Params: []string{"details", "only_missing", "ratios", "actions", "year", "quarter"},
	},
	"db_audit": {
		Key: "db_audit", Label: "DB Health Audit", Script: "db_daily_audit.py", LogPrefix: "db_audit",
		Description: "Audit integritas database, self-healing, dan vacuuming",
		Category:    "daily", Params: []string{"autofix", "vacuum", "notify", "alert_only", "audit_date"},
	},
	"pipeline": {
		Key: "pipeline", Label: "Full Pipeline", Script: "run_daily.sh", LogPrefix: "cron_daily",
		Description: "Pipeline Harian Lengkap (Backfill + DB Audit + ETL + Audit + Scanner + Backup)",
		Category:    "pipeline", IsShell: true, Params: []string{},
	},
	"backup": {
		Key: "backup", Label: "Database Backup", Script: "auto_backup.sh", LogPrefix: "backup",
		Description: "Pencadangan database PostgreSQL (pg_dump) & sync Google Drive",
		Category:    "maintenance", IsShell: true, Params: []string{},
	},
}

type runningJob struct {
	cmd       *exec.Cmd
	pid       int
	startedAt time.Time
	logFile   string
	doneCh    chan int
}

type lastRunResult struct {
	finishedAt string
	exitCode   int
	status     string
	logFile    string
}

// JobManager tracks and triggers SMtracker automation scripts.
type JobManager struct {
	SMtrackerDir string
	PythonBin    string
	mu           sync.Mutex
	running      map[string]*runningJob
	lastResults  map[string]lastRunResult
}

// NewJobManager initializes a JobManager rooted at smtrackerDir.
func NewJobManager(smtrackerDir string) *JobManager {
	if strings.TrimSpace(smtrackerDir) == "" {
		smtrackerDir = "/root/SMtracker"
	}
	py := filepath.Join(smtrackerDir, "venv", "bin", "python3")
	if _, err := os.Stat(py); err != nil {
		py = "python3"
	}
	return &JobManager{
		SMtrackerDir: smtrackerDir,
		PythonBin:    py,
		running:      make(map[string]*runningJob),
		lastResults:  make(map[string]lastRunResult),
	}
}

// BuildJobCommand builds the non-interactive CLI invocation matching routers/system.py.
func (m *JobManager) BuildJobCommand(jobName string, cfg JobConfig, req JobRunRequest) []string {
	if cfg.IsShell {
		return []string{"/bin/bash", filepath.Join(m.SMtrackerDir, cfg.Script)}
	}
	cmd := []string{m.PythonBin, cfg.Script}

	switch jobName {
	case "backfill":
		uni := "all"
		if req.Universe != nil && *req.Universe != "" {
			uni = *req.Universe
		}
		cmd = append(cmd, "--universe", uni)
		if req.Days != nil && *req.Days > 1 {
			cmd = append(cmd, "--days", strconv.Itoa(*req.Days))
		}
		if req.Force {
			cmd = append(cmd, "--force")
		}
		if req.RateLimit != nil && *req.RateLimit > 0 {
			cmd = append(cmd, "--rate-limit", strconv.FormatFloat(*req.RateLimit, 'f', -1, 64))
		}
	case "backfill_monthly":
		uni := "all"
		if req.Universe != nil && *req.Universe != "" {
			uni = *req.Universe
		}
		cmd = append(cmd, "--universe", uni)
		if req.StatusOnly {
			cmd = append(cmd, "--status")
		} else {
			months := "last3"
			if req.Months != nil && *req.Months != "" {
				months = *req.Months
			}
			cmd = append(cmd, "--months", months)
			if req.Force {
				cmd = append(cmd, "--force")
			}
			if req.RateLimit != nil && *req.RateLimit > 0 {
				cmd = append(cmd, "--rate-limit", strconv.FormatFloat(*req.RateLimit, 'f', -1, 64))
			}
			cmd = append(cmd, "-y")
		}
	case "etl_flow":
		uni := "watchlist"
		if req.Universe != nil && *req.Universe != "" {
			uni = *req.Universe
		}
		cmd = append(cmd, "--universe", uni, "--window")
		wins := req.Windows
		if len(wins) == 0 {
			wins = []int{20, 60}
		}
		for _, w := range wins {
			cmd = append(cmd, strconv.Itoa(w))
		}
		if req.Workers != nil && *req.Workers > 0 {
			cmd = append(cmd, "--workers", strconv.Itoa(*req.Workers))
		}
	case "backfill_fundamentals":
		hasSel := req.Details || req.Ratios || req.Actions
		if !hasSel || req.Details {
			cmd = append(cmd, "--details")
		}
		if req.OnlyMissing {
			cmd = append(cmd, "--only-missing")
		}
		if req.Ratios {
			y := time.Now().Year()
			if req.Year != nil && *req.Year >= 2000 {
				y = *req.Year
			}
			q := 4
			if req.Quarter != nil && *req.Quarter >= 1 && *req.Quarter <= 4 {
				q = *req.Quarter
			}
			cmd = append(cmd, "--ratios", "--year", strconv.Itoa(y), "--quarter", strconv.Itoa(q))
		}
		if req.Actions {
			cmd = append(cmd, "--actions")
		}
	case "db_audit":
		if req.AlertOnly {
			cmd = append(cmd, "--alert-only")
		} else {
			if req.Autofix {
				cmd = append(cmd, "--autofix")
			}
			if req.Vacuum {
				cmd = append(cmd, "--vacuum")
			}
			if req.Notify {
				cmd = append(cmd, "--notify")
			}
		}
		if req.AuditDate != nil && *req.AuditDate != "" {
			cmd = append(cmd, "--date", *req.AuditDate)
		}
	}
	return cmd
}

func findExternalScriptPID(scriptName string) *int {
	if scriptName == "" {
		return nil
	}
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil
	}
	selfPID := os.Getpid()
	for _, entry := range entries {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil || pid == selfPID {
			continue
		}
		raw, err := os.ReadFile(filepath.Join("/proc", entry.Name(), "cmdline"))
		if err != nil {
			continue
		}
		cmd := strings.ReplaceAll(string(raw), "\x00", " ")
		if strings.Contains(cmd, scriptName) {
			return &pid
		}
	}
	return nil
}

// CheckStatus returns live status for a single job.
func (m *JobManager) CheckStatus(jobName string) JobStatusItem {
	cfg := JobConfigs[jobName]
	m.mu.Lock()
	if info, ok := m.running[jobName]; ok {
		select {
		case code := <-info.doneCh:
			st := "failed"
			if code == 0 {
				st = "success"
			}
			m.lastResults[jobName] = lastRunResult{
				finishedAt: time.Now().Format(time.RFC3339),
				exitCode:   code,
				status:     st,
				logFile:    info.logFile,
			}
			delete(m.running, jobName)
		default:
			elapsed := int(time.Since(info.startedAt).Seconds())
			pid := info.pid
			m.mu.Unlock()
			return JobStatusItem{
				Key: jobName, Label: cfg.Label, Description: cfg.Description, Category: cfg.Category,
				Status: "running", PID: &pid, StartedAt: info.startedAt.Format(time.RFC3339),
				ElapsedSeconds: &elapsed, LogFile: info.logFile,
				SupportsUniverse: cfg.SupportsUniverse, SupportsWindow: cfg.SupportsWindow, Params: cfg.Params,
			}
		}
	}
	last, hasLast := m.lastResults[jobName]
	m.mu.Unlock()

	if extPID := findExternalScriptPID(cfg.Script); extPID != nil {
		return JobStatusItem{
			Key: jobName, Label: cfg.Label, Description: cfg.Description, Category: cfg.Category,
			Status: "running", PID: extPID, IsExternal: true,
			SupportsUniverse: cfg.SupportsUniverse, SupportsWindow: cfg.SupportsWindow, Params: cfg.Params,
		}
	}

	status := "idle"
	var lastFin, lastLog *string
	var lastCode *int
	if hasLast {
		status = last.status
		lastFin = &last.finishedAt
		c := last.exitCode
		lastCode = &c
		if last.logFile != "" {
			lf := last.logFile
			lastLog = &lf
		}
	}
	return JobStatusItem{
		Key: jobName, Label: cfg.Label, Description: cfg.Description, Category: cfg.Category,
		Status: status, PID: nil, LastFinishedAt: lastFin, LastExitCode: lastCode, LastLogFile: lastLog,
		SupportsUniverse: cfg.SupportsUniverse, SupportsWindow: cfg.SupportsWindow, Params: cfg.Params,
	}
}

// StartJob launches a background job writing output to <SMtrackerDir>/debug/.
func (m *JobManager) StartJob(jobName string, cfg JobConfig, args []string) (int, string, error) {
	logDir := filepath.Join(m.SMtrackerDir, "debug")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return 0, "", err
	}
	logName := fmt.Sprintf("%s_%s.log", cfg.LogPrefix, time.Now().Format("20060102"))
	logPath := filepath.Join(logDir, logName)

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return 0, "", err
	}
	banner := fmt.Sprintf("\n\n%s\n🚀 MANUAL TRIGGER [%s] PADA %s\nPerintah : %s\n%s\n\n",
		strings.Repeat("=", 60), strings.ToUpper(jobName),
		time.Now().Format("2006-01-02 15:04:05"), strings.Join(args, " "), strings.Repeat("=", 60))
	_, _ = f.WriteString(banner)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = m.SMtrackerDir
	cmd.Stdout = f
	cmd.Stderr = f
	cmd.Env = append(os.Environ(), "PYTHONPATH="+filepath.Join(m.SMtrackerDir, "src"))
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	if err := cmd.Start(); err != nil {
		_ = f.Close()
		return 0, "", err
	}
	doneCh := make(chan int, 1)
	go func() {
		err := cmd.Wait()
		_ = f.Close()
		code := 0
		if err != nil {
			code = 1
			if exitErr, ok := err.(*exec.ExitError); ok {
				code = exitErr.ExitCode()
			}
		}
		doneCh <- code
	}()

	m.mu.Lock()
	m.running[jobName] = &runningJob{
		cmd: cmd, pid: cmd.Process.Pid, startedAt: time.Now(), logFile: logName, doneCh: doneCh,
	}
	m.mu.Unlock()
	return cmd.Process.Pid, logName, nil
}

// StopJob sends SIGTERM to a running job process.
func (m *JobManager) StopJob(jobName string, pid int) error {
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		return err
	}
	m.mu.Lock()
	delete(m.running, jobName)
	m.lastResults[jobName] = lastRunResult{
		finishedAt: time.Now().Format(time.RFC3339),
		exitCode:   -15,
		status:     "stopped",
	}
	m.mu.Unlock()
	return nil
}
