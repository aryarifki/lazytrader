package foreign_flow_pro

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/aryarifki/lazytrader/backend/internal/domains/terminal/shared"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler exposes the 3 Foreign Flow Pro HTTP endpoints.
type Handler struct {
	Repo  *Repository
	Cache *cache.Store
}

// Routes mounts /api/v1/terminal/foreign-flow/*.
func Routes(db *pgxpool.Pool, store *cache.Store) chi.Router {
	h := &Handler{Repo: &Repository{DB: db}, Cache: store}
	r := chi.NewRouter()
	r.Get("/screener", h.GetScreener)
	r.Get("/series/{ticker}", h.GetSeries)
	r.Get("/ihsg", h.GetIHSG)
	return r
}

func (h *Handler) GetScreener(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	universe := q.Get("universe")
	if universe == "" {
		universe = "ff"
	}
	window, err := shared.ParseIntQueryRange(r, "window", 40, 1, 10000)
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	fpBadge := q.Get("fp_badge")
	transBucket := q.Get("trans_bucket")
	badge := q.Get("badge")
	limit, err := shared.ParseIntQueryRange(r, "limit", 200, 1, 1000)
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// FastAPI @cached("ff_screener", 120) uses positional/sorted kwargs equivalent.
	cacheKey := fmt.Sprintf("ff_screener:universe=%s:window=%d:fp_badge=%s:trans_bucket=%s:badge=%s:limit=%d",
		universe, window, fpBadge, transBucket, badge, limit)

	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, cacheKey, 120*time.Second, func(ctx context.Context) (any, int, string, error) {
		rows, err := h.Repo.ListScreener(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}

		results := make([]ScreenerItem, 0, len(rows))
		accumKey := fmt.Sprintf("%sAccum%dd", universe, window)
		priceChgKey := fmt.Sprintf("priceChange%dd", window)

		for _, row := range rows {
			if row.Ticker == "IHSG_Aggregate" {
				continue
			}
			summary := row.ScreenerSummary
			if summary == nil {
				summary = map[string]any{}
			}
			badges := row.Badges
			if badges == nil {
				badges = []string{}
			}

			if badge != "" && !containsStr(badges, badge) {
				continue
			}

			curFP, fpValid := mapOptionalString(summary, "fpBadge", "LFP")
			if fpBadge != "" && fpBadge != "ALL" && (!fpValid || curFP == nil || *curFP != fpBadge) {
				continue
			}

			avgTrans := mapFloat(summary, "avgTransValue10d", 0.0)
			if transBucket == "gt35" && avgTrans < 35.0 {
				continue
			}
			if transBucket == "10to35" && !(10.0 <= avgTrans && avgTrans <= 35.0) {
				continue
			}
			if transBucket == "lt10" && avgTrans >= 10.0 {
				continue
			}

			var accumVal float64
			switch universe {
			case "ff", "fco", "kf", "xl", "semiritel":
				accumVal = mapFloat(summary, accumKey, 0.0)
			}

			chgVal := 0.0
			if rawChg, exists := summary[priceChgKey]; exists {
				if f, ok := toFloat(rawChg); ok {
					chgVal = f
				}
			} else {
				chgVal = mapFloat(summary, "changePct", 0.0)
			}

			name := ""
			if row.Name != nil && *row.Name != "" {
				name = *row.Name
			} else {
				name = mapString(summary, "name", "")
			}

			results = append(results, ScreenerItem{
				Ticker:              row.Ticker,
				Name:                name,
				Close:               mapFloat(summary, "close", 0.0),
				ChangePct:           mapFloat(summary, "changePct", 0.0),
				FPBadge:             curFP,
				FPScore:             mapFloat(summary, "fpScoreLatest", 0.0),
				AvgTransValue10d:    avgTrans,
				AccumVal:            accumVal,
				PriceChangeWindow:   chgVal,
				Badges:              badges,
				DailyClose:          mapArrayOrEmpty(summary, "dailyClose"),
				DailyValue:          mapArrayOrEmpty(summary, "dailyValue"),
				DailyFFValue:        mapArrayOrEmpty(summary, "dailyFFValue"),
				DailyFCoValue:       mapArrayOrEmpty(summary, "dailyFCoValue"),
				DailyKFValue:        mapArrayOrEmpty(summary, "dailyKFValue"),
				DailyXLValue:        mapArrayOrEmpty(summary, "dailyXLValue"),
				DailySemiRitelValue: mapArrayOrEmpty(summary, "dailySemiRitelValue"),
				MaxDailyFF60d:       summary["maxDailyFF60d"],
				MinDailyFF60d:       summary["minDailyFF60d"],
				Summary:             summary,
			})
		}

		sort.SliceStable(results, func(i, j int) bool {
			return math.Abs(results[i].AccumVal) > math.Abs(results[j].AccumVal)
		})
		if len(results) > limit {
			results = results[:limit]
		}

		return map[string]any{
			"success":  true,
			"universe": universe,
			"window":   window,
			"total":    len(results),
			"screener": results,
		}, http.StatusOK, "", nil
	})
	if err != nil {
		httpx.WriteError(w, status, detail)
		return
	}
	shared.WriteCachedTerminalJSON(w, body)
}

