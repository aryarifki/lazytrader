package bandarmology

import (
	"math"
	"sort"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
)

var profileOrder = []string{
	"smart_foreign", "local_institutional", "market_maker",
	"bandar_gorengan", "retail", "lainnya",
}

var retailBrokerCodes = map[string]bool{
	"YP": true, "PD": true, "XC": true, "KK": true, "NI": true,
	"XL": true, "SQ": true, "XA": true, "AZ": true,
}

// BrokerSubtype mirrors bandarmology_service.broker_subtype.
func BrokerSubtype(pType string, netVal, freq float64) string {
	if shared.ParticipantLabel(pType) != "FOREIGN" {
		return "-"
	}
	absNet := math.Abs(netVal)
	f := math.Max(freq, 1.0)
	avgVal := absNet / f
	if avgVal >= 500_000_000.0 || (absNet >= 5_000_000_000.0 && f <= 500.0) {
		return "Institutional"
	}
	if f >= 2000.0 || avgVal <= 100_000_000.0 {
		return "Speculative"
	}
	return "Mixed"
}

// PriceAtOrBefore returns the last PriceRow on or before ts.
func PriceAtOrBefore(prices []PriceRow, ts time.Time) *PriceRow {
	var best *PriceRow
	for i := range prices {
		if !prices[i].Date.After(ts) {
			if best == nil || !prices[i].Date.Before(best.Date) {
				best = &prices[i]
			}
		}
	}
	return best
}

// ReturnToDate mirrors bandarmology_service.return_to_date.
func ReturnToDate(prices []PriceRow, ts time.Time, periods int) *float64 {
	sub := make([]PriceRow, 0, len(prices))
	for _, p := range prices {
		if !p.Date.After(ts) {
			sub = append(sub, p)
		}
	}
	sort.Slice(sub, func(i, j int) bool { return sub[i].Date.Before(sub[j].Date) })
	if len(sub) <= periods {
		return nil
	}
	last := sub[len(sub)-1].Close
	base := sub[len(sub)-periods-1].Close
	if last == nil || base == nil || *base == 0 {
		return nil
	}
	v := *last / *base - 1.0
	return &v
}

// FlowRowAt returns the latest FlowRow for ticker on or before ts.
func FlowRowAt(flows []FlowRow, ticker string, ts time.Time) *FlowRow {
	var best *FlowRow
	for i := range flows {
		if flows[i].Ticker == ticker && !flows[i].Date.After(ts) {
			if best == nil || !flows[i].Date.Before(best.Date) {
				best = &flows[i]
			}
		}
	}
	return best
}

// LatestActivityDate returns the latest ActivityRow date for ticker on or before ts.
func LatestActivityDate(acts []ActivityRow, ticker string, ts time.Time) *time.Time {
	var best *time.Time
	for i := range acts {
		if acts[i].Ticker == ticker && !acts[i].Date.After(ts) {
			if best == nil || acts[i].Date.After(*best) {
				d := acts[i].Date
				best = &d
			}
		}
	}
	return best
}

// SparklineValues renders 5-char +/-/0 sparkline for a broker up to endTS.
func SparklineValues(acts []ActivityRow, brokerCode string, endTS time.Time, days int) string {
	sub := make([]ActivityRow, 0, days)
	for _, a := range acts {
		if a.BrokerCode == brokerCode && !a.Date.After(endTS) {
			sub = append(sub, a)
		}
	}
	if len(sub) == 0 {
		return "-----"
	}
	sort.Slice(sub, func(i, j int) bool { return sub[i].Date.Before(sub[j].Date) })
	if len(sub) > days {
		sub = sub[len(sub)-days:]
	}
	var sb strings.Builder
	for _, a := range sub {
		switch {
		case a.NetValue > 0:
			sb.WriteByte('+')
		case a.NetValue < 0:
			sb.WriteByte('-')
		default:
			sb.WriteByte('0')
		}
	}
	return sb.String()
}

