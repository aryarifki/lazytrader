package mmd

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

// Handler exposes the 7 MM Detector Pro HTTP endpoints.
type Handler struct {
	Repo  *Repository
	Cache *cache.Store
}

// Routes mounts /api/v1/terminal/mmd/*.
func Routes(db *pgxpool.Pool, store *cache.Store) chi.Router {
	h := &Handler{Repo: &Repository{DB: db}, Cache: store}
	r := chi.NewRouter()
	r.Get("/tickers", h.GetTickers)
	r.Get("/explorer/scan", h.ScanExplorer)
	r.Get("/{ticker}/analysis", h.GetTickerAnalysis)
	r.Get("/{ticker}/t5", h.GetTickerT5)
	r.Get("/{ticker}/group-performance", h.GetGroupPerformance)
	r.Get("/{ticker}/broker-drilldown", h.GetBrokerDrilldown)
	r.Get("/{ticker}", h.GetTickerDetail)
	return r
}

func (h *Handler) GetTickers(w http.ResponseWriter, r *http.Request) {
	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, "mmd_tickers", 300*time.Second, func(ctx context.Context) (any, int, string, error) {
		rows, err := h.Repo.ListTickersMeta(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		return map[string]any{
			"success": true,
			"total":   len(rows),
			"tickers": rows,
		}, http.StatusOK, "", nil
	})
	if err != nil || (status != 0 && status != http.StatusOK) {
		httpx.WriteError(w, status, detail)
		return
	}
	shared.WriteCachedTerminalJSON(w, body)
}

