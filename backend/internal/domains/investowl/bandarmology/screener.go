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
)

// GetScreener handles GET /api/bandar/screener.
func (h *Handler) GetScreener(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	universeMode := q.Get("universe_mode")
	if universeMode == "" {
		universeMode = "watchlist"
	}
	horizon := 10
	if raw := strings.TrimSpace(q.Get("horizon")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			horizon = n
		}
	}

	tickers := shared.GetDynamicUniverse(r.Context(), h.Repo.DB, universeMode)
	prices, err := h.Repo.ReadPrices(r.Context(), tickers, nil, nil)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	acts, err := h.Repo.ReadBrokerActivity(r.Context(), tickers, nil, nil)
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
	rows := alphaRowsToMaps(shared.BrokerAlphaScan(actIn, pxIn, horizon, 5, 0))
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"data": rows})
}

// GetScreenerV2 handles deprecated GET /api/bandar/screener-v2.
func (h *Handler) GetScreenerV2(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	universeMode := q.Get("universe_mode")
	if universeMode == "" {
		universeMode = "lq45"
	}
	rawDate := strings.TrimSpace(q.Get("analysis_date"))
	windowDays := 20
	if raw := strings.TrimSpace(q.Get("window_days")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			windowDays = n
		}
	}
	rawTickers := strings.TrimSpace(q.Get("tickers"))

	ck := bandarCacheKey("bandar_screener_v2", "screener", fmt.Sprintf("%s|%s|%d|%s", universeMode, rawDate, windowDays, rawTickers))
	if h.tryCache(w, r, ck, 30, 300) {
		return
	}

	var tickerList []string
	if rawTickers != "" {
		for _, p := range strings.Split(rawTickers, ",") {
			if s := strings.ToUpper(strings.TrimSpace(p)); s != "" {
				tickerList = append(tickerList, s)
			}
		}
	} else {
		tickerList = shared.GetDynamicUniverse(r.Context(), h.Repo.DB, universeMode)
	}
	if len(tickerList) == 0 {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"data": []any{}, "meta": map[string]any{}})
		return
	}
	if len(tickerList) > 100 {
		tickerList = tickerList[:100]
	}

	var analysisTS time.Time
	if rawDate != "" {
		t, err := time.Parse("2006-01-02", rawDate)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid analysis_date")
			return
		}
		analysisTS = t
	} else {
		analysisTS = toDateOnly(time.Now())
	}
	windowStart := analysisTS.AddDate(0, 0, -windowDays)
	cutoffFlow := analysisTS.AddDate(0, 0, -(windowDays + 60))

	allFlows, _ := h.Repo.ReadBrokerFlow(r.Context(), tickerList, &cutoffFlow, &analysisTS)
	allActs, _ := h.Repo.ReadBrokerActivity(r.Context(), tickerList, &cutoffFlow, &analysisTS)
	allPrices, _ := h.Repo.ReadPrices(r.Context(), tickerList, &cutoffFlow, &analysisTS)

	actInAll := make([]shared.ActivityInput, len(allActs))
	for i, a := range allActs {
		actInAll[i] = shared.ActivityInput{Date: fmtDateISO(a.Date), Ticker: a.Ticker, BrokerCode: a.BrokerCode, NetValue: a.NetValue}
	}
	pxInAll := make([]shared.PriceInput, 0, len(allPrices))
	for _, p := range allPrices {
		if p.Close != nil && *p.Close > 0 {
			pxInAll = append(pxInAll, shared.PriceInput{Date: fmtDateISO(p.Date), Ticker: p.Ticker, Close: *p.Close})
		}
	}
	scan10D := shared.BrokerAlphaScan(actInAll, pxInAll, 10, 5, 0)

	flByTicker := make(map[string][]FlowRow)
	for _, f := range allFlows {
		flByTicker[f.Ticker] = append(flByTicker[f.Ticker], f)
	}
	actByTicker := make(map[string][]ActivityRow)
	for _, a := range allActs {
		actByTicker[a.Ticker] = append(actByTicker[a.Ticker], a)
	}
	pxByTicker := make(map[string][]PriceRow)
	for _, p := range allPrices {
		pxByTicker[p.Ticker] = append(pxByTicker[p.Ticker], p)
	}

	type screenerV2Item struct {
		Ticker          string  `json:"ticker"`
		Signal          string  `json:"signal"`
		ConvictionScore float64 `json:"conviction_score"`
		ForeignNet      float64 `json:"foreign_net"`
		NetValue        float64 `json:"net_value"`
		BandarAvgPrice  float64 `json:"bandar_avg_price"`
		TotalValue      float64 `json:"total_value"`
		CurrentPrice    float64 `json:"current_price"`
		TopBuyer        string  `json:"top_buyer"`
		DataDate        string  `json:"data_date"`
	}
	results := make([]screenerV2Item, 0, len(tickerList))

	for _, tk := range tickerList {
		tFlow := flByTicker[tk]
		if len(tFlow) == 0 {
			continue
		}
		sort.Slice(tFlow, func(i, j int) bool { return tFlow[i].Date.Before(tFlow[j].Date) })
		latestFlow := tFlow[len(tFlow)-1]

		rawSig := "NEUTRAL"
		if latestFlow.BandarSignal != nil && *latestFlow.BandarSignal != "" {
			rawSig = *latestFlow.BandarSignal
		}
		sigTitle := shared.TitleWords(strings.ReplaceAll(rawSig, "_", " "))

		f5Start := max(0, len(tFlow)-5)
		foreign5D := 0.0
		totVals5D := make([]float64, 0, 5)
		flowInputs := make([]shared.FlowInput, len(tFlow))
		for i, fl := range tFlow {
			flowInputs[i] = shared.FlowInput{Date: fmtDateISO(fl.Date), ForeignNetBroker: shared.Finite(fl.ForeignNetBroker, 0)}
			if i >= f5Start {
				foreign5D += shared.Finite(fl.ForeignNetBroker, 0)
				totVals5D = append(totVals5D, shared.Finite(fl.TotalValue, 0))
			}
		}

		tPx := pxByTicker[tk]
		pxInputs := make([]shared.PriceInput, 0, len(tPx))
		for _, p := range tPx {
			if p.Close != nil && *p.Close > 0 {
				pxInputs = append(pxInputs, shared.PriceInput{Date: fmtDateISO(p.Date), Ticker: p.Ticker, Close: *p.Close})
			}
		}
		conv := shared.ConvictionScore(rawSig, &foreign5D, scan10D, tk, pxInputs, flowInputs, totVals5D)

		tAct := actByTicker[tk]
		var actWin []ActivityRow
		if len(tAct) > 0 {
			maxAD := tAct[0].Date
			for _, a := range tAct[1:] {
				if a.Date.After(maxAD) {
					maxAD = a.Date
				}
			}
			cut := maxAD.AddDate(0, 0, -windowDays)
			for _, a := range tAct {
				if !a.Date.Before(cut) {
					actWin = append(actWin, a)
				}
			}
		}
		dfSmart := SmartDailyFromActivity(actWin)
		totalNetVal := 0.0
		for _, sd := range dfSmart {
			totalNetVal += sd.SmartNet
		}

		topB, topS := TopNetBrokerSummary(tAct, tk, nil, 150)
		topBuyer := "-"
		if len(topB) > 0 {
			topBuyer = topB[0].BrokerCode
		}
		var bVal5, bLot5, mktVal, mktLot float64
		for i, b := range topB {
			if i < 5 {
				bVal5 += b.BuyValue
				bLot5 += shared.Finite(b.BuyLot, 0)
			}
			mktVal += b.BuyValue + b.SellValue
			mktLot += shared.Finite(b.BuyLot, 0) + shared.Finite(b.SellLot, 0)
		}
		for _, s := range topS {
			mktVal += s.BuyValue + s.SellValue
			mktLot += shared.Finite(s.BuyLot, 0) + shared.Finite(s.SellLot, 0)
		}
		bandarAvg := 0.0
		if bLot5 > 0 {
			bandarAvg = bVal5 / (bLot5 * 100.0)
		}
		pxLatest := PriceAtOrBefore(tPx, analysisTS)
		officialClose := 0.0
		if pxLatest != nil && pxLatest.Close != nil {
			officialClose = *pxLatest.Close
		}
		currPrice := officialClose
		if currPrice <= 0 && mktLot > 0 {
			currPrice = mktVal / (mktLot * 100.0)
		}

		results = append(results, screenerV2Item{
			Ticker:          tk,
			Signal:          sigTitle,
			ConvictionScore: conv.Score,
			ForeignNet:      foreign5D,
			NetValue:        totalNetVal,
			BandarAvgPrice:  bandarAvg,
			TotalValue:      shared.Finite(latestFlow.TotalValue, 0),
			CurrentPrice:    currPrice,
			TopBuyer:        topBuyer,
			DataDate:        fmtDateISO(latestFlow.Date),
		})
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].ConvictionScore > results[j].ConvictionScore
	})

	h.storeCache(w, r, ck, map[string]any{
		"data": results,
		"meta": map[string]any{
			"window_start":  fmtDateISO(windowStart),
			"analysis_date": fmtDateISO(analysisTS),
		},
	}, 30, 300)
}
