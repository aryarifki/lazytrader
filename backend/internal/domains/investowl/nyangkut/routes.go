package nyangkut

import (
	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Mount registers all Dashboard Nyangkut routes and their aliases on the root router.
func Mount(r chi.Router, db *pgxpool.Pool, store *cache.Store) {
	h := &Handler{
		Svc:   &Service{DB: db},
		Cache: store,
	}

	r.Get("/api/v1/stock/search", h.SearchStocks)
	r.Get("/api/v1/stock/{ticker}/nyangkut-tracker", h.GetStockTracker)
	r.Get("/api/nyangkut/stock/{ticker}", h.GetStockTracker)
	r.Get("/api/v1/nyangkut/leaderboard", h.GetLeaderboard)
	r.Get("/api/nyangkut/leaderboard", h.GetLeaderboard)
	r.Get("/api/v1/nyangkut/summary", h.GetSummary)
	r.Get("/api/nyangkut/summary", h.GetSummary)
	r.Get("/api/nyangkut/search", h.SearchStocks)
	r.Post("/api/nyangkut/refresh", h.RefreshCache)
}
