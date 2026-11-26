package ipohunter

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/domains/terminal/shared"
)

// AfterInfo captures After-IPO D0+later window state for a ticker.
type AfterInfo struct {
	Days        []DayRow
	Gap         *GapRange
	LastDate    *time.Time
	D0Date      *time.Time
	InD0        bool
	NLater      int
	ListingDate *time.Time
	NoD0        bool
	FromListing bool
	FirstDate   *time.Time
}

// BrokerNet stores a single broker's net volume.
type BrokerNet struct {
	Broker string
	Net    float64
}

// BrokerTotal stores gross buy and sell volumes for a broker.
type BrokerTotal struct {
	Buy  float64
	Sell float64
}

// Aggregation stores broker flow and price stats over analysis days.
type Aggregation struct {
	Ticker        string
	Net           map[string]float64
	BrokerTotals  map[string]BrokerTotal
	TotBuy        float64
	TotSell       float64
	TotNetVolume  float64
	TotalVolume   float64
	Buyers        []BrokerNet
	Sellers       []BrokerNet
	AvgDailyValue *float64
	PriceChange   *float64
}

// Decision holds bandar and retail bandar classifications.
type Decision struct {
	ActionWeight           float64
	VolumeRatioAdjusted    *float64
	BandarAction           ActionLabel
	RetailBandarBroker     *string
	RetailBandarBrokers    []string
	RetailBandarAutoResult *string
	RetailBandarReasons    map[string][]string
	EffectiveRetailList    []string
	FinalRetailSet         map[string]struct{}
	BrokerCategory         map[string]string
	RetailRatio            *float64
	RetailAction           ActionLabel
}

// BuildAfterDays constructs the After-IPO window merging D0 baseline and daily rows <= hiDate.
func BuildAfterDays(
	_ string,
	master *MasterRecord,
	d0Row *DayRow,
	sortedDaily []DayRow,
	hiDate *time.Time,
	holidays map[int]struct{},
) AfterInfo {
	var listing *time.Time
	if master != nil {
		listing = master.ListingDate
	}
	var firstDate *time.Time
	if len(sortedDaily) > 0 {
		fd := sortedDaily[0].Date
		firstDate = &fd
	}
	out := AfterInfo{
		Days:        make([]DayRow, 0, len(sortedDaily)+1),
		ListingDate: listing,
		FirstDate:   firstDate,
	}
	if d0Row != nil {
		d0dt := d0Row.Date
		out.D0Date = &d0dt
	}

	if d0Row == nil {
		inWin := listing != nil && firstDate != nil && !listing.After(*firstDate)
		out.FromListing = inWin
		out.NoD0 = !inWin
		base := make([]DayRow, 0, len(sortedDaily))
		for _, d := range sortedDaily {
			if inWin && d.Date.Before(*listing) {
				continue
			}
			if hiDate != nil && d.Date.After(*hiDate) {
				continue
			}
			base = append(base, d)
		}
		out.Days = base
		out.NLater = len(base)
		if len(base) > 0 {
			ld := base[len(base)-1].Date
			out.LastDate = &ld
		}
		return out
	}

	d0Date := d0Row.Date
	later := make([]DayRow, 0, len(sortedDaily))
	for _, d := range sortedDaily {
		if d.Date.After(d0Date) && (hiDate == nil || !d.Date.After(*hiDate)) {
			later = append(later, d)
		}
	}

	start := d0Date
	if listing != nil {
		start = *listing
	}
	d0Copy := *d0Row
	d0Copy.IsD0 = true
	d0Copy.DayCount = WeekdaysBetween(start, d0Date, holidays)
	out.Days = append(out.Days, d0Copy)
	out.InD0 = true

	out.Days = append(out.Days, later...)
	out.NLater = len(later)
	if len(later) > 0 {
		ld := later[len(later)-1].Date
		out.LastDate = &ld
		gapDays := int(later[0].Date.Sub(d0Date).Hours() / 24)
		if gapDays > 4 {
			out.Gap = &GapRange{
				From: FormatDateISO(d0Date.AddDate(0, 0, 1)),
				To:   FormatDateISO(later[0].Date.AddDate(0, 0, -1)),
			}
		}
	} else {
		ld := d0Date
		out.LastDate = &ld
	}
	return out
}