func (h *Handler) ScanExplorer(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	category := q.Get("category")
	if category == "" {
		category = "normal"
	}
	criteria := q.Get("criteria")
	minVolVal, hasMinVol, err := parseOptionalFloatQuery(r, "min_volume_val")
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	maxVolVal, _, err := parseOptionalFloatQuery(r, "max_volume_val")
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	priceAction := q.Get("price_action")
	limit, err := shared.ParseIntQueryRange(r, "limit", 50, 1, 200)
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	cacheKey := fmt.Sprintf(
		"mmd_explorer_scan:category=%s:criteria=%s:min_volume_val=%s:max_volume_val=%s:price_action=%s:limit=%d",
		category, criteria, q.Get("min_volume_val"), q.Get("max_volume_val"), priceAction, limit,
	)
	_ = maxVolVal

	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, cacheKey, 60*time.Second, func(ctx context.Context) (any, int, string, error) {
		if category == "bigPlayer" {
			rows, err := h.Repo.ListBigPlayerT5Summary(ctx, limit*2)
			if err != nil {
				return nil, http.StatusInternalServerError, err.Error(), err
			}
			results := make([]ExplorerBigPlayerItem, 0, limit)
			for _, row := range rows {
				var lastVol int64
				if row.LastVolume != nil {
					lastVol = *row.LastVolume
				}
				var lastClose float64
				if row.LastClose != nil {
					lastClose = *row.LastClose
				}
				lastTavg := lastClose
				if row.LastTavg != nil && *row.LastTavg != 0 {
					lastTavg = *row.LastTavg
				}
				dailyVal := float64(lastVol) * lastTavg
				if hasMinVol && minVolVal != 0 && dailyVal < minVolVal {
					continue
				}
				results = append(results, ExplorerBigPlayerItem{
					Ticker:      row.Ticker,
					LastClose:   lastClose,
					LastVolume:  lastVol,
					PriceChange: 0.0,
					Badge:       Badge{Label: "Big Player Active", Cls: "big-acc", Color: "#2FD98A"},
					Detail:      fmt.Sprintf("%d Transaksi T5 (Terakhir: %s)", row.T5Count, row.LastT5Date),
				})
				if len(results) >= limit {
					break
				}
			}
			return map[string]any{
				"success":  true,
				"category": category,
				"total":    len(results),
				"results":  results,
			}, http.StatusOK, "", nil
		}

		tickers, err := h.Repo.ListLatestCandlesForExplorer(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		results := make([]ExplorerScanItem, 0, len(tickers))
		critUpper := strings.ToUpper(criteria)

		for _, t := range tickers {
			var lastClose float64
			if t.Close != nil {
				lastClose = *t.Close
			}
			prevClose := lastClose
			if t.Prev != nil && *t.Prev != 0 {
				prevClose = *t.Prev
			} else if t.Open != nil && *t.Open != 0 {
				prevClose = *t.Open
			}
			var lastVol int64
			if t.Volume != nil {
				lastVol = *t.Volume
			}
			lastTavg := lastClose
			if t.Tavg != nil && *t.Tavg != 0 {
				lastTavg = *t.Tavg
			}
			dailyVal := 0.0
			if lastVol > 0 {
				dailyVal = float64(lastVol) * lastTavg
			}
			if hasMinVol && minVolVal != 0 && dailyVal < minVolVal {
				continue
			}

			priceChange := 0.0
			if prevClose > 0 {
				priceChange = shared.PyRound(((lastClose-prevClose)/prevClose)*100.0, 2)
			}

			sumB := topNSumFloat(t.BuyBrokers, 5)
			sumS := topNSumFloat(t.SellBrokers, 5)
			diff := sumB - sumS
			denom := (sumB + sumS) / 2.0
			rawRatio := 0.0
			if denom > 0 {
				rawRatio = (diff / denom) * 100.0
			}
			ratio := shared.PyRound(rawRatio, 1)
			if math.Abs(ratio) == 0.0 {
				ratio = 0.0
			}
			badge := ClassifyBandarVal(ratio)

			if critUpper != "" && critUpper != "ALL" {
				if !matchesExplorerCriteria(category, critUpper, badge.Label, priceChange, ratio) {
					continue
				}
			}

			lotStr := formatCommaInt(lastVol / 100)
			results = append(results, ExplorerScanItem{
				Ticker:      t.Ticker,
				CompanyName: t.CompanyName,
				LastClose:   lastClose,
				LastVolume:  lastVol,
				PriceChange: priceChange,
				Ratio:       ratio,
				Badge:       badge,
				Detail:      fmt.Sprintf("Ratio: %+.1f%% | Chg: %+.2f%% | Vol: %s lot", ratio, priceChange, lotStr),
			})
		}

		sort.SliceStable(results, func(i, j int) bool {
			return math.Abs(results[i].Ratio) > math.Abs(results[j].Ratio)
		})
		if len(results) > limit {
			results = results[:limit]
		}
		return map[string]any{
			"success":  true,
			"category": category,
			"total":    len(results),
			"results":  results,
		}, http.StatusOK, "", nil
	})
	if err != nil || (status != 0 && status != http.StatusOK) {
		httpx.WriteError(w, status, detail)
		return
	}
	shared.WriteCachedTerminalJSON(w, body)
}

