package ipohunter

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

func pctPtr(num, denom float64) *float64 {
	if denom == 0 {
		return nil
	}
	v := (num / denom) * 100.0
	return &v
}

func categoryLabelForTBS(broker string, finalRetailSet map[string]struct{}, brokerCategory map[string]string) string {
	cat := brokerCategory[broker]
	if cat == "uw" {
		return "Underwriter"
	}
	if _, ok := finalRetailSet[broker]; ok {
		return "Retail"
	}
	switch cat {
	case "semi":
		return "Semi-Ritel"
	case "institusi":
		return "Institusi"
	case "nonRitel":
		return "Non Ritel"
	case "bandar":
		return "Dipakai Bandar"
	default:
		return "Others"
	}
}

func tbsStyleKey(broker string, finalRetailSet, bandarSet map[string]struct{}, brokerCategory map[string]string) string {
	cat := brokerCategory[broker]
	if cat == "uw" {
		return "uw"
	}
	if cat == "institusi" {
		return "institusi"
	}
	if _, ok := finalRetailSet[broker]; ok {
		return "retail"
	}
	if cat == "nonRitel" {
		return "nonRitel"
	}
	_, isBandar := bandarSet[broker]
	if cat == "semi" || isBandar {
		return "semi"
	}
	return "others"
}

