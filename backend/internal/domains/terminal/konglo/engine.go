package konglo

import (
	"sort"
	"strings"

	"github.com/aryarifki/lazytrader/backend/internal/domains/terminal/shared"
)

// FundamentalTickers is the canonical 50 fundamental benchmark tickers from konglo_engine.py.
var FundamentalTickers = []string{
	"BBCA", "BBRI", "BMRI", "TLKM", "ASII", "ANTM", "BBNI", "UNTR", "TINS", "INDF",
	"MAPI", "KLBF", "INCO", "ISAT", "MEDC", "CPIN", "PGAS", "ITMG", "JPFA", "ICBP",
	"GGRM", "AKRA", "SMGR", "PTBA", "UNVR", "BBTN", "ERAA", "HRTA", "PWON", "INTP",
	"PGEO", "BRIS", "EXCL", "ELSA", "CMRY", "MYOR", "JSMR", "BDMN", "HMSP", "SIDO",
	"ACES", "HRUM", "BTPS", "BNGA", "ULTJ", "SMRA", "POWR", "MPMX", "BIRD", "BSSR",
}

// CanonicalKongloOrder is the canonical display ordering for conglomerate groups from konglo_engine.py.
var CanonicalKongloOrder = []string{
	"Djarum Group", "SSI Group", "Salim Group", "Lippo Group", "Sinarmas Group",
	"MNC Group", "Astra Group", "Saratoga Group", "Emtek Group", "Barito Group",
	"Bakrie Group", "Triputra Group", "Aguan Group", "Kresna Group", "Northstar Group",
	"Thohir Group", "Hapsoro Group", "Tancorp Group", "Ciputra Group", "LTK Group",
	"Rajawali Group", "Mayapada Group", "Haji Isam Group", "Sinarmas Land", "Hashim Group",
	"Harita Group", "Panin Group", "CT Corp Group", "Alfa Group",
}