func (h *Handler) GetTickerAnalysis(w http.ResponseWriter, r *http.Request) {
	tickerClean := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	q := r.URL.Query()
	dateFrom := q.Get("date_from")
	dateTo := q.Get("date_to")
	autoRetail, err := parseBoolQuery(r, "auto_retail", true)
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	retailOverride := q.Get("retail_override")

	cacheKey := fmt.Sprintf("mmd_ticker_analysis:%s:date_from=%s:date_to=%s:auto_retail=%t:retail_override=%s",
		tickerClean, dateFrom, dateTo, autoRetail, retailOverride)

	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, cacheKey, 60*time.Second, func(ctx context.Context) (any, int, string, error) {
		meta, err := h.Repo.GetTickerMeta(ctx, tickerClean)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		if meta == nil {
			return nil, http.StatusNotFound, fmt.Sprintf("Ticker '%s' not found.", tickerClean), nil
		}
		allCandles, err := h.Repo.GetAllCandles(ctx, tickerClean)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		if len(allCandles) == 0 {
			return nil, http.StatusNotFound, fmt.Sprintf("No candle data for '%s'.", tickerClean), nil
		}

		filtered := filterCandlesByDate(allCandles, dateFrom, dateTo)
		t5Rows, err := h.Repo.ListT5Transactions(ctx, tickerClean, 50)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		groupPerf, err := h.Repo.GetGroupPerformance(ctx, tickerClean)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		macroRepo := h.Repo.GetLatestMacroRepoRatio(ctx)

		res, ok := ComputeMMDAnalysis(AnalysisInput{
			Ticker:             tickerClean,
			Candles:            filtered,
			AllHistoryCandles:  allCandles,
			T5Records:          t5Rows,
			Meta:               meta,
			GroupPerf:          groupPerf,
			KSEIMacroRepoRatio: macroRepo,
			AutoRetail:         autoRetail,
			RetailOverride:     retailOverride,
			Today:              time.Now(),
		})
		if !ok {
			return map[string]any{
				"success": true,
				"ticker":  tickerClean,
				"analysis": map[string]any{
					"error":             "No candle data",
					"group_performance": groupPerf,
				},
			}, http.StatusOK, "", nil
		}
		return map[string]any{
			"success":  true,
			"ticker":   tickerClean,
			"analysis": res,
		}, http.StatusOK, "", nil
	})
	if err != nil || (status != 0 && status != http.StatusOK) {
		httpx.WriteError(w, status, detail)
		return
	}
	shared.WriteCachedTerminalJSON(w, body)
}

func (h *Handler) GetTickerT5(w http.ResponseWriter, r *http.Request) {
	tickerClean := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	cacheKey := "mmd_ticker_t5:" + tickerClean
	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, cacheKey, 60*time.Second, func(ctx context.Context) (any, int, string, error) {
		txs, err := h.Repo.ListT5Transactions(ctx, tickerClean, 100)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		return map[string]any{
			"success":      true,
			"ticker":       tickerClean,
			"total":        len(txs),
			"transactions": txs,
			"t5_records":   txs,
		}, http.StatusOK, "", nil
	})
	if err != nil || (status != 0 && status != http.StatusOK) {
		httpx.WriteError(w, status, detail)
		return
	}
	shared.WriteCachedTerminalJSON(w, body)
}

func (h *Handler) GetGroupPerformance(w http.ResponseWriter, r *http.Request) {
	tickerClean := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	cacheKey := "mmd_group_performance:" + tickerClean
	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, cacheKey, 60*time.Second, func(ctx context.Context) (any, int, string, error) {
		perf, err := h.Repo.GetGroupPerformance(ctx, tickerClean)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		return map[string]any{
			"success":     true,
			"ticker":      tickerClean,
			"performance": perf,
		}, http.StatusOK, "", nil
	})
	if err != nil || (status != 0 && status != http.StatusOK) {
		httpx.WriteError(w, status, detail)
		return
	}
	shared.WriteCachedTerminalJSON(w, body)
}

