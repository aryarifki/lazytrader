package mmd

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/domains/terminal/shared"
)

// AnalysisInput holds all DB inputs and options for ComputeMMDAnalysis.
type AnalysisInput struct {
	Ticker             string
	Candles            []DailyCandle
	AllHistoryCandles  []DailyCandle
	T5Records          []T5Record
	Meta               *TickerMetaFull
	GroupPerf          *GroupPerformance
	KSEIMacroRepoRatio float64
	AutoRetail         bool
	RetailOverride     string
	Today              time.Time
}

type brokerNet struct {
	Broker string
	Net    int64
}

type brokerTotal struct {
	Buy  int64
	Sell int64
}

type sliceResume struct {
	Action           Badge
	PriceChange      float64
	PengendaliAction Badge
	RetailAction     Badge
}

func pyFloorDiv(a, b int64) int64 {
	q := a / b
	r := a % b
	if r != 0 && (a < 0) != (b < 0) {
		q--
	}
	return q
}

func sortedJSONBKeys(m map[string]float64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if len(keys[i]) != len(keys[j]) {
			return len(keys[i]) < len(keys[j])
		}
		return keys[i] < keys[j]
	})
	return keys
}

func formatCommaInt(v int64) string {
	neg := v < 0
	var uv uint64
	if neg {
		uv = uint64(-v)
	} else {
		uv = uint64(v)
	}
	s := strconv.FormatUint(uv, 10)
	n := len(s)
	var b strings.Builder
	if neg {
		b.WriteByte('-')
	}
	for i := 0; i < n; i++ {
		if i > 0 && (n-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func formatSignedCommaInt(v int64) string {
	if v == 0 {
		return "0"
	}
	if v > 0 {
		return "+" + formatCommaInt(v)
	}
	return formatCommaInt(v)
}

func formatCommaFloat(v float64, digits int) string {
	neg := v < 0
	raw := strconv.FormatFloat(math.Abs(v), 'f', digits, 64)
	dot := strings.IndexByte(raw, '.')
	intPart := raw
	fracPart := ""
	if dot >= 0 {
		intPart = raw[:dot]
		fracPart = raw[dot:]
	}
	n := len(intPart)
	var b strings.Builder
	if neg {
		b.WriteByte('-')
	}
	for i := 0; i < n; i++ {
		if i > 0 && (n-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteByte(intPart[i])
	}
	b.WriteString(fracPart)
	return b.String()
}

func formatPyFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func candlePrevOpenClose(c DailyCandle) float64 {
	if c.Prev != nil && *c.Prev != 0 {
		return *c.Prev
	}
	if c.Open != nil && *c.Open != 0 {
		return *c.Open
	}
	if c.Close != nil {
		return *c.Close
	}
	return 0
}

func candleTavgOrClose(c DailyCandle) float64 {
	if c.Tavg != nil && *c.Tavg != 0 {
		return *c.Tavg
	}
	if c.Close != nil {
		return *c.Close
	}
	return 0
}

func cumulativeDiffAtN(buyers, sellers []brokerNet, n int) int64 {
	var buySum, sellSum int64
	for i := 0; i < n && i < len(buyers); i++ {
		buySum += buyers[i].Net
	}
	for i := 0; i < n && i < len(sellers); i++ {
		if sellers[i].Net < 0 {
			sellSum += -sellers[i].Net
		} else {
			sellSum += sellers[i].Net
		}
	}
	return buySum - sellSum
}

func aggregateBrokerSlice(candles []DailyCandle, trackTotals bool) (
	netMap map[string]int64,
	brokerTotals map[string]brokerTotal,
	totBuy int64,
	totSell int64,
	totNetVol int64,
	totalVol int64,
	sumDailyVal float64,
) {
	netMap = make(map[string]int64, 64)
	if trackTotals {
		brokerTotals = make(map[string]brokerTotal, 64)
	}
	for _, c := range candles {
		seen := make(map[string]bool, len(c.BuyBrokers)+len(c.SellBrokers))
		for b := range c.BuyBrokers {
			seen[b] = true
		}
		for b := range c.SellBrokers {
			seen[b] = true
		}
		for b := range seen {
			bv := int64(c.BuyBrokers[b])
			sv := int64(c.SellBrokers[b])
			totBuy += bv
			totSell += sv
			netMap[b] += bv - sv
			if trackTotals {
				bt := brokerTotals[b]
				bt.Buy += bv
				bt.Sell += sv
				brokerTotals[b] = bt
			}
		}
		var vol int64
		if c.Volume != nil {
			vol = *c.Volume
		}
		sumDailyVal += float64(vol) * candleTavgOrClose(c)
	}
	var absNetSum int64
	for _, v := range netMap {
		if v < 0 {
			absNetSum += -v
		} else {
			absNetSum += v
		}
	}
	totNetVol = pyFloorDiv(absNetSum, 2)
	totalVol = pyFloorDiv(totBuy+totSell, 2)
	return netMap, brokerTotals, totBuy, totSell, totNetVol, totalVol, sumDailyVal
}

func sortBuyersSellers(netMap map[string]int64) (buyers []brokerNet, sellers []brokerNet) {
	buyers = make([]brokerNet, 0, 16)
	sellers = make([]brokerNet, 0, 16)
	for b, net := range netMap {
		if net > 0 {
			buyers = append(buyers, brokerNet{Broker: b, Net: net})
		} else if net < 0 {
			sellers = append(sellers, brokerNet{Broker: b, Net: net})
		}
	}
	sort.SliceStable(buyers, func(i, j int) bool {
		if buyers[i].Net != buyers[j].Net {
			return buyers[i].Net > buyers[j].Net
		}
		return buyers[i].Broker < buyers[j].Broker
	})
	if len(buyers) > 10 {
		buyers = buyers[:10]
	}
	sort.SliceStable(sellers, func(i, j int) bool {
		if sellers[i].Net != sellers[j].Net {
			return sellers[i].Net < sellers[j].Net
		}
		return sellers[i].Broker < sellers[j].Broker
	})
	if len(sellers) > 10 {
		sellers = sellers[:10]
	}
	return buyers, sellers
}

func computeSliceResume(candles []DailyCandle) sliceResume {
	defBadge := Badge{Label: "—", Cls: "undetected", Color: "#94A3B8"}
	defRetail := Badge{Label: "Do Nothing", Cls: "undetected", Color: "#94A3B8"}
	if len(candles) == 0 {
		return sliceResume{Action: defBadge, PriceChange: 0.0, PengendaliAction: defBadge, RetailAction: defRetail}
	}

	firstC := candles[0]
	lastC := candles[len(candles)-1]
	prevClose := candlePrevOpenClose(firstC)
	lastClose := 0.0
	if lastC.Close != nil {
		lastClose = *lastC.Close
	}
	priceChange := 0.0
	if prevClose > 0 {
		priceChange = shared.PyRound(((lastClose-prevClose)/prevClose)*100.0, 2)
	}

	netMap, _, _, _, totNetVol, totalVol, _ := aggregateBrokerSlice(candles, false)
	buyers, sellers := sortBuyersSellers(netMap)

	maxRows := len(buyers)
	if len(sellers) > maxRows {
		maxRows = len(sellers)
	}
	if maxRows > 10 {
		maxRows = 10
	}

	accDistRows := make([]int64, 0, maxRows)
	for i := 1; i <= maxRows; i++ {
		accDistRows = append(accDistRows, cumulativeDiffAtN(buyers, sellers, i))
	}

	var nilaiTop3, nilaiTop5, nilaiAvg10 int64
	if maxRows >= 3 {
		nilaiTop3 = cumulativeDiffAtN(buyers, sellers, 3)
	}
	if maxRows >= 5 {
		nilaiTop5 = cumulativeDiffAtN(buyers, sellers, 5)
	}
	if len(accDistRows) > 0 {
		var s int64
		for _, d := range accDistRows {
			s += d
		}
		nilaiAvg10 = pyFloorDiv(s, int64(len(accDistRows)))
	}

	avgNilai := pyFloorDiv(nilaiTop3+nilaiTop5+nilaiAvg10, 3)
	avgRatio := 0.0
	if totNetVol > 0 {
		avgRatio = shared.PyRound(float64(avgNilai)/float64(totNetVol)*100.0, 1)
	}
	vr := 1.0
	if totNetVol > 0 {
		vr = float64(totalVol) / float64(totNetVol)
	}
	action := ClassifyBandarVal(float64(avgRatio * VolumeRatioWeight(vr)))

	var instNet, retailNet int64
	for b := range InstitusiBrokers {
		instNet += netMap[b]
	}
	for b := range RetailBrokers {
		retailNet += netMap[b]
	}
	pengendaliRatio := 0.0
	retailRatio := 0.0
	if totNetVol > 0 {
		pengendaliRatio = shared.PyRound(float64(instNet)/float64(totNetVol)*100.0, 1)
		retailRatio = shared.PyRound(float64(retailNet)/float64(totNetVol)*100.0, 1)
	}

	return sliceResume{
		Action:           action,
		PriceChange:      priceChange,
		PengendaliAction: ClassifyBandarVal(pengendaliRatio),
		RetailAction:     ClassifyRetailAction(retailRatio),
	}
}

func detectRetailBandar(buyers []brokerNet, brokerTotals map[string]brokerTotal) *string {
	if len(buyers) == 0 {
		return nil
	}
	var picked string
	for i := 0; i < 3 && i < len(buyers); i++ {
		if RetailBrokers[buyers[i].Broker] {
			picked = buyers[i].Broker
			break
		}
	}
	if picked == "" {
		return nil
	}
	bt := brokerTotals[picked]
	bSellAbs := bt.Sell
	if bSellAbs < 0 {
		bSellAbs = -bSellAbs
	}
	bTrans := bt.Buy + bSellAbs

	var allRetailTrans int64
	for r := range RetailBrokers {
		rt := brokerTotals[r]
		rSell := rt.Sell
		if rSell < 0 {
			rSell = -rSell
		}
		allRetailTrans += rt.Buy + rSell
	}
	if allRetailTrans > 0 && float64(bTrans)/float64(allRetailTrans) >= 0.35 {
		return &picked
	}
	return nil
}

// ComputeMMDAnalysis computes the full MM Detector Pro analysis payload.
func ComputeMMDAnalysis(in AnalysisInput) (*AnalysisResult, bool) {
	if len(in.Candles) == 0 {
		return nil, false
	}

	netMap, brokerTotals, totBuy, totSell, totNetVol, totalVol, sumDailyVal := aggregateBrokerSlice(in.Candles, true)
	avgDailyVal := sumDailyVal / float64(len(in.Candles))

	firstC := in.Candles[0]
	lastC := in.Candles[len(in.Candles)-1]
	prevClose := candlePrevOpenClose(firstC)
	lastClose := 0.0
	if lastC.Close != nil {
		lastClose = *lastC.Close
	}
	priceChange := 0.0
	if prevClose > 0 {
		priceChange = shared.PyRound(((lastClose-prevClose)/prevClose)*100.0, 2)
	}

	buyers, sellers := sortBuyersSellers(netMap)
	maxRows := len(buyers)
	if len(sellers) > maxRows {
		maxRows = len(sellers)
	}
	if maxRows > 10 {
		maxRows = 10
	}
	accDistRows := make([]int64, 0, maxRows)
	for i := 1; i <= maxRows; i++ {
		accDistRows = append(accDistRows, cumulativeDiffAtN(buyers, sellers, i))
	}

	autoRB := detectRetailBandar(buyers, brokerTotals)
	var activeRB *string
	overrideClean := strings.ToUpper(strings.TrimSpace(in.RetailOverride))
	if overrideClean != "" {
		activeRB = &overrideClean
	} else if in.AutoRetail {
		activeRB = autoRB
	}

	var nilaiTop3, nilaiTop5, nilaiAvg10 int64
	if maxRows >= 3 {
		nilaiTop3 = cumulativeDiffAtN(buyers, sellers, 3)
	}
	if maxRows >= 5 {
		nilaiTop5 = cumulativeDiffAtN(buyers, sellers, 5)
	}
	if len(accDistRows) > 0 {
		var s int64
		for _, d := range accDistRows {
			s += d
		}
		nilaiAvg10 = pyFloorDiv(s, int64(len(accDistRows)))
	}

	ratioTop3, ratioTop5, ratioAvg10, avgRatio := 0.0, 0.0, 0.0, 0.0
	avgNilai := pyFloorDiv(nilaiTop3+nilaiTop5+nilaiAvg10, 3)
	if totNetVol > 0 {
		ratioTop3 = shared.PyRound(float64(nilaiTop3)/float64(totNetVol)*100.0, 1)
		ratioTop5 = shared.PyRound(float64(nilaiTop5)/float64(totNetVol)*100.0, 1)
		ratioAvg10 = shared.PyRound(float64(nilaiAvg10)/float64(totNetVol)*100.0, 1)
		avgRatio = shared.PyRound(float64(avgNilai)/float64(totNetVol)*100.0, 1)
	}

	vr := 1.0
	if totNetVol > 0 {
		vr = float64(totalVol) / float64(totNetVol)
	}
	actionCls := ClassifyBandarVal(float64(avgRatio * VolumeRatioWeight(vr)))

	buyTop3, sellTop3 := sumNetAbs(buyers, 3), sumAbsNet(sellers, 3)
	buyTop5, sellTop5 := sumNetAbs(buyers, 5), sumAbsNet(sellers, 5)
	bTenLen := minInt(len(buyers), 10)
	sTenLen := minInt(len(sellers), 10)
	var buyAvg10, sellAvg10 int64
	if len(buyers) > 0 {
		buyAvg10 = pyFloorDiv(sumNetAbs(buyers, 10), int64(maxInt(bTenLen, 1)))
	}
	if len(sellers) > 0 {
		sellAvg10 = pyFloorDiv(sumAbsNet(sellers, 10), int64(maxInt(sTenLen, 1)))
	}
	avgBuy := pyFloorDiv(buyTop3+buyTop5+buyAvg10, 3)
	avgSell := pyFloorDiv(sellTop3+sellTop5+sellAvg10, 3)

	bandarTable := []BandarTableRow{
		{Label: "Top 3", Buy: buyTop3, Sell: sellTop3, Nilai: nilaiTop3, Ratio: ratioTop3, Badge: ClassifyBandarVal(ratioTop3)},
		{Label: "Top 5", Buy: buyTop5, Sell: sellTop5, Nilai: nilaiTop5, Ratio: ratioTop5, Badge: ClassifyBandarVal(ratioTop5)},
		{Label: "Average T10", Buy: buyAvg10, Sell: sellAvg10, Nilai: nilaiAvg10, Ratio: ratioAvg10, Badge: ClassifyBandarVal(ratioAvg10)},
		{Label: "Aggregate Score", Buy: avgBuy, Sell: avgSell, Nilai: avgNilai, Ratio: avgRatio, Badge: actionCls, IsAction: true},
	}

	topBuyers := make([]TopBrokerRow, 0, len(buyers))
	for idx, b := range buyers {
		valJt := shared.PyRound(float64(float64(b.Net)*100.0*lastClose)/1e6, 1)
		topBuyers = append(topBuyers, TopBrokerRow{
			Broker: b.Broker, Net: b.Net, VolumeFmt: formatCommaInt(b.Net),
			ValueJt: valJt, Kategori: CategoryLabelForBroker(b.Broker, idx, "buy", actionCls.Label, activeRB),
		})
	}
	topSellers := make([]TopBrokerRow, 0, len(sellers))
	for idx, s := range sellers {
		absNet := s.Net
		if absNet < 0 {
			absNet = -absNet
		}
		valJt := shared.PyRound(float64(float64(absNet)*100.0*lastClose)/1e6, 1)
		topSellers = append(topSellers, TopBrokerRow{
			Broker: s.Broker, Net: s.Net, VolumeFmt: formatCommaInt(absNet),
			ValueJt: valJt, Kategori: CategoryLabelForBroker(s.Broker, idx, "sell", actionCls.Label, activeRB),
		})
	}

	top10Combined := buildTop10Combined(buyers, sellers, accDistRows, topBuyers, topSellers, lastClose, maxRows)
	retailRows, retailRatio, retailBadge := buildRetailSection(brokerTotals, totNetVol, totalVol, activeRB)
	pengendaliStats := buildPengendaliSection(netMap, brokerTotals, totNetVol, totalVol, activeRB, in.T5Records, in.KSEIMacroRepoRatio)

	sortedAll := make([]DailyCandle, len(in.AllHistoryCandles))
	if len(sortedAll) == 0 {
		sortedAll = append([]DailyCandle{}, in.Candles...)
	} else {
		copy(sortedAll, in.AllHistoryCandles)
	}
	sort.SliceStable(sortedAll, func(i, j int) bool {
		return sortedAll[i].DateISO < sortedAll[j].DateISO
	})

	prevRes, prevCtx := computePreviousResults(sortedAll)
	today := in.Today
	if today.IsZero() {
		today = time.Now()
	}
	aiScore, korelasiRaw, aiComp, dailyResumes := computeAIScore(
		in.Ticker, sortedAll, priceChange, avgDailyVal, actionCls, in.GroupPerf, prevCtx,
	)
	activeCommentary, bigPlayerSummary := generateCommentaries(
		in, sortedAll, dailyResumes, buyers, sellers, topBuyers,
		pengendaliStats, actionCls, retailBadge, aiScore, korelasiRaw,
		vr, avgDailyVal, totNetVol, totalVol, lastClose, priceChange, today,
	)

	t5Out := in.T5Records
	if t5Out == nil {
		t5Out = []T5Record{}
	}
	manualOverrideOut := ""
	if in.RetailOverride != "" {
		manualOverrideOut = in.RetailOverride
	}

	return &AnalysisResult{
		Ticker:                     in.Ticker,
		TotalDays:                  len(in.Candles),
		DateFrom:                   firstC.DateISO,
		DateTo:                     lastC.DateISO,
		LastPrice:                  lastClose,
		PriceChange:                priceChange,
		TotalBuy:                   totBuy,
		TotalSell:                  totSell,
		TotalVolume:                totalVol,
		TotNetVolume:               totNetVol,
		AvgDailyValue:              avgDailyVal,
		Buyers:                     topBuyers,
		Sellers:                    topSellers,
		AccDistRows:                accDistRows,
		Top10Combined:              top10Combined,
		BandarTable:                bandarTable,
		BandarAction:               actionCls,
		RetailRows:                 retailRows,
		RetailRatio:                retailRatio,
		RetailAction:               retailBadge,
		RetailBandarAutoResult:     autoRB,
		ActiveRetailBandar:         activeRB,
		RetailBandarAutoEnabled:    in.AutoRetail,
		RetailBandarManualOverride: manualOverrideOut,
		EffectiveRetailBrokers:     EffectiveRetailBrokers,
		PengendaliStats:            pengendaliStats,
		Resume:                     actionCls,
		AIScore:                    aiScore,
		Korelasi:                   korelasiRaw,
		AIScoreComponents:          aiComp,
		PreviousResults:            prevRes,
		AICommentary:               activeCommentary,
		Commentary:                 activeCommentary,
		BigPlayerSummary:           bigPlayerSummary,
		T5Records:                  t5Out,
		GroupPerformance:           in.GroupPerf,
	}, true
}

func buildTop10Combined(
	buyers, sellers []brokerNet, accDistRows []int64,
	topBuyers, topSellers []TopBrokerRow, lastClose float64, maxRows int,
) []Top10CombinedRow {
	rows := make([]Top10CombinedRow, 0, maxRows)
	for i := 0; i < maxRows; i++ {
		var b *brokerNet
		var s *brokerNet
		if i < len(buyers) {
			b = &buyers[i]
		}
		if i < len(sellers) {
			s = &sellers[i]
		}
		var diff int64
		if i < len(accDistRows) {
			diff = accDistRows[i]
		}
		row := Top10CombinedRow{
			Rank: i + 1, BuyerBroker: "—", BuyerNetFmt: "—", BuyerKategori: "—",
			SellerBroker: "—", SellerNetFmt: "—", SellerKategori: "—",
			AccDist: diff, AccDistFmt: formatSignedCommaInt(diff),
		}
		if b != nil {
			row.BuyerBroker = b.Broker
			row.BuyerNet = b.Net
			row.BuyerNetFmt = formatCommaInt(b.Net)
			row.BuyerValJt = shared.PyRound(float64(float64(b.Net)*100.0*lastClose)/1e6, 1)
			if i < len(topBuyers) {
				row.BuyerKategori = topBuyers[i].Kategori
			}
		}
		if s != nil {
			absS := s.Net
			if absS < 0 {
				absS = -absS
			}
			row.SellerBroker = s.Broker
			row.SellerNet = absS
			row.SellerNetFmt = formatCommaInt(absS)
			row.SellerValJt = shared.PyRound(float64(float64(absS)*100.0*lastClose)/1e6, 1)
			if i < len(topSellers) {
				row.SellerKategori = topSellers[i].Kategori
			}
		}
		rows = append(rows, row)
	}
	return rows
}

func buildRetailSection(
	brokerTotals map[string]brokerTotal, totNetVol, totalVol int64, activeRB *string,
) ([]RetailRow, float64, Badge) {
	rows := make([]RetailRow, 0, len(CoreRetailList))
	var sumRetailNet int64
	totTurnover := maxInt64(totalVol*2, 1)

	for _, rCode := range CoreRetailList {
		bt := brokerTotals[rCode]
		sellAbs := bt.Sell
		if sellAbs < 0 {
			sellAbs = -sellAbs
		}
		trans := bt.Buy + sellAbs
		net := bt.Buy - bt.Sell
		pctTrans, pctNet := 0.0, 0.0
		if totalVol > 0 {
			pctTrans = shared.PyRound(float64(trans)/float64(totTurnover)*100.0, 1)
		}
		if totNetVol > 0 {
			pctNet = shared.PyRound(float64(net)/float64(totNetVol)*100.0, 1)
		}
		isRB := activeRB != nil && rCode == *activeRB
		if !isRB {
			sumRetailNet += net
		}
		rows = append(rows, RetailRow{
			Broker: rCode, Trans: trans, PctTrans: pctTrans, Net: net, PctNet: pctNet, IsRetailBandar: isRB,
		})
	}
	retailRatio := 0.0
	if totNetVol > 0 {
		retailRatio = shared.PyRound(float64(sumRetailNet)/float64(totNetVol)*100.0, 1)
	}
	return rows, retailRatio, ClassifyRetailAction(retailRatio)
}

func buildPengendaliSection(
	netMap map[string]int64, brokerTotals map[string]brokerTotal,
	totNetVol, totalVol int64, activeRB *string, t5 []T5Record, macroRepo float64,
) PengendaliStats {
	var instVol, semiVol, ritelVol int64
	for b, bt := range brokerTotals {
		if InstitusiBrokers[b] {
			instVol += bt.Buy + bt.Sell
		}
		if SemiRetailBrokers[b] {
			semiVol += bt.Buy + bt.Sell
		}
		if RetailBrokers[b] && (activeRB == nil || b != *activeRB) {
			ritelVol += bt.Buy + bt.Sell
		}
	}
	if activeRB != nil {
		if bt, ok := brokerTotals[*activeRB]; ok {
			instVol += bt.Buy + bt.Sell
		}
	}
	otherVol := totalVol*2 - (instVol + semiVol + ritelVol)
	if otherVol < 0 {
		otherVol = 0
	}
	totTurnover := maxInt64(totalVol*2, 1)

	var transferLots float64
	for _, rec := range t5 {
		if rec.Action != nil {
			act := strings.ToUpper(*rec.Action)
			if act == "TRANSFER" || act == "REPO" {
				transferLots += math.Abs(rec.Lot)
			}
		}
	}
	repoVal := 0.0
	if transferLots > 0 && totalVol > 0 {
		repoVal = shared.PyRound((transferLots/float64(totalVol))*100.0, 1)
	} else if macroRepo > 0 {
		repoVal = shared.PyRound(macroRepo*100.0, 1)
	}

	otherVal := shared.PyRound(float64(otherVol)/float64(totTurnover)*100.0, 1)
	if repoVal > 0 {
		otherVal = math.Max(0.0, shared.PyRound(otherVal-math.Min(otherVal, float64(repoVal*0.2)), 1))
	}

	var instNet, semiNet, ritelNet int64
	for b := range InstitusiBrokers {
		instNet += netMap[b]
	}
	if activeRB != nil && !InstitusiBrokers[*activeRB] {
		instNet += netMap[*activeRB]
	}
	for b := range SemiRetailBrokers {
		semiNet += netMap[b]
	}
	for b := range RetailBrokers {
		if activeRB == nil || b != *activeRB {
			ritelNet += netMap[b]
		}
	}
	nonRitelNet := instNet + semiNet

	denom := totNetVol
	if denom <= 0 {
		denom = 1
	}
	pctSub := func(num int64) float64 {
		if totNetVol <= 0 {
			return 0.0
		}
		return shared.PyRound(float64(num)/float64(denom)*100.0, 1)
	}
	pctMain := func(num int64) float64 {
		return shared.PyRound(float64(num)/float64(totTurnover)*100.0, 1)
	}
	instRatio := pctMain(instVol)
	pengendaliNetRatio := 0.0
	if totNetVol > 0 {
		pengendaliNetRatio = shared.PyRound(float64(instNet)/float64(totNetVol)*100.0, 1)
	}

	return PengendaliStats{
		Institusi: PengendaliGroupStat{Main: instRatio, Sub: pctSub(instNet)},
		SemiRitel: PengendaliGroupStat{Main: pctMain(semiVol), Sub: pctSub(semiNet)},
		Ritel:     PengendaliGroupStat{Main: pctMain(ritelVol), Sub: pctSub(ritelNet)},
		NonRitel:  PengendaliGroupStat{Main: pctMain(instVol + semiVol), Sub: pctSub(nonRitelNet)},
		Repo:      ValStat{Val: repoVal},
		Lainnya:   ValStat{Val: otherVal},
		Ratio:     instRatio,
		Action:    ClassifyBandarVal(pengendaliNetRatio),
	}
}

func sumNetAbs(items []brokerNet, n int) int64 {
	var s int64
	for i := 0; i < n && i < len(items); i++ {
		s += items[i].Net
	}
	return s
}

func sumAbsNet(items []brokerNet, n int) int64 {
	var s int64
	for i := 0; i < n && i < len(items); i++ {
		if items[i].Net < 0 {
			s += -items[i].Net
		} else {
			s += items[i].Net
		}
	}
	return s
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
