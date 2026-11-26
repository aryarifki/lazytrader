package bandarmology

import (
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
	"github.com/go-chi/chi/v5"
)

// GetTickerDetail handles GET /api/bandar/detail/{ticker}.
func (h *Handler) GetTickerDetail(w http.ResponseWriter, r *http.Request) {
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

	ck := bandarCacheKey("bandar", "ticker_detail", fmt.Sprintf("%s|%s|%d|%d", ticker, rawDate, windowDays, horizon))
	if h.tryCache(w, r, ck, 30, 1800) {
		return
	}

	var analysisTS time.Time
	if rawDate != "" {
		parsed, err := time.Parse("2006-01-02", rawDate)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid analysis_date")
			return
		}
		analysisTS = parsed
	} else if maxD := h.Repo.GetMaxActivityDate(r.Context(), ticker); maxD != nil {
		analysisTS = *maxD
	} else {
		analysisTS = toDateOnly(time.Now())
	}

	fetchStart := analysisTS.AddDate(0, 0, -max(windowDays+60, 180))
	fetchEnd := analysisTS
	priceStart := analysisTS.AddDate(0, 0, -max(windowDays+90, 240))

	actDF, err := h.Repo.ReadBrokerActivity(r.Context(), []string{ticker}, &fetchStart, &fetchEnd)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(actDF) == 0 {
		h.storeCache(w, r, ck, map[string]any{"error": "No broker history for " + ticker}, 30, 1800)
		return
	}

	priceDF, err := h.Repo.ReadPrices(r.Context(), []string{ticker}, &priceStart, &fetchEnd)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	brokerDF, err := h.Repo.ReadBrokerFlow(r.Context(), []string{ticker}, &fetchStart, &fetchEnd)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(brokerDF) == 0 {
		h.storeCache(w, r, ck, map[string]any{"error": "No broker history for " + ticker}, 30, 1800)
		return
	}

	windowStart := analysisTS.AddDate(0, 0, -windowDays)
	var priceWin []PriceRow
	for _, p := range priceDF {
		if !p.Date.Before(windowStart) && !p.Date.After(analysisTS) {
			priceWin = append(priceWin, p)
		}
	}
	var brokerWin []FlowRow
	for _, f := range brokerDF {
		if !f.Date.Before(windowStart) && !f.Date.After(analysisTS) {
			brokerWin = append(brokerWin, f)
		}
	}
	var actWin []ActivityRow
	for _, a := range actDF {
		if !a.Date.Before(windowStart) && !a.Date.After(analysisTS) {
			actWin = append(actWin, a)
		}
	}

	pxRow := PriceAtOrBefore(priceDF, analysisTS)
	signalRow := FlowRowAt(brokerDF, ticker, analysisTS)
	actDate := LatestActivityDate(actDF, ticker, analysisTS)
	topBuy, topSell := TopNetBrokerSummary(actDF, ticker, actDate, 6)

	dailySmart := SmartDailyFromActivity(actWin)
	profileEntries := ProfileFlowFromActivity(actWin)

	// Convert to shared inputs for BrokerAlphaScan & ConvictionScore
	allActIn := make([]shared.ActivityInput, len(actDF))
	for i, a := range actDF {
		allActIn[i] = shared.ActivityInput{Date: fmtDateISO(a.Date), Ticker: a.Ticker, BrokerCode: a.BrokerCode, NetValue: a.NetValue}
	}
	allPxIn := make([]shared.PriceInput, 0, len(priceDF))
	for _, p := range priceDF {
		if p.Close != nil && *p.Close > 0 {
			allPxIn = append(allPxIn, shared.PriceInput{Date: fmtDateISO(p.Date), Ticker: p.Ticker, Close: *p.Close})
		}
	}
	scan10D := shared.BrokerAlphaScan(allActIn, allPxIn, 10, 5, 0)

	var closeVal *float64
	if pxRow != nil && pxRow.Close != nil {
		v := *pxRow.Close
		closeVal = &v
	}
	ret5D := ReturnToDate(priceDF, analysisTS, 5)
	ret10D := ReturnToDate(priceDF, analysisTS, 10)

	sort.Slice(brokerWin, func(i, j int) bool { return brokerWin[i].Date.Before(brokerWin[j].Date) })
	foreign5D := 0.0
	if len(brokerWin) > 0 {
		for _, fl := range brokerWin[max(0, len(brokerWin)-5):] {
			foreign5D += shared.Finite(fl.ForeignNetBroker, 0)
		}
	}
	var smartCum *float64
	if len(dailySmart) > 0 {
		v := dailySmart[len(dailySmart)-1].CumulativeNet
		smartCum = &v
	}

	rawSigStr := "None"
	rawSigBlank := ""
	signalScore := 0.0
	if signalRow != nil {
		if signalRow.BandarSignal != nil {
			rawSigStr = *signalRow.BandarSignal
			rawSigBlank = *signalRow.BandarSignal
		}
		if signalRow.BandarSignalScore != nil {
			signalScore = *signalRow.BandarSignalScore
		}
	}
	sigTitle := shared.TitleWords(strings.ReplaceAll(rawSigStr, "_", " "))

	sort.Slice(brokerDF, func(i, j int) bool { return brokerDF[i].Date.Before(brokerDF[j].Date) })
	flowIn := make([]shared.FlowInput, len(brokerDF))
	totVals5D := make([]float64, 0, 5)
	fStart5 := max(0, len(brokerDF)-5)
	for i, fl := range brokerDF {
		flowIn[i] = shared.FlowInput{Date: fmtDateISO(fl.Date), ForeignNetBroker: shared.Finite(fl.ForeignNetBroker, 0)}
		if i >= fStart5 {
			totVals5D = append(totVals5D, shared.Finite(fl.TotalValue, 0))
		}
	}
	conv := shared.ConvictionScore(rawSigBlank, &foreign5D, scan10D, ticker, allPxIn, flowIn, totVals5D)
	alerts := ContradictionAlerts(rawSigBlank, ret5D, ret10D, &foreign5D, smartCum)

	var sig10D []shared.AlphaScanRow
	for _, s := range scan10D {
		if s.Significant {
			sig10D = append(sig10D, s)
		}
	}
	var verdict string
	if len(sig10D) == 0 {
		r5, r10 := 0.0, 0.0
		if ret5D != nil {
			r5 = *ret5D
		}
		if ret10D != nil {
			r10 = *ret10D
		}
		verdict = fmt.Sprintf("%s shows %s with %+.2f%% over 5D and %+.2f%% over 10D. The current read is directional, but broker-specific 10D validation is not yet statistically strong.",
			ticker, sigTitle, r5*100.0, r10*100.0)
	} else {
		sort.SliceStable(sig10D, func(i, j int) bool {
			pi, pj := math.Inf(1), math.Inf(1)
			if !math.IsNaN(sig10D[i].PValueOneSided) {
				pi = sig10D[i].PValueOneSided
			}
			if !math.IsNaN(sig10D[j].PValueOneSided) {
				pj = sig10D[j].PValueOneSided
			}
			if pi != pj {
				return pi < pj
			}
			return sig10D[i].MeanFwdReturn > sig10D[j].MeanFwdReturn
		})
		best := sig10D[0]
		winPct := int(math.RoundToEven(best.WinRate * 100.0))
		verdict = fmt.Sprintf("%s shows %s. Broker %s is the strongest 10D validation: %d events, mean return %+.2f%%, win rate %d%%, p-value %.4f.",
			ticker, sigTitle, best.BrokerCode, best.NEvents, best.MeanFwdReturn*100.0, winPct, best.PValueOneSided)
	}

	rawSummary := make([]map[string]any, 0, 6)
	for _, sidePair := range []struct {
		side string
		rows []ActivityRow
	}{
		{"Buy", topBuy[:min(len(topBuy), 3)]},
		{"Sell", topSell[:min(len(topSell), 3)]},
	} {
		for _, row := range sidePair.rows {
			rawSummary = append(rawSummary, map[string]any{
				"side":   sidePair.side,
				"broker": row.BrokerCode,
				"type":   shared.ParticipantLabel(row.ParticipantType),
				"net":    row.NetValue,
				"spark":  SparklineValues(actDF, row.BrokerCode, analysisTS, 5),
			})
		}
	}
	brokerSummaryRows, _ := EnrichBrokerSummaryAffiliations(r.Context(), h.Repo, ticker, rawSummary)
	affiliatedDetected := DetectWindowAffiliations(r.Context(), h.Repo, ticker, actWin)

	perfAll := PricePerformanceTable(priceDF)
	wantPeriods := map[string]bool{"1D": true, "1W": true, "1M": true, "3M": true, "6M": true, "YTD": true}
	perfRows := make([]map[string]any, 0, 6)
	for _, pr := range perfAll {
		if wantPeriods[pr["timeframe"].(string)] {
			perfRows = append(perfRows, map[string]any{"period": pr["timeframe"], "value": pr["return"]})
		}
	}

	sort.SliceStable(profileEntries, func(i, j int) bool { return profileEntries[i].Net > profileEntries[j].Net })
	profileRows := make([]map[string]any, 0, min(len(profileEntries), 6))
	for i := 0; i < min(len(profileEntries), 6); i++ {
		profileRows = append(profileRows, map[string]any{"label": profileEntries[i].Label, "net": profileEntries[i].Net})
	}

	smartDailyRows := make([]map[string]any, 0, len(dailySmart))
	for _, sd := range dailySmart {
		smartDailyRows = append(smartDailyRows, map[string]any{
			"date":           fmtDateMidnightSpace(sd.Date),
			"smart_net":      sd.SmartNet,
			"cumulative_net": sd.CumulativeNet,
		})
	}

	// Bandar VWAP of top 5 net-buying brokers in window
	brkTotals := make(map[string]float64)
	for _, a := range actWin {
		brkTotals[a.BrokerCode] += a.NetValue
	}
	type bt struct {
		code string
		net  float64
	}
	var posBrks []bt
	for c, v := range brkTotals {
		if v > 0 {
			posBrks = append(posBrks, bt{c, v})
		}
	}
	sort.SliceStable(posBrks, func(i, j int) bool {
		if posBrks[i].net != posBrks[j].net {
			return posBrks[i].net > posBrks[j].net
		}
		return posBrks[i].code < posBrks[j].code
	})
	topBuyersList := make([]string, 0, min(len(posBrks), 5))
	topBuyerSet := make(map[string]bool, 5)
	for i := 0; i < min(len(posBrks), 5); i++ {
		topBuyersList = append(topBuyersList, posBrks[i].code)
		topBuyerSet[posBrks[i].code] = true
	}

	sort.Slice(priceWin, func(i, j int) bool { return priceWin[i].Date.Before(priceWin[j].Date) })
	priceChartRows := make([]map[string]any, 0, len(priceWin))
	var bandarVWAPLatest *float64
	for _, row := range priceWin {
		var vwapPtr *float64
		if len(topBuyersList) > 0 && len(actWin) > 0 {
			var totVal, totLot float64
			for _, a := range actWin {
				if !a.Date.After(row.Date) && topBuyerSet[a.BrokerCode] {
					totVal += a.BuyValue
					totLot += shared.Finite(a.BuyLot, 0)
				}
			}
			if totLot > 0 {
				v := shared.Round(totVal/(totLot*100.0), 2)
				vwapPtr = &v
				bandarVWAPLatest = &v
			}
		}
		var volPtr *float64
		if row.Volume != nil {
			vf := float64(*row.Volume)
			volPtr = &vf
		}
		priceChartRows = append(priceChartRows, map[string]any{
			"date":        fmtDateMidnightSpace(row.Date),
			"close":       row.Close,
			"volume":      volPtr,
			"bandar_vwap": vwapPtr,
		})
	}
	var bandarPnLPct *float64
	if bandarVWAPLatest != nil && closeVal != nil && *bandarVWAPLatest > 0 {
		v := shared.Round(((*closeVal-*bandarVWAPLatest)/(*bandarVWAPLatest))*100.0, 2)
		bandarPnLPct = &v
	}

	retailSummary := buildRetailSummary(actWin, dailySmart)

	signalOverlay := make([]map[string]any, 0, len(brokerWin))
	for _, br := range brokerWin {
		bSig := "None"
		if br.BandarSignal != nil {
			bSig = *br.BandarSignal
		}
		signalOverlay = append(signalOverlay, map[string]any{
			"date":   fmtDateMidnightSpace(br.Date),
			"signal": shared.TitleWords(strings.ReplaceAll(bSig, "_", " ")),
			"score":  br.BandarSignalScore,
		})
	}

	compName, compGroup := h.Repo.GetCompanyMeta(r.Context(), ticker)
	topBuyerMap := map[string]any{"broker": "-", "net": nil}
	if len(topBuy) > 0 {
		topBuyerMap = map[string]any{"broker": topBuy[0].BrokerCode, "net": topBuy[0].NetValue}
	}
	topSellerMap := map[string]any{"broker": "-", "net": nil}
	if len(topSell) > 0 {
		topSellerMap = map[string]any{"broker": topSell[0].BrokerCode, "net": topSell[0].NetValue}
	}
	var actDateStr *string
	if actDate != nil {
		s := fmtDateISO(*actDate)
		actDateStr = &s
	}

	result := map[string]any{
		"ticker":        ticker,
		"company":       map[string]any{"name": compName, "group": compGroup},
		"analysis_date": fmtDateISO(analysisTS),
		"window_start":  fmtDateISO(windowStart),
		"window_days":   windowDays,
		"close":         closeVal,
		"ret_5d":        ret5D,
		"ret_10d":       ret10D,
		"signal":        sigTitle,
		"signal_raw":    rawSigBlank,
		"signal_score":  signalScore,
		"foreign_5d":    foreign5D,
		"smart_cumulative": smartCum,
		"conviction_score": conv.Score,
		"conviction_breakdown": map[string]any{
			"score":               conv.Score,
			"p_value":             conv.PValue,
			"causality_component": conv.CausalityComponent,
			"signal_component":    conv.SignalComponent,
			"foreign_component":   conv.ForeignComponent,
			"broker_component":    conv.BrokerComponent,
			"broker_note":         conv.BrokerNote,
		},
		"top_buyer":                  topBuyerMap,
		"top_seller":                 topSellerMap,
		"alerts":                     alerts,
		"verdict":                    verdict,
		"broker_summary":             brokerSummaryRows,
		"affiliated_brokers_detected": affiliatedDetected,
		"price_performance":          perfRows,
		"profile_flow":               profileRows,
		"profile_broker_detail":      ProfileBrokerDetailTable(actWin),
		"smart_daily":                smartDailyRows,
		"price_chart":                priceChartRows,
		"signal_overlay":             signalOverlay,
		"activity_date":              actDateStr,
		"bandar_vwap":                bandarVWAPLatest,
		"bandar_pnl_pct":             bandarPnLPct,
		"top_accumulating_brokers":   topBuyersList,
		"retail_summary":             retailSummary,
	}
	h.storeCache(w, r, ck, result, 30, 1800)
}

