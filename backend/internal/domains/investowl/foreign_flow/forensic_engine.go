package foreign_flow

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func classifyCounterparty(participantType string, taxCategory *string) string {
	cat := ""
	if taxCategory != nil {
		cat = strings.ToLower(strings.TrimSpace(*taxCategory))
	}
	if cat == "retail" || cat == "semi_ritel" {
		return "retail"
	}
	if participantType == "Asing" {
		return "foreign"
	}
	return "local_inst"
}

func evaluateDivergence(priceChangePct, netFlowWindow, totalTurnoverWindow float64) (string, string, string, string) {
	if totalTurnoverWindow < 0 {
		totalTurnoverWindow = 0
	}
	deadband := math.Max(50_000_000.0, 0.005*totalTurnoverWindow)
	if math.Abs(netFlowWindow) < deadband {
		return "NEUTRAL", "Netral / Konsolidasi", "neutral",
			"Arus modal asing relatif kecil dan berada di bawah ambang materialitas sepanjang window."
	}
	if netFlowWindow > 0 {
		if priceChangePct <= 2.0 {
			return "SILENT_ABSORPTION", "Silent Absorption (Serok Bawah)", "positive",
				"Asing mengakumulasi bersih sementara harga masih tertekan atau sideways."
		}
		return "MARKUP_CONFIRMATION", "Markup Confirmation", "positive",
			"Kenaikan harga dikonfirmasi oleh akumulasi bersih modal asing."
	}
	if netFlowWindow < 0 {
		if priceChangePct >= -2.0 {
			return "HIDDEN_DISTRIBUTION", "Hidden Distribution", "warning",
				"Harga tampak ditahan/naik, namun modal asing mencatat distribusi bersih."
		}
		return "EXODUS_CAPITULATION", "Exodus / Kapitulasi Asing", "negative",
			"Penurunan harga disertai tekanan jual bersih dari investor asing."
	}
	return "NEUTRAL", "Netral / Seimbang", "neutral",
		"Arus modal asing berada pada posisi seimbang sepanjang window."
}

func evaluateCounterpartyVerdict(foreignTotal, localInstTotal, retailTotal float64) (string, string) {
	switch {
	case foreignTotal > 0 && localInstTotal > 0 && retailTotal < 0:
		return "Big Money Kompak Akumulasi ← Ritel", "positive"
	case foreignTotal < 0 && localInstTotal < 0 && retailTotal > 0:
		return "Big Money Kompak Keluar → Ritel", "negative"
	case foreignTotal > 0 && retailTotal < 0 && math.Abs(retailTotal) >= math.Abs(localInstTotal):
		return "Asing Akumulasi ← Ritel Melepas", "positive"
	case foreignTotal > 0 && localInstTotal < 0:
		return "Asing Akumulasi ← Institusi Lokal Melepas", "positive"
	case foreignTotal < 0 && retailTotal > 0 && retailTotal >= localInstTotal:
		return "Asing Distribusi → Ritel Menampung", "negative"
	case foreignTotal < 0 && localInstTotal > 0:
		return "Asing Distribusi → Institusi Lokal Menyerap", "warning"
	default:
		return "Tarik-Ulur Asing vs Institusi Lokal", "neutral"
	}
}

type forensicMacroRow struct {
	Date             time.Time
	ForeignNetBroker float64
	TotalValue       float64
	Close            float64
	PriceValue       float64
}

type forensicActRow struct {
	Date            string
	BrokerCode      string
	ParticipantType string
	BuyValue        float64
	SellValue       float64
	NetValue        float64
	BuyLot          float64
	SellLot         float64
	BrokerName      string
	TaxCategory     *string
}

type kseiNormRow struct {
	PeriodDate          string  `json:"period_date"`
	PctAsing            float64 `json:"pct_asing"`
	PctAsingGt1         float64 `json:"pct_asing_gt1"`
	PctAsingLt1         float64 `json:"pct_asing_lt1"`
	PctRitelPure        float64 `json:"pct_ritel_pure"`
	PctLocalInstitution float64 `json:"pct_local_institution"`
}

