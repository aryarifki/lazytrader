package market_report

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var progressLineRe = regexp.MustCompile(`\[PROGRESS\]\s+(\d+)/(\d+)\s+\(([\d\.]+)%\)\s*-\s*(.*)`)

func defaultCommandRunner(ctx context.Context, bin string, args []string, extraEnv map[string]string, onLine func(string)) (int, error) {
	cmd := exec.CommandContext(ctx, bin, args...)
	env := os.Environ()
	for k, v := range extraEnv {
		env = append(env, k+"="+v)
	}
	cmd.Env = env

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return -1, err
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return -1, err
	}

	sc := bufio.NewScanner(stdout)
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 1024*1024)
	for sc.Scan() {
		if onLine != nil {
			onLine(sc.Text())
		}
	}

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), nil
		}
		return -1, err
	}
	return 0, nil
}

// StartGeneration triggers asynchronous report generation subprocess.
func (e *Engine) StartGeneration(session int, targetRPM float64, maxWorkers int) map[string]any {
	e.mu.Lock()
	if e.state.IsRunning {
		stCopy := e.state
		e.mu.Unlock()
		return map[string]any{
			"success": false,
			"message": "Pemindaian sedang berjalan. Harap tunggu hingga selesai.",
			"status":  stCopy,
		}
	}

	nowStr := getWIBNow()
	e.state.IsRunning = true
	e.state.Status = "running"
	e.state.Session = session
	e.state.Progress = 0
	e.state.Total = 962
	e.state.Pct = 0.0
	e.state.Error = nil
	e.state.StartedAt = &nowStr
	e.state.CompletedAt = nil
	e.state.Message = fmt.Sprintf("Memulai pemindaian universe untuk Laporan Sesi %d...", session)
	e.appendLogLocked(strings.Repeat("=", 60), true)
	e.appendLogLocked(fmt.Sprintf("🚀 Memulai generator Laporan Pasar Sesi %d (Rate: %g rpm, Workers: %d)", session, targetRPM, maxWorkers), true)
	e.mu.Unlock()

	go e.runScanSubprocess(session, targetRPM, maxWorkers)

	return map[string]any{
		"success": true,
		"message": fmt.Sprintf("Generator Laporan Sesi %d berhasil dimulai.", session),
		"session": session,
	}
}

func (e *Engine) runScanSubprocess(session int, targetRPM float64, maxWorkers int) {
	args := []string{
		"-u",
		e.scriptPath,
		"--session", strconv.Itoa(session),
		"--rate-limit", strconv.FormatFloat(targetRPM, 'f', -1, 64),
		"--workers", strconv.Itoa(maxWorkers),
	}
	extraEnv := map[string]string{
		"PYTHONPATH":               e.smtrackerSrcDir,
		"PLAYWRIGHT_BROWSERS_PATH": "/root/.cache/ms-playwright",
		"PYTHONUNBUFFERED":         "1",
	}
	e.addLog(fmt.Sprintf("Menjalankan subprocess: %s %s", e.pythonBin, strings.Join(args, " ")), true)

	var lastErrMsg string
	rc, err := e.runCmd(context.Background(), e.pythonBin, args, extraEnv, func(raw string) {
		line := strings.TrimSpace(raw)
		if line == "" {
			return
		}
		e.mu.Lock()
		defer e.mu.Unlock()
		e.appendLogLocked(line, true)
		if m := progressLineRe.FindStringSubmatch(line); len(m) == 5 {
			prog, _ := strconv.Atoi(m[1])
			tot, _ := strconv.Atoi(m[2])
			pct, _ := strconv.ParseFloat(m[3], 64)
			submsg := m[4]
			e.state.Progress = prog
			e.state.Total = tot
			e.state.Pct = pct
			e.state.Message = fmt.Sprintf("Memindai harga %d/%d emiten (%.1f%%)... %s", prog, tot, pct, submsg)
		} else if strings.HasPrefix(line, "[ERROR]") {
			msg := strings.TrimSpace(strings.TrimPrefix(line, "[ERROR]"))
			lastErrMsg = msg
			e.state.Error = &msg
		}
	})

	e.mu.Lock()
	defer e.mu.Unlock()
	if err != nil {
		errStr := err.Error()
		e.state.IsRunning = false
		e.state.Status = "error"
		e.state.Error = &errStr
		e.state.Message = fmt.Sprintf("Terjadi kesalahan: %s", errStr)
		e.appendLogLocked(fmt.Sprintf("❌ Exception: %s", errStr), true)
		return
	}
	if rc == 0 {
		doneAt := getWIBNow()
		e.loadLatestFromDiskLocked()
		e.state.IsRunning = false
		e.state.Status = "completed"
		e.state.Pct = 100.0
		e.state.CompletedAt = &doneAt
		e.state.Message = fmt.Sprintf("Laporan Sesi %d berhasil dibuat!", session)
		e.appendLogLocked(fmt.Sprintf("🎉 Selesai! Laporan Sesi %d berhasil dibuat.", session), true)
		return
	}

	errDetail := lastErrMsg
	if errDetail == "" {
		errDetail = fmt.Sprintf("Subprocess berhenti dengan exit code %d", rc)
	}
	e.state.IsRunning = false
	e.state.Status = "error"
	e.state.Error = &errDetail
	e.state.Message = fmt.Sprintf("Gagal membuat laporan: %s", errDetail)
	e.appendLogLocked(fmt.Sprintf("❌ Proses scan gagal (%s).", errDetail), true)
}

