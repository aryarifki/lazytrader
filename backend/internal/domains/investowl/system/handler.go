package system

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler serves /api/logs endpoints.
type Handler struct {
	Manager  *JobManager
	Cache    *cache.Store
	DB       *pgxpool.Pool
	Host     string
	Port     int
	CFZoneID string
}

// GetDailyLogs handles GET /api/logs/daily.
func (h *Handler) GetDailyLogs(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	logsMap := make(map[string]string, len(JobConfigs))
	for key, cfg := range JobConfigs {
		prefix := cfg.LogPrefix
		if prefix == "" {
			prefix = key
		}
		logsMap[key] = GetLatestLog(h.Manager.SMtrackerDir, prefix)
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"timestamp": time.Now().Format(time.RFC3339),
		"logs":      logsMap,
	})
}

// GetAllJobsStatus handles GET /api/logs/jobs/status.
func (h *Handler) GetAllJobsStatus(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	results := make(map[string]JobStatusItem, len(JobConfigs))
	for key := range JobConfigs {
		results[key] = h.Manager.CheckStatus(key)
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"timestamp": time.Now().Format(time.RFC3339),
		"jobs":      results,
	})
}

// TriggerJob handles POST /api/logs/jobs/{job_name}/run.
func (h *Handler) TriggerJob(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	jobName := chi.URLParam(r, "job_name")
	cfg, ok := JobConfigs[jobName]
	if !ok {
		httpx.WriteError(w, http.StatusNotFound, fmt.Sprintf("Job '%s' tidak terdaftar.", jobName))
		return
	}

	req := DefaultJobRunRequest()
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "Invalid JSON request body")
			return
		}
	}

	st := h.Manager.CheckStatus(jobName)
	if st.Status == "running" {
		pidVal := 0
		if st.PID != nil {
			pidVal = *st.PID
		}
		httpx.WriteError(w, http.StatusConflict, fmt.Sprintf("Job '%s' sedang berjalan (PID: %d). Harap tunggu hingga selesai atau hentikan proses.", jobName, pidVal))
		return
	}

	if jobName != "pipeline" {
		pStatus := h.Manager.CheckStatus("pipeline")
		if pStatus.Status == "running" {
			httpx.WriteError(w, http.StatusConflict, "Full Pipeline sedang berjalan. Tidak disarankan menjalankan individual job bersamaan.")
			return
		}
	}

	args := h.Manager.BuildJobCommand(jobName, cfg, req)
	if req.DryRun {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{
			"status":   "dry_run",
			"job":      jobName,
			"label":    cfg.Label,
			"command":  strings.Join(args, " "),
			"cmd_args": args,
			"cwd":      h.Manager.SMtrackerDir,
			"message":  "Simulasi berhasil! Sintaks perintah valid dan tidak ada proses yang dijalankan ke database.",
		})
		return
	}

	pid, logFile, err := h.Manager.StartJob(jobName, cfg, args)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal menjalankan subprocess: %v", err))
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":   "started",
		"job":      jobName,
		"pid":      pid,
		"log_file": logFile,
		"command":  strings.Join(args, " "),
	})
}

// StopJob handles POST /api/logs/jobs/{job_name}/stop.
func (h *Handler) StopJob(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	jobName := chi.URLParam(r, "job_name")
	st := h.Manager.CheckStatus(jobName)
	if st.PID == nil {
		httpx.WriteError(w, http.StatusBadRequest, fmt.Sprintf("Tidak ada proses aktif untuk job '%s'.", jobName))
		return
	}
	if err := h.Manager.StopJob(jobName, *st.PID); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal menghentikan PID %d: %v", *st.PID, err))
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status": "stopped",
		"job":    jobName,
		"pid":    *st.PID,
	})
}

// GetLogHistory handles GET /api/logs/history.
func (h *Handler) GetLogHistory(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	job := strings.TrimSpace(r.URL.Query().Get("job"))
	if job == "" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "Query parameter 'job' is required")
		return
	}
	files := ListHistoryFiles(h.Manager.SMtrackerDir, job)
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"job":   job,
		"files": files,
	})
}