type t5RecentRow struct {
	Date         string  `json:"date"`
	InvestorName string  `json:"investor_name"`
	Broker       string  `json:"broker"`
	Action       string  `json:"action"`
	ChangeLot    float64 `json:"change_lot"`
	ValueM       float64 `json:"value_m"`
	Price        float64 `json:"price"`
}

// ComputeForeignForensic is an exact 1:1 port of services/foreign_forensic.py.
func ComputeForeignForensic(ctx context.Context, db *pgxpool.Pool, ticker string, lookbackDays int) (map[string]any, error) {
	tk := strings.ToUpper(strings.TrimSpace(ticker))

	var metaName, metaSector, metaGroups *string
	err := db.QueryRow(ctx, `
		SELECT t.name, t.sector, cg.groups
		FROM public.tickers t
		LEFT JOIN (
			SELECT ticker, STRING_AGG(DISTINCT group_name, ' / ' ORDER BY group_name) AS groups
			FROM ownership_staging.v_conglomerate_members
			WHERE ticker = $1
			GROUP BY ticker
		) cg ON cg.ticker = t.ticker
		WHERE t.ticker = $1
	`, tk).Scan(&metaName, &metaSector, &metaGroups)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	companyName := fmt.Sprintf("PT %s Tbk", tk)
	if metaName != nil && strings.TrimSpace(*metaName) != "" {
		companyName = strings.TrimSpace(*metaName)
	}
	sectorName := ""
	if metaSector != nil {
		sectorName = strings.TrimSpace(*metaSector)
	}
	groupName := "Independen / Non-Afiliasi"
	if metaGroups != nil && strings.TrimSpace(*metaGroups) != "" {
		groupName = strings.TrimSpace(*metaGroups)
	} else if sectorName != "" {
		groupName = sectorName
	}

	macroQ, err := db.Query(ctx, `
		SELECT bf.date,
		       bf.foreign_net_broker::double precision,
		       bf.total_value::double precision,
		       p.close::double precision,
		       p.value::double precision AS price_value
		FROM public.broker_flow bf
		LEFT JOIN public.prices p ON bf.ticker = p.ticker AND bf.date = p.date
		WHERE bf.ticker = $1
		ORDER BY bf.date DESC
		LIMIT $2
	`, tk, lookbackDays)
	if err != nil {
		return nil, err
	}
	defer macroQ.Close()

	macroDesc := make([]forensicMacroRow, 0, lookbackDays)
	for macroQ.Next() {
		var dt time.Time
		var fnb, tv, cl, pv *float64
		if err := macroQ.Scan(&dt, &fnb, &tv, &cl, &pv); err != nil {
			return nil, err
		}
		macroDesc = append(macroDesc, forensicMacroRow{
			Date:             dt,
			ForeignNetBroker: shared.Finite(fnb, 0),
			TotalValue:       shared.Finite(tv, 0),
			Close:            shared.Finite(cl, 0),
			PriceValue:       shared.Finite(pv, 0),
		})
	}
	if err := macroQ.Err(); err != nil {
		return nil, err
	}
	if len(macroDesc) < 10 {
		return nil, fmt.Errorf("Riwayat transaksi emiten '%s' belum mencukupi untuk analitik forensik (min. 10 hari bursa).", tk)
	}

	nMacro := len(macroDesc)
	macroAsc := make([]forensicMacroRow, nMacro)
	for i := 0; i < nMacro; i++ {
		macroAsc[i] = macroDesc[nMacro-1-i]
	}

	dates := make([]string, nMacro)
	foreignNet := make([]float64, nMacro)
	totalValues := make([]float64, nMacro)
	priceValues := make([]float64, nMacro)
	closes := make([]float64, nMacro)

	lastClose := 0.0
	firstPosClose := 0.0
	for i, r := range macroAsc {
		dates[i] = r.Date.Format("2006-01-02")
		foreignNet[i] = r.ForeignNetBroker
		totalValues[i] = r.TotalValue
		priceValues[i] = r.PriceValue
		if r.Close > 0 {
			lastClose = r.Close
		}
		closes[i] = lastClose
		if firstPosClose == 0 && lastClose > 0 {
			firstPosClose = lastClose
		}
	}
	if firstPosClose > 0 {
		for i := range closes {
			if closes[i] <= 0 {
				closes[i] = firstPosClose
			}
		}
	}

	startDate := macroAsc[0].Date
	endDate := macroAsc[nMacro-1].Date
	latestClose := closes[nMacro-1]

	actQ, err := db.Query(ctx, `
		SELECT ba.date, ba.broker_code, ba.participant_type,
		       ba.buy_value::double precision, ba.sell_value::double precision,
		       ba.net_value::double precision, ba.buy_lot::double precision,
		       ba.sell_lot::double precision,
		       bt.name AS broker_name, bt.category AS tax_category
		FROM public.broker_activity ba
		LEFT JOIN public.broker_taxonomy bt ON bt.code = ba.broker_code
		WHERE ba.ticker = $1
		  AND ba.date >= $2
		  AND ba.date <= $3
		ORDER BY ba.date ASC
	`, tk, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer actQ.Close()

	foreignDailyBuyVal := make(map[string]float64, nMacro)
	foreignDailyBuyLot := make(map[string]float64, nMacro)
	foreignDailyTurnover := make(map[string]float64, nMacro)
	type cpBucketVals struct {
		foreign, localInst, retail float64
	}
	cpDaily := make(map[string]*cpBucketVals, nMacro)
	for _, d := range dates {
		cpDaily[d] = &cpBucketVals{}
	}

	type brokerAgg struct {
		code, name                         string
		netVal, netLot, buyVal, sellVal    float64
		buyLot, sellLot                    float64
		vwapBuyVal, vwapBuyLot             float64
		vwapSellVal, vwapSellLot           float64
		dates                              map[string]struct{}
	}
	foreignBrokerAgg := make(map[string]*brokerAgg)
	microRowsAsing := make([]MicroRow, 0, 128)

	for actQ.Next() {
		var dt time.Time
		var bCode, pType, bNamePtr, taxCat *string
		var bValPtr, sValPtr, nValPtr, bLotPtr, sLotPtr *float64
		if err := actQ.Scan(&dt, &bCode, &pType, &bValPtr, &sValPtr, &nValPtr, &bLotPtr, &sLotPtr, &bNamePtr, &taxCat); err != nil {
			return nil, err
		}
		dStr := dt.Format("2006-01-02")
		code := ""
		if bCode != nil {
			code = strings.TrimSpace(*bCode)
		}
		ptype := ""
		if pType != nil {
			ptype = strings.TrimSpace(*pType)
		}
		bName := code
		if bNamePtr != nil && strings.TrimSpace(*bNamePtr) != "" {
			bName = strings.TrimSpace(*bNamePtr)
		}
		buyVal := shared.Finite(bValPtr, 0)
		sellVal := shared.Finite(sValPtr, 0)
		netVal := shared.Finite(nValPtr, 0)
		buyLot := shared.Finite(bLotPtr, 0)
		sellLot := shared.Finite(sLotPtr, 0)

		bucket := classifyCounterparty(ptype, taxCat)
		if cp, ok := cpDaily[dStr]; ok {
			switch bucket {
			case "foreign":
				cp.foreign += netVal
			case "retail":
				cp.retail += netVal
			default:
				cp.localInst += netVal
			}
		}
		if bucket == "foreign" {
			if _, ok := cpDaily[dStr]; ok {
				if buyLot > 0 && buyVal > 0 {
					foreignDailyBuyVal[dStr] += buyVal
					foreignDailyBuyLot[dStr] += buyLot
				}
				foreignDailyTurnover[dStr] += buyVal + sellVal
			}
			microRowsAsing = append(microRowsAsing, MicroRow{Date: dStr, BrokerCode: code, NetValue: netVal})

			agg, ok := foreignBrokerAgg[code]
			if !ok {
				agg = &brokerAgg{code: code, name: bName, dates: make(map[string]struct{})}
				foreignBrokerAgg[code] = agg
			}
			agg.netVal += netVal
			agg.netLot += buyLot - sellLot
			agg.buyVal += buyVal
			agg.sellVal += sellVal
			agg.buyLot += buyLot
			agg.sellLot += sellLot
			if buyLot > 0 && buyVal > 0 {
				agg.vwapBuyVal += buyVal
				agg.vwapBuyLot += buyLot
			}
			if sellLot > 0 && sellVal > 0 {
				agg.vwapSellVal += sellVal
				agg.vwapSellLot += sellLot
			}
			agg.dates[dStr] = struct{}{}
		}
	}
	if err := actQ.Err(); err != nil {
		return nil, err
	}

	// Module 1: Cost Basis
	cumForeignNet := make([]float64, nMacro)
	foreignVWAP := make([]*float64, nMacro)
	var runNet, cumBuyVal, cumBuyShares, peakCumNet float64
	for i, dStr := range dates {
		runNet += foreignNet[i]
		rNet := shared.Round(runNet, 2)
		cumForeignNet[i] = rNet
		if rNet > peakCumNet {
			peakCumNet = rNet
		}
		cumBuyVal += foreignDailyBuyVal[dStr]
		cumBuyShares += foreignDailyBuyLot[dStr] * 100.0
		if cumBuyShares > 0 && cumBuyVal > 0 {
			vw := cumBuyVal / cumBuyShares
			if !math.IsNaN(vw) && !math.IsInf(vw, 0) {
				v := shared.Round(vw, 2)
				foreignVWAP[i] = &v
			}
		}
	}

	var totalBuyLot, totalSellLot, totalVwapBuyVal, totalVwapBuyLot, totalVwapSellVal, totalVwapSellLot float64
	aggSlice := make([]*brokerAgg, 0, len(foreignBrokerAgg))
	for _, agg := range foreignBrokerAgg {
		aggSlice = append(aggSlice, agg)
		totalBuyLot += agg.buyLot
		totalSellLot += agg.sellLot
		totalVwapBuyVal += agg.vwapBuyVal
		totalVwapBuyLot += agg.vwapBuyLot
		totalVwapSellVal += agg.vwapSellVal
		totalVwapSellLot += agg.vwapSellLot
	}

	avgBuyPrice := 0.0
	if totalVwapBuyLot > 0 && totalVwapBuyVal > 0 {
		avgBuyPrice = shared.Round(totalVwapBuyVal/(totalVwapBuyLot*100.0), 2)
	}
	avgSellPrice := 0.0
	if totalVwapSellLot > 0 && totalVwapSellVal > 0 {
		avgSellPrice = shared.Round(totalVwapSellVal/(totalVwapSellLot*100.0), 2)
	}
	priceVsVWAPPct := 0.0
	if avgBuyPrice > 0 && latestClose > 0 {
		priceVsVWAPPct = shared.Round(((latestClose-avgBuyPrice)/avgBuyPrice)*100.0, 2)
	}

	var totalNetValue, totalTurnoverWindow float64
	for i := range foreignNet {
		totalNetValue += foreignNet[i]
		totalTurnoverWindow += totalValues[i]
	}
	var floatingPnLPtr *float64
	if totalNetValue > 0 {
		v := priceVsVWAPPct
		floatingPnLPtr = &v
	}
	totalNetLot := shared.Round(totalBuyLot-totalSellLot, 2)
	retentionPct := 0.0
	if peakCumNet > 0 {
		retentionPct = shared.Round(math.Max(0.0, math.Min(100.0, (cumForeignNet[nMacro-1]/peakCumNet)*100.0)), 1)
	}
	positionStatus := "Underwater (Rugi)"
	switch {
	case totalNetValue < 0:
		positionStatus = "Distribusi / Net Sell"
	case totalNetValue == 0:
		positionStatus = "Netral / Seimbang"
	case priceVsVWAPPct >= 0:
		positionStatus = "Floating Profit"
	}

	formatBroker := func(it *brokerAgg) ForensicBrokerEntry {
		bAvg := 0.0
		if it.vwapBuyLot > 0 && it.vwapBuyVal > 0 {
			bAvg = shared.Round(it.vwapBuyVal/(it.vwapBuyLot*100.0), 2)
		}
		sAvg := 0.0
		if it.vwapSellLot > 0 && it.vwapSellVal > 0 {
			sAvg = shared.Round(it.vwapSellVal/(it.vwapSellLot*100.0), 2)
		}
		var pnl *float64
		if it.netVal > 0 {
			p := 0.0
			if bAvg > 0 && latestClose > 0 {
				p = shared.Round(((latestClose-bAvg)/bAvg)*100.0, 2)
			}
			pnl = &p
		}
		return ForensicBrokerEntry{
			BrokerCode:     it.code,
			BrokerName:     it.name,
			NetValue:       shared.Round(it.netVal, 2),
			NetLot:         shared.Round(it.netLot, 0),
			BuyValue:       shared.Round(it.buyVal, 2),
			SellValue:      shared.Round(it.sellVal, 2),
			AvgBuyPrice:    bAvg,
			AvgSellPrice:   sAvg,
			FloatingPnLPct: pnl,
			ActiveDays:     len(it.dates),
		}
	}

	holdersRaw := make([]*brokerAgg, 0)
	sellersRaw := make([]*brokerAgg, 0)
	for _, agg := range aggSlice {
		if agg.netVal > 0 {
			holdersRaw = append(holdersRaw, agg)
		} else if agg.netVal < 0 {
			sellersRaw = append(sellersRaw, agg)
		}
	}
	sort.SliceStable(holdersRaw, func(i, j int) bool { return holdersRaw[i].netVal > holdersRaw[j].netVal })
	sort.SliceStable(sellersRaw, func(i, j int) bool { return sellersRaw[i].netVal < sellersRaw[j].netVal })
	if len(holdersRaw) > 5 {
		holdersRaw = holdersRaw[:5]
	}
	if len(sellersRaw) > 5 {
		sellersRaw = sellersRaw[:5]
	}
	topHolders := make([]ForensicBrokerEntry, 0, len(holdersRaw))
	for _, item := range holdersRaw {
		topHolders = append(topHolders, formatBroker(item))
	}
	topSellers := make([]ForensicBrokerEntry, 0, len(sellersRaw))
	for _, item := range sellersRaw {
		topSellers = append(topSellers, formatBroker(item))
	}

	// Module 2: Impact & Divergence (gracefully optional terminal_engine schema)
	var fpBadge *string
	var fpScorePtr, avgTrans10DPtr *float64
	tePartMap := make(map[string]float64)

	var bTxt *string
	var sNum, tNum *float64
	if err := db.QueryRow(ctx, `
		SELECT screener_summary->>'fpBadge',
		       NULLIF(screener_summary->>'fpScoreLatest', '')::double precision,
		       NULLIF(screener_summary->>'avgTransValue10d', '')::double precision
		FROM terminal_engine.foreign_flow_screener
		WHERE ticker = $1
	`, tk).Scan(&bTxt, &sNum, &tNum); err == nil {
		if bTxt != nil && *bTxt != "" {
			fpBadge = bTxt
		}
		if sNum != nil {
			v := shared.Round(*sNum, 1)
			fpScorePtr = &v
		}
		if tNum != nil {
			v := shared.Round(*tNum, 2)
			avgTrans10DPtr = &v
		}
	}

	if teRows, err := db.Query(ctx, `
		SELECT date_iso, partisipasi_asing::double precision
		FROM terminal_engine.foreign_flow_daily
		WHERE ticker = $1 AND date_iso >= $2 AND date_iso <= $3
		ORDER BY date_iso ASC
	`, tk, startDate, endDate); err == nil {
		for teRows.Next() {
			var dISO time.Time
			var pPtr *float64
			if teRows.Scan(&dISO, &pPtr) == nil && pPtr != nil && !math.IsNaN(*pPtr) {
				tePartMap[dISO.Format("2006-01-02")] = shared.Round(*pPtr, 2)
			}
		}
		teRows.Close()
	}

	partSeries := make([]float64, nMacro)
	var sumPart float64
	for i, dStr := range dates {
		if v, ok := tePartMap[dStr]; ok {
			partSeries[i] = v
		} else {
			pxVal := priceValues[i]
			if pxVal <= 0 {
				pxVal = totalValues[i]
			}
			if pxVal > 0 {
				est := math.Min(100.0, (foreignDailyTurnover[dStr]/(2.0*pxVal))*100.0)
				partSeries[i] = shared.Round(est, 2)
			}
		}
		sumPart += partSeries[i]
	}
	avgPartPct := shared.Round(sumPart/float64(nMacro), 2)
	resolvedBadge := "LFP"
	if fpBadge != nil {
		resolvedBadge = *fpBadge
	} else if avgPartPct >= 30.0 {
		resolvedBadge = "HFP"
	} else if avgPartPct >= 10.0 {
		resolvedBadge = "MFP"
	}

	retPairsX := make([]float64, 0, nMacro)
	retPairsY := make([]float64, 0, nMacro)
	for i := 1; i < nMacro; i++ {
		if closes[i-1] > 0 && closes[i] > 0 {
			dailyRet := (closes[i] - closes[i-1]) / closes[i-1]
			denom := priceValues[i]
			if denom <= 0 {
				denom = totalValues[i]
			}
			flowNorm := foreignNet[i]
			if denom > 0 {
				flowNorm = foreignNet[i] / denom
			}
			retPairsX = append(retPairsX, flowNorm)
			retPairsY = append(retPairsY, dailyRet)
		}
	}
	foreignPriceCorr := math.Max(-1.0, math.Min(1.0, shared.PearsonCorr(retPairsX, retPairsY)))
	priceChangePct := 0.0
	if closes[0] > 0 {
		priceChangePct = shared.Round(((closes[nMacro-1]-closes[0])/closes[0])*100.0, 2)
	}

	buyDays, sellDays := 0, 0
	var absSum, maxSingleAbs float64
	for _, v := range foreignNet {
		if v > 0 {
			buyDays++
		} else if v < 0 {
			sellDays++
		}
		av := math.Abs(v)
		absSum += av
		if av > maxSingleAbs {
			maxSingleAbs = av
		}
	}
	maxDomDays := buyDays
	if sellDays > maxDomDays {
		maxDomDays = sellDays
	}
	consistencyPct := shared.Round((float64(maxDomDays)/float64(nMacro))*100.0, 1)

	streakDays := 0
	for i := nMacro - 1; i >= 0; i-- {
		x := foreignNet[i]
		if x > 0 && streakDays >= 0 {
			streakDays++
		} else if x < 0 && streakDays <= 0 {
			streakDays--
		} else {
			break
		}
	}

	maxSingleShare := 0.0
	if absSum > 0 {
		maxSingleShare = shared.Round((maxSingleAbs/absSum)*100.0, 1)
	}
	crossingWarning := absSum > 0 && nMacro >= 10 && maxSingleShare > 45.0 && maxSingleAbs >= 1_000_000_000.0

	divState, divLabel, divTone, divSummary := evaluateDivergence(priceChangePct, totalNetValue, totalTurnoverWindow)
	fMean, fStd := shared.MeanStd(foreignNet)
	zscore := 0.0
	if fStd > 0 {
		zscore = (foreignNet[nMacro-1] - fMean) / fStd
		if math.IsNaN(zscore) || math.IsInf(zscore, 0) {
			zscore = 0.0
		}
	}

	// Module 3: Counterparty Tug-of-War
	cpF := make([]float64, nMacro)
	cpLI := make([]float64, nMacro)
	cpR := make([]float64, nMacro)
	cumF := make([]float64, nMacro)
	cumLI := make([]float64, nMacro)
	cumR := make([]float64, nMacro)
	var accF, accLI, accR float64
	for i, dStr := range dates {
		row := cpDaily[dStr]
		fV := shared.Round(row.foreign, 2)
		liV := shared.Round(row.localInst, 2)
		rV := shared.Round(row.retail, 2)
		cpF[i], cpLI[i], cpR[i] = fV, liV, rV
		accF += fV
		accLI += liV
		accR += rV
		cumF[i] = shared.Round(accF, 2)
		cumLI[i] = shared.Round(accLI, 2)
		cumR[i] = shared.Round(accR, 2)
	}
	matchupVerdict, matchupTone := evaluateCounterpartyVerdict(accF, accLI, accR)
	brokerNetSlice := make([]float64, 0, len(aggSlice))
	for _, a := range aggSlice {
		brokerNetSlice = append(brokerNetSlice, a.netVal)
	}
	hhiScore := ComputeForeignHHI(brokerNetSlice)
	heatmap := ComputeBrokerHeatmap(microRowsAsing, 15)

	// Module 4: KSEI Reality Check & T5 >1%
	kseiHist := make([]kseiNormRow, 0, 6)
	if kr, err := db.Query(ctx, `
		SELECT period_date, pct_asing::double precision, pct_asing_gt1::double precision,
		       pct_asing_lt1::double precision, pct_ritel_pure::double precision,
		       pct_local_institution::double precision
		FROM ownership_staging.v_composition_ksei_normalized
		WHERE ticker = $1
		ORDER BY period_date DESC
		LIMIT 6
	`, tk); err == nil {
		descK := make([]kseiNormRow, 0, 6)
		for kr.Next() {
			var dt time.Time
			var pa, pgt, plt, pr, pli *float64
			if kr.Scan(&dt, &pa, &pgt, &plt, &pr, &pli) == nil {
				descK = append(descK, kseiNormRow{
					PeriodDate:          dt.Format("2006-01-02"),
					PctAsing:            shared.Round(shared.Finite(pa, 0), 2),
					PctAsingGt1:         shared.Round(shared.Finite(pgt, 0), 2),
					PctAsingLt1:         shared.Round(shared.Finite(plt, 0), 2),
					PctRitelPure:        shared.Round(shared.Finite(pr, 0), 2),
					PctLocalInstitution: shared.Round(shared.Finite(pli, 0), 2),
				})
			}
		}
		kr.Close()
		for i := len(descK) - 1; i >= 0; i-- {
			kseiHist = append(kseiHist, descK[i])
		}
	}

	t5Recent := make([]t5RecentRow, 0, 8)
	if tr, err := db.Query(ctx, `
		SELECT movement_date, investor_name, broker, action,
		       change_lot::double precision, value_m::double precision, price::double precision
		FROM ownership_staging.v_t5_transactions
		WHERE ticker = $1
		ORDER BY movement_date DESC
		LIMIT 8
	`, tk); err == nil {
		for tr.Next() {
			var dt time.Time
			var inv, brk, act *string
			var clot, valM, prc *float64
			if tr.Scan(&dt, &inv, &brk, &act, &clot, &valM, &prc) == nil {
				invS, brkS, actS := "-", "-", "-"
				if inv != nil && *inv != "" {
					invS = *inv
				}
				if brk != nil && *brk != "" {
					brkS = *brk
				}
				if act != nil && *act != "" {
					actS = *act
				}
				t5Recent = append(t5Recent, t5RecentRow{
					Date:         dt.Format("2006-01-02"),
					InvestorName: invS,
					Broker:       brkS,
					Action:       actS,
					ChangeLot:    shared.Round(shared.Finite(clot, 0), 0),
					ValueM:       shared.Round(shared.Finite(valM, 0), 2),
					Price:        shared.Round(shared.Finite(prc, 0), 0),
				})
			}
		}
		tr.Close()
	}

	var latestPeriod *string
	var latestPctAsing, pctAsingGt1, pctAsingLt1, latestPctRitel, delta1m, delta3m *float64
	if len(kseiHist) > 0 {
		last := kseiHist[len(kseiHist)-1]
		latestPeriod = &last.PeriodDate
		latestPctAsing = &last.PctAsing
		pctAsingGt1 = &last.PctAsingGt1
		pctAsingLt1 = &last.PctAsingLt1
		latestPctRitel = &last.PctRitelPure
	}
	if len(kseiHist) >= 2 {
		v := shared.Round(kseiHist[len(kseiHist)-1].PctAsing-kseiHist[len(kseiHist)-2].PctAsing, 2)
		delta1m = &v
	}
	if len(kseiHist) >= 4 {
		v := shared.Round(kseiHist[len(kseiHist)-1].PctAsing-kseiHist[len(kseiHist)-4].PctAsing, 2)
		delta3m = &v
	}

	alignmentStatus := "Divergen (Rotasi/Arbitrase Broker)"
	switch {
	case delta1m == nil || totalNetValue == 0:
		alignmentStatus = "Data KSEI Terbatas"
	case math.Abs(*delta1m) < 0.05:
		alignmentStatus = "KSEI Relatif Stabil"
	case (totalNetValue > 0 && *delta1m > 0) || (totalNetValue < 0 && *delta1m < 0):
		alignmentStatus = "Searah (Konfirmasi KSEI)"
	}

	displaySector := "-"
	if sectorName != "" {
		displaySector = sectorName
	}

	return map[string]any{
		"ticker": tk,
		"company": map[string]any{
			"name":   companyName,
			"group":  groupName,
			"sector": displaySector,
		},
		"lookback_days": lookbackDays,
		"latest_date":   dates[nMacro-1],
		"features": map[string]any{
			"foreign_hhi":         shared.Round(hhiScore, 4),
			"foreign_zscore":      shared.Round(zscore, 2),
			"foreign_power_badge": resolvedBadge,
			"foreign_power_score": fpScorePtr,
			"foreign_avg_cost":    avgBuyPrice,
			"floating_pnl_pct":    floatingPnLPtr,
			"divergence_state":    divState,
			"streak_days":         streakDays,
		},
		"timeseries": map[string]any{
			"dates":             dates,
			"close_prices":      closes,
			"foreign_net":       foreignNet,
			"cum_foreign_net":   cumForeignNet,
			"foreign_vwap":      foreignVWAP,
			"participation_pct": partSeries,
		},
		"models": map[string]any{
			"broker_heatmap": heatmap,
		},
		"cost_basis": map[string]any{
			"latest_close":            shared.Round(latestClose, 2),
			"avg_buy_price":           avgBuyPrice,
			"avg_sell_price":          avgSellPrice,
			"floating_pnl_pct":        floatingPnLPtr,
			"price_vs_vwap_pct":       priceVsVWAPPct,
			"total_net_value":         shared.Round(totalNetValue, 2),
			"total_net_lot":           totalNetLot,
			"peak_cum_net_value":      shared.Round(peakCumNet, 2),
			"inventory_retention_pct": retentionPct,
			"position_status":         positionStatus,
			"top_holders":             topHolders,
			"top_sellers":             topSellers,
		},
		"impact_divergence": map[string]any{
			"divergence_state":         divState,
			"divergence_label":         divLabel,
			"divergence_tone":          divTone,
			"summary_text":             divSummary,
			"foreign_power_badge":      resolvedBadge,
			"foreign_power_score":      fpScorePtr,
			"avg_trans_10d_m":          avgTrans10DPtr,
			"avg_participation_pct":    avgPartPct,
			"foreign_price_corr":       foreignPriceCorr,
			"price_change_pct":         priceChangePct,
			"net_flow_window":          shared.Round(totalNetValue, 2),
			"buy_days":                 buyDays,
			"sell_days":                sellDays,
			"total_days":               nMacro,
			"streak_days":              streakDays,
			"consistency_pct":          consistencyPct,
			"max_single_day_share_pct": maxSingleShare,
			"crossing_warning":         crossingWarning,
		},
		"counterparty": map[string]any{
			"dates":          dates,
			"foreign_net":    cpF,
			"local_inst_net": cpLI,
			"retail_net":     cpR,
			"cum_foreign":    cumF,
			"cum_local_inst": cumLI,
			"cum_retail":     cumR,
			"totals": map[string]any{
				"foreign":             shared.Round(accF, 2),
				"local_inst":          shared.Round(accLI, 2),
				"retail":              shared.Round(accR, 2),
				"unreported_residual": shared.Round(-(accF + accLI + accR), 2),
			},
			"matchup_verdict": matchupVerdict,
			"matchup_tone":    matchupTone,
			"foreign_hhi":     shared.Round(hhiScore, 4),
		},
		"ksei_reality": map[string]any{
			"latest_period":         latestPeriod,
			"latest_pct_asing":      latestPctAsing,
			"pct_asing_gt1":         pctAsingGt1,
			"pct_asing_lt1":         pctAsingLt1,
			"latest_pct_ritel_pure": latestPctRitel,
			"delta_1m_pp":           delta1m,
			"delta_3m_pp":           delta3m,
			"alignment_status":      alignmentStatus,
			"history":               kseiHist,
			"t5_recent":             t5Recent,
		},
	}, nil
}
