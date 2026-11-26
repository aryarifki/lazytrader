package konglo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/aryarifki/lazytrader/backend/internal/domains/terminal/shared"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	shortDaysID = [7]string{"Min", "Sen", "Sel", "Rab", "Kam", "Jum", "Sab"}
	fullDaysID  = [7]string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	monthsID    = [12]string{"Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"}
)

// Handler exposes the 4 Konglo Screener HTTP endpoints.
type Handler struct {
	Repo  *Repository
	Cache *cache.Store
}

// Routes mounts /api/v1/terminal/konglo/*.
func Routes(db *pgxpool.Pool, store *cache.Store) chi.Router {
	h := &Handler{Repo: &Repository{DB: db}, Cache: store}
	r := chi.NewRouter()
	r.Get("/report", h.GetReport)
	r.Get("/groups", h.GetGroups)
	r.Get("/summary", h.GetSummary)
	r.Get("/group/{group_name}", h.GetGroupDetail)
	return r
}

func (h *Handler) GetReport(w http.ResponseWriter, r *http.Request) {
	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, "konglo_full_report", 60*time.Second, func(ctx context.Context) (any, int, string, error) {
		stocks, err := h.Repo.ListReportStocks(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		groupOrder, groupsMap, err := h.Repo.ListGroupsUnordered(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		ts, err := h.Repo.GetReportTimestamps(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}

		report := CalculateKongloReport(stocks, groupOrder, groupsMap)
		if ts.IHSGSummary != nil {
			closeVal := 5096.83
			if v, ok := toFloat(ts.IHSGSummary["close"]); ok {
				closeVal = v
			}
			pctVal := 0.79
			if v, ok := toFloat(ts.IHSGSummary["changePct"]); ok {
				pctVal = v
			}
			report.Summary.IHSGClose = shared.FormatComma2(closeVal)
			report.Summary.IHSGPct = fmt.Sprintf("%+.2f", pctVal)
		}

		report.Metadata = buildReportMetadata(ts.MaxTradeDate, ts.MaxUpdated, report.Summary.TotalActiveStocks, report.Summary.TotalStocks)

		return map[string]any{
			"success": true,
			"report":  report,
		}, http.StatusOK, "", nil
	})
	if err != nil {
		httpx.WriteError(w, status, detail)
		return
	}
	shared.WriteCachedTerminalJSON(w, body)
}

func (h *Handler) GetGroups(w http.ResponseWriter, r *http.Request) {
	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, "konglo_groups", 600*time.Second, func(ctx context.Context) (any, int, string, error) {
		groups, err := h.Repo.ListGroupsOrdered(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		return map[string]any{
			"success":      true,
			"total_groups": len(groups),
			"groups":       groups,
		}, http.StatusOK, "", nil
	})
	if err != nil {
		httpx.WriteError(w, status, detail)
		return
	}
	shared.WriteCachedTerminalJSON(w, body)
}

func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	limit, err := shared.ParseIntQueryRange(r, "limit", 200, 1, 2000)
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	cacheKey := fmt.Sprintf("konglo_summary:group=%s:limit=%d", group, limit)

	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, cacheKey, 120*time.Second, func(ctx context.Context) (any, int, string, error) {
		stocks, err := h.Repo.ListStockSummary(ctx, group, limit)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		return map[string]any{
			"success": true,
			"total":   len(stocks),
			"stocks":  stocks,
		}, http.StatusOK, "", nil
	})
	if err != nil {
		httpx.WriteError(w, status, detail)
		return
	}
	shared.WriteCachedTerminalJSON(w, body)
}

func (h *Handler) GetGroupDetail(w http.ResponseWriter, r *http.Request) {
	rawName := chi.URLParam(r, "group_name")
	groupName, err := url.PathUnescape(rawName)
	if err != nil {
		groupName = rawName
	}

	cacheKey := "konglo_group_detail:" + groupName
	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, cacheKey, 120*time.Second, func(ctx context.Context) (any, int, string, error) {
		group, err := h.Repo.GetGroupByPattern(ctx, groupName)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		if group == nil {
			return nil, http.StatusNotFound, fmt.Sprintf("Group '%s' not found.", groupName), nil
		}

		members, err := h.Repo.ListGroupMembers(ctx, group.Tickers)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}

		var totalMC, totalVal, pctSum float64
		for _, m := range members {
			if m.MarketCap != nil {
				totalMC += *m.MarketCap
			}
			if m.Val != nil {
				totalVal += *m.Val
			}
			last := floatOrZero(m.Last)
			prev := floatOrZero(m.Prev)
			if prev > 0 {
				pctSum += (last - prev) / prev * 100.0
			}
		}
		avgChg := 0.0
		if len(members) > 0 {
			avgChg = pctSum / float64(len(members))
		}

		return map[string]any{
			"success":            true,
			"group_name":         group.GroupName,
			"total_members":      len(members),
			"total_market_cap":   totalMC,
			"total_value_traded": totalVal,
			"average_change_pct": shared.PyRound(avgChg, 2),
			"members":            members,
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

func buildReportMetadata(maxTradeDate, maxUpdated *time.Time, totalActive, totalRegistered int) *ReportMetadata {
	var targetDate time.Time
	switch {
	case maxTradeDate != nil:
		targetDate = *maxTradeDate
	case maxUpdated != nil:
		targetDate = *maxUpdated
	default:
		targetDate = time.Date(2026, 11, 2, 0, 0, 0, 0, time.UTC)
	}

	dow := int(targetDate.Weekday())
	dName := shortDaysID[dow]
	fdName := fullDaysID[dow]
	mName := monthsID[int(targetDate.Month())-1]

	dateDisplay := fmt.Sprintf("%s, %02d %s %04d", fdName, targetDate.Day(), mName, targetDate.Year())
	timeDisplay := "16:00"
	lastUpdatedStr := fmt.Sprintf("%s, %02d/%02d/%04d 16:00 WIB", dName, targetDate.Day(), int(targetDate.Month()), targetDate.Year())

	tsStr := targetDate.Format("2006-01-02") + "T16:00:00+07:00"
	if maxUpdated != nil {
		tsStr = shared.FormatISO8601UTC(*maxUpdated)
	}

	return &ReportMetadata{
		Timestamp:      tsStr,
		TradeDate:      targetDate.Format("2006-01-02"),
		DateDisplay:    dateDisplay,
		TimeDisplay:    timeDisplay,
		LastUpdatedStr: lastUpdatedStr,
		SourceStatus:   "ok",
		SourceTitle:    fmt.Sprintf("Data live aktif — %d saham aktif (dari %d emiten)", totalActive, totalRegistered),
		SourceSub:      fmt.Sprintf("Data per penutupan %s | Auto-update setiap 5 menit", lastUpdatedStr),
	}
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
