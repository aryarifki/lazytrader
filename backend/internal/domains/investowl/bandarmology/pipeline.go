package bandarmology

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
)

func (h *Handler) pythonExecutable() (pyBin, invDir string) {
	invDir = strings.TrimSpace(h.InvDir)
	if invDir == "" {
		invDir = "/root/SMtracker"
	}
	pyBin = strings.TrimSpace(h.PythonBin)
	if pyBin == "" {
		cand := filepath.Join(invDir, "venv", "bin", "python3")
		if _, err := os.Stat(cand); err == nil {
			pyBin = cand
		} else {
			pyBin = "python3"
		}
	}
	return pyBin, invDir
}

func (h *Handler) pythonEnv(invDir string) []string {
	return append(os.Environ(), "PYTHONPATH="+filepath.Join(invDir, "src")+":"+invDir)
}

// RunPipeline handles POST /api/bandar/pipeline/run (approved session required).
func (h *Handler) RunPipeline(w http.ResponseWriter, r *http.Request) {
	universeMode := strings.TrimSpace(r.URL.Query().Get("universe_mode"))
	if universeMode == "" {
		universeMode = "watchlist"
	}
	pyBin, invDir := h.pythonExecutable()

	cmd := exec.Command(pyBin, "-c",
		"import sys; from idx_bandarmology import pipeline; pipeline.run(universe_mode=sys.argv[1])",
		universeMode)
	cmd.Dir = invDir
	cmd.Env = h.pythonEnv(invDir)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	go func() { _ = cmd.Wait() }()

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":        "started",
		"universe_mode": universeMode,
	})
}

// BackfillPipeline handles POST /api/bandar/pipeline/backfill (approved session required).
func (h *Handler) BackfillPipeline(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	tickers := strings.TrimSpace(q.Get("tickers"))
	if tickers == "" {
		tickers = "BBCA"
	}
	start := strings.TrimSpace(q.Get("start"))
	if start == "" {
		start = "2024-01-01"
	}
	end := strings.TrimSpace(q.Get("end"))
	pyBin, invDir := h.pythonExecutable()

	cmd := exec.Command(pyBin, "-c",
		"import sys; from idx_bandarmology import pipeline; pipeline.backfill_broker_history(tickers=sys.argv[1].split(','), start=sys.argv[2], end=(sys.argv[3] or None))",
		tickers, start, end)
	cmd.Dir = invDir
	cmd.Env = h.pythonEnv(invDir)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	go func() { _ = cmd.Wait() }()

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":  "started",
		"tickers": strings.Split(tickers, ","),
		"start":   start,
	})
}

// RefreshUniverse handles POST /api/bandar/universe/refresh.
func (h *Handler) RefreshUniverse(w http.ResponseWriter, r *http.Request) {
	pyBin, invDir := h.pythonExecutable()
	cmd := exec.Command(pyBin, "-c",
		"from idx_bandarmology import universe; print(universe.refresh_master_tickers(force=True))")
	cmd.Dir = invDir
	cmd.Env = h.pythonEnv(invDir)
	out, err := cmd.Output()
	if err != nil {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{
			"status": "error",
			"count":  0,
			"error":  err.Error(),
		})
		return
	}
	lines := strings.Fields(strings.TrimSpace(string(out)))
	cnt := 0
	if len(lines) > 0 {
		cnt, _ = strconv.Atoi(lines[len(lines)-1])
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":  "success",
		"count":   cnt,
		"message": fmt.Sprintf("Berhasil memperbarui %d emiten dari IDX.", cnt),
	})
}

// ClearBandarCache handles POST /api/bandar/clear-cache and GET /api/bandar/refresh-cache.
func (h *Handler) ClearBandarCache(w http.ResponseWriter, r *http.Request) {
	purgeCF := true
	if raw := strings.TrimSpace(r.URL.Query().Get("purge_cloudflare")); raw != "" {
		if b, err := strconv.ParseBool(raw); err == nil {
			purgeCF = b
		}
	}
	cleared := 0
	if h.Cache != nil {
		cleared = h.Cache.SafePurge(r.Context(), "bandar*")
	}
	var cfRes map[string]any
	if purgeCF {
		cfRes = cache.PurgeCloudflare(r.Context(), h.CFZoneID, nil, true)
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":             "ok",
		"message":            "Bandarmology cache cleared successfully and synchronized with Cloudflare Edge.",
		"cleared_redis_keys": cleared,
		"cloudflare_sync":    cfRes,
	})
}
