package signal

import (
	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/aryarifki/lazytrader/backend/internal/core/flags"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Routes mounts /api/signal endpoints.
func Routes(db *pgxpool.Pool, cacheStore *cache.Store, flagStore *flags.Store) chi.Router {
	h := &Handler{
		Repo:  &Repository{DB: db},
		Cache: cacheStore,
		Flags: flagStore,
	}

	r := chi.NewRouter()
	r.Post("/refresh", h.Refresh)
	r.Get("/daily", h.GetDaily)
	r.Get("/backtest/stats", h.GetBacktestStats)
	r.Get("/backtest/history", h.GetBacktestHistory)
	r.Get("/outcomes/summary", h.GetOutcomesSummary)
	r.Get("/outcomes/export", h.ExportOutcomes)
	r.Get("/active", h.GetActive)
	r.Get("/confluence-screener", h.GetConfluenceScreener)
	return r
}
