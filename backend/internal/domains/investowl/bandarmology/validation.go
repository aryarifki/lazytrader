package bandarmology

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
	"github.com/go-chi/chi/v5"
)

func alphaRowToMap(r shared.AlphaScanRow) map[string]any {
	return map[string]any{
		"ticker":            r.Ticker,
		"broker_code":       r.BrokerCode,
		"n_events":          r.NEvents,
		"mean_fwd_return":   shared.FinitePtr(&r.MeanFwdReturn),
		"median_fwd_return": shared.FinitePtr(&r.MedianFwdReturn),
		"win_rate":          shared.FinitePtr(&r.WinRate),
		"avg_net_value":     shared.FinitePtr(&r.AvgNetValue),
		"total_net_value":   shared.FinitePtr(&r.TotalNetValue),
		"se_hac":            shared.FinitePtr(&r.SEHAC),
		"t_stat":            shared.FinitePtr(&r.TStat),
		"p_value_one_sided": shared.FinitePtr(&r.PValueOneSided),
		"status":            r.Status,
		"p_value_fdr":       shared.FinitePtr(&r.PValueFDR),
		"significant":       r.Significant,
	}
}

func alphaRowsToMaps(rows []shared.AlphaScanRow) []map[string]any {
	out := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		out = append(out, alphaRowToMap(r))
	}
	return out
}

func buildEventStudyCharts(evRows []EventStudyRow, v2 bool) ([]map[string]any, []map[string]any, []map[string]any) {
	if len(evRows) == 0 {
		return []map[string]any{}, []map[string]any{}, []map[string]any{}
	}
	xs := []int{0, 1, 3, 5, 10}
	labels := []string{"Signal", "+1D", "+3D", "+5D", "+10D"}

	ribbon := make([]map[string]any, 0, len(xs))
	for i, h := range xs {
		var vals []float64
		for _, r := range evRows {
			if p := r.Values[h]; p != nil {
				vals = append(vals, *p)
			}
		}
		var medPtr, q25Ptr, q75Ptr *float64
		if len(vals) > 0 {
			m := shared.LinearPercentile(vals, 50)
			l := shared.LinearPercentile(vals, 25)
			u := shared.LinearPercentile(vals, 75)
			medPtr, q25Ptr, q75Ptr = &m, &l, &u
		}
		ribbon = append(ribbon, map[string]any{
			"day":    labels[i],
			"median": medPtr,
			"range":  []*float64{q25Ptr, q75Ptr},
		})
	}

	paths := make([]map[string]any, 0, len(evRows))
	table := make([]map[string]any, 0, len(evRows))
	for idx, r := range evRows {
		pData := make(map[string]*float64, len(xs))
		for i, h := range xs {
			pData[labels[i]] = r.Values[h]
		}
		paths = append(paths, map[string]any{
			"id":   fmt.Sprintf("event_%d", idx),
			"data": pData,
		})

		if v2 {
			table = append(table, map[string]any{
				"ticker":      r.Ticker,
				"signal_date": fmtDateMidnightSpace(r.SignalDate),
				"signal":      r.BandarSignal,
				"t_plus_1d":   r.Values[1],
				"t_plus_3d":   r.Values[3],
				"t_plus_5d":   r.Values[5],
				"t_plus_10d":  r.Values[10],
			})
		} else {
			table = append(table, map[string]any{
				"ticker":       r.Ticker,
				"signal_date":  fmtDateMidnightSpace(r.SignalDate),
				"signal":       r.BandarSignal,
				"signal_score": r.BandarSignalScore,
				"t_plus_0d":    r.Values[0],
				"t_plus_1d":    r.Values[1],
				"t_plus_3d":    r.Values[3],
				"t_plus_5d":    r.Values[5],
				"t_plus_10d":   r.Values[10],
			})
		}
	}
	return ribbon, paths, table
}

