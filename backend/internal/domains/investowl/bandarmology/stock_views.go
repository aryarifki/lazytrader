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

func activityRowsToDicts(rows []ActivityRow, includeGross bool) []map[string]any {
	out := make([]map[string]any, 0, len(rows))
	for _, a := range rows {
		m := map[string]any{
			"date":             fmtDateMidnightT(a.Date),
			"ticker":           a.Ticker,
			"broker_code":      a.BrokerCode,
			"participant_type": a.ParticipantType,
			"buy_value":        a.BuyValue,
			"sell_value":       a.SellValue,
			"net_value":        a.NetValue,
			"buy_lot":          a.BuyLot,
			"sell_lot":         a.SellLot,
			"frequency":        a.Frequency,
			"buy_avg_price":    a.BuyAvgPrice,
			"sell_avg_price":   a.SellAvgPrice,
			"fetched_at":       fmtFetchedAtOffset(a.FetchedAt),
		}
		if includeGross {
			m["gross_value"] = a.BuyValue + a.SellValue
		}
		out = append(out, m)
	}
	return out
}

// GetStockMetrics handles GET /api/bandar/stocks/{ticker}/metrics.
func (h *Handler) GetStockMetrics(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	q := r.URL.Query()
	rawDate := strings.TrimSpace(q.Get("date"))
	window := 20
	if raw := strings.TrimSpace(q.Get("window")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			window = n
		}
	}

	ck := bandarCacheKey("bandar_metrics", "metrics", fmt.Sprintf("%s|%s|%d", ticker, rawDate, window))
	if h.tryCache(w, r, ck, 30, 1800) {
		return
	}

	var fetchStart time.Time
	var fetchEnd *time.Time
	if rawDate != "" {
		parsed, err := time.Parse("2006-01-02", rawDate)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid date")
			return
		}
		fetchEnd = &parsed
		fetchStart = parsed.AddDate(0, 0, -max(window+60, 400))
	} else {
		fetchStart = toDateOnly(time.Now()).AddDate(0, 0, -400)
	}

	priceDF, err := h.Repo.ReadPrices(r.Context(), []string{ticker}, &fetchStart, fetchEnd)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	flowDF, err := h.Repo.ReadBrokerFlow(r.Context(), []string{ticker}, &fetchStart, fetchEnd)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	actDF, err := h.Repo.ReadBrokerActivity(r.Context(), []string{ticker}, &fetchStart, fetchEnd)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(flowDF) == 0 {
		httpx.WriteError(w, http.StatusNotFound, fmt.Sprintf("No broker data for %s", ticker))
		return
	}

	var ts time.Time
	if fetchEnd != nil {
		ts = *fetchEnd
	} else {
		ts = flowDF[len(flowDF)-1].Date
	}
	winStart := ts.AddDate(0, 0, -window)

	var flowWin []FlowRow
	for _, f := range flowDF {
		if !f.Date.Before(winStart) && !f.Date.After(ts) {
			flowWin = append(flowWin, f)
		}
	}
	var actWin []ActivityRow
	for _, a := range actDF {
		if !a.Date.Before(winStart) && !a.Date.After(ts) {
			actWin = append(actWin, a)
		}
	}

	sort.Slice(flowWin, func(i, j int) bool { return flowWin[i].Date.Before(flowWin[j].Date) })
	var foreign5D *float64
	if len(flowWin) > 0 {
		s := 0.0
		for _, fl := range flowWin[max(0, len(flowWin)-5):] {
			s += shared.Finite(fl.ForeignNetBroker, 0)
		}
		foreign5D = &s
	}

	topBuy, topSell := TopNetBrokerSummary(actDF, ticker, &ts, 6)
	var latestSigPtr *string
	rawSig := ""
	if len(flowWin) > 0 {
		latestSigPtr = flowWin[len(flowWin)-1].BandarSignal
		if latestSigPtr != nil {
			rawSig = *latestSigPtr
		}
	}

	actIn := make([]shared.ActivityInput, len(actDF))
	for i, a := range actDF {
		actIn[i] = shared.ActivityInput{Date: fmtDateISO(a.Date), Ticker: a.Ticker, BrokerCode: a.BrokerCode, NetValue: a.NetValue}
	}
	pxIn := make([]shared.PriceInput, 0, len(priceDF))
	for _, p := range priceDF {
		if p.Close != nil && *p.Close > 0 {
			pxIn = append(pxIn, shared.PriceInput{Date: fmtDateISO(p.Date), Ticker: p.Ticker, Close: *p.Close})
		}
	}
	scan := shared.BrokerAlphaScan(actIn, pxIn, 10, 5, 0)

	sort.Slice(flowDF, func(i, j int) bool { return flowDF[i].Date.Before(flowDF[j].Date) })
	flIn := make([]shared.FlowInput, len(flowDF))
	tot5D := make([]float64, 0, 5)
	for i, f := range flowDF {
		flIn[i] = shared.FlowInput{Date: fmtDateISO(f.Date), ForeignNetBroker: shared.Finite(f.ForeignNetBroker, 0)}
		if i >= max(0, len(flowDF)-5) {
			tot5D = append(tot5D, shared.Finite(f.TotalValue, 0))
		}
	}
	conv := shared.ConvictionScore(rawSig, foreign5D, scan, ticker, pxIn, flIn, tot5D)

	var smartCum *float64
	if len(actWin) > 0 {
		hasSmart := false
		sum := 0.0
		for _, a := range actWin {
			if shared.IsSmartProfile(shared.BrokerProfileOf(a.BrokerCode)) {
				hasSmart = true
				sum += a.NetValue
			}
		}
		if hasSmart {
			smartCum = &sum
		}
	}

	res := map[string]any{
		"ticker":        ticker,
		"analysis_date": fmtDateISO(ts),
		"window_start":  fmtDateISO(winStart),
		"signal":        latestSigPtr,
		"conviction": map[string]any{
			"score":   shared.Round(conv.Score, 1),
			"p_value": conv.PValue,
			"components": map[string]any{
				"causality": conv.CausalityComponent,
				"signal":    conv.SignalComponent,
				"foreign":   conv.ForeignComponent,
				"broker":    shared.Round(conv.BrokerComponent, 1),
			},
			"broker_note": conv.BrokerNote,
		},
		"ret_5d":           ReturnToDate(priceDF, ts, 5),
		"ret_10d":          ReturnToDate(priceDF, ts, 10),
		"foreign_net_5d":   foreign5D,
		"smart_cumulative": smartCum,
		"top_buyers":       activityRowsToDicts(topBuy, false),
		"top_sellers":      activityRowsToDicts(topSell, false),
	}
	h.storeCache(w, r, ck, res, 30, 1800)
}

