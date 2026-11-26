package ownership

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
)

func returnAt(values []float64, periods int) *float64 {
	if periods <= 0 || len(values) <= periods {
		return nil
	}
	now := values[len(values)-1]
	before := values[len(values)-1-periods]
	return pctChange(&now, &before)
}

func ytdReturn(history []IndexHistoryPoint) *float64 {
	if len(history) == 0 {
		return nil
	}
	last := history[len(history)-1]
	if last.Date == nil || len(*last.Date) < 4 {
		return nil
	}
	yearPrefix := (*last.Date)[:4]
	for _, row := range history {
		if row.Date != nil && strings.HasPrefix(*row.Date, yearPrefix) {
			now := last.Close
			before := row.Close
			return pctChange(&now, &before)
		}
	}
	return nil
}

// GetMarketOverview handles GET /api/v1/ownership/market.
func (h *Handler) GetMarketOverview(w http.ResponseWriter, r *http.Request) {
	historyLimit := 260
	if raw := strings.TrimSpace(r.URL.Query().Get("history_limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 30 || n > 600 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "history_limit must be between 30 and 600")
			return
		}
		historyLimit = n
	}

	ck := ownershipCacheKey("ownership_market_v2", fmt.Sprintf("%d", historyLimit))
	if h.tryCache(w, r, ck, 1800) {
		return
	}

	ctx := r.Context()
	fullHistory, err := indexHistory(ctx, h.DB, "COMPOSITE", 320)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(fullHistory) == 0 {
		httpx.WriteError(w, http.StatusNotFound, "Data IHSG tidak tersedia")
		return
	}

	closes := make([]float64, len(fullHistory))
	highs := []float64{}
	for i, row := range fullHistory {
		closes[i] = row.Close
		if row.Close > 0 {
			highs = append(highs, row.Close)
		}
	}
	latest := fullHistory[len(fullHistory)-1]
	histStart := len(fullHistory) - historyLimit
	if histStart < 0 {
		histStart = 0
	}
	history := fullHistory[histStart:]

	var high52w, low52w *float64
	if len(highs) > 0 {
		winStart := len(highs) - 252
		if winStart < 0 {
			winStart = 0
		}
		win := highs[winStart:]
		hi, lo := win[0], win[0]
		for _, v := range win[1:] {
			if v > hi {
				hi = v
			}
			if v < lo {
				lo = v
			}
		}
		hiR := roundEven(hi, 4)
		loR := roundEven(lo, 4)
		high52w = &hiR
		low52w = &loR
	}

	// Sector index performance
	secRows, err := h.DB.Query(ctx, `
		SELECT DISTINCT index_code, index_name
		FROM idx_index_summary
		WHERE index_code <> 'COMPOSITE'
	`)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	type codeName struct {
		code string
		name string
	}
	secPairs := []codeName{}
	for secRows.Next() {
		var code string
		var name *string
		if scanErr := secRows.Scan(&code, &name); scanErr != nil {
			secRows.Close()
			httpx.WriteError(w, http.StatusInternalServerError, scanErr.Error())
			return
		}
		nm := code
		if name != nil && *name != "" {
			nm = *name
		}
		secPairs = append(secPairs, codeName{code: code, name: nm})
	}
	secRows.Close()
	if err := secRows.Err(); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var latestDateObj *time.Time
	if latest.Date != nil {
		if t, parseErr := time.Parse("2006-01-02", *latest.Date); parseErr == nil {
			latestDateObj = &t
		}
	}

	sectors := []map[string]any{}
	for _, pair := range secPairs {
		sh, shErr := indexHistory(ctx, h.DB, pair.code, 90)
		if shErr != nil {
			httpx.WriteError(w, http.StatusInternalServerError, shErr.Error())
			return
		}
		if len(sh) < 2 {
			continue
		}
		lastItem := sh[len(sh)-1]
		if latestDateObj != nil && lastItem.Date != nil {
			if lastObj, pErr := time.Parse("2006-01-02", *lastItem.Date); pErr == nil {
				if latestDateObj.Sub(lastObj).Hours()/24 > 10 {
					continue
				}
			}
		}
		vals := make([]float64, len(sh))
		for i, rItem := range sh {
			vals[i] = rItem.Close
		}
		sectors = append(sectors, map[string]any{
			"code":          pair.code,
			"name":          pair.name,
			"close":         roundEven(vals[len(vals)-1], 4),
			"change_1d_pct": returnAt(vals, 1),
			"change_1w_pct": returnAt(vals, 5),
			"change_1m_pct": returnAt(vals, 21),
			"date":          lastItem.Date,
		})
	}
	sort.SliceStable(sectors, func(i, j int) bool {
		ci, _ := sectors[i]["change_1m_pct"].(*float64)
		cj, _ := sectors[j]["change_1m_pct"].(*float64)
		vi := -9999.0
		if ci != nil {
			vi = *ci
		}
		vj := -9999.0
		if cj != nil {
			vj = *cj
		}
		return vi > vj
	})

	// Market breadth from latest two trading days
	var latestPriceDate, prevPriceDate *time.Time
	if err := h.DB.QueryRow(ctx, `SELECT MAX(date) FROM prices`).Scan(&latestPriceDate); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if latestPriceDate != nil {
		_ = h.DB.QueryRow(ctx, `SELECT MAX(date) FROM prices WHERE date < $1`, *latestPriceDate).Scan(&prevPriceDate)
	}

	var adv, dec, unch, totBreadth int
	if latestPriceDate != nil && prevPriceDate != nil {
		bRows, bErr := h.DB.Query(ctx, `
			SELECT p1.close::double precision, p2.close::double precision
			FROM prices p1
			JOIN prices p2 ON p2.ticker = p1.ticker
			WHERE p1.date = $1 AND p2.date = $2
			  AND p1.close IS NOT NULL AND p2.close IS NOT NULL
		`, *latestPriceDate, *prevPriceDate)
		if bErr != nil {
			httpx.WriteError(w, http.StatusInternalServerError, bErr.Error())
			return
		}
		for bRows.Next() {
			var nowV, beforeV float64
			if scanErr := bRows.Scan(&nowV, &beforeV); scanErr != nil {
				bRows.Close()
				httpx.WriteError(w, http.StatusInternalServerError, scanErr.Error())
				return
			}
			totBreadth++
			if nowV > beforeV {
				adv++
			} else if nowV < beforeV {
				dec++
			} else {
				unch++
			}
		}
		bRows.Close()
		if err := bRows.Err(); err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	foreignFlow := map[string]any{
		"date": isoDatePtr(latestPriceDate),
		"1d":   nil,
		"1w":   nil,
		"1m":   nil,
	}
	if latestPriceDate != nil {
		for _, pair := range []struct {
			k    string
			days int
		}{{"1d", 1}, {"1w", 7}, {"1m", 30}} {
			startDate := latestPriceDate.AddDate(0, 0, -(pair.days - 1))
			var val float64
			// Keep verbatim (value / volume) bigint division semantics matching Python ownership.py.
			qErr := h.DB.QueryRow(ctx, `
				SELECT COALESCE(SUM((foreign_buy - foreign_sell) *
					CASE
						WHEN volume > 0 AND value > 0 THEN (value / volume)
						ELSE COALESCE(close, 0)
					END
				), 0)::double precision
				FROM prices
				WHERE date >= $1 AND date <= $2
			`, startDate, *latestPriceDate).Scan(&val)
			if qErr != nil {
				httpx.WriteError(w, http.StatusInternalServerError, qErr.Error())
				return
			}
			foreignFlow[pair.k] = roundEven(val, 2)
		}
	}

	var monthlyCount int
	var monthlyDate *time.Time
	if err := h.DB.QueryRow(ctx, `SELECT COUNT(DISTINCT ticker)::int, MAX(period_date) FROM monthly_ownership_metrics`).Scan(&monthlyCount, &monthlyDate); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var totalTickers int
	if err := h.DB.QueryRow(ctx, `SELECT COUNT(*)::int FROM tickers`).Scan(&totalTickers); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ihsgComp := getIHSGComposition(ctx, h.DB)
	hasComposition := ihsgComp != nil
	hasTop20 := false
	if ihsgComp != nil {
		if top20, ok := ihsgComp["top20"].([]map[string]any); ok && len(top20) > 0 {
			hasTop20 = true
		}
	}
	covNote := "Komposisi kepemilikan IHSG dan Top 20 konstituen memerlukan dataset " +
		"shares outstanding/market cap per emiten yang belum ada di read layer ini. " +
		"Angka IHSG di atas berasal dari idx_index_summary."
	if hasComposition {
		covNote = "Komposisi makro IHSG dan Top 20 konstituen bersumber dari baseline agregat KSEI / CT."
	}

	resp := map[string]any{
		"index": map[string]any{
			"code":           "COMPOSITE",
			"name":           "IHSG",
			"date":           latest.Date,
			"close":          latest.Close,
			"change_1d_pct":  returnAt(closes, 1),
			"return_1w_pct":  returnAt(closes, 5),
			"return_1m_pct":  returnAt(closes, 21),
			"return_3m_pct":  returnAt(closes, 63),
			"return_6m_pct":  returnAt(closes, 126),
			"return_1y_pct":  returnAt(closes, 252),
			"return_ytd_pct": ytdReturn(fullHistory),
			"high_52w":       high52w,
			"low_52w":        low52w,
			"market_cap":     latest.MarketCap,
			"value":          latest.Value,
			"volume":         latest.Volume,
		},
		"history":      history,
		"sectors":      sectors,
		"breadth": map[string]any{
			"advancers": adv,
			"decliners": dec,
			"unchanged": unch,
			"total":     totBreadth,
		},
		"foreign_flow": foreignFlow,
		"composition":  ihsgComp,
		"ownership_coverage": map[string]any{
			"covered_tickers":       monthlyCount,
			"total_tickers":         totalTickers,
			"period_date":           isoDatePtr(monthlyDate),
			"composition_available": hasComposition,
			"top20_available":       hasTop20,
			"note":                  covNote,
		},
		"data_quality": map[string]any{
			"index_source":          "idx_index_summary",
			"breadth_source":        "prices",
			"foreign_flow_source":   "broker_flow",
			"ownership_composition": hasComposition,
			"notes": []string{
				"IHSG dan sektor dihitung dari idx_index_summary.",
				"Breadth dihitung dari perubahan harga harian seluruh emiten di prices.",
				"Foreign flow adalah agregat foreign_net_broker, bukan kepemilikan.",
				"Komposisi makro dan Top 20 konstituen dihitung dari data historis agregat KSEI.",
			},
		},
	}
	h.storeCache(w, r, ck, resp, 1800)
}