func buildRetailSeriesAndBrokers(actWin []ActivityRow, dailySmart []SmartDailyRow) ([]map[string]any, []map[string]any, float64, float64) {
	retailDailyMap := make(map[time.Time]float64)
	type brkAgg struct{ net, buy, sell float64 }
	retByBroker := make(map[string]*brkAgg)
	hasRetail := false

	for _, a := range actWin {
		if !retailBrokerCodes[a.BrokerCode] {
			continue
		}
		hasRetail = true
		retailDailyMap[a.Date] += a.NetValue
		agg, ok := retByBroker[a.BrokerCode]
		if !ok {
			agg = &brkAgg{}
			retByBroker[a.BrokerCode] = agg
		}
		agg.net += a.NetValue
		agg.buy += a.BuyValue
		agg.sell += a.SellValue
	}

	retailDailyRows := make([]map[string]any, 0)
	retailCum, smartCum := 0.0, 0.0
	if hasRetail {
		smartMap := make(map[time.Time]float64, len(dailySmart))
		dateSet := make(map[time.Time]struct{})
		for d := range retailDailyMap {
			dateSet[d] = struct{}{}
		}
		for _, sd := range dailySmart {
			smartMap[sd.Date] = sd.SmartNet
			dateSet[sd.Date] = struct{}{}
		}
		dates := make([]time.Time, 0, len(dateSet))
		for d := range dateSet {
			dates = append(dates, d)
		}
		sort.Slice(dates, func(i, j int) bool { return dates[i].Before(dates[j]) })

		for _, d := range dates {
			rn := retailDailyMap[d]
			sn := smartMap[d]
			retailCum += rn
			smartCum += sn
			retailDailyRows = append(retailDailyRows, map[string]any{
				"date":       fmtDateMidnightSpace(d),
				"retail_net": rn,
				"smart_net":  sn,
				"retail_cum": retailCum,
				"smart_cum":  smartCum,
			})
		}
	}

	topRetail := make([]map[string]any, 0, len(retByBroker))
	for c, agg := range retByBroker {
		topRetail = append(topRetail, map[string]any{
			"broker": c,
			"net":    agg.net,
			"buy":    agg.buy,
			"sell":   agg.sell,
		})
	}
	sort.SliceStable(topRetail, func(i, j int) bool {
		return topRetail[i]["buy"].(float64) > topRetail[j]["buy"].(float64)
	})
	return retailDailyRows, topRetail, retailCum, smartCum
}