// LastCloseIn returns the latest positive Close on non-D0 days.
func LastCloseIn(days []DayRow) *float64 {
	var bestDt *time.Time
	var bestC *float64
	for i := range days {
		d := &days[i]
		if d.IsD0 {
			continue
		}
		if d.Close != nil && *d.Close > 0 && (bestDt == nil || !d.Date.Before(*bestDt)) {
			dt := d.Date
			c := *d.Close
			bestDt = &dt
			bestC = &c
		}
	}
	return bestC
}

// AfterVWAP computes volume-weighted average price using Tavg or positive OHLC/4 fallback.
func AfterVWAP(days []DayRow) *float64 {
	sv := 0.0
	sp := 0.0
	for i := range days {
		d := &days[i]
		if d.Volume == nil || *d.Volume == 0 {
			continue
		}
		var ap float64
		if d.Tavg != nil && *d.Tavg > 0 {
			ap = *d.Tavg
		} else {
			if d.Open == nil || *d.Open <= 0 || d.High == nil || *d.High <= 0 || d.Low == nil || *d.Low <= 0 || d.Close == nil || *d.Close <= 0 {
				continue
			}
			ap = (*d.Open + *d.High + *d.Low + *d.Close) / 4.0
		}
		vol := float64(*d.Volume)
		sv += vol
		sp += vol * ap
	}
	if sv == 0 {
		return nil
	}
	res := sp / sv
	return &res
}

// BuildAggregation accumulates broker buy/sell totals, Top 10 buyers/sellers, avg daily value, and price change.
func BuildAggregation(ticker string, days []DayRow) *Aggregation {
	if len(days) == 0 {
		return nil
	}

	net := make(map[string]float64, 64)
	brokerTotals := make(map[string]BrokerTotal, 64)
	totBuy := 0.0
	totSell := 0.0

	for i := range days {
		d := &days[i]
		seen := make(map[string]struct{}, len(d.BuyBrokers)+len(d.SellBrokers))
		for b := range d.BuyBrokers {
			seen[b] = struct{}{}
		}
		for b := range d.SellBrokers {
			seen[b] = struct{}{}
		}
		brokers := make([]string, 0, len(seen))
		for b := range seen {
			brokers = append(brokers, b)
		}
		sort.Strings(brokers)

		for _, b := range brokers {
			bv := d.BuyBrokers[b]
			sv := d.SellBrokers[b]
			totBuy += bv
			totSell += sv
			net[b] += bv - sv
			bt := brokerTotals[b]
			bt.Buy += bv
			bt.Sell += sv
			brokerTotals[b] = bt
		}
	}

	totNetVolume := 0.0
	buyers := make([]BrokerNet, 0, 32)
	sellers := make([]BrokerNet, 0, 32)
	for b, v := range net {
		totNetVolume += math.Abs(v)
		if v > 0 {
			buyers = append(buyers, BrokerNet{Broker: b, Net: v})
		} else if v < 0 {
			sellers = append(sellers, BrokerNet{Broker: b, Net: v})
		}
	}
	totNetVolume /= 2.0

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

	totalVolume := (totBuy + totSell) / 2.0

	// Match L509-521: skips only Volume == nil (NOT Volume == 0), uses Tavg when != nil even if 0.
	sumDailyValue := 0.0
	daysWithValue := 0
	for i := range days {
		d := &days[i]
		if d.Volume == nil {
			continue
		}
		var ap float64
		if d.Tavg != nil {
			ap = *d.Tavg
		} else {
			if d.Open == nil || d.High == nil || d.Low == nil || d.Close == nil {
				continue
			}
			ap = (*d.Open + *d.High + *d.Low + *d.Close) / 4.0
		}
		sumDailyValue += float64(*d.Volume) * 100.0 * ap
		daysWithValue++
	}
	var avgDailyValue *float64
	if daysWithValue > 0 {
		v := sumDailyValue / float64(daysWithValue)
		avgDailyValue = &v
	}

	var firstPrev *float64
	var lastClose *float64
	for i := range days {
		if days[i].Prev != nil && firstPrev == nil {
			firstPrev = days[i].Prev
		}
		if days[i].Close != nil {
			lastClose = days[i].Close
		}
	}
	var priceChange *float64
	if firstPrev != nil && lastClose != nil && *firstPrev != 0 {
		pc := ((*lastClose - *firstPrev) / *firstPrev) * 100.0
		priceChange = &pc
	}

	return &Aggregation{
		Ticker:        ticker,
		Net:           net,
		BrokerTotals:  brokerTotals,
		TotBuy:        totBuy,
		TotSell:       totSell,
		TotNetVolume:  totNetVolume,
		TotalVolume:   totalVolume,
		Buyers:        buyers,
		Sellers:       sellers,
		AvgDailyValue: avgDailyValue,
		PriceChange:   priceChange,
	}
}

