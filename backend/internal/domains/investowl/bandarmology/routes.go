package bandarmology

import (
	"github.com/aryarifki/lazytrader/backend/internal/core/auth"
	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Deps bundles dependencies for /api/bandar routes.
type Deps struct {
	DB       *pgxpool.Pool
	Cache    *cache.Store
	Guard    *auth.Guard
	CFZoneID string
}

// Routes mounts all 30 InvestOwl /api/bandar endpoints, registering static segments before wildcards.
func Routes(d Deps) chi.Router {
	h := &Handler{
		Repo:     &Repository{DB: d.DB},
		Cache:    d.Cache,
		CFZoneID: d.CFZoneID,
	}

	r := chi.NewRouter()

	// Universe, System, Summary & Screener
	r.Get("/universe", h.GetUniverse)
	r.Post("/universe/refresh", h.RefreshUniverse)
	r.Get("/universe/{mode}", h.GetUniverseByMode)
	r.Get("/tickers", h.GetMasterTickers)
	r.Get("/system-status", h.GetSystemStatus)
	r.Get("/daily-summary", h.GetDailySummary)
	r.Get("/dates/{ticker}", h.GetTickerDates)

	// Pipeline & Cache
	r.Group(func(protected chi.Router) {
		protected.Use(d.Guard.RequireApproved)
		protected.Post("/pipeline/run", h.RunPipeline)
		protected.Post("/pipeline/backfill", h.BackfillPipeline)
	})
	r.Post("/clear-cache", h.ClearBandarCache)
	r.Get("/refresh-cache", h.ClearBandarCache)

	// Screener & Validation (static /validation/broker-scan before /validation/{ticker})
	r.Get("/screener", h.GetScreener)
	r.Get("/screener-v2", h.GetScreenerV2)
	r.Get("/validation/broker-scan", h.GetValidationBrokerScan)
	r.Get("/validation/{ticker}", h.GetValidation)
	r.Get("/validation-v2/{ticker}", h.GetValidationV2)

	// Detail, Broker-Flow & Causality
	r.Get("/detail/{ticker}", h.GetTickerDetail)
	r.Get("/broker-flow/{ticker}", h.GetBrokerFlowDetail)
	r.Get("/causality/{ticker}", h.GetCausalityInsight)

	// Stock-specific endpoints
	r.Get("/stocks/{ticker}/metrics", h.GetStockMetrics)
	r.Get("/stocks/{ticker}/smart-flow", h.GetStockSmartFlow)
	r.Get("/stocks/{ticker}/broker-compare", h.GetStockBrokerCompare)
	r.Get("/stocks/{ticker}/causality", h.GetStockCausality)
	r.Get("/stocks/{ticker}/event-study", h.GetStockEventStudy)
	r.Get("/stocks/{ticker}/raw", h.GetStockRaw)
	r.Get("/stocks/{ticker}/raw-tables", h.GetStockRawTables)
	r.Get("/stocks/{ticker}/broker-profiles", h.GetStockBrokerProfiles)
	r.Get("/stocks/{ticker}/price-performance", h.GetStockPricePerformance)
	r.Get("/stocks/{ticker}/broker-distribution", h.GetStockBrokerDistribution)
	r.Get("/stocks/{ticker}/retail-flow", h.GetStockRetailFlow)

	return r
}