func buildRetailSummary(actWin []ActivityRow, dailySmart []SmartDailyRow) map[string]any {
	rows, topRet, retCum, smCum := buildRetailSeriesAndBrokers(actWin, dailySmart)
	if len(topRet) > 5 {
		topRet = topRet[:5]
	}

	sentiment := "Konsolidasi / Netral"
	desc := "Arus modal ritel dan institusi seimbang tanpa arah akumulasi dominan."
	tone := "neutral"
	switch {
	case smCum > 0 && retCum < 0:
		sentiment = "Akumulasi Murni (Ritel Jual ke Smart Money)"
		desc = "Smart money agresif menyerap saham saat investor ritel melepas posisi (Cut Loss / Panic Sell)."
		tone = "positive"
	case smCum < 0 && retCum > 0:
		sentiment = "Distribusi Trap (Ritel Menampung Distribusi)"
		desc = "Investor ritel memborong saham saat smart money melepas posisi (Potensi jebakan distribusi)."
		tone = "negative"
	case smCum > 0 && retCum > 0:
		sentiment = "Inflow Bersama (Smart Money & Ritel Beli)"
		desc = "Smart money dan ritel berada di posisi beli bersamaan (Momentum Mark-up aktif)."
		tone = "primary"
	case smCum < 0 && retCum < 0:
		sentiment = "Outflow Bersama (Smart Money & Ritel Keluar)"
		desc = "Seluruh partisipan melepas posisi (Fase Markdown / Likuidasi)."
		tone = "warning"
	}

	return map[string]any{
		"retail_cum_net":     retCum,
		"smart_cum_net":      smCum,
		"sentiment":          sentiment,
		"sentiment_desc":     desc,
		"sentiment_tone":     tone,
		"top_retail_brokers": topRet,
		"retail_daily":       rows,
	}
}