// ViewSpecificLog handles GET /api/logs/view.
func (h *Handler) ViewSpecificLog(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	rawName := r.URL.Query().Get("filename")
	if strings.TrimSpace(rawName) == "" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "Query parameter 'filename' is required")
		return
	}
	cleanName, ok := SanitizeLogFilename(rawName)
	if !ok {
		httpx.WriteError(w, http.StatusBadRequest, "Nama file log tidak valid.")
		return
	}
	target := filepath.Join(logDir(h.Manager.SMtrackerDir), cleanName)
	info, err := os.Stat(target)
	if err != nil || info.IsDir() {
		httpx.WriteError(w, http.StatusNotFound, fmt.Sprintf("File log '%s' tidak ditemukan.", cleanName))
		return
	}
	content, lineCount, err := readLastLines(target, 1000)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Error membaca file %s: %v", cleanName, err))
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"filename":    cleanName,
		"lines_count": lineCount,
		"content":     content,
	})
}

// ClearLogs handles DELETE /api/logs/clear.
func (h *Handler) ClearLogs(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	dir := logDir(h.Manager.SMtrackerDir)
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{
			"status":        "success",
			"message":       "Direktori log tidak ditemukan.",
			"deleted_count": 0,
		})
		return
	}

	q := r.URL.Query()
	filename := strings.TrimSpace(q.Get("filename"))
	job := strings.TrimSpace(q.Get("job"))
	clearAll := false
	if raw := strings.TrimSpace(q.Get("clear_all")); raw != "" {
		b, _ := strconv.ParseBool(raw)
		clearAll = b
	}

	deleted := make([]string, 0)
	switch {
	case filename != "" && filename != "latest":
		cleanName, ok := SanitizeLogFilename(filename)
		if !ok {
			httpx.WriteError(w, http.StatusBadRequest, "Nama file log tidak valid.")
			return
		}
		target := filepath.Join(dir, cleanName)
		info, err := os.Stat(target)
		if err != nil || info.IsDir() {
			httpx.WriteError(w, http.StatusNotFound, fmt.Sprintf("File log '%s' tidak ditemukan.", cleanName))
			return
		}
		if err := os.Remove(target); err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal menghapus file %s: %v", cleanName, err))
			return
		}
		deleted = append(deleted, cleanName)

	case job != "":
		prefix := job
		if cfg, ok := JobConfigs[job]; ok && cfg.LogPrefix != "" {
			prefix = cfg.LogPrefix
		}
		for _, fp := range matchingLogFiles(dir, prefix) {
			if info, err := os.Stat(fp); err == nil && !info.IsDir() {
				if os.Remove(fp) == nil {
					deleted = append(deleted, filepath.Base(fp))
				}
			}
		}

	case clearAll:
		matches, _ := filepath.Glob(filepath.Join(dir, "*.log"))
		for _, fp := range matches {
			if info, err := os.Stat(fp); err == nil && !info.IsDir() {
				if os.Remove(fp) == nil {
					deleted = append(deleted, filepath.Base(fp))
				}
			}
		}

	default:
		httpx.WriteError(w, http.StatusBadRequest, "Parameter filename, job, atau clear_all harus disertakan.")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":        "success",
		"deleted_count": len(deleted),
		"deleted_files": deleted,
		"message":       fmt.Sprintf("Berhasil menghapus %d file log.", len(deleted)),
	})
}

// GetCacheStats handles GET /api/logs/cache/stats.
func (h *Handler) GetCacheStats(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	httpx.WriteJSON(w, http.StatusOK, h.Cache.Stats(r.Context()))
}

// PurgeCache handles POST /api/logs/cache/purge using SafePurge (never FLUSHDB).
func (h *Handler) PurgeCache(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	q := r.URL.Query()
	pattern := strings.TrimSpace(q.Get("pattern"))
	if pattern == "" {
		pattern = "*"
	}
	purgeCF := true
	if raw := strings.TrimSpace(q.Get("purge_cloudflare")); raw != "" {
		if b, err := strconv.ParseBool(raw); err == nil {
			purgeCF = b
		}
	}

	deleted := h.Cache.SafePurge(r.Context(), pattern)
	var cfRes map[string]any
	if purgeCF {
		cfRes = cache.PurgeCloudflare(r.Context(), h.CFZoneID, nil, true)
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":       "success",
		"deleted_keys": deleted,
		"pattern":      pattern,
		"cloudflare":   cfRes,
		"message":      "Universal cache purged successfully across Memory L1, Redis L2, and Cloudflare CDN.",
	})
}

// WarmupCache handles POST /api/logs/cache/warmup.
func (h *Handler) WarmupCache(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	TriggerWarmup(h.DB, h.Host, h.Port)
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":  "success",
		"message": "Pre-warming cache dimulai di background.",
	})
}