// TopNetBrokerSummary mirrors analysis.top_net_broker_summary.
func TopNetBrokerSummary(acts []ActivityRow, ticker string, tradeDate *time.Time, topN int) ([]ActivityRow, []ActivityRow) {
	if len(acts) == 0 {
		return nil, nil
	}
	var target time.Time
	if tradeDate != nil {
		target = *tradeDate
	} else {
		for _, a := range acts {
			if a.Date.After(target) {
				target = a.Date
			}
		}
	}
	buyers := make([]ActivityRow, 0)
	sellers := make([]ActivityRow, 0)
	tk := strings.ToUpper(strings.TrimSpace(ticker))
	for _, a := range acts {
		if a.Ticker != tk || !a.Date.Equal(target) {
			continue
		}
		if a.NetValue > 0 {
			buyers = append(buyers, a)
		} else if a.NetValue < 0 {
			sellers = append(sellers, a)
		}
	}
	sort.SliceStable(buyers, func(i, j int) bool { return buyers[i].NetValue > buyers[j].NetValue })
	sort.SliceStable(sellers, func(i, j int) bool { return sellers[i].NetValue < sellers[j].NetValue })
	if len(buyers) > topN {
		buyers = buyers[:topN]
	}
	if len(sellers) > topN {
		sellers = sellers[:topN]
	}
	return buyers, sellers
}

// ContradictionAlerts mirrors bandarmology_service.contradiction_alerts.
func ContradictionAlerts(signal string, ret5d, ret10d, foreign5d, smartCum *float64) []string {
	raw := strings.ToUpper(strings.TrimSpace(signal))
	alerts := make([]string, 0, 4)
	if shared.DistSignals[raw] && ((ret5d != nil && *ret5d > 0) || (ret10d != nil && *ret10d > 0)) {
		alerts = append(alerts, "Distribution while price is still rising — potential unfinished distribution or new buyer absorption. Monitor volume.")
	}
	if shared.AccSignals[raw] && ret5d != nil && *ret5d < 0 {
		alerts = append(alerts, "Accumulation signal with negative 5D return — accumulation may be early, failed, or absorbed by larger supply.")
	}
	if foreign5d != nil && *foreign5d < 0 && shared.AccSignals[raw] {
		alerts = append(alerts, "Aggregate accumulation conflicts with foreign net selling — check whether the move is driven by local brokers.")
	}
	if smartCum != nil && *smartCum < 0 && shared.AccSignals[raw] {
		alerts = append(alerts, "Signal is accumulation but smart-money cumulative flow is negative in the selected window.")
	}
	return alerts
}

// SmartDailyRow represents one row of smart-money daily + cumulative net.
type SmartDailyRow struct {
	Date          time.Time
	SmartNet      float64
	CumulativeNet float64
}

// SmartDailyFromActivity mirrors bandarmology_service.smart_daily_from_activity.
func SmartDailyFromActivity(acts []ActivityRow) []SmartDailyRow {
	if len(acts) == 0 {
		return nil
	}
	byDate := make(map[time.Time]float64)
	for _, a := range acts {
		if shared.IsSmartProfile(shared.BrokerProfileOf(a.BrokerCode)) {
			byDate[a.Date] += a.NetValue
		}
	}
	if len(byDate) == 0 {
		return nil
	}
	dates := make([]time.Time, 0, len(byDate))
	for d := range byDate {
		dates = append(dates, d)
	}
	sort.Slice(dates, func(i, j int) bool { return dates[i].Before(dates[j]) })
	out := make([]SmartDailyRow, len(dates))
	cum := 0.0
	for i, d := range dates {
		cum += byDate[d]
		out[i] = SmartDailyRow{Date: d, SmartNet: byDate[d], CumulativeNet: cum}
	}
	return out
}

// ProfileFlowEntry is shared by detail & broker-flow profile summaries.
type ProfileFlowEntry struct {
	Profile     string
	Label       string
	Description string
	Net         float64
	Buy         float64
	Sell        float64
	TopBrokers  []map[string]any
}

// ProfileFlowFromActivity mirrors bandarmology_service.profile_flow_from_activity & analysis.broker_profile_flow_table.
func ProfileFlowFromActivity(acts []ActivityRow) []ProfileFlowEntry {
	if len(acts) == 0 {
		return nil
	}
	type pbKey struct {
		profile, broker, ptype string
	}
	type pbAgg struct {
		net, buy, sell float64
	}
	grouped := make(map[pbKey]*pbAgg)
	for _, a := range acts {
		prof := shared.BrokerProfileOf(a.BrokerCode)
		k := pbKey{prof, a.BrokerCode, a.ParticipantType}
		g, ok := grouped[k]
		if !ok {
			g = &pbAgg{}
			grouped[k] = g
		}
		g.net += a.NetValue
		g.buy += a.BuyValue
		g.sell += a.SellValue
	}

	out := make([]ProfileFlowEntry, 0, len(profileOrder))
	for _, prof := range profileOrder {
		var members []map[string]any
		var totNet, totBuy, totSell float64
		for k, g := range grouped {
			if k.profile != prof {
				continue
			}
			totNet += g.net
			totBuy += g.buy
			totSell += g.sell
			members = append(members, map[string]any{
				"broker_code":      k.broker,
				"participant_type": k.ptype,
				"net":              g.net,
			})
		}
		if len(members) == 0 {
			continue
		}
		sort.SliceStable(members, func(i, j int) bool {
			return math.Abs(members[i]["net"].(float64)) > math.Abs(members[j]["net"].(float64))
		})
		if len(members) > 6 {
			members = members[:6]
		}
		meta := shared.ProfileMeta[prof]
		out = append(out, ProfileFlowEntry{
			Profile:     prof,
			Label:       meta.Label,
			Description: meta.Description,
			Net:         totNet,
			Buy:         totBuy,
			Sell:        totSell,
			TopBrokers:  members,
		})
	}
	return out
}