// GetStockSmartFlow handles GET /api/bandar/stocks/{ticker}/smart-flow.
func (h *Handler) GetStockSmartFlow(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	window := 20
	if raw := strings.TrimSpace(r.URL.Query().Get("window")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			window = n
		}
	}
	act, err := h.Repo.ReadBrokerActivity(r.Context(), []string{ticker}, nil, nil)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(act) == 0 {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"data": []any{}})
		return
	}
	end := act[len(act)-1].Date
	for _, a := range act {
		if a.Date.After(end) {
			end = a.Date
		}
	}
	start := end.AddDate(0, 0, -window)
	var win []ActivityRow
	for _, a := range act {
		if !a.Date.Before(start) && !a.Date.After(end) {
			win = append(win, a)
		}
	}
	daily := SmartDailyFromActivity(win)
	rows := make([]map[string]any, 0, len(daily))
	for _, d := range daily {
		rows = append(rows, map[string]any{
			"date":           fmtDateISO(d.Date),
			"smart_net":      d.SmartNet,
			"cumulative_net": d.CumulativeNet,
		})
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"data": rows})
}

// GetStockBrokerCompare handles GET /api/bandar/stocks/{ticker}/broker-compare.
func (h *Handler) GetStockBrokerCompare(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	q := r.URL.Query()
	window := 20
	if raw := strings.TrimSpace(q.Get("window")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			window = n
		}
	}
	mode := q.Get("mode")
	if mode == "" {
		mode = "cumulative"
	}

	act, err := h.Repo.ReadBrokerActivity(r.Context(), []string{ticker}, nil, nil)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(act) == 0 {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"data": []any{}})
		return
	}
	end := act[len(act)-1].Date
	for _, a := range act {
		if a.Date.After(end) {
			end = a.Date
		}
	}
	start := end.AddDate(0, 0, -window)

	dateSet := make(map[time.Time]struct{})
	codeSet := make(map[string]struct{})
	grid := make(map[time.Time]map[string]float64)
	for _, a := range act {
		if a.Date.Before(start) || a.Date.After(end) {
			continue
		}
		dateSet[a.Date] = struct{}{}
		codeSet[a.BrokerCode] = struct{}{}
		if grid[a.Date] == nil {
			grid[a.Date] = make(map[string]float64)
		}
		grid[a.Date][a.BrokerCode] += a.NetValue
	}
	if len(dateSet) == 0 {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"data": []any{}})
		return
	}
	dates := make([]time.Time, 0, len(dateSet))
	for d := range dateSet {
		dates = append(dates, d)
	}
	sort.Slice(dates, func(i, j int) bool { return dates[i].Before(dates[j]) })

	codes := make([]string, 0, len(codeSet))
	for c := range codeSet {
		codes = append(codes, c)
	}
	sort.Strings(codes)

	isCum := strings.EqualFold(strings.TrimSpace(mode), "cumulative")
	running := make(map[string]float64, len(codes))
	data := make([]map[string]any, 0, len(dates))
	for _, d := range dates {
		rec := make(map[string]any, len(codes)+1)
		rec["date"] = fmtDateMidnightSpace(d)
		m := grid[d]
		for _, c := range codes {
			v := m[c]
			if isCum {
				running[c] += v
				v = running[c]
			}
			rec[c] = v
		}
		data = append(data, rec)
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"data": data})
}