func (h *Handler) GetBrokerDrilldown(w http.ResponseWriter, r *http.Request) {
	tickerClean := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	q := r.URL.Query()
	rawBroker := q.Get("broker")
	if rawBroker == "" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "Query parameter 'broker' is required")
		return
	}
	if len(rawBroker) < 2 || len(rawBroker) > 5 {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "Query parameter 'broker' must have length between 2 and 5")
		return
	}
	brokerClean := strings.ToUpper(strings.TrimSpace(rawBroker))
	dateFrom := q.Get("date_from")
	dateTo := q.Get("date_to")

	cacheKey := fmt.Sprintf("mmd_broker_drilldown:%s:broker=%s:date_from=%s:date_to=%s",
		tickerClean, brokerClean, dateFrom, dateTo)

	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, cacheKey, 60*time.Second, func(ctx context.Context) (any, int, string, error) {
		firmName, err := h.Repo.GetBrokerFirmName(ctx, brokerClean)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		rows, err := h.Repo.ListBrokerActivity(ctx, tickerClean, brokerClean, dateFrom, dateTo)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}

		daily := make([]BrokerDrilldownDay, 0, len(rows))
		var totBLot, totSLot, totBVal, totSVal float64
		for _, row := range rows {
			bLot, sLot, bVal, sVal := derefFloat(row.BuyLot), derefFloat(row.SellLot), derefFloat(row.BuyVal), derefFloat(row.SellVal)
			totBLot += bLot
			totSLot += sLot
			totBVal += bVal
			totSVal += sVal

			bAvg, sAvg := 0.0, 0.0
			if bLot > 0 {
				bAvg = shared.PyRound(bVal/(bLot*100.0), 1)
			}
			if sLot > 0 {
				sAvg = shared.PyRound(sVal/(sLot*100.0), 1)
			}
			daily = append(daily, BrokerDrilldownDay{
				Date: row.Date, BuyLot: bLot, BuyVal: bVal, BuyAvg: bAvg,
				SellLot: sLot, SellVal: sVal, SellAvg: sAvg,
				NetLot: bLot - sLot, NetVal: bVal - sVal,
			})
		}

		avgBuyPrice, avgSellPrice := 0.0, 0.0
		if totBLot > 0 {
			avgBuyPrice = shared.PyRound(totBVal/(totBLot*100.0), 1)
		}
		if totSLot > 0 {
			avgSellPrice = shared.PyRound(totSVal/(totSLot*100.0), 1)
		}

		var dfPtr, dtPtr *string
		if _, has := q["date_from"]; has {
			dfPtr = &dateFrom
		}
		if _, has := q["date_to"]; has {
			dtPtr = &dateTo
		}

		return map[string]any{
			"success":     true,
			"ticker":      tickerClean,
			"broker_code": brokerClean,
			"broker_name": firmName,
			"category":    ClassifyBrokerCategory(brokerClean),
			"date_from":   dfPtr,
			"date_to":     dtPtr,
			"summary": BrokerDrilldownSummary{
				TotalBuyLot: totBLot, TotalBuyValueM: shared.PyRound(totBVal/1e6, 2), TotalBuyValueB: shared.PyRound(totBVal/1e9, 2),
				AvgBuyPrice: avgBuyPrice,
				TotalSellLot: totSLot, TotalSellValueM: shared.PyRound(totSVal/1e6, 2), TotalSellValueB: shared.PyRound(totSVal/1e9, 2),
				AvgSellPrice: avgSellPrice,
				NetLot:       totBLot - totSLot,
				NetValueM:    shared.PyRound((totBVal-totSVal)/1e6, 2),
				NetValueB:    shared.PyRound((totBVal-totSVal)/1e9, 2),
				DaysActive:   len(daily),
			},
			"daily_breakdown": daily,
		}, http.StatusOK, "", nil
	})
	if err != nil || (status != 0 && status != http.StatusOK) {
		httpx.WriteError(w, status, detail)
		return
	}
	shared.WriteCachedTerminalJSON(w, body)
}

func (h *Handler) GetTickerDetail(w http.ResponseWriter, r *http.Request) {
	tickerClean := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	limit, err := shared.ParseIntQueryRange(r, "limit", 35, 1, 100)
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	cacheKey := "mmd_ticker_detail:" + tickerClean + ":limit=" + strconv.Itoa(limit)
	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, cacheKey, 60*time.Second, func(ctx context.Context) (any, int, string, error) {
		meta, err := h.Repo.GetTickerMeta(ctx, tickerClean)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		if meta == nil {
			return nil, http.StatusNotFound, fmt.Sprintf("Ticker '%s' not found in MM Detector.", tickerClean), nil
		}
		candles, err := h.Repo.GetLatestCandles(ctx, tickerClean, limit)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}

		var cumBandar int64
		series := make([]BandarSeriesPoint, 0, len(candles))
		for _, c := range candles {
			var netVol int64
			if c.NetVolume != nil {
				netVol = *c.NetVolume
			}
			cumBandar += netVol
			var cl *float64
			if c.Close != nil && *c.Close != 0 {
				v := *c.Close
				cl = &v
			}
			series = append(series, BandarSeriesPoint{
				Date: c.DateISO, NetVolume: netVol, CumulativeBandar: cumBandar, Close: cl,
			})
		}

		return map[string]any{
			"success":                    true,
			"ticker":                     tickerClean,
			"meta":                       meta,
			"total_candles":              len(candles),
			"candles":                    candles,
			"bandar_accumulation_series": series,
		}, http.StatusOK, "", nil
	})
	if err != nil || (status != 0 && status != http.StatusOK) {
		httpx.WriteError(w, status, detail)
		return
	}
	shared.WriteCachedTerminalJSON(w, body)
}