func (h *Handler) runValidation(r *http.Request, v2 bool) (any, int, error) {
	ticker := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	q := r.URL.Query()
	rawDate := strings.TrimSpace(q.Get("analysis_date"))

	windowDays := 20
	if raw := strings.TrimSpace(q.Get("window_days")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			windowDays = n
		}
	}
	horizon := 10
	if raw := strings.TrimSpace(q.Get("horizon")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			horizon = n
		}
	}
	minEvents := 5
	if raw := strings.TrimSpace(q.Get("min_events")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			minEvents = n
		}
	}
	minNetBuy := 0.0
	if raw := strings.TrimSpace(q.Get("min_net_buy")); raw != "" {
		if v, err := strconv.ParseFloat(raw, 64); err == nil {
			minNetBuy = v
		}
	}

	var endTS *time.Time
	if rawDate != "" {
		t, err := time.Parse("2006-01-02", rawDate)
		if err != nil {
			return nil, http.StatusUnprocessableEntity, fmt.Errorf("invalid analysis_date")
		}
		endTS = &t
	}

	prices, err := h.Repo.ReadPrices(r.Context(), []string{ticker}, nil, endTS)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	acts, err := h.Repo.ReadBrokerActivity(r.Context(), []string{ticker}, nil, endTS)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	flows, err := h.Repo.ReadBrokerFlow(r.Context(), []string{ticker}, nil, endTS)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	actIn := make([]shared.ActivityInput, len(acts))
	for i, a := range acts {
		actIn[i] = shared.ActivityInput{Date: fmtDateISO(a.Date), Ticker: a.Ticker, BrokerCode: a.BrokerCode, NetValue: a.NetValue}
	}
	pxIn := make([]shared.PriceInput, 0, len(prices))
	for _, p := range prices {
		if p.Close != nil && *p.Close > 0 {
			pxIn = append(pxIn, shared.PriceInput{Date: fmtDateISO(p.Date), Ticker: p.Ticker, Close: *p.Close})
		}
	}
	scanRows := alphaRowsToMaps(shared.BrokerAlphaScan(actIn, pxIn, horizon, minEvents, minNetBuy*1e9))

	accMap := map[string]bool{
		"STRONG_ACCUMULATION": true, "ACCUMULATION": true, "NET_BUY": true,
		"AKUMULASI_KUAT": true, "AKUMULASI": true,
	}
	evRows := EventStudyTable(prices, flows, []int{1, 3, 5, 10}, &windowDays, accMap, endTS, 5)
	chart, paths, table := buildEventStudyCharts(evRows, v2)

	var scanPayload any = scanRows
	if v2 {
		scanPayload = map[string]any{"ticker": scanRows, "all": scanRows}
	}
	return map[string]any{
		"broker_scan": scanPayload,
		"event_study": map[string]any{
			"chart": chart,
			"paths": paths,
			"table": table,
		},
	}, http.StatusOK, nil
}

// GetValidation handles GET /api/bandar/validation/{ticker}.
func (h *Handler) GetValidation(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ck := bandarCacheKey("bandar_validation", "validation_insight", fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s",
		strings.ToUpper(chi.URLParam(r, "ticker")), q.Get("analysis_date"), q.Get("window_days"),
		q.Get("horizon"), q.Get("min_events"), q.Get("min_net_buy"), q.Get("universe_mode")))
	if h.tryCache(w, r, ck, 30, 1800) {
		return
	}
	res, status, err := h.runValidation(r, false)
	if err != nil {
		httpx.WriteError(w, status, err.Error())
		return
	}
	h.storeCache(w, r, ck, res, 30, 1800)
}