// CalculateKongloReport computes all 7 tabs for Daily Konglo Report verbatim.
func CalculateKongloReport(stocks []StockRow, groupOrder []string, groupsMap map[string][]string) *Report {
	stockByTicker := make(map[string]ReportStockItem, len(stocks))
	activeStocks := make([]ReportStockItem, 0, len(stocks))

	for _, s := range stocks {
		tk := strings.ToUpper(strings.TrimSpace(s.Ticker))
		last := floatOrZero(s.Last)
		prev := last
		if s.Prev != nil && *s.Prev != 0 {
			prev = *s.Prev
		}
		var chgPct float64
		if prev > 0 {
			chgPct = (last - prev) / prev * 100.0
		} else {
			chgPct = floatOrZero(s.Chg)
		}
		shares := floatOrZero(s.SharesListed)
		mcChgTrillion := 0.0
		if shares > 0 {
			mcChgTrillion = shared.PyRound((shares*(last-prev))/1e12, 4)
		}
		val := floatOrZero(s.Val)
		var vol float64
		if s.Vol != nil {
			vol = float64(*s.Vol)
		}
		var freq int64
		if s.Freq != nil {
			freq = *s.Freq
		}
		name := tk
		if s.Name != nil {
			name = *s.Name
		}
		group := "Others"
		if s.GroupName != nil && *s.GroupName != "" {
			group = *s.GroupName
		}

		item := ReportStockItem{
			Ticker: tk,
			Name:   name,
			Group:  group,
			Prev:   prev,
			Close:  last,
			Pct:    chgPct,
			Val:    val,
			Vol:    vol,
			Freq:   freq,
			Shares: shares,
			MCChg:  mcChgTrillion,
		}
		stockByTicker[tk] = item
		if last > 0 {
			activeStocks = append(activeStocks, item)
		}
	}

	// 1. Canonical + insertion group ordering
	inCanonical := make(map[string]bool, len(CanonicalKongloOrder))
	orderedGroups := make([]string, 0, len(groupsMap))
	for _, g := range CanonicalKongloOrder {
		inCanonical[g] = true
		if _, ok := groupsMap[g]; ok {
			orderedGroups = append(orderedGroups, g)
		}
	}
	for _, g := range groupOrder {
		if !inCanonical[g] {
			orderedGroups = append(orderedGroups, g)
		}
	}

	kongloStocks := make([]ReportStockItem, 0, 400)
	uniqueKongloStocks := make([]ReportStockItem, 0, 400)
	seenKonglo := make(map[string]bool, 400)
	missingKonglo := make([]MissingKonglo, 0)

	for _, gname := range orderedGroups {
		for _, tk := range groupsMap[gname] {
			p, ok := stockByTicker[tk]
			if !ok || p.Close <= 0 {
				missingKonglo = append(missingKonglo, MissingKonglo{Ticker: tk, Group: gname})
				continue
			}
			item := p
			item.Group = gname
			kongloStocks = append(kongloStocks, item)
			if !seenKonglo[tk] {
				seenKonglo[tk] = true
				uniqueKongloStocks = append(uniqueKongloStocks, item)
			}
		}
	}

	// 2. Konglo Overview
	kongloRising, kongloFalling, kongloFlat := 0, 0, 0
	var kongloMCSum float64
	for _, s := range uniqueKongloStocks {
		switch {
		case s.Pct > 0:
			kongloRising++
		case s.Pct < 0:
			kongloFalling++
		default:
			kongloFlat++
		}
		kongloMCSum += s.MCChg
	}
	kongloTotalMCChg := shared.PyRound(kongloMCSum, 2)

	// 3 & 4. Top 5 & Top 15
	byPctDesc := sortedCopy(uniqueKongloStocks, func(a, b ReportStockItem) bool { return a.Pct > b.Pct })
	byMCDesc := sortedCopy(uniqueKongloStocks, func(a, b ReportStockItem) bool { return a.MCChg > b.MCChg })
	byPctAsc := sortedCopy(uniqueKongloStocks, func(a, b ReportStockItem) bool { return a.Pct < b.Pct })
	byMCAsc := sortedCopy(uniqueKongloStocks, func(a, b ReportStockItem) bool { return a.MCChg < b.MCChg })

	// 5. Group Stats & Top 5 Konglo
	groupStats := make([]GroupStat, 0, len(orderedGroups))
	for _, gname := range orderedGroups {
		gs := make([]ReportStockItem, 0)
		for _, s := range kongloStocks {
			if s.Group == gname {
				gs = append(gs, s)
			}
		}
		if len(gs) == 0 {
			continue
		}
		var totMCChg, totVal, totPct, totMCClose, totMCPrev float64
		var totFreq int64
		for _, s := range gs {
			totMCChg += s.MCChg
			totVal += s.Val
			totFreq += s.Freq
			totPct += s.Pct
			totMCClose += (s.Shares * s.Close) / 1e12
			totMCPrev += (s.Shares * s.Prev) / 1e12
		}
		avgPct := shared.PyRound(totPct/float64(len(gs)), 2)
		mcPct := 0.0
		if totMCPrev > 0 {
			mcPct = shared.PyRound(totMCChg/totMCPrev*100.0, 2)
		}
		membersSorted := sortedCopy(gs, func(a, b ReportStockItem) bool { return a.Pct > b.Pct })
		top2 := make([]Top2Entry, 0, 2)
		for i := 0; i < len(membersSorted) && i < 2; i++ {
			top2 = append(top2, Top2Entry{Ticker: membersSorted[i].Ticker, Pct: membersSorted[i].Pct})
		}

		groupStats = append(groupStats, GroupStat{
			Group:        gname,
			TotalMembers: len(gs),
			MCChg:        shared.PyRound(totMCChg, 2),
			TotalMC:      shared.PyRound(totMCClose, 2),
			MCPrev:       shared.PyRound(totMCPrev, 2),
			MCPct:        mcPct,
			TotalVal:     shared.PyRound(totVal, 2),
			TotalFreq:    totFreq,
			AvgPct:       avgPct,
			Top2:         top2,
			Members:      membersSorted,
		})
	}

	top5Konglo := sortedGroupCopy(groupStats, func(a, b GroupStat) bool { return a.AvgPct > b.AvgPct })
	if len(top5Konglo) > 5 {
		top5Konglo = top5Konglo[:5]
	}

	// 6. Best & Worst Leaders
	bestMC := LeaderCard{Group: "-", Val: 0, Count: 0}
	worstMC := LeaderCard{Group: "-", Val: 0, Count: 0}
	bestPct := LeaderCard{Group: "-", Val: 0, Count: 0}
	worstPct := LeaderCard{Group: "-", Val: 0, Count: 0}
	if len(groupStats) > 0 {
		sBMC := sortedGroupCopy(groupStats, func(a, b GroupStat) bool { return a.MCChg > b.MCChg })[0]
		sWMC := sortedGroupCopy(groupStats, func(a, b GroupStat) bool { return a.MCChg < b.MCChg })[0]
		sBPct := sortedGroupCopy(groupStats, func(a, b GroupStat) bool { return a.AvgPct > b.AvgPct })[0]
		sWPct := sortedGroupCopy(groupStats, func(a, b GroupStat) bool { return a.AvgPct < b.AvgPct })[0]

		bestMC = LeaderCard{Group: sBMC.Group, Val: sBMC.MCChg, Count: sBMC.TotalMembers}
		worstMC = LeaderCard{Group: sWMC.Group, Val: sWMC.MCChg, Count: sWMC.TotalMembers}
		bestPct = LeaderCard{Group: sBPct.Group, Val: sBPct.AvgPct, Count: sBPct.TotalMembers}
		worstPct = LeaderCard{Group: sWPct.Group, Val: sWPct.AvgPct, Count: sWPct.TotalMembers}
	}

	// 7. Fundamental 50 Report
	fundaStocks := make([]ReportStockItem, 0, len(FundamentalTickers))
	missingFunda := make([]string, 0)
	for _, tk := range FundamentalTickers {
		p, ok := stockByTicker[tk]
		if !ok || p.Close <= 0 {
			missingFunda = append(missingFunda, tk)
			continue
		}
		fundaStocks = append(fundaStocks, p)
	}
	fundaRising, fundaFalling, fundaFlat := 0, 0, 0
	var fundaMCSum, fundaPrevMC float64
	for _, s := range fundaStocks {
		switch {
		case s.Pct > 0:
			fundaRising++
		case s.Pct < 0:
			fundaFalling++
		default:
			fundaFlat++
		}
		fundaMCSum += s.MCChg
		fundaPrevMC += (s.Shares * s.Prev) / 1e12
	}
	fundaTotMCChg := shared.PyRound(fundaMCSum, 2)
	fundaMCPct := 0.0
	if fundaPrevMC > 0 {
		fundaMCPct = shared.PyRound(fundaTotMCChg/fundaPrevMC*100.0, 2)
	}
	fPctUp := sortedCopy(fundaStocks, func(a, b ReportStockItem) bool { return a.Pct > b.Pct })
	fMCUp := sortedCopy(fundaStocks, func(a, b ReportStockItem) bool { return a.MCChg > b.MCChg })
	fPctDown := sortedCopy(fundaStocks, func(a, b ReportStockItem) bool { return a.Pct < b.Pct })
	fMCDown := sortedCopy(fundaStocks, func(a, b ReportStockItem) bool { return a.MCChg < b.MCChg })

	// 8. Overall Market Summary
	var totalTurnoverVal float64
	stocksUp, stocksDown, stocksFlat := 0, 0, 0
	for _, s := range activeStocks {
		totalTurnoverVal += s.Val
		switch {
		case s.Pct > 0:
			stocksUp++
		case s.Pct < 0:
			stocksDown++
		default:
			stocksFlat++
		}
	}

	allGroupStats := sortedGroupCopy(groupStats, func(a, b GroupStat) bool { return a.Group < b.Group })

	return &Report{
		Summary: ReportSummary{
			StocksUp:          stocksUp,
			StocksDown:        stocksDown,
			StocksFlat:        stocksFlat,
			TotalActiveStocks: len(activeStocks),
			TotalStocks:       len(stocks),
			TotalTurnoverM:    shared.PyRound(totalTurnoverVal/1e9, 2),
		},
		KongloCards: KongloCards{
			Overview: OverviewCard{
				TotalKongloStocks: len(uniqueKongloStocks),
				Rising:            kongloRising,
				Falling:           kongloFalling,
				Flat:              kongloFlat,
				TotalMCChg:        kongloTotalMCChg,
				MissingCount:      len(missingKonglo),
			},
			BestMC:   bestMC,
			BestPct:  bestPct,
			WorstMC:  worstMC,
			WorstPct: worstPct,
		},
		Top5Summary: Top5Summary{
			Top5PctUp:   takeSlice(byPctDesc, 5),
			Top5MCUp:    takeSlice(byMCDesc, 5),
			Top5PctDown: takeSlice(byPctAsc, 5),
			Top5MCDown:  takeSlice(byMCAsc, 5),
		},
		Top15Tables: Top15Tables{
			PctUp:   takeSlice(byPctDesc, 15),
			MCUp:    takeSlice(byMCDesc, 15),
			PctDown: takeSlice(byPctAsc, 15),
			MCDown:  takeSlice(byMCAsc, 15),
		},
		Top5Konglo:    top5Konglo,
		AllGroupStats: allGroupStats,
		FundamentalSummary: FundamentalSummary{
			TotalStocks:  len(fundaStocks) + len(missingFunda),
			MissingCount: len(missingFunda),
			Rising:       fundaRising,
			Falling:      fundaFalling,
			Flat:         fundaFlat,
			TotalMCChg:   fundaTotMCChg,
			MCPct:        fundaMCPct,
			Top5PctUp:    takeSlice(fPctUp, 5),
			Top5MCUp:     takeSlice(fMCUp, 5),
			Top5PctDown:  takeSlice(fPctDown, 5),
			Top5MCDown:   takeSlice(fMCDown, 5),
		},
		FundamentalReport: fPctUp,
		MissingKonglo:     missingKonglo,
		MissingFunda:      missingFunda,
	}
}

func floatOrZero(p *float64) float64 {
	if p == nil {
		return 0.0
	}
	return *p
}

func sortedCopy(in []ReportStockItem, less func(a, b ReportStockItem) bool) []ReportStockItem {
	cp := make([]ReportStockItem, len(in))
	copy(cp, in)
	sort.SliceStable(cp, func(i, j int) bool {
		return less(cp[i], cp[j])
	})
	return cp
}

func sortedGroupCopy(in []GroupStat, less func(a, b GroupStat) bool) []GroupStat {
	cp := make([]GroupStat, len(in))
	copy(cp, in)
	sort.SliceStable(cp, func(i, j int) bool {
		return less(cp[i], cp[j])
	})
	return cp
}

func takeSlice(in []ReportStockItem, n int) []ReportStockItem {
	if len(in) <= n {
		return in
	}
	return in[:n]
}
