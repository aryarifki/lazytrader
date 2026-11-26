package signal

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
)

type rawConfluenceCache struct {
	Data []ConfluenceItem `json:"data"`
	Meta map[string]any   `json:"meta"`
}

type flowScreenerRow struct {
	Date             string
	Ticker           string
	BandarSignal     string
	ForeignNetBroker float64
	TotalValue       float64
}

type actScreenerRow struct {
	Date       string
	Ticker     string
	BrokerCode string
	BuyValue   float64
	BuyLot     float64
	NetValue   float64
}

var (
	smartForeignSet = map[string]bool{"AK": true, "BK": true, "KZ": true, "CS": true}
	gorenganSet     = map[string]bool{"MG": true, "SQ": true, "CP": true, "EP": true, "XC": true, "YP": true}
)

// GetConfluenceScreener mirrors GET /api/signal/confluence-screener.
func (r *Repository) GetConfluenceScreener(
	ctx context.Context,
	cacheStore *cache.Store,
	universeMode string,
	analysisDate *time.Time,
	windowDays int,
	tickersCSV string,
	minConviction float64,
	minLiquidity float64,
	tagFilter string,
	signalFilter string,
) (map[string]any, error) {
	var aDateStr string
	if analysisDate != nil {
		aDateStr = analysisDate.Format("2006-01-02")
	}
	var rawKeyTickers string
	if strings.TrimSpace(tickersCSV) != "" {
		rawKeyTickers = tickersCSV
	} else {
		rawKeyTickers = "None"
	}
	if aDateStr == "" {
		aDateStr = "None"
	}
	rawCacheKey := fmt.Sprintf("confluence_raw_v3:%s|%s|%d|%s", universeMode, aDateStr, windowDays, rawKeyTickers)

	var rawResults []ConfluenceItem
	var metaInfo map[string]any

	loadedFromRaw := false
	if cacheStore != nil {
		if b, ok := cacheStore.Get(ctx, rawCacheKey); ok {
			var entry rawConfluenceCache
			if json.Unmarshal(b, &entry) == nil && entry.Meta != nil {
				rawResults = entry.Data
				metaInfo = entry.Meta
				loadedFromRaw = true
			}
		}
	}

	if !loadedFromRaw {
		var tickerList []string
		if strings.TrimSpace(tickersCSV) != "" {
			for _, p := range strings.Split(tickersCSV, ",") {
				if s := strings.ToUpper(strings.TrimSpace(p)); s != "" {
					tickerList = append(tickerList, s)
				}
			}
		} else {
			tickerList = shared.GetDynamicUniverse(ctx, r.DB, universeMode)
		}
		if len(tickerList) == 0 {
			return map[string]any{"count": 0, "meta": map[string]any{}, "data": []ConfluenceItem{}}, nil
		}
		if len(tickerList) > 100 {
			tickerList = tickerList[:100]
		}

		now := time.Now()
		analysisTS := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		if analysisDate != nil {
			analysisTS = *analysisDate
		}
		cutoffDate := analysisTS.AddDate(0, 0, -(windowDays + 30))
		endDateStr := analysisTS.Format("2006-01-02")

		// 1. Read broker_flow, broker_activity, prices pushed down to SQL
		flowRows, err := r.DB.Query(ctx, `
			SELECT date, ticker, COALESCE(bandar_signal, 'NEUTRAL'),
			       foreign_net_broker::double precision, total_value::double precision
			FROM broker_flow
			WHERE ticker = ANY($1) AND date >= $2 AND date <= $3
			ORDER BY ticker, date ASC
		`, tickerList, cutoffDate, analysisTS)
		if err != nil {
			return nil, err
		}
		flowByTicker := make(map[string][]flowScreenerRow)
		for flowRows.Next() {
			var dt time.Time
			var tk, sig string
			var fn, tv *float64
			if err := flowRows.Scan(&dt, &tk, &sig, &fn, &tv); err != nil {
				flowRows.Close()
				return nil, err
			}
			flowByTicker[tk] = append(flowByTicker[tk], flowScreenerRow{
				Date:             dt.Format("2006-01-02"),
				Ticker:           tk,
				BandarSignal:     sig,
				ForeignNetBroker: shared.Finite(fn, 0),
				TotalValue:       shared.Finite(tv, 0),
			})
		}
		flowRows.Close()

		actRows, err := r.DB.Query(ctx, `
			SELECT date, ticker, broker_code,
			       buy_value::double precision, buy_lot::double precision, net_value::double precision
			FROM broker_activity
			WHERE ticker = ANY($1) AND date >= $2 AND date <= $3
			ORDER BY ticker, date ASC, net_value DESC
		`, tickerList, cutoffDate, analysisTS)
		if err != nil {
			return nil, err
		}
		actByTicker := make(map[string][]actScreenerRow)
		allActInputs := make([]shared.ActivityInput, 0, 2048)
		for actRows.Next() {
			var dt time.Time
			var tk, brk string
			var bVal, bLot, nVal *float64
			if err := actRows.Scan(&dt, &tk, &brk, &bVal, &bLot, &nVal); err != nil {
				actRows.Close()
				return nil, err
			}
			dStr := dt.Format("2006-01-02")
			nv := shared.Finite(nVal, 0)
			actByTicker[tk] = append(actByTicker[tk], actScreenerRow{
				Date: dStr, Ticker: tk, BrokerCode: brk,
				BuyValue: shared.Finite(bVal, 0), BuyLot: shared.Finite(bLot, 0), NetValue: nv,
			})
			allActInputs = append(allActInputs, shared.ActivityInput{
				Date: dStr, Ticker: tk, BrokerCode: brk, NetValue: nv,
			})
		}
		actRows.Close()

		pxRows, err := r.DB.Query(ctx, `
			SELECT date, ticker, close::double precision
			FROM prices
			WHERE ticker = ANY($1) AND date >= $2 AND date <= $3
			ORDER BY ticker, date ASC
		`, tickerList, cutoffDate, analysisTS)
		if err != nil {
			return nil, err
		}
		pxByTicker := make(map[string][]shared.PriceInput)
		allPxInputs := make([]shared.PriceInput, 0, 1024)
		for pxRows.Next() {
			var dt time.Time
			var tk string
			var cl *float64
			if err := pxRows.Scan(&dt, &tk, &cl); err != nil {
				pxRows.Close()
				return nil, err
			}
			if cl != nil && !math.IsNaN(*cl) && !math.IsInf(*cl, 0) {
				pi := shared.PriceInput{Date: dt.Format("2006-01-02"), Ticker: tk, Close: *cl}
				pxByTicker[tk] = append(pxByTicker[tk], pi)
				allPxInputs = append(allPxInputs, pi)
			}
		}
		pxRows.Close()

		scan10D := shared.BrokerAlphaScan(allActInputs, allPxInputs, 10, 5, 0)

		// Metadata & Free Float
		type metaRec struct {
			compName  string
			grpName   string
			freeFloat float64
		}
		metaMap := make(map[string]metaRec, len(tickerList))
		mRows, err := r.DB.Query(ctx, `
			SELECT t.ticker, t.name,
			       COALESCE(cg.groups, t.sector, '-'),
			       s.free_float_pct::double precision
			FROM tickers t
			LEFT JOIN (
				SELECT ticker, STRING_AGG(DISTINCT group_name, ' / ' ORDER BY group_name) AS groups
				FROM ownership_staging.v_conglomerate_members
				WHERE ticker = ANY($1)
				GROUP BY ticker
			) cg ON cg.ticker = t.ticker
			LEFT JOIN (
				SELECT ticker,
				       CASE
				           WHEN COALESCE(SUM(CASE WHEN (name ILIKE 'Masyarakat%' OR name ILIKE 'Publik%') THEN pct ELSE 0 END), 0) > 0
				           THEN LEAST(100.0, COALESCE(SUM(CASE WHEN (name ILIKE 'Masyarakat%' OR name ILIKE 'Publik%') THEN pct ELSE 0 END), 0))
				           ELSE GREATEST(0.0, LEAST(100.0, 100.0 - COALESCE(SUM(CASE WHEN (is_controlling = true OR pct >= 5.0) AND NOT (name ILIKE 'Masyarakat%' OR name ILIKE 'Publik%') THEN pct ELSE 0 END), 0)))
				       END AS free_float_pct
				FROM company_shareholders
				WHERE ticker = ANY($1)
				GROUP BY ticker
			) s ON t.ticker = s.ticker
			WHERE t.ticker = ANY($1)
		`, tickerList)
		if err != nil {
			return nil, err
		}
		for mRows.Next() {
			var tk string
			var cName, gName *string
			var ff *float64
			if err := mRows.Scan(&tk, &cName, &gName, &ff); err != nil {
				mRows.Close()
				return nil, err
			}
			cn := tk
			if cName != nil && *cName != "" {
				cn = *cName
			}
			gn := "-"
			if gName != nil && *gName != "" {
				gn = *gName
			}
			ffVal := 0.0
			if ff != nil {
				ffVal = shared.Round(*ff, 1)
			}
			metaMap[tk] = metaRec{compName: cn, grpName: gn, freeFloat: ffVal}
		}
		mRows.Close()

		rawResults = make([]ConfluenceItem, 0, len(tickerList))
		for _, tk := range tickerList {
			tFlow := flowByTicker[tk]
			if len(tFlow) == 0 {
				continue
			}
			tAct := actByTicker[tk]
			tPx := pxByTicker[tk]

			latestFlow := tFlow[len(tFlow)-1]
			rawSig := latestFlow.BandarSignal
			if rawSig == "" {
				rawSig = "NEUTRAL"
			}
			sig := shared.TitleWords(strings.ReplaceAll(rawSig, "_", " "))

			currPx := 0.0
			if len(tPx) > 0 {
				currPx = tPx[len(tPx)-1].Close
			}

			f5Start := max(0, len(tFlow)-5)
			foreign5D := 0.0
			totVals5D := make([]float64, 0, 5)
			flowInputs := make([]shared.FlowInput, len(tFlow))
			for i, fl := range tFlow {
				flowInputs[i] = shared.FlowInput{Date: fl.Date, ForeignNetBroker: fl.ForeignNetBroker}
				if i >= f5Start {
					foreign5D += fl.ForeignNetBroker
					totVals5D = append(totVals5D, fl.TotalValue)
				}
			}

			conv := shared.ConvictionScore(rawSig, &foreign5D, scan10D, tk, tPx, flowInputs, totVals5D)
			convScore := conv.Score

			bandarAvg := 0.0
			topBuyer := "-"
			topSeller := "-"
			smartForeignNet := 0.0
			gorenganNet := 0.0

			if len(tAct) > 0 {
				latestActDate := tAct[len(tAct)-1].Date
				var buyers, sellers []actScreenerRow
				for _, a := range tAct {
					if a.Date != latestActDate {
						continue
					}
					if a.NetValue > 0 {
						buyers = append(buyers, a)
					} else if a.NetValue < 0 {
						sellers = append(sellers, a)
					}
					if smartForeignSet[a.BrokerCode] {
						smartForeignNet += a.NetValue
					}
					if gorenganSet[a.BrokerCode] {
						gorenganNet += a.NetValue
					}
				}
				sort.SliceStable(buyers, func(i, j int) bool { return buyers[i].NetValue > buyers[j].NetValue })
				sort.SliceStable(sellers, func(i, j int) bool { return sellers[i].NetValue < sellers[j].NetValue })

				if len(buyers) > 0 {
					topBuyer = buyers[0].BrokerCode
					top5 := buyers
					if len(top5) > 5 {
						top5 = top5[:5]
					}
					var bVal, bLot float64
					for _, b := range top5 {
						bVal += b.BuyValue
						bLot += b.BuyLot
					}
					if bLot > 0 {
						bandarAvg = shared.Round(bVal/(bLot*100.0), 1)
					}
				}
				if len(sellers) > 0 {
					topSeller = sellers[0].BrokerCode
				}
			}

			if currPx == 0 && bandarAvg > 0 {
				currPx = bandarAvg
			}
			distPct := 0.0
			if bandarAvg > 0 {
				distPct = shared.Round(((currPx-bandarAvg)/bandarAvg)*100.0, 2)
			}

			m, okM := metaMap[tk]
			if !okM {
				m = metaRec{compName: tk, grpName: "-"}
			}

			fnStart := max(0, len(tFlow)-20)
			fnSeries := make([]float64, 0, len(tFlow)-fnStart)
			for _, fl := range tFlow[fnStart:] {
				fnSeries = append(fnSeries, fl.ForeignNetBroker)
			}
			latestFN := fnSeries[len(fnSeries)-1]
			fnMean, fnStd := shared.MeanStd(fnSeries)
			foreignZ := 0.0
			if len(fnSeries) >= 5 && fnStd > 0 {
				foreignZ = (latestFN - fnMean) / fnStd
			}
			foreignHMM := "Neutral / Transition"
			if latestFN > 0 {
				if foreignZ > 1.0 {
					foreignHMM = "Strong Accumulation"
				} else if foreignZ >= -0.3 {
					foreignHMM = "Accumulation"
				} else {
					foreignHMM = "Slowing Accumulation"
				}
			} else if latestFN < 0 {
				if foreignZ < -1.0 {
					foreignHMM = "Strong Distribution"
				} else if foreignZ <= 0.3 {
					foreignHMM = "Distribution"
				} else {
					foreignHMM = "Slowing Distribution"
				}
			}

			quickTags := make([]string, 0, 5)
			rawUpper := strings.ToUpper(rawSig)
			if convScore >= 70 && strings.Contains(rawUpper, "ACCUMULATION") && distPct >= -5 && distPct <= 5 {
				quickTags = append(quickTags, "GOLDEN_CONFLUENCE")
			}
			if strings.Contains(rawUpper, "ACCUMULATION") && distPct < 0 {
				quickTags = append(quickTags, "STEALTH_ACCUMULATION")
			}
			if gorenganNet > 500_000_000 && math.Abs(distPct) < 15 {
				quickTags = append(quickTags, "GORENGAN_MOMENTUM")
			}
			if currPx > bandarAvg && strings.Contains(rawUpper, "DISTRIBUTION") {
				quickTags = append(quickTags, "DISTRIBUTION_TRAP")
			}
			if m.freeFloat > 0 && m.freeFloat <= 25.0 {
				quickTags = append(quickTags, "LOW_FLOAT")
			}

			rawResults = append(rawResults, ConfluenceItem{
				Ticker:                 tk,
				CompanyName:            m.compName,
				GroupName:              m.grpName,
				CurrentPrice:           currPx,
				TotalValue:             latestFlow.TotalValue,
				Signal:                 sig,
				ConvictionScore:        convScore,
				BandarAvgPrice:         bandarAvg,
				BandarPriceDistancePct: distPct,
				Foreign5D:              foreign5D,
				SmartForeignNet:        smartForeignNet,
				GorenganNet:            gorenganNet,
				ForeignHMMState:        foreignHMM,
				FreeFloatEst:           m.freeFloat,
				TopBuyer:               topBuyer,
				TopSeller:              topSeller,
				QuickTags:              quickTags,
			})
		}

		sort.SliceStable(rawResults, func(i, j int) bool {
			return rawResults[i].ConvictionScore > rawResults[j].ConvictionScore
		})
		metaInfo = map[string]any{
			"universe_mode": universeMode,
			"window_days":   windowDays,
			"analysis_date": endDateStr,
			"cached":        false,
		}
		if cacheStore != nil {
			cachedMeta := map[string]any{
				"universe_mode": universeMode,
				"window_days":   windowDays,
				"analysis_date": endDateStr,
				"cached":        true,
			}
			if b, err := json.Marshal(rawConfluenceCache{Data: rawResults, Meta: cachedMeta}); err == nil {
				cacheStore.Set(ctx, rawCacheKey, b, 3600*time.Second)
			}
		}
	}

	filtered := make([]ConfluenceItem, 0, len(rawResults))
	tagUp := strings.ToUpper(strings.TrimSpace(tagFilter))
	sigLow := strings.ToLower(strings.TrimSpace(signalFilter))
	for _, x := range rawResults {
		if minConviction > 0 && x.ConvictionScore < minConviction {
			continue
		}
		if minLiquidity > 0 && x.TotalValue < minLiquidity {
			continue
		}
		if tagUp != "" {
			hasTag := false
			for _, tg := range x.QuickTags {
				if tg == tagUp {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}
		if sigLow != "" && !strings.Contains(strings.ToLower(x.Signal), sigLow) {
			continue
		}
		filtered = append(filtered, x)
	}

	return map[string]any{
		"count": len(filtered),
		"meta":  metaInfo,
		"data":  filtered,
	}, nil
}