// GetStockRaw handles GET /api/bandar/stocks/{ticker}/raw.
func (h *Handler) GetStockRaw(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	window := 20
	if raw := strings.TrimSpace(r.URL.Query().Get("window")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			window = n
		}
	}
	flowDF, err := h.Repo.ReadBrokerFlow(r.Context(), []string{ticker}, nil, nil)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	actDF, err := h.Repo.ReadBrokerActivity(r.Context(), []string{ticker}, nil, nil)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var end time.Time
	if len(flowDF) > 0 {
		end = flowDF[len(flowDF)-1].Date
	} else if len(actDF) > 0 {
		end = actDF[len(actDF)-1].Date
	} else {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"flow": []any{}, "activity": []any{}})
		return
	}
	start := end.AddDate(0, 0, -window)

	flowOut := make([]map[string]any, 0)
	for _, f := range flowDF {
		if f.Date.Before(start) || f.Date.After(end) {
			continue
		}
		flowOut = append(flowOut, map[string]any{
			"date":                fmtDateISO(f.Date),
			"ticker":              f.Ticker,
			"bandar_signal":       f.BandarSignal,
			"bandar_signal_score": f.BandarSignalScore,
			"foreign_net_broker":  f.ForeignNetBroker,
			"local_net_broker":    f.LocalNetBroker,
			"gov_net_broker":      f.GovNetBroker,
			"foreign_net_flow":    f.ForeignNetFlow,
			"domestic_net_flow":   f.DomesticNetFlow,
			"total_value":         f.TotalValue,
			"foreign_signal":      f.ForeignSignal,
			"conclusion_broker":   f.ConclusionBroker,
			"conclusion_flow":     f.ConclusionFlow,
			"fetched_at":          fmtFetchedAtOffset(f.FetchedAt),
		})
	}
	actOut := make([]map[string]any, 0)
	for _, a := range actDF {
		if a.Date.Before(start) || a.Date.After(end) {
			continue
		}
		actOut = append(actOut, map[string]any{
			"date":             fmtDateISO(a.Date),
			"ticker":           a.Ticker,
			"broker_code":      a.BrokerCode,
			"participant_type": a.ParticipantType,
			"buy_value":        a.BuyValue,
			"sell_value":       a.SellValue,
			"net_value":        a.NetValue,
			"buy_lot":          a.BuyLot,
			"sell_lot":         a.SellLot,
			"frequency":        a.Frequency,
			"buy_avg_price":    a.BuyAvgPrice,
			"sell_avg_price":   a.SellAvgPrice,
			"fetched_at":       fmtFetchedAtOffset(a.FetchedAt),
		})
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"flow": flowOut, "activity": actOut})
}