// ProfileBrokerDetailTable mirrors bandarmology_service.profile_broker_detail_table.
func ProfileBrokerDetailTable(acts []ActivityRow) []map[string]any {
	if len(acts) == 0 {
		return []map[string]any{}
	}
	type pbKey struct {
		profile, broker, ptype string
	}
	type pbAgg struct {
		buy, sell, net, freq float64
		dates                map[string]struct{}
	}
	grouped := make(map[pbKey]*pbAgg)
	for _, a := range acts {
		k := pbKey{shared.BrokerProfileOf(a.BrokerCode), a.BrokerCode, a.ParticipantType}
		g, ok := grouped[k]
		if !ok {
			g = &pbAgg{dates: make(map[string]struct{})}
			grouped[k] = g
		}
		g.buy += a.BuyValue
		g.sell += a.SellValue
		g.net += a.NetValue
		g.freq += shared.Finite(a.Frequency, 0)
		g.dates[fmtDateISO(a.Date)] = struct{}{}
	}

	type sortable struct {
		prof, broker string
		net          float64
		row          map[string]any
	}
	items := make([]sortable, 0, len(grouped))
	for k, g := range grouped {
		label := k.profile
		if m, ok := shared.ProfileMeta[k.profile]; ok {
			label = m.Label
		}
		avgTx := math.Abs(g.net) / math.Max(g.freq, 1.0)
		items = append(items, sortable{
			prof:   k.profile,
			broker: k.broker,
			net:    g.net,
			row: map[string]any{
				"profile":      label,
				"profile_key":  k.profile,
				"broker":       k.broker,
				"type":         shared.ParticipantLabel(k.ptype),
				"buy":          g.buy,
				"sell":         g.sell,
				"net":          g.net,
				"freq":         g.freq,
				"days":         len(g.dates),
				"avg_value_tx": avgTx,
			},
		})
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].prof != items[j].prof {
			return items[i].prof < items[j].prof
		}
		return items[i].net > items[j].net
	})
	rows := make([]map[string]any, len(items))
	for i, it := range items {
		rows[i] = it.row
	}
	return rows
}