// TriggerAutoRenewToken triggers SMtracker Two-Tier automatic token renewal in background.
func (e *Engine) TriggerAutoRenewToken() map[string]any {
	e.mu.Lock()
	if e.state.IsRunning {
		e.mu.Unlock()
		return map[string]any{
			"success": false,
			"message": "Pemindaian atau proses pembaruan token sedang berjalan.",
		}
	}
	e.state.IsRunning = true
	e.state.Status = "running"
	e.state.Message = "Menjalankan pembaruan token otomatis via SMtracker (Two-Tier)..."
	e.appendLogLocked(strings.Repeat("=", 60), true)
	e.appendLogLocked("🔐 [SMtracker] Memulai Two-Tier token auto-renewal (metode backfill_daily)...", true)
	e.mu.Unlock()

	go func() {
		pyCode := fmt.Sprintf("import sys, os; sys.path.insert(0, %q); from idx_bandarmology import broker_api; ok = broker_api.ensure_token_valid(); sys.exit(0 if ok else 1)", e.smtrackerSrcDir)
		args := []string{"-u", "-c", pyCode}
		extraEnv := map[string]string{
			"PYTHONPATH":               e.smtrackerSrcDir,
			"PLAYWRIGHT_BROWSERS_PATH": "/root/.cache/ms-playwright",
			"PYTHONUNBUFFERED":         "1",
		}
		rc, err := e.runCmd(context.Background(), e.pythonBin, args, extraEnv, func(raw string) {
			if line := strings.TrimSpace(raw); line != "" {
				e.addLog(line, true)
			}
		})

		e.mu.Lock()
		defer e.mu.Unlock()
		e.state.IsRunning = false
		if err != nil {
			errStr := err.Error()
			e.state.Status = "error"
			e.state.Error = &errStr
			e.state.Message = fmt.Sprintf("Kesalahan internal auto-renew: %s", errStr)
			e.appendLogLocked(fmt.Sprintf("❌ [SMtracker] Exception saat auto-renew: %s", errStr), true)
			return
		}
		if rc == 0 {
			e.state.Status = "idle"
			e.state.Error = nil
			e.state.Message = "Token Stockbit berhasil diperbarui otomatis oleh SMtracker!"
			e.appendLogLocked("✨ [SMtracker] Berhasil memperbarui dan memvalidasi token otomatis.", true)
			return
		}
		errMsg := "Auto-renew SMtracker gagal atau batas waktu autentikasi HP habis."
		e.state.Status = "error"
		e.state.Error = &errMsg
		e.state.Message = errMsg
		e.appendLogLocked("❌ [SMtracker] Pembaruan token otomatis gagal atau waktu habis.", true)
	}()

	return map[string]any{
		"success": true,
		"message": "Pembaruan token otomatis via SMtracker dimulai. Periksa notifikasi di aplikasi HP Anda jika diminta.",
	}
}