// GetStockRawTables handles GET /api/bandar/stocks/{ticker}/raw-tables.
func (h *Handler) GetStockRawTables(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	q := r.URL.Query()
	rawDate := strings.TrimSpace(q.Get("analysis_date"))
	windowDays := 20
	if raw := strings.TrimSpace(q.Get("window_days")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			windowDays = n
		}
	}

	ck := bandarCacheKey("bandar_rawtables", "get_raw_tables", fmt.Sprintf("%s|%s|%d", ticker, rawDate, windowDays))
	if h.tryCache(w, r, ck, 30, 1800) {
		return
	}

	var analysisTS time.Time
	if rawDate != "" {
		t, err := time.Parse("2006-01-02", rawDate)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid analysis_date")
			return
		}
		analysisTS = t
	} else if maxD := h.Repo.GetMaxActivityDate(r.Context(), ticker); maxD != nil {
		analysisTS = *maxD
	} else {
		analysisTS = toDateOnly(time.Now())
	}

	fetchStart := analysisTS.AddDate(0, 0, -(windowDays + 14))
	windowStart := analysisTS.AddDate(0, 0, -windowDays)

	flowDF, _ := h.Repo.ReadBrokerFlow(r.Context(), []string{ticker}, &fetchStart, &analysisTS)
	var flowWin []FlowRow
	for _, f := range flowDF {
		if !f.Date.Before(windowStart) && !f.Date.After(analysisTS) {
			flowWin = append(flowWin, f)
		}
	}
	sort.SliceStable(flowWin, func(i, j int) bool { return flowWin[i].Date.After(flowWin[j].Date) })
	flowList := make([]map[string]any, 0, len(flowWin))
	for _, f := range flowWin {
		sigRaw := "NEUTRAL"
		if f.BandarSignal != nil && *f.BandarSignal != "" {
			sigRaw = *f.BandarSignal
		}
		sig := shared.TitleWords(strings.ReplaceAll(sigRaw, "_", " "))
		if strings.EqualFold(sig, "NET BUY") {
			sig = "Net Buy"
		} else if strings.EqualFold(sig, "NET SELL") {
			sig = "Net Sell"
		}
		flowList = append(flowList, map[string]any{
			"date":        fmtDateISO(f.Date),
			"signal":      sig,
			"score":       shared.Finite(f.BandarSignalScore, 0),
			"foreign_net": shared.Finite(f.ForeignNetBroker, 0),
			"local_net":   shared.Finite(f.LocalNetBroker, 0),
			"total_value": shared.Finite(f.TotalValue, 0),
		})
	}

	actDF, _ := h.Repo.ReadBrokerActivity(r.Context(), []string{ticker}, &fetchStart, &analysisTS)
	var actWin []ActivityRow
	for _, a := range actDF {
		if !a.Date.Before(windowStart) && !a.Date.After(analysisTS) {
			actWin = append(actWin, a)
		}
	}
	sort.SliceStable(actWin, func(i, j int) bool {
		if !actWin[i].Date.Equal(actWin[j].Date) {
			return actWin[i].Date.After(actWin[j].Date)
		}
		gi := actWin[i].BuyValue + actWin[i].SellValue
		gj := actWin[j].BuyValue + actWin[j].SellValue
		return gi > gj
	})
	actList := make([]map[string]any, 0, len(actWin))
	for _, a := range actWin {
		tLabel := strings.ToUpper(a.ParticipantType)
		switch a.ParticipantType {
		case "Asing":
			tLabel = "FOREIGN"
		case "Lokal":
			tLabel = "LOCAL"
		case "Pemerintah":
			tLabel = "GOV"
		}
		actList = append(actList, map[string]any{
			"date":   fmtDateISO(a.Date),
			"broker": a.BrokerCode,
			"type":   tLabel,
			"buy":    a.BuyValue,
			"sell":   a.SellValue,
			"net":    a.NetValue,
			"freq":   int(shared.Finite(a.Frequency, 0)),
		})
	}

	h.storeCache(w, r, ck, map[string]any{"flow": flowList, "activity": actList}, 30, 1800)
}