// BrokerDistributionDataRange mirrors bandarmology_service.broker_distribution_data_range.
func BrokerDistributionDataRange(acts []ActivityRow, distStart, distEnd time.Time) map[string]any {
	type bpKey struct{ broker, ptype string }
	type bpAgg struct {
		buyVal, sellVal, netVal, freq, buyLot, sellLot float64
		vwapBuyVal, vwapBuyLot                         float64
		vwapSellVal, vwapSellLot                       float64
		bAvgSum                                        float64
		bAvgCnt                                        int
		sAvgSum                                        float64
		sAvgCnt                                        int
		hasBuyLot, hasSellLot                          bool
	}
	grouped := make(map[bpKey]*bpAgg)
	count := 0
	for _, a := range acts {
		if a.Date.Before(distStart) || a.Date.After(distEnd) {
			continue
		}
		count++
		k := bpKey{a.BrokerCode, a.ParticipantType}
		g, ok := grouped[k]
		if !ok {
			g = &bpAgg{}
			grouped[k] = g
		}
		g.buyVal += a.BuyValue
		g.sellVal += a.SellValue
		g.netVal += a.NetValue
		g.freq += shared.Finite(a.Frequency, 0)
		bLot := shared.Finite(a.BuyLot, 0)
		sLot := shared.Finite(a.SellLot, 0)
		if a.BuyLot != nil {
			g.hasBuyLot = true
		}
		if a.SellLot != nil {
			g.hasSellLot = true
		}
		g.buyLot += bLot
		g.sellLot += sLot
		if bLot > 0 {
			g.vwapBuyVal += a.BuyValue
			g.vwapBuyLot += bLot
		}
		if sLot > 0 {
			g.vwapSellVal += a.SellValue
			g.vwapSellLot += sLot
		}
		if a.BuyAvgPrice != nil && *a.BuyAvgPrice > 0 {
			g.bAvgSum += *a.BuyAvgPrice
			g.bAvgCnt++
		}
		if a.SellAvgPrice != nil && *a.SellAvgPrice > 0 {
			g.sAvgSum += *a.SellAvgPrice
			g.sAvgCnt++
		}
	}
	if count == 0 {
		return map[string]any{"buyers": []any{}, "sellers": []any{}, "edges": []any{}}
	}

	type sideRec struct {
		broker, pLabel               string
		buyVal, sellVal, netVal, freq float64
		buyLot, sellLot              *float64
		buyAvg, sellAvg              *float64
		rem                          float64
	}
	var buyers, sellers []sideRec
	for k, g := range grouped {
		var bLotPtr, sLotPtr, bAvgPtr, sAvgPtr *float64
		if g.hasBuyLot {
			v := g.buyLot
			bLotPtr = &v
		}
		if g.hasSellLot {
			v := g.sellLot
			sLotPtr = &v
		}
		if g.vwapBuyLot > 0 {
			v := g.vwapBuyVal / (g.vwapBuyLot * 100.0)
			bAvgPtr = &v
		} else if g.bAvgCnt > 0 {
			v := g.bAvgSum / float64(g.bAvgCnt)
			bAvgPtr = &v
		}
		if g.vwapSellLot > 0 {
			v := g.vwapSellVal / (g.vwapSellLot * 100.0)
			sAvgPtr = &v
		} else if g.sAvgCnt > 0 {
			v := g.sAvgSum / float64(g.sAvgCnt)
			sAvgPtr = &v
		}

		rec := sideRec{
			broker: k.broker, pLabel: shared.ParticipantLabel(k.ptype),
			buyVal: g.buyVal, sellVal: g.sellVal, netVal: g.netVal, freq: g.freq,
			buyLot: bLotPtr, sellLot: sLotPtr, buyAvg: bAvgPtr, sellAvg: sAvgPtr,
			rem: math.Abs(g.netVal),
		}
		if g.netVal > 0 {
			buyers = append(buyers, rec)
		} else if g.netVal < 0 {
			sellers = append(sellers, rec)
		}
	}
	sort.SliceStable(buyers, func(i, j int) bool { return buyers[i].netVal > buyers[j].netVal })
	sort.SliceStable(sellers, func(i, j int) bool { return sellers[i].netVal < sellers[j].netVal })

	nB8 := min(len(buyers), 8)
	nS8 := min(len(sellers), 8)
	bRem := make([]float64, nB8)
	sRem := make([]float64, nS8)
	for i := 0; i < nB8; i++ {
		bRem[i] = buyers[i].netVal
	}
	for i := 0; i < nS8; i++ {
		sRem[i] = math.Abs(sellers[i].netVal)
	}
	edges := make([]map[string]any, 0)
	sIdx := 0
	for bIdx := 0; bIdx < nB8; bIdx++ {
		bLeft := bRem[bIdx]
		for bLeft > 1e-9 && sIdx < nS8 {
			sLeft := sRem[sIdx]
			if sLeft <= 1e-9 {
				sIdx++
				continue
			}
			matched := math.Min(bLeft, sLeft)
			edges = append(edges, map[string]any{
				"buyer_code":    buyers[bIdx].broker,
				"buyer_type":    buyers[bIdx].pLabel,
				"seller_code":   sellers[sIdx].broker,
				"seller_type":   sellers[sIdx].pLabel,
				"matched_value": matched,
			})
			bLeft -= matched
			sRem[sIdx] = sLeft - matched
			if sRem[sIdx] <= 1e-9 {
				sIdx++
			}
		}
	}

	buyersOut := make([]map[string]any, 0, min(len(buyers), 10))
	for i := 0; i < min(len(buyers), 10); i++ {
		b := buyers[i]
		buyersOut = append(buyersOut, map[string]any{
			"broker":    b.broker,
			"type":      b.pLabel,
			"buy_value": b.buyVal,
			"sell_value": b.sellVal,
			"net_value": b.netVal,
			"freq":      b.freq,
			"buy_lot":   b.buyLot,
			"buy_avg":   b.buyAvg,
		})
	}
	sellersOut := make([]map[string]any, 0, min(len(sellers), 10))
	for i := 0; i < min(len(sellers), 10); i++ {
		s := sellers[i]
		sellersOut = append(sellersOut, map[string]any{
			"broker":     s.broker,
			"type":       s.pLabel,
			"buy_value":  s.buyVal,
			"sell_value": s.sellVal,
			"net_value":  s.netVal,
			"freq":       s.freq,
			"sell_lot":   s.sellLot,
			"sell_avg":   s.sellAvg,
		})
	}
	return map[string]any{
		"buyers":     buyersOut,
		"sellers":    sellersOut,
		"edges":      edges,
		"dist_start": fmtDateISO(distStart),
		"dist_end":   fmtDateISO(distEnd),
		"dist_date":  fmtDateISO(distEnd),
	}
}