// GetValidationV2 handles GET /api/bandar/validation-v2/{ticker}.
func (h *Handler) GetValidationV2(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ck := bandarCacheKey("bandar_validation", "validation_insight_v2", fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s",
		strings.ToUpper(chi.URLParam(r, "ticker")), q.Get("analysis_date"), q.Get("window_days"),
		q.Get("horizon"), q.Get("min_events"), q.Get("min_net_buy"), q.Get("universe_mode")))
	if h.tryCache(w, r, ck, 30, 1800) {
		return
	}
	res, status, err := h.runValidation(r, true)
	if err != nil {
		httpx.WriteError(w, status, err.Error())
		return
	}
	h.storeCache(w, r, ck, res, 30, 1800)
}

// GetValidationBrokerScan handles GET /api/bandar/validation/broker-scan.
func (h *Handler) GetValidationBrokerScan(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ticker := strings.ToUpper(strings.TrimSpace(q.Get("ticker")))
	if ticker == "" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "ticker is required")
		return
	}
	horizon := 10
	if raw := strings.TrimSpace(q.Get("horizon")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			horizon = n
		}
	}
	minEvents := 5
	if raw := strings.TrimSpace(q.Get("min_events")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			minEvents = n
		}
	}
	minNetB := 0.0
	if raw := strings.TrimSpace(q.Get("min_net_b")); raw != "" {
		if v, err := strconv.ParseFloat(raw, 64); err == nil {
			minNetB = v
		}
	}

	prices, err := h.Repo.ReadPrices(r.Context(), []string{ticker}, nil, nil)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	acts, err := h.Repo.ReadBrokerActivity(r.Context(), []string{ticker}, nil, nil)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	actIn := make([]shared.ActivityInput, len(acts))
	for i, a := range acts {
		actIn[i] = shared.ActivityInput{Date: fmtDateISO(a.Date), Ticker: a.Ticker, BrokerCode: a.BrokerCode, NetValue: a.NetValue}
	}
	pxIn := make([]shared.PriceInput, 0, len(prices))
	for _, p := range prices {
		if p.Close != nil && *p.Close > 0 {
			pxIn = append(pxIn, shared.PriceInput{Date: fmtDateISO(p.Date), Ticker: p.Ticker, Close: *p.Close})
		}
	}
	rows := alphaRowsToMaps(shared.BrokerAlphaScan(actIn, pxIn, horizon, minEvents, minNetB*1e9))
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"data": rows})
}

// GetStockEventStudy handles GET /api/bandar/stocks/{ticker}/event-study.
func (h *Handler) GetStockEventStudy(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	q := r.URL.Query()
	rawHor := q.Get("horizons")
	if rawHor == "" {
		rawHor = "1,3,5,10"
	}
	var horizons []int
	for _, p := range strings.Split(rawHor, ",") {
		if n, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
			horizons = append(horizons, n)
		}
	}
	if len(horizons) == 0 {
		horizons = []int{1, 3, 5, 10}
	}
	lookbackDays := 20
	if raw := strings.TrimSpace(q.Get("lookback_days")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			lookbackDays = n
		}
	}

	prices, err := h.Repo.ReadPrices(r.Context(), []string{ticker}, nil, nil)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	flows, err := h.Repo.ReadBrokerFlow(r.Context(), []string{ticker}, nil, nil)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sigSet := map[string]bool{
		"AKUMULASI_KUAT": true, "AKUMULASI": true,
		"STRONG_ACCUMULATION": true, "ACCUMULATION": true, "NET_BUY": true,
	}
	evRows := EventStudyTable(prices, flows, horizons, &lookbackDays, sigSet, nil, 5)
	out := make([]map[string]any, 0, len(evRows))
	sortedH := append([]int{0}, horizons...)
	sort.Ints(sortedH)
	for _, ev := range evRows {
		m := map[string]any{
			"ticker":              ev.Ticker,
			"signal_date":         fmtDateMidnightT(ev.SignalDate),
			"bandar_signal":       ev.BandarSignal,
			"bandar_signal_score": ev.BandarSignalScore,
		}
		for _, hVal := range sortedH {
			m[fmt.Sprintf("t_plus_%dd", hVal)] = ev.Values[hVal]
		}
		out = append(out, m)
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"data": out})
}
