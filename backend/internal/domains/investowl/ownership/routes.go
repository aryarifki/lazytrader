package ownership

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ownershipCacheKey(prefix, raw string) string {
	h := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%s:%s", prefix, hex.EncodeToString(h[:])[:16])
}

// Routes mounts all 14 /api/v1/ownership endpoints with static paths before parameter routes.
func Routes(db *pgxpool.Pool, store *cache.Store) chi.Router {
	h := &Handler{
		DB:    db,
		Cache: store,
	}

	r := chi.NewRouter()
	r.Get("/stock/{ticker}", h.GetStockOwnership)

	r.Get("/investors", h.SearchInvestors)
	r.Get("/investors/detail", h.GetInvestorDetail)
	r.Get("/investors/ct", h.SearchCTInvestors)
	r.Get("/investors/ct/detail", h.GetCTInvestorDetail)

	r.Get("/market", h.GetMarketOverview)

	r.Get("/groups", h.GetGroups)
	r.Post("/groups/sync-catalog", h.SyncGroupsCatalog)
	r.Get("/groups/{group_id}", h.GetGroupDetail)

	r.Get("/explorer", h.GetExplorer)
	r.Get("/explorer/ct", h.GetCTExplorer)

	r.Get("/staging/overview", h.GetStagingOverview)
	r.Get("/staging/stock/{ticker}", h.GetStagingStockProfile)
	r.Get("/staging/parity/{ticker}", h.GetStagingParity)
	return r
}