// ComputeActionWeight calculates CT adjusted volume ratio and weight factor.
func ComputeActionWeight(net map[string]float64, totBuy, totSell float64, dayCount int) (float64, float64) {
	sumAbsNet := 0.0
	for _, v := range net {
		sumAbsNet += math.Abs(v)
	}
	if sumAbsNet == 0 {
		return math.Inf(1), 0.25
	}
	vr := (totBuy + totSell) / sumAbsNet
	if vr > 7.0 {
		if dayCount > 15 {
			vr *= 0.7
		} else if dayCount > 7 {
			vr *= 0.8
		}
	}
	return vr, VolumeRatioWeight(vr)
}

// CumulativeDiffAtN sums top-n buyers minus abs top-n sellers.
func CumulativeDiffAtN(buyers, sellers []BrokerNet, n int) float64 {
	cumBuy := 0.0
	cumSell := 0.0
	for i := 1; i <= n; i++ {
		if i <= len(buyers) {
			cumBuy += buyers[i-1].Net
		}
		if i <= len(sellers) {
			cumSell += math.Abs(sellers[i-1].Net)
		}
	}
	return cumBuy - cumSell
}

// ComputeBandarActionClassification classifies net top-broker dominance weighted by action weight.
func ComputeBandarActionClassification(buyers, sellers []BrokerNet, totNetVolume, actionWeight float64) ActionLabel {
	nilaiTop3 := CumulativeDiffAtN(buyers, sellers, 3)
	nilaiTop5 := CumulativeDiffAtN(buyers, sellers, 5)
	sumSelisih := 0.0
	for n := 1; n <= 10; n++ {
		sumSelisih += CumulativeDiffAtN(buyers, sellers, n)
	}
	nilaiAvg10 := sumSelisih / 10.0
	avgNilai := (nilaiTop3 + nilaiTop5 + nilaiAvg10) / 3.0
	if totNetVolume == 0 {
		return ClassifyBandar(nil)
	}
	avgRatio := (avgNilai / totNetVolume) * 100.0
	classificationRatio := avgRatio * actionWeight
	return ClassifyBandar(&classificationRatio)
}

func rbPositionWeight(pos int) float64 {
	if pos <= 3 {
		return 1.0
	}
	if pos <= 6 {
		return 0.75
	}
	return 0.5
}