// BrokerSummaryTable mirrors bandarmology_service.broker_summary_table.
func BrokerSummaryTable(dist map[string]any) []map[string]any {
	buyers, _ := dist["buyers"].([]map[string]any)
	sellers, _ := dist["sellers"].([]map[string]any)
	maxLen := max(len(buyers), len(sellers))
	if maxLen > 10 {
		maxLen = 10
	}
	rows := make([]map[string]any, 0, maxLen)
	for i := 0; i < maxLen; i++ {
		row := make(map[string]any, 10)
		if i < len(buyers) {
			b := buyers[i]
			row["buy_broker"] = b["broker"]
			row["buy_type"] = b["type"]
			row["buy_value"] = b["buy_value"]
			row["buy_lot"] = b["buy_lot"]
			row["buy_avg"] = b["buy_avg"]
		} else {
			row["buy_broker"] = ""
			row["buy_type"] = ""
			row["buy_value"] = nil
			row["buy_lot"] = nil
			row["buy_avg"] = nil
		}
		if i < len(sellers) {
			s := sellers[i]
			sv, _ := s["sell_value"].(float64)
			row["sell_broker"] = s["broker"]
			row["sell_type"] = s["type"]
			row["sell_value"] = math.Abs(sv)
			row["sell_lot"] = s["sell_lot"]
			row["sell_avg"] = s["sell_avg"]
		} else {
			row["sell_broker"] = ""
			row["sell_type"] = ""
			row["sell_value"] = nil
			row["sell_lot"] = nil
			row["sell_avg"] = nil
		}
		rows = append(rows, row)
	}
	return rows
}

// BrokerCompareData mirrors bandarmology_service.broker_compare_data (values in Rp B).
func BrokerCompareData(acts []ActivityRow, selected []string, flowMode string) []map[string]any {
	if len(acts) == 0 || len(selected) == 0 {
		return []map[string]any{}
	}
	selSet := make(map[string]bool, len(selected))
	for _, c := range selected {
		selSet[c] = true
	}
	dateSet := make(map[time.Time]struct{})
	byDateBroker := make(map[time.Time]map[string]float64)
	for _, a := range acts {
		if !selSet[a.BrokerCode] {
			continue
		}
		dateSet[a.Date] = struct{}{}
		m, ok := byDateBroker[a.Date]
		if !ok {
			m = make(map[string]float64)
			byDateBroker[a.Date] = m
		}
		m[a.BrokerCode] += a.NetValue
	}
	if len(dateSet) == 0 {
		return []map[string]any{}
	}
	dates := make([]time.Time, 0, len(dateSet))
	for d := range dateSet {
		dates = append(dates, d)
	}
	sort.Slice(dates, func(i, j int) bool { return dates[i].Before(dates[j]) })

	cum := strings.EqualFold(strings.TrimSpace(flowMode), "cumulative")
	running := make(map[string]float64, len(selected))
	out := make([]map[string]any, 0, len(dates))
	for _, d := range dates {
		item := map[string]any{"date": fmtDateISO(d)}
		m := byDateBroker[d]
		for _, c := range selected {
			v := m[c]
			if cum {
				running[c] += v
				v = running[c]
			}
			item[c] = v / 1e9
		}
		out = append(out, item)
	}
	return out
}

