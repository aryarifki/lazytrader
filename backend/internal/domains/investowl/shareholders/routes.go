package shareholders

import (
	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Routes mounts all 6 /api/bandar/shareholders endpoints with static paths before /{ticker}.
func Routes(db *pgxpool.Pool, store *cache.Store) chi.Router {
	h := &Handler{
		Svc:   &Service{DB: db},
		Cache: store,
	}

	r := chi.NewRouter()
	r.Get("/top-cross-holdings", h.GetTopCrossHoldings)
	r.Get("/search", h.SearchShareholders)
	r.Get("/groups", h.GetConglomerateGroups)
	r.Get("/network-graph", h.GetOwnershipNetworkGraph)
	r.Post("/refresh", h.RefreshCache)
	r.Get("/{ticker}", h.GetStockShareholders)
	return r
}
