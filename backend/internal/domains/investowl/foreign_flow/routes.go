package foreign_flow

import (
	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/forensic"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Routes mounts /api/foreign-flow endpoints. Static /sectors/* routes precede /{ticker}.
func Routes(db *pgxpool.Pool, cacheStore *cache.Store, forensicSvc *forensic.Service) chi.Router {
	h := &Handler{
		DB:    db,
		Cache: cacheStore,
		Sectors: &SectorService{
			DB:       db,
			Cache:    cacheStore,
			Forensic: forensicSvc,
		},
	}

	r := chi.NewRouter()
	r.Get("/sectors/summary", h.GetSectorsSummary)
	r.Get("/sectors/history", h.GetSectorsHistory)
	r.Get("/{ticker}", h.GetTickerForensic)
	return r
}