// PricePerformanceTable mirrors analysis.price_performance_table.
func PricePerformanceTable(prices []PriceRow) []map[string]any {
	if len(prices) == 0 {
		return []map[string]any{}
	}
	px := append([]PriceRow(nil), prices...)
	sort.Slice(px, func(i, j int) bool { return px[i].Date.Before(px[j].Date) })
	latest := px[len(px)-1]
	windows := []struct {
		label   string
		periods int
	}{
		{"1D", 1}, {"1W", 5}, {"1M", 21}, {"3M", 63}, {"6M", 126},
		{"1Y", 252}, {"3Y", 756}, {"5Y", 1260}, {"10Y", 2520},
	}
	rows := make([]map[string]any, 0, len(windows)+1)
	for _, w := range windows {
		var ret *float64
		if len(px) > w.periods && latest.Close != nil && px[len(px)-w.periods-1].Close != nil {
			base := *px[len(px)-w.periods-1].Close
			if base != 0 {
				v := *latest.Close/base - 1.0
				ret = &v
			}
		}
		rows = append(rows, map[string]any{"timeframe": w.label, "return": ret})
	}
	var ytd *float64
	yr := latest.Date.Year()
	var sameYr []PriceRow
	for _, p := range px {
		if p.Date.Year() == yr {
			sameYr = append(sameYr, p)
		}
	}
	if len(sameYr) > 1 && latest.Close != nil && sameYr[0].Close != nil && *sameYr[0].Close != 0 {
		v := *latest.Close/(*sameYr[0].Close) - 1.0
		ytd = &v
	}
	// Insert YTD at index 5
	res := make([]map[string]any, 0, len(rows)+1)
	res = append(res, rows[:5]...)
	res = append(res, map[string]any{"timeframe": "YTD", "return": ytd})
	res = append(res, rows[5:]...)
	return res
}

// EventStudyRow holds one signal event normalized to t=0 = 100.0.
type EventStudyRow struct {
	Ticker            string
	SignalDate        time.Time
	BandarSignal      string
	BandarSignalScore *float64
	Values            map[int]*float64
}

// EventStudyTable mirrors analysis.event_study_table.
func EventStudyTable(prices []PriceRow, flows []FlowRow, horizons []int, lookbackDays *int, signals map[string]bool, endDate *time.Time, minEventSpacing int) []EventStudyRow {
	px := make([]PriceRow, 0, len(prices))
	for _, p := range prices {
		if p.Close == nil || *p.Close <= 0 {
			continue
		}
		if endDate != nil && p.Date.After(*endDate) {
			continue
		}
		px = append(px, p)
	}
	fl := make([]FlowRow, 0, len(flows))
	for _, f := range flows {
		if endDate != nil && f.Date.After(*endDate) {
			continue
		}
		fl = append(fl, f)
	}
	if len(px) == 0 || len(fl) == 0 {
		return nil
	}
	sort.Slice(px, func(i, j int) bool { return px[i].Date.Before(px[j].Date) })
	sort.Slice(fl, func(i, j int) bool { return fl[i].Date.Before(fl[j].Date) })

	if lookbackDays != nil {
		maxD := fl[len(fl)-1].Date
		cutoff := maxD.AddDate(0, 0, -*lookbackDays)
		kept := make([]FlowRow, 0, len(fl))
		for _, f := range fl {
			if !f.Date.Before(cutoff) {
				kept = append(kept, f)
			}
		}
		fl = kept
	}

	idxByDate := make(map[string]int, len(px))
	for i, p := range px {
		idxByDate[fmtDateISO(p.Date)] = i
	}

	lastIdxBySig := make(map[string]int)
	rows := make([]EventStudyRow, 0)
	for _, f := range fl {
		sig := ""
		if f.BandarSignal != nil {
			sig = *f.BandarSignal
		}
		if len(signals) > 0 && !signals[sig] {
			continue
		}
		startIdx, ok := idxByDate[fmtDateISO(f.Date)]
		if !ok {
			continue
		}
		startClose := *px[startIdx].Close
		if startClose <= 0 {
			continue
		}
		if prevIdx, seen := lastIdxBySig[sig]; minEventSpacing > 0 && seen && (startIdx-prevIdx) < minEventSpacing {
			continue
		}
		lastIdxBySig[sig] = startIdx

		vals := make(map[int]*float64, len(horizons)+1)
		h0 := 100.0
		vals[0] = &h0
		for _, h := range horizons {
			idx := startIdx + h
			if idx >= 0 && idx < len(px) && px[idx].Close != nil {
				v := (*px[idx].Close / startClose) * 100.0
				vals[h] = &v
			} else {
				vals[h] = nil
			}
		}
		rows = append(rows, EventStudyRow{
			Ticker:            f.Ticker,
			SignalDate:        f.Date,
			BandarSignal:      sig,
			BandarSignalScore: f.BandarSignalScore,
			Values:            vals,
		})
	}
	return rows
}