func topNSumFloat(m map[string]float64, n int) float64 {
	if len(m) == 0 {
		return 0.0
	}
	vals := make([]float64, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	sort.Slice(vals, func(i, j int) bool { return vals[i] > vals[j] })
	var sum float64
	for i := 0; i < n && i < len(vals); i++ {
		sum += vals[i]
	}
	return sum
}

func matchesExplorerCriteria(category, critUpper, badgeLbl string, priceChange, ratio float64) bool {
	if category == "normal" {
		switch critUpper {
		case "BIG ACC":
			return badgeLbl == "Big Acc" || badgeLbl == "Massive Acc"
		case "NORMAL ACC":
			return badgeLbl == "Normal Acc"
		case "SMALL ACC":
			return badgeLbl == "Small Acc"
		case "NETRAL":
			return badgeLbl == "Undetected"
		case "DISTRIBUSI":
			return strings.Contains(badgeLbl, "Dist")
		}
		return true
	}
	if category == "advanced" {
		switch critUpper {
		case "MARK UP":
			return priceChange > 0 && ratio >= 7.5
		case "MARK DOWN":
			return priceChange < 0 && ratio >= -7.5
		case "SIDEWAYS ACC":
			return math.Abs(priceChange) <= 1.5 && ratio >= 7.5
		case "SIDEWAYS DIST":
			return math.Abs(priceChange) <= 1.5 && ratio <= -7.5
		case "ACC":
			return ratio >= 15.0
		case "DIST":
			return ratio <= -15.0
		}
	}
	return true
}

func filterCandlesByDate(all []DailyCandle, dateFrom, dateTo string) []DailyCandle {
	if dateFrom != "" && dateTo != "" {
		out := make([]DailyCandle, 0, len(all))
		for _, c := range all {
			if c.DateISO >= dateFrom && c.DateISO <= dateTo {
				out = append(out, c)
			}
		}
		if len(out) == 0 {
			return all
		}
		return out
	}
	if dateFrom != "" {
		out := make([]DailyCandle, 0, len(all))
		for _, c := range all {
			if c.DateISO >= dateFrom {
				out = append(out, c)
			}
		}
		return out
	}
	if dateTo != "" {
		out := make([]DailyCandle, 0, len(all))
		for _, c := range all {
			if c.DateISO <= dateTo {
				out = append(out, c)
			}
		}
		return out
	}
	return all
}

func parseOptionalFloatQuery(r *http.Request, name string) (float64, bool, error) {
	raw, exists := r.URL.Query()[name]
	if !exists || len(raw) == 0 || strings.TrimSpace(raw[0]) == "" {
		return 0, false, nil
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(raw[0]), 64)
	if err != nil {
		return 0, false, fmt.Errorf("Input should be a valid number for '%s'", name)
	}
	return v, true, nil
}

func parseBoolQuery(r *http.Request, name string, def bool) (bool, error) {
	raw := strings.TrimSpace(r.URL.Query().Get(name))
	if raw == "" {
		return def, nil
	}
	switch strings.ToLower(raw) {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("Input should be a valid boolean for '%s'", name)
	}
}

func derefFloat(p *float64) float64 {
	if p == nil {
		return 0.0
	}
	return *p
}