// GetStockRetailFlow handles GET /api/bandar/stocks/{ticker}/retail-flow.
func (h *Handler) GetStockRetailFlow(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	windowDays := 30
	if raw := strings.TrimSpace(r.URL.Query().Get("window_days")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			windowDays = n
		}
	}

	act, err := h.Repo.ReadBrokerActivity(r.Context(), []string{ticker}, nil, nil)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(act) == 0 {
		httpx.WriteError(w, http.StatusNotFound, fmt.Sprintf("Tidak ada data aktivitas broker untuk ticker %s", ticker))
		return
	}
	maxD := act[len(act)-1].Date
	for _, a := range act {
		if a.Date.After(maxD) {
			maxD = a.Date
		}
	}
	winStart := maxD.AddDate(0, 0, -windowDays)
	var sub []ActivityRow
	for _, a := range act {
		if !a.Date.Before(winStart) && !a.Date.After(maxD) {
			sub = append(sub, a)
		}
	}

	dailySmart := SmartDailyFromActivity(sub)
	rows, topRet, retCum, smCum := buildRetailSeriesAndBrokers(sub, dailySmart)
	if len(topRet) > 8 {
		topRet = topRet[:8]
	}

	sent, tone := "Konsolidasi Netral", "neutral"
	switch {
	case smCum > 0 && retCum < 0:
		sent, tone = "Akumulasi Murni (Ritel Cut Loss / Panik)", "positive"
	case smCum < 0 && retCum > 0:
		sent, tone = "Distribusi Trap (Ritel Borong Saham)", "negative"
	case smCum > 0 && retCum > 0:
		sent, tone = "Inflow Bersama (Smart Money & Ritel Beli)", "primary"
	case smCum < 0 && retCum < 0:
		sent, tone = "Outflow Bersama (Likuidasi Pasar)", "warning"
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"ticker":                ticker,
		"window_days":           windowDays,
		"retail_cumulative_net": retCum,
		"smart_cumulative_net":  smCum,
		"sentiment":             sent,
		"sentiment_tone":        tone,
		"retail_brokers":        topRet,
		"daily_series":          rows,
	})
}