// DetermineRetailBandarOtomatis implements Criteria 1, 3, 4, 6, 7 for auto-detecting bandar retail nominee.
func DetermineRetailBandarOtomatis(
	net map[string]float64,
	brokerTotals map[string]BrokerTotal,
	buyers, sellers []BrokerNet,
	effRetail []string,
	totNetVolume, totVolume float64,
) (*string, map[string][]string) {
	retailSet := make(map[string]struct{}, len(effRetail))
	for _, b := range effRetail {
		retailSet[b] = struct{}{}
	}
	reasons := make(map[string][]string)
	reasonOrder := make([]string, 0, len(effRetail))
	addNominee := func(broker, reason string) {
		if _, exists := reasons[broker]; !exists {
			reasonOrder = append(reasonOrder, broker)
		}
		reasons[broker] = append(reasons[broker], reason)
	}

	// Criterion 1
	totals := make(map[string]float64, len(effRetail))
	for _, b := range effRetail {
		bt := brokerTotals[b]
		totals[b] = bt.Buy + bt.Sell
	}
	for _, b := range effRetail {
		othersSum := 0.0
		for _, x := range effRetail {
			if x != b {
				othersSum += totals[x]
			}
		}
		coef, ok := RetailBandarCoef[b]
		if !ok {
			coef = 1.0
		}
		if totals[b] > othersSum*coef {
			addNominee(b, fmt.Sprintf("C1: Transaksi %s > %.1fx total retail lain", b, coef))
		}
	}

	// Criterion 3
	buyWeight := 0.0
	for i, x := range buyers {
		if _, ok := retailSet[x.Broker]; ok {
			buyWeight += rbPositionWeight(i + 1)
		}
	}
	sellWeight := 0.0
	for i, x := range sellers {
		if _, ok := retailSet[x.Broker]; ok {
			sellWeight += rbPositionWeight(i + 1)
		}
	}
	if buyWeight > 4.0 {
		lim := min(5, len(sellers))
		for _, x := range sellers[:lim] {
			if _, ok := retailSet[x.Broker]; ok {
				addNominee(x.Broker, "C3: Dominasi retail di Top Buyer, muncul di Top 5 Seller")
			}
		}
	}
	if sellWeight > 4.0 {
		lim := min(5, len(buyers))
		for _, x := range buyers[:lim] {
			if _, ok := retailSet[x.Broker]; ok {
				addNominee(x.Broker, "C3: Dominasi retail di Top Seller, muncul di Top 5 Buyer")
			}
		}
	}

	// Criterion 4
	type rankItem struct {
		broker string
		val    float64
		rank   int
	}
	for _, side := range []struct {
		name string
		list []BrokerNet
	}{{"Buyer", buyers}, {"Seller", sellers}} {
		retailInList := make([]rankItem, 0, len(effRetail))
		for i, x := range side.list {
			if _, ok := retailSet[x.Broker]; ok {
				retailInList = append(retailInList, rankItem{
					broker: x.Broker,
					val:    math.Abs(x.Net),
					rank:   i + 1,
				})
			}
		}
		if len(retailInList) >= 2 {
			sort.SliceStable(retailInList, func(i, j int) bool {
				return retailInList[i].val > retailInList[j].val
			})
			first, second := retailInList[0], retailInList[1]
			if first.rank <= 3 && second.val > 0 && first.val >= second.val*3.0 {
				addNominee(first.broker, fmt.Sprintf("C4: Top %d %s & >= 3x retail berikutnya", first.rank, side.name))
			}
		}
	}

	// Criterion 6
	if totNetVolume != 0 {
		var buyerMembers []BrokerNet
		var sellerMembers []BrokerNet
		for _, b := range effRetail {
			v := net[b]
			if v > 0 {
				buyerMembers = append(buyerMembers, BrokerNet{Broker: b, Net: v})
			} else if v < 0 {
				sellerMembers = append(sellerMembers, BrokerNet{Broker: b, Net: v})
			}
		}
		if len(buyerMembers) == 1 && (math.Abs(buyerMembers[0].Net)/totNetVolume*100.0) > 10.0 {
			addNominee(buyerMembers[0].Broker, "C6: Satu-satunya retail net buyer (>10% net vol)")
		}
		if len(sellerMembers) == 1 && (math.Abs(sellerMembers[0].Net)/totNetVolume*100.0) > 10.0 {
			addNominee(sellerMembers[0].Broker, "C6: Satu-satunya retail net seller (>10% net vol)")
		}
	}

	// Criterion 7
	if totVolume != 0 {
		for _, b := range effRetail {
			thr, ok := RetailBandarTransThreshold[b]
			if !ok {
				continue
			}
			pct := (totals[b] / totVolume) * 100.0
			if pct > thr {
				addNominee(b, fmt.Sprintf("C7: Porsi transaksi %.1f%% > %.1f%%", shared.PyRound(pct, 1), thr))
			}
		}
	}

	if len(reasons) == 0 {
		return nil, map[string][]string{}
	}

	var best string
	bestVal := math.Inf(-1)
	for _, b := range reasonOrder {
		val := math.Abs(net[b])
		if val > bestVal {
			bestVal = val
			best = b
		}
	}
	return &best, reasons
}