// GetStockBrokerProfiles handles GET /api/bandar/stocks/{ticker}/broker-profiles.
func (h *Handler) GetStockBrokerProfiles(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	lookbackDays := 20
	if raw := strings.TrimSpace(r.URL.Query().Get("lookback_days")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			lookbackDays = n
		}
	}
	act, err := h.Repo.ReadBrokerActivity(r.Context(), []string{ticker}, nil, nil)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(act) == 0 {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"data": []any{}})
		return
	}
	maxD := act[len(act)-1].Date
	for _, a := range act {
		if a.Date.After(maxD) {
			maxD = a.Date
		}
	}
	cutoff := maxD.AddDate(0, 0, -lookbackDays)
	var sub []ActivityRow
	for _, a := range act {
		if !a.Date.Before(cutoff) {
			sub = append(sub, a)
		}
	}
	entries := ProfileFlowFromActivity(sub)
	out := make([]map[string]any, 0, len(entries))
	for _, pe := range entries {
		out = append(out, map[string]any{
			"profile":        pe.Profile,
			"label":          pe.Label,
			"characteristic": pe.Description,
			"net":            pe.Net,
			"buy":            pe.Buy,
			"sell":           pe.Sell,
			"top_brokers":    pe.TopBrokers,
		})
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"data": out})
}

// GetStockPricePerformance handles GET /api/bandar/stocks/{ticker}/price-performance.
func (h *Handler) GetStockPricePerformance(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	prices, err := h.Repo.ReadPrices(r.Context(), []string{ticker}, nil, nil)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"data": PricePerformanceTable(prices)})
}

// GetStockBrokerDistribution handles GET /api/bandar/stocks/{ticker}/broker-distribution.
func (h *Handler) GetStockBrokerDistribution(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	q := r.URL.Query()
	rawDate := strings.TrimSpace(q.Get("trade_date"))
	topN := 12
	if raw := strings.TrimSpace(q.Get("top_n")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			topN = n
		}
	}

	act, err := h.Repo.ReadBrokerActivity(r.Context(), []string{ticker}, nil, nil)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(act) == 0 {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"data": []any{}})
		return
	}
	var target time.Time
	if rawDate != "" {
		parsed, err := time.Parse("2006-01-02", rawDate)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid trade_date")
			return
		}
		target = parsed
	} else {
		target = act[len(act)-1].Date
		for _, a := range act {
			if a.Date.After(target) {
				target = a.Date
			}
		}
	}
	var sub []ActivityRow
	for _, a := range act {
		if a.Ticker == ticker && a.Date.Equal(target) {
			sub = append(sub, a)
		}
	}
	sort.SliceStable(sub, func(i, j int) bool {
		return (sub[i].BuyValue + sub[i].SellValue) > (sub[j].BuyValue + sub[j].SellValue)
	})
	if len(sub) > topN {
		sub = sub[:topN]
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"data": activityRowsToDicts(sub, true)})
}