func (h *Handler) GetSeries(w http.ResponseWriter, r *http.Request) {
	tickerClean := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	limit, err := shared.ParseIntQueryRange(r, "limit", 250, 10, 1500)
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	cacheKey := "ff_ticker_series:" + tickerClean + ":limit=" + strconv.Itoa(limit)
	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, cacheKey, 60*time.Second, func(ctx context.Context) (any, int, string, error) {
		meta, err := h.Repo.GetTickerMeta(ctx, tickerClean)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		points, err := h.Repo.GetTickerDaily(ctx, tickerClean, limit)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		if len(points) == 0 {
			return nil, http.StatusNotFound, fmt.Sprintf("No Foreign Flow data found for '%s'.", tickerClean), nil
		}

		name := tickerClean
		badges := []string{}
		if meta != nil {
			if meta.Name != nil && *meta.Name != "" {
				name = *meta.Name
			} else if meta.Name == nil {
				// FastAPI: meta["name"] if meta else ticker_clean -> returns null if meta exists with NULL name!
				name = ""
			}
			if meta.Badges != nil {
				badges = meta.Badges
			}
		}

		var nameVal any = name
		if meta != nil && meta.Name == nil {
			nameVal = nil
		}

		analytics := ProcessForeignFlowAnalytics(points)
		return map[string]any{
			"success":      true,
			"ticker":       tickerClean,
			"name":         nameVal,
			"badges":       badges,
			"total_points": len(points),
			"analytics":    analytics,
		}, http.StatusOK, "", nil
	})
	if err != nil {
		httpx.WriteError(w, status, detail)
		return
	}
	if status != http.StatusOK {
		httpx.WriteError(w, status, detail)
		return
	}
	shared.WriteCachedTerminalJSON(w, body)
}

func (h *Handler) GetIHSG(w http.ResponseWriter, r *http.Request) {
	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, "ff_ihsg_summary", 300*time.Second, func(ctx context.Context) (any, int, string, error) {
		payload, updatedAt, err := h.Repo.GetIHSGSummary(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}

		summary, ok := payload["summary"]
		if !ok || summary == nil {
			summary = map[string]any{}
		}
		points, ok := payload["points"]
		if !ok || points == nil {
			points = []any{}
		}
		totalPoints := payload["total_points"]
		if totalPoints == nil {
			totalPoints = 0
		}
		lastUpdate := payload["last_update"]
		if lastUpdate == nil {
			lastUpdate = ""
		}

		var updatedStr *string
		if updatedAt != nil {
			s := shared.FormatISO8601UTC(*updatedAt)
			updatedStr = &s
		}

		return map[string]any{
			"success":      true,
			"summary":      summary,
			"points":       points,
			"total_points": totalPoints,
			"last_update":  lastUpdate,
			"updated_at":   updatedStr,
		}, http.StatusOK, "", nil
	})
	if err != nil {
		httpx.WriteError(w, status, detail)
		return
	}
	shared.WriteCachedTerminalJSON(w, body)
}

func containsStr(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func mapOptionalString(m map[string]any, key, def string) (*string, bool) {
	v, exists := m[key]
	if !exists {
		s := def
		return &s, true
	}
	if v == nil {
		return nil, false
	}
	if s, ok := v.(string); ok {
		return &s, true
	}
	s := def
	return &s, true
}

func mapString(m map[string]any, key, def string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return def
	}
	if s, ok := v.(string); ok {
		return s
	}
	return def
}

func mapFloat(m map[string]any, key string, def float64) float64 {
	v, ok := m[key]
	if !ok || v == nil {
		return def
	}
	if f, ok := toFloat(v); ok {
		return f
	}
	return def
}

func mapArrayOrEmpty(m map[string]any, key string) any {
	v, ok := m[key]
	if !ok || v == nil {
		return []any{}
	}
	return v
}

func toFloat(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}
