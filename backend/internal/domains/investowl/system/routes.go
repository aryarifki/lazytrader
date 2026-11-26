package system

import (
	"github.com/aryarifki/lazytrader/backend/internal/core/auth"
	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Deps groups dependencies for /api/logs routes.
type Deps struct {
	DB            *pgxpool.Pool
	Cache         *cache.Store
	TerminalCache *cache.Store
	Guard         *auth.Guard
	SMtrackerDir  string
	Host          string
	Port          int
	CFZoneID      string
}

// Routes mounts /api/logs endpoints matching InvestOwl routers/system.py auth gates.
func Routes(d Deps) chi.Router {
	h := &Handler{
		Manager:       NewJobManager(d.SMtrackerDir),
		Cache:         d.Cache,
		TerminalCache: d.TerminalCache,
		DB:            d.DB,
		Host:          d.Host,
		Port:          d.Port,
		CFZoneID:      d.CFZoneID,
	}

	r := chi.NewRouter()
	r.Group(func(admin chi.Router) {
		admin.Use(d.Guard.RequireAdmin)
		admin.Get("/daily", h.GetDailyLogs)
		admin.Get("/jobs/status", h.GetAllJobsStatus)
		admin.Post("/jobs/{job_name}/run", h.TriggerJob)
		admin.Post("/jobs/{job_name}/stop", h.StopJob)
		admin.Get("/history", h.GetLogHistory)
		admin.Get("/view", h.ViewSpecificLog)
		admin.Delete("/clear", h.ClearLogs)
		admin.Get("/cache/stats", h.GetCacheStats)
	})
	r.Group(func(approved chi.Router) {
		approved.Use(d.Guard.RequireApproved)
		approved.Post("/cache/purge", h.PurgeCache)
		approved.Post("/cache/warmup", h.WarmupCache)
	})
	return r
}