// ComputeDecision evaluates bandar and retail classifications and broker categories.
func ComputeDecision(
	agg *Aggregation,
	dayCount int,
	uwSet map[string]struct{},
	autoRetail bool,
	retailOverride string,
) Decision {
	vr, weight := ComputeActionWeight(agg.Net, agg.TotBuy, agg.TotSell, dayCount)
	var vrAdj *float64
	if !math.IsInf(vr, 0) {
		v := vr
		vrAdj = &v
	}
	bandarAction := ComputeBandarActionClassification(agg.Buyers, agg.Sellers, agg.TotNetVolume, weight)

	effRetailList := EffectiveRetailBrokers(agg.Ticker, uwSet)
	retailSet := make(map[string]struct{}, len(effRetailList))
	for _, b := range effRetailList {
		retailSet[b] = struct{}{}
	}
	retailBandarAuto, rbReasons := DetermineRetailBandarOtomatis(
		agg.Net, agg.BrokerTotals, agg.Buyers, agg.Sellers, effRetailList, agg.TotNetVolume, agg.TotBuy+agg.TotSell,
	)

	manualClean := strings.ToUpper(strings.TrimSpace(retailOverride))
	retailBandarBrokers := make([]string, 0, 2)
	hasBandar := make(map[string]struct{}, 2)
	if manualClean != "" {
		if _, ok := retailSet[manualClean]; ok {
			retailBandarBrokers = append(retailBandarBrokers, manualClean)
			hasBandar[manualClean] = struct{}{}
		}
	}
	if autoRetail && retailBandarAuto != nil {
		if _, dup := hasBandar[*retailBandarAuto]; !dup {
			retailBandarBrokers = append(retailBandarBrokers, *retailBandarAuto)
			hasBandar[*retailBandarAuto] = struct{}{}
		}
	}

	var retailBandarBroker *string
	if manualClean != "" {
		if _, ok := retailSet[manualClean]; ok {
			m := manualClean
			retailBandarBroker = &m
		}
	}
	if retailBandarBroker == nil && autoRetail && retailBandarAuto != nil {
		a := *retailBandarAuto
		retailBandarBroker = &a
	}

	finalRetailSet := make(map[string]struct{}, len(retailSet))
	for b := range retailSet {
		if _, excl := hasBandar[b]; !excl {
			finalRetailSet[b] = struct{}{}
		}
	}

	retailNet := 0.0
	for b, v := range agg.Net {
		if _, ok := finalRetailSet[b]; ok {
			retailNet += v
		}
	}
	var retailRatio *float64
	if agg.TotNetVolume != 0 {
		rr := (retailNet / agg.TotNetVolume) * 100.0
		retailRatio = &rr
	}
	retailAction := ClassifyRetailAction(retailRatio)

	brokerCategory := make(map[string]string, len(agg.Buyers)+len(agg.Sellers))
	for _, x := range agg.Buyers {
		if _, isB := hasBandar[x.Broker]; isB {
			brokerCategory[x.Broker] = "bandar"
		} else {
			brokerCategory[x.Broker] = CategoryOf(x.Broker, agg.Ticker, uwSet)
		}
	}
	for _, x := range agg.Sellers {
		if _, isB := hasBandar[x.Broker]; isB {
			brokerCategory[x.Broker] = "bandar"
		} else {
			brokerCategory[x.Broker] = CategoryOf(x.Broker, agg.Ticker, uwSet)
		}
	}

	return Decision{
		ActionWeight:           weight,
		VolumeRatioAdjusted:    vrAdj,
		BandarAction:           bandarAction,
		RetailBandarBroker:     retailBandarBroker,
		RetailBandarBrokers:    retailBandarBrokers,
		RetailBandarAutoResult: retailBandarAuto,
		RetailBandarReasons:    rbReasons,
		EffectiveRetailList:    effRetailList,
		FinalRetailSet:         finalRetailSet,
		BrokerCategory:         brokerCategory,
		RetailRatio:            retailRatio,
		RetailAction:           retailAction,
	}
}