// BuildAnalysisCards produces all 7 Analisis cards matching Creative Trader and Python engine outputs.
func BuildAnalysisCards(
	ticker, mode string,
	master *MasterRecord,
	agg *Aggregation,
	decision Decision,
	analysisDays []DayRow,
	afterInfo *AfterInfo,
	holidays map[int]struct{},
	maxHolidayYear int,
) AnalysisCards {
	uwList := master.Underwriters
	uwSet := UWSetOf(master)
	bt := agg.BrokerTotals
	tv := agg.TotalVolume
	tn := agg.TotNetVolume
	denom := agg.TotBuy + agg.TotSell

	bandarSet := make(map[string]struct{}, len(decision.RetailBandarBrokers))
	for _, b := range decision.RetailBandarBrokers {
		bandarSet[b] = struct{}{}
	}
	bandarSorted := append([]string(nil), decision.RetailBandarBrokers...)
	sort.Strings(bandarSorted)

	lastClose := LastCloseIn(analysisDays)
	vwap := AfterVWAP(analysisDays)
	valuePrice := lastClose
	if mode == "after" && vwap != nil && *vwap > 0 {
		valuePrice = vwap
	}

	// 1. Top Buyer & Top Seller
	formatTBSRow := func(x BrokerNet) TBSRow {
		var valJt *float64
		if valuePrice != nil {
			v := math.Abs(x.Net) * 100.0 * (*valuePrice) / 1e6
			valJt = &v
		}
		_, isUW := uwSet[x.Broker]
		return TBSRow{
			Broker:        x.Broker,
			Net:           x.Net,
			ValueJt:       valJt,
			CategoryLabel: categoryLabelForTBS(x.Broker, decision.FinalRetailSet, decision.BrokerCategory),
			StyleKey:      tbsStyleKey(x.Broker, decision.FinalRetailSet, bandarSet, decision.BrokerCategory),
			IsUW:          isUW,
		}
	}
	topBuyers := make([]TBSRow, 0, len(agg.Buyers))
	for _, x := range agg.Buyers {
		topBuyers = append(topBuyers, formatTBSRow(x))
	}
	topSellers := make([]TBSRow, 0, len(agg.Sellers))
	for _, x := range agg.Sellers {
		topSellers = append(topSellers, formatTBSRow(x))
	}
	maxRows := min(10, max(len(topBuyers), len(topSellers)))
	accDist := make([]float64, 0, maxRows)
	for i := 0; i < maxRows; i++ {
		accDist = append(accDist, CumulativeDiffAtN(agg.Buyers, agg.Sellers, i+1))
	}

	// 2. Underwriter Card
	uwRows := make([]UnderwriterCardRow, 0, len(uwList))
	uwSumTrans := 0.0
	uwSumNet := 0.0
	for _, u := range uwList {
		t := bt[u.Broker]
		trans := t.Buy + math.Abs(t.Sell)
		netV := t.Buy - t.Sell
		uwSumTrans += trans
		uwSumNet += netV
		uwRows = append(uwRows, UnderwriterCardRow{
			Broker:   u.Broker,
			Role:     u.Role,
			Net:      netV,
			PctNet:   pctPtr(netV, tn),
			Trans:    trans,
			PctTrans: pctPtr(trans, tv),
		})
	}
	var lotIPO *float64
	if master.OfferingShares != nil {
		v := float64(*master.OfferingShares) / 100.0
		lotIPO = &v
	}
	var pctOfIPO *float64
	if lotIPO != nil && *lotIPO != 0 {
		v := (uwSumNet / *lotIPO) * 100.0
		pctOfIPO = &v
	}
	uwStatus := "NET SELL"
	if uwSumNet >= 0 {
		uwStatus = "NET BUY"
	}
	uwCard := UnderwriterCard{
		HasUW:         len(uwList) > 0,
		TotalNet:      uwSumNet,
		TotalPctNet:   pctPtr(uwSumNet, tn),
		TotalTrans:    uwSumTrans,
		TotalPctTrans: pctPtr(uwSumTrans, tv),
		LotIPO:        lotIPO,
		PctOfIPO:      pctOfIPO,
		Status:        uwStatus,
		Rows:          uwRows,
	}

	// 3. Total Transaksi per Kategori Card
	catDefs := []struct {
		Key   string
		Label string
	}{
		{"uw", "Underwriter"},
		{"retail", "Retail"},
		{"bandar", "Ritel dipakai Bandar"},
		{"semi", "Semi Retail"},
		{"institusi", "Institusi"},
		{"nonRitel", "Non Retail"},
		{"repoOthers", "Repo + Others"},
	}
	catSums := make(map[string]float64, len(catDefs))
	catNets := make(map[string]float64, len(catDefs))
	for b, t := range bt {
		c := CategoryOf(b, ticker, uwSet)
		_, isBandar := bandarSet[b]
		_, isRepo := RepoBrokers[b]
		if c != "uw" && isBandar {
			c = "bandar"
		} else if c != "uw" && (isRepo || c == "zombie") {
			c = "repoOthers"
		}
		if c == "zombie" {
			c = "repoOthers"
		}
		catSums[c] += t.Buy + math.Abs(t.Sell)
		catNets[c] += t.Buy - t.Sell
	}
	catRows := make([]CategoryCardRow, 0, len(catDefs))
	totCatTrans := 0.0
	totCatNet := 0.0
	for _, def := range catDefs {
		if def.Key == "uw" && len(uwList) == 0 {
			continue
		}
		if def.Key == "bandar" && len(bandarSet) == 0 {
			continue
		}
		sVal := catSums[def.Key]
		nVal := catNets[def.Key]
		totCatTrans += sVal
		totCatNet += nVal
		brokers := []string{}
		if def.Key == "bandar" {
			brokers = bandarSorted
		}
		catRows = append(catRows, CategoryCardRow{
			Key:           def.Key,
			Label:         def.Label,
			BandarBrokers: brokers,
			Net:           nVal,
			PctNet:        pctPtr(nVal, tn),
			Trans:         sVal,
			PctTrans:      pctPtr(sVal, denom),
		})
	}
	categoryCard := CategoryCard{
		Rows: catRows,
		Total: CategoryCardTotal{
			Net:      totCatNet,
			PctNet:   pctPtr(totCatNet, tn),
			Trans:    totCatTrans,
			PctTrans: pctPtr(totCatTrans, denom),
		},
	}

	// 4. Shared Pie Categories (Buyer vs Seller Donut + After IPO Distribution)
	pieNets := make(map[string]float64, len(PieKeys))
	pieBuy := make(map[string]float64, len(PieKeys))
	pieSell := make(map[string]float64, len(PieKeys))
	pieNBuy := make(map[string]float64, len(PieKeys))
	pieNSell := make(map[string]float64, len(PieKeys))
	type brokerAmt struct {
		broker string
		val    float64
	}
	nbBuyList := make(map[string][]brokerAmt, len(PieKeys))
	nbSellList := make(map[string][]brokerAmt, len(PieKeys))

	for b, t := range bt {
		c := CategoryOf(b, ticker, uwSet)
		_, isBandar := bandarSet[b]
		_, isRepo := RepoBrokers[b]
		if c != "uw" && isBandar {
			c = "bandar"
		} else if c != "uw" && (isRepo || c == "zombie") {
			c = "repo"
		}
		if _, valid := PieMeta[c]; !valid || c == "ipo" || c == "others" {
			c = "repo"
		}
		bv := t.Buy
		sv := math.Abs(t.Sell)
		nv := bv - sv
		pieNets[c] += nv
		pieBuy[c] += bv
		pieSell[c] += sv
		if nv > 0 {
			pieNBuy[c] += nv
			nbBuyList[c] = append(nbBuyList[c], brokerAmt{broker: b, val: nv})
		} else if nv < 0 {
			pieNSell[c] += -nv
			nbSellList[c] = append(nbSellList[c], brokerAmt{broker: b, val: -nv})
		}
	}

	activePieKeys := make([]string, 0, len(PieKeys))
	for _, k := range PieKeys {
		if (k == "uw" && len(uwList) == 0) || (k == "bandar" && len(bandarSet) == 0) {
			continue
		}
		activePieKeys = append(activePieKeys, k)
	}

	buildDonutSide := func(valsMap map[string]float64, topMap map[string][]brokerAmt) DonutSide {
		tot := 0.0
		for _, k := range activePieKeys {
			tot += valsMap[k]
		}
		slices := make([]DonutSlice, 0, len(activePieKeys))
		for _, k := range activePieKeys {
			v := valsMap[k]
			if v <= 0 {
				continue
			}
			topStr := ""
			if topMap != nil {
				items := append([]brokerAmt(nil), topMap[k]...)
				sort.SliceStable(items, func(i, j int) bool {
					if items[i].val != items[j].val {
						return items[i].val > items[j].val
					}
					return items[i].broker < items[j].broker
				})
				if len(items) > 3 {
					items = items[:3]
				}
				parts := make([]string, 0, len(items))
				for _, it := range items {
					parts = append(parts, fmt.Sprintf("%s %s", it.broker, FormatIntComma(JSRound(it.val))))
				}
				topStr = strings.Join(parts, ", ")
			}
			pct := 0.0
			if tot != 0 {
				pct = (v / tot) * 100.0
			}
			meta := PieMeta[k]
			slices = append(slices, DonutSlice{
				Key:        k,
				Label:      meta[0],
				Color:      meta[1],
				Value:      v,
				Pct:        pct,
				TopBrokers: topStr,
			})
		}
		return DonutSide{Total: tot, Slices: slices}
	}

	buyerSellerPie := BuyerSellerPieCard{
		Net: PieModePair{
			Buyer:  buildDonutSide(pieNBuy, nbBuyList),
			Seller: buildDonutSide(pieNSell, nbSellList),
		},
		Gross: PieModePair{
			Buyer:  buildDonutSide(pieBuy, nil),
			Seller: buildDonutSide(pieSell, nil),
		},
	}

	// 5. After IPO Distribution Card
	buildDistSide := func(wantBuy bool) DistSide {
		if lotIPO == nil || *lotIPO <= 0 {
			return DistSide{Available: false, Slices: []DistSlice{}, Total: 0.0}
		}
		type catVal struct {
			key string
			val float64
		}
		side := make([]catVal, 0, len(PieKeys))
		for _, k := range PieKeys {
			if k == "uw" && len(uwList) == 0 {
				continue
			}
			v := pieNets[k]
			if (wantBuy && v > 0) || (!wantBuy && v < 0) {
				side = append(side, catVal{key: k, val: math.Abs(v)})
			}
		}
		sort.SliceStable(side, func(i, j int) bool {
			return side[i].val > side[j].val
		})
		if len(side) == 0 {
			emptyVal := true
			return DistSide{Available: true, Empty: &emptyVal, Slices: []DistSlice{}, Total: *lotIPO}
		}
		top3 := side
		var rest []catVal
		if len(side) > 3 {
			top3 = side[:3]
			rest = side[3:]
		}
		ipoLotVal := *lotIPO
		slices := []DistSlice{{
			Key:    "ipo",
			Label:  "Jumlah IPO",
			Color:  PieMeta["ipo"][1],
			Value:  ipoLotVal,
			PctIPO: 100.0,
		}}
		totRing := ipoLotVal
		for _, cv := range top3 {
			totRing += cv.val
			meta := PieMeta[cv.key]
			slices = append(slices, DistSlice{
				Key:    cv.key,
				Label:  meta[0],
				Color:  meta[1],
				Value:  cv.val,
				PctIPO: (cv.val / ipoLotVal) * 100.0,
			})
		}
		if len(rest) > 0 {
			restVal := 0.0
			subNames := make([]string, 0, len(rest))
			for _, cv := range rest {
				restVal += cv.val
				subNames = append(subNames, PieMeta[cv.key][0])
			}
			totRing += restVal
			slices = append(slices, DistSlice{
				Key:        "others",
				Label:      PieMeta["others"][0],
				Color:      PieMeta["others"][1],
				Value:      restVal,
				PctIPO:     (restVal / ipoLotVal) * 100.0,
				SubBrokers: strings.Join(subNames, ", "),
			})
		}
		for i := range slices {
			if totRing != 0 {
				slices[i].PctRing = (slices[i].Value / totRing) * 100.0
			}
		}
		emptyFalse := false
		return DistSide{
			Available: true,
			Empty:     &emptyFalse,
			IPOLot:    &ipoLotVal,
			Total:     totRing,
			Slices:    slices,
		}
	}

	afterIPODist := AfterIPODistCard{
		Enabled: mode == "after",
		Sell:    buildDistSide(false),
		Buy:     buildDistSide(true),
	}

	// 6. Retail Action Card
	raRows := make([]RetailActionRow, 0, len(decision.EffectiveRetailList))
	raSumTrans := 0.0
	raSumNet := 0.0
	for _, b := range decision.EffectiveRetailList {
		t := bt[b]
		trans := t.Buy + math.Abs(t.Sell)
		netV := t.Buy - t.Sell
		isAuto := decision.RetailBandarAutoResult != nil && b == *decision.RetailBandarAutoResult
		if !isAuto {
			raSumTrans += trans
			raSumNet += netV
		}
		_, isActive := bandarSet[b]
		raRows = append(raRows, RetailActionRow{
			Broker:         b,
			Net:            netV,
			PctNet:         pctPtr(netV, tn),
			Trans:          trans,
			PctTrans:       pctPtr(trans, tv),
			IsAutoBandar:   isAuto,
			IsActiveBandar: isActive,
		})
	}
	retailActionCard := RetailActionCard{
		RetailRatio:            decision.RetailRatio,
		RetailAction:           decision.RetailAction,
		RetailBandarAutoResult: decision.RetailBandarAutoResult,
		RetailBandarBroker:     decision.RetailBandarBroker,
		RetailBandarBrokers:    decision.RetailBandarBrokers,
		RetailBandarReasons:    decision.RetailBandarReasons,
		EffectiveRetailList:    decision.EffectiveRetailList,
		Rows:                   raRows,
		TotalExclAuto: RetailTotalExclAuto{
			SumTrans: raSumTrans,
			SumNet:   raSumNet,
			PctTrans: pctPtr(raSumTrans, tv),
			PctNet:   pctPtr(raSumNet, tn),
		},
	}

	// 7. Info IPO Card
	var ipoVsClose *int64
	if lastClose != nil && master.IPOPrice != 0 {
		v := JSRound((*lastClose/master.IPOPrice - 1.0) * 100.0)
		ipoVsClose = &v
	}
	var avgVsClose *int64
	if lastClose != nil && vwap != nil && *vwap != 0 {
		v := JSRound((*lastClose/(*vwap) - 1.0) * 100.0)
		avgVsClose = &v
	}
	var endDt time.Time
	var hasEndDt bool
	if afterInfo != nil && afterInfo.LastDate != nil {
		endDt = *afterInfo.LastDate
		hasEndDt = true
	} else if len(analysisDays) > 0 {
		endDt = analysisDays[len(analysisDays)-1].Date
		hasEndDt = true
	}
	var tradingSince *int
	holidayWarn := false
	if master.ListingDate != nil && hasEndDt && !endDt.Before(*master.ListingDate) {
		ts := TradingDays(*master.ListingDate, endDt, holidays)
		tradingSince = &ts
		holidayWarn = endDt.UTC().Year() > maxHolidayYear
	}

	var lotTotal *float64
	if master.TotalListedShares != nil {
		v := float64(*master.TotalListedShares) / 100.0
		lotTotal = &v
	}
	var nilaiMiliar *float64
	if master.NilaiIPO != nil {
		v := *master.NilaiIPO / 1e9
		nilaiMiliar = &v
	}
	var listingStr *string
	if master.ListingDate != nil {
		s := FormatDateISO(*master.ListingDate)
		listingStr = &s
	}
	compName := ticker
	if master.CompanyName != nil && *master.CompanyName != "" {
		compName = *master.CompanyName
	}
	var d0DateStr *string
	if afterInfo != nil && afterInfo.D0Date != nil {
		s := FormatDateISO(*afterInfo.D0Date)
		d0DateStr = &s
	}
	nLater := len(analysisDays)
	if afterInfo != nil {
		nLater = afterInfo.NLater
	}
	var gap *GapRange
	if afterInfo != nil {
		gap = afterInfo.Gap
	}

	infoCard := InfoCard{
		Ticker:              ticker,
		CompanyName:         compName,
		ListingDate:         listingStr,
		IPOPrice:            master.IPOPrice,
		LastClose:           lastClose,
		VWAP:                vwap,
		IPOVsClosePct:       ipoVsClose,
		AvgVsClosePct:       avgVsClose,
		LotOffer:            lotIPO,
		LotTotal:            lotTotal,
		PctIPO:              master.PctIPO,
		NilaiMiliar:         nilaiMiliar,
		Penjatahan:          master.Penjatahan,
		Papan:               master.Papan,
		JumlahSID:           master.JumlahSID,
		TradingDaysSinceIPO: tradingSince,
		HolidayWarn:         holidayWarn,
		MaxHolidayYear:      maxHolidayYear,
		Gap:                 gap,
		NoD0:                afterInfo != nil && afterInfo.NoD0,
		InD0:                afterInfo != nil && afterInfo.InD0,
		D0Date:              d0DateStr,
		NLater:              nLater,
	}

	return AnalysisCards{
		TopBuyerSeller: TopBuyerSellerCard{
			Buyers:     topBuyers,
			Sellers:    topSellers,
			AccDist:    accDist,
			ValuePrice: valuePrice,
		},
		UnderwriterCard:      uwCard,
		CategoryCard:         categoryCard,
		AfterIPODistribution: afterIPODist,
		BuyerSellerPie:       buyerSellerPie,
		RetailActionCard:     retailActionCard,
		InfoCard:             infoCard,
	}
}
