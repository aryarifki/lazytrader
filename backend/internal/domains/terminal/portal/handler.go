package portal

import (
	"net/http"

	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/go-chi/chi/v5"
)

// Routes registers /api/v1/terminal/portal/* endpoints (explicitly excluding bei1persen).
func Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/status", GetStatus)
	r.Get("/entitlements", GetEntitlements)
	return r
}

func GetStatus(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":               "online",
		"system_name":          "LazyTrader Fusion Suite",
		"version":              "1.0.0-fusion",
		"active_systems_count": 6,
		"active_systems": []map[string]string{
			{"id": "mm-detector", "name": "MM Detector Pro", "tier": "pro_max", "status": "active"},
			{"id": "foreign-flow", "name": "Foreign Flow Pro", "tier": "explorer", "status": "active"},
			{"id": "foreign-flow-mobile", "name": "Foreign Flow Mobile", "tier": "explorer", "status": "active"},
			{"id": "ipo-hunter", "name": "IPO Hunter Pro", "tier": "standard", "status": "active"},
			{"id": "konglo-screener", "name": "Konglo Screener", "tier": "standard", "status": "active"},
			{"id": "mm-detector-lite", "name": "MM Detector Lite", "tier": "mobile_pro_max", "status": "active"},
		},
	})
}

func GetEntitlements(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"entitlements": map[string]map[string]string{
			"mm_detector":         {"tier": "pro_max", "expires": "2030-12-31T23:59:59.000Z", "status": "active"},
			"foreign_flow":        {"tier": "explorer", "expires": "2030-12-31T23:59:59.000Z", "status": "active"},
			"foreign_flow_mobile": {"tier": "explorer", "expires": "2030-12-31T23:59:59.000Z", "status": "active"},
			"ipo_hunter":          {"tier": "standard", "expires": "2030-12-31T23:59:59.000Z", "status": "active"},
			"konglo_screener":     {"tier": "standard", "expires": "2030-12-31T23:59:59.000Z", "status": "active"},
			"mm_detector_lite":    {"tier": "mobile_pro_max", "expires": "2030-12-31T23:59:59.000Z", "status": "active"},
		},
		"readInboxIds": []string{},
	})
}
