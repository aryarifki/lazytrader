package ipohunter

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// BuildMetaData computes GET /meta from masters, distinct dates/tickers, and optional D0 date.
func BuildMetaData(masters []MasterRecord, dates []time.Time, withData map[string]struct{}, d0Date *time.Time) MetaResponse {
	pool := make([]MasterRecord, 0, len(masters))
	for _, m := range masters {
		if m.ListingDate == nil {
			continue
		}
		if len(withData) > 0 {
			if _, ok := withData[m.Ticker]; !ok {
				continue
			}
		}
		pool = append(pool, m)
	}
	sort.SliceStable(pool, func(i, j int) bool {
		di := DateOrdinal(*pool[i].ListingDate)
		dj := DateOrdinal(*pool[j].ListingDate)
		if di != dj {
			return di > dj
		}
		return pool[i].Ticker > pool[j].Ticker
	})
	limit := min(40, len(pool))
	quickPicks := make([]QuickPickItem, 0, limit)
	for _, m := range pool[:limit] {
		cname := m.Ticker
		if m.CompanyName != nil && *m.CompanyName != "" {
			cname = *m.CompanyName
		}
		quickPicks = append(quickPicks, QuickPickItem{
			Ticker:      m.Ticker,
			CompanyName: cname,
			ListingDate: FormatDateISO(*m.ListingDate),
			IPOPrice:    m.IPOPrice,
		})
	}

	uwCount := make(map[string]int, 128)
	yearsSet := make(map[int]struct{}, 16)
	tickers := make([]MetaTickerItem, 0, len(masters))
	for _, m := range masters {
		var ldStr *string
		if m.ListingDate != nil {
			yearsSet[m.ListingDate.UTC().Year()] = struct{}{}
			s := FormatDateISO(*m.ListingDate)
			ldStr = &s
		}
		for _, u := range m.Underwriters {
			uwCount[u.Broker]++
		}
		cname := ""
		if m.CompanyName != nil {
			cname = *m.CompanyName
		}
		tickers = append(tickers, MetaTickerItem{
			Ticker:      m.Ticker,
			CompanyName: cname,
			ListingDate: ldStr,
			IPOPrice:    m.IPOPrice,
		})
	}

	years := make([]int, 0, len(yearsSet))
	for y := range yearsSet {
		years = append(years, y)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(years)))

	uwKeys := make([]string, 0, len(uwCount))
	for k := range uwCount {
		uwKeys = append(uwKeys, k)
	}
	sort.Strings(uwKeys)
	uwrList := make([]UnderwriterCountItem, 0, len(uwKeys))
	for _, k := range uwKeys {
		uwrList = append(uwrList, UnderwriterCountItem{Broker: k, Count: uwCount[k]})
	}

	availDates := make([]string, 0, len(dates))
	for _, d := range dates {
		availDates = append(availDates, FormatDateISO(d))
	}
	var earliest, latest, d0Str *string
	if len(availDates) > 0 {
		earliest = &availDates[0]
		latest = &availDates[len(availDates)-1]
	}
	if d0Date != nil {
		s := FormatDateISO(*d0Date)
		d0Str = &s
	}

	return MetaResponse{
		TotalIPOs:      len(masters),
		D0Date:         d0Str,
		EarliestDate:   earliest,
		LatestDate:     latest,
		AvailableDates: availDates,
		QuickPicks:     quickPicks,
		Years:          years,
		Underwriters:   uwrList,
		Tickers:        tickers,
	}
}

func valBucketOK(bucket string, nilai *float64) bool {
	if bucket == "" {
		return true
	}
	if nilai == nil {
		return false
	}
	m := *nilai / 1e9
	switch bucket {
	case "lt100":
		return m < 100
	case "100_500":
		return 100 <= m && m < 500
	case "500_1000":
		return 500 <= m && m < 1000
	case "1_5t":
		return 1000 <= m && m < 5000
	case "gt5t":
		return m >= 5000
	default:
		return true
	}
}

func listSortValue(m MasterRecord, sortBy string) (numVal float64, strVal string, isStr bool, ok bool) {
	switch sortBy {
	case "kode":
		return 0, m.Ticker, true, true
	case "nama":
		if m.CompanyName == nil || *m.CompanyName == "" {
			return 0, "", true, false
		}
		return 0, *m.CompanyName, true, true
	case "papan":
		if m.Papan != nil {
			return 0, *m.Papan, true, true
		}
		return 0, "", true, true
	case "uw":
		if len(m.Underwriters) > 0 {
			return 0, m.Underwriters[0].Broker, true, true
		}
		return 0, "", true, true
	case "nilai":
		if m.NilaiIPO == nil {
			return 0, "", false, false
		}
		return *m.NilaiIPO, "", false, true
	case "price":
		return m.IPOPrice, "", false, true
	case "sid":
		if m.JumlahSID == nil {
			return 0, "", false, false
		}
		return float64(*m.JumlahSID), "", false, true
	case "penjatahan":
		if m.Penjatahan == nil {
			return 0, "", false, false
		}
		return *m.Penjatahan, "", false, true
	case "perf2w", "perf_2w":
		if m.Perf2W == nil {
			return 0, "", false, false
		}
		return *m.Perf2W, "", false, true
	case "perf1m", "perf_1m":
		if m.Perf1M == nil {
			return 0, "", false, false
		}
		return *m.Perf1M, "", false, true
	case "perf3m", "perf_3m":
		if m.Perf3M == nil {
			return 0, "", false, false
		}
		return *m.Perf3M, "", false, true
	default:
		if strings.HasPrefix(sortBy, "d") && len(sortBy) > 1 {
			if idx, err := strconv.Atoi(sortBy[1:]); err == nil {
				idx--
				if idx >= 0 && idx < len(m.DayReturns) && m.DayReturns[idx] != nil {
					return *m.DayReturns[idx], "", false, true
				}
				return 0, "", false, false
			}
		}
		if m.ListingDate == nil {
			return 0, "", false, false
		}
		return float64(DateOrdinal(*m.ListingDate)), "", false, true
	}
}

// BuildIPOListData filters, sorts (nulls last), and summarizes IPO List rows.
func BuildIPOListData(masters []MasterRecord, p ListParams) ListResponse {
	qUp := strings.ToUpper(strings.TrimSpace(p.Query))
	uwUp := strings.ToUpper(strings.TrimSpace(p.UW))

	filtered := make([]MasterRecord, 0, len(masters))
	for _, m := range masters {
		if uwUp != "" {
			matched := false
			for _, u := range m.Underwriters {
				if u.Broker == uwUp && (!p.LeadOnly || u.Role == "LEAD") {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		if !valBucketOK(p.ValBucket, m.NilaiIPO) {
			continue
		}
		if p.Year != nil && (m.ListingDate == nil || m.ListingDate.UTC().Year() != *p.Year) {
			continue
		}
		if qUp != "" {
			cname := ""
			if m.CompanyName != nil {
				cname = strings.ToUpper(*m.CompanyName)
			}
			if !strings.Contains(m.Ticker, qUp) && !strings.Contains(cname, qUp) {
				continue
			}
		}
		filtered = append(filtered, m)
	}

	nonNull := make([]MasterRecord, 0, len(filtered))
	nulls := make([]MasterRecord, 0)
	for _, m := range filtered {
		if _, _, _, ok := listSortValue(m, p.SortBy); ok {
			nonNull = append(nonNull, m)
		} else {
			nulls = append(nulls, m)
		}
	}
	desc := p.SortDir < 0
	sort.SliceStable(nonNull, func(i, j int) bool {
		ni, si, isStr, _ := listSortValue(nonNull[i], p.SortBy)
		nj, sj, _, _ := listSortValue(nonNull[j], p.SortBy)
		if isStr {
			if desc {
				return si > sj
			}
			return si < sj
		}
		if desc {
			return ni > nj
		}
		return ni < nj
	})
	ordered := append(nonNull, nulls...)

	// Kinerja summary across ordered slice
	colKeys := make([]string, 0, 13)
	for i := 1; i <= 10; i++ {
		colKeys = append(colKeys, fmt.Sprintf("d%d", i))
	}
	colKeys = append(colKeys, "perf2w", "perf1m", "perf3m")
	colVals := make(map[string][]float64, 13)
	for _, m := range ordered {
		for i := 0; i < 10 && i < len(m.DayReturns); i++ {
			if m.DayReturns[i] != nil {
				k := fmt.Sprintf("d%d", i+1)
				colVals[k] = append(colVals[k], *m.DayReturns[i])
			}
		}
		if m.Perf2W != nil {
			colVals["perf2w"] = append(colVals["perf2w"], *m.Perf2W)
		}
		if m.Perf1M != nil {
			colVals["perf1m"] = append(colVals["perf1m"], *m.Perf1M)
		}
		if m.Perf3M != nil {
			colVals["perf3m"] = append(colVals["perf3m"], *m.Perf3M)
		}
	}
	rataRata := make(map[string]*int64, 13)
	kemungkinanNaik := make(map[string]*int64, 13)
	for _, k := range colKeys {
		vals := colVals[k]
		if len(vals) == 0 {
			rataRata[k] = nil
			kemungkinanNaik[k] = nil
			continue
		}
		sum := 0.0
		pos := 0
		for _, v := range vals {
			sum += v
			if v > 0 {
				pos++
			}
		}
		avg := JSRound(sum / float64(len(vals)))
		up := JSRound(float64(pos) / float64(len(vals)) * 100.0)
		rataRata[k] = &avg
		kemungkinanNaik[k] = &up
	}

	rowsOut := make([]ListRow, 0, len(ordered))
	for _, m := range ordered {
		chips := make([]ListUnderwriterChip, 0, len(m.Underwriters))
		for _, u := range m.Underwriters {
			chips = append(chips, ListUnderwriterChip{
				Broker:   u.Broker,
				Role:     u.Role,
				Category: CategoryOf(u.Broker, m.Ticker, nil),
			})
		}
		cname := ""
		if m.CompanyName != nil {
			cname = *m.CompanyName
		}
		var ldStr *string
		if m.ListingDate != nil {
			s := FormatDateISO(*m.ListingDate)
			ldStr = &s
		}
		var nilaiM *float64
		if m.NilaiIPO != nil {
			v := *m.NilaiIPO / 1e9
			nilaiM = &v
		}
		dayRet := m.DayReturns
		if dayRet == nil {
			dayRet = []*float64{}
		}
		rowsOut = append(rowsOut, ListRow{
			Ticker:       m.Ticker,
			CompanyName:  cname,
			ListingDate:  ldStr,
			IPOPrice:     m.IPOPrice,
			NilaiIPO:     m.NilaiIPO,
			NilaiMiliar:  nilaiM,
			Papan:        m.Papan,
			Underwriters: chips,
			JumlahSID:    m.JumlahSID,
			Penjatahan:   m.Penjatahan,
			DayReturns:   dayRet,
			Perf2W:       m.Perf2W,
			Perf1M:       m.Perf1M,
			Perf3M:       m.Perf3M,
		})
	}

	return ListResponse{
		Total:    len(rowsOut),
		TotalAll: len(masters),
		Rows:     rowsOut,
		KinerjaSummary: KinerjaSummary{
			RataRata:        rataRata,
			KemungkinanNaik: kemungkinanNaik,
		},
	}
}

// ComputeScreenerRow computes all 11 screener metrics for a single ticker in Periodic or After-IPO mode.
// When after==true, afterDays is the ticker's daily rows strictly after D0 (and <= hiDate) plus D0 via buildAfterDays.
// When after==false, baseVols provides up to 20 positive volumes strictly before loDate.
func ComputeScreenerRow(
	kode string,
	master *MasterRecord,
	d0Row *DayRow,
	periodDays []DayRow,
	baseVols []float64,
	after bool,
	hiDate *time.Time,
	holidays map[int]struct{},
) *ScreenerRow {
	var useDays []DayRow
	if after {
		useDays = BuildAfterDays(kode, master, d0Row, periodDays, hiDate, holidays).Days
	} else {
		useDays = periodDays
	}

	agg := BuildAggregation(kode, useDays)
	if agg == nil {
		return nil
	}

	tv := agg.TotalVolume
	uwSet := UWSetOf(master)
	nets := map[string]float64{
		"uw": 0, "retail": 0, "semi": 0, "institusi": 0, "nonRitel": 0, "repo": 0,
	}
	for b, t := range agg.BrokerTotals {
		c := CategoryOf(b, kode, uwSet)
		_, isRepo := RepoBrokers[b]
		if c != "uw" && (isRepo || c == "zombie") {
			c = "repo"
		}
		if _, exists := nets[c]; !exists {
			c = "repo"
		}
		nets[c] += t.Buy - t.Sell
	}
	pc := func(v float64) *float64 {
		if tv == 0 {
			return nil
		}
		res := (v / tv) * 100.0
		return &res
	}

	sv := 0.0
	nd := 0
	for i := range useDays {
		d := &useDays[i]
		if d.Volume == nil {
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
		sv += float64(*d.Volume) * 100.0 * ap
		dc := d.DayCount
		if dc <= 0 {
			dc = 1
		}
		nd += dc
	}

	var volRatio *float64
	if after {
		var lastDay *DayRow
		for i := range useDays {
			d := &useDays[i]
			if !d.IsD0 && d.Volume != nil && *d.Volume > 0 {
				lastDay = d
			}
		}
		if lastDay != nil {
			pv := 0.0
			pd := 0
			for i := range useDays {
				d := &useDays[i]
				if d == lastDay || !d.Date.Before(lastDay.Date) {
					continue
				}
				if d.Volume != nil && *d.Volume > 0 {
					pv += float64(*d.Volume)
					dc := d.DayCount
					if dc <= 0 {
						dc = 1
					}
					pd += dc
				}
			}
			if pd > 0 && pv > 0 {
				vr := float64(*lastDay.Volume) / (pv / float64(pd))
				volRatio = &vr
			}
		}
	} else if len(baseVols) > 0 {
		curSum := 0.0
		curN := 0
		for i := range useDays {
			d := &useDays[i]
			if d.Volume != nil && *d.Volume > 0 {
				curSum += float64(*d.Volume)
				curN++
			}
		}
		if curN > 0 {
			baseSum := 0.0
			for _, v := range baseVols {
				baseSum += v
			}
			baseMean := baseSum / float64(len(baseVols))
			if baseMean > 0 {
				vr := (curSum / float64(curN)) / baseMean
				volRatio = &vr
			}
		}
	}

	last := LastCloseIn(useDays)
	vw := AfterVWAP(useDays)
	var perfIpo *float64
	if last != nil && master.IPOPrice != 0 {
		v := (*last/master.IPOPrice - 1.0) * 100.0
		perfIpo = &v
	}
	var perfAvg *float64
	if last != nil && vw != nil && *vw != 0 {
		v := (*last/(*vw) - 1.0) * 100.0
		perfAvg = &v
	}
	chg := agg.PriceChange
	if after {
		chg = perfIpo
	}

	var uwPtr *float64
	if len(uwSet) > 0 {
		uwPtr = pc(nets["uw"])
	}
	var netRetail *float64
	if tv != 0 {
		v := -nets["retail"] / tv * 100.0
		netRetail = &v
	}
	var ratrans *float64
	if nd != 0 {
		v := sv / float64(nd) / 1e9
		ratrans = &v
	}
	cname := ""
	if master.CompanyName != nil {
		cname = *master.CompanyName
	}
	var ldStr *string
	if master.ListingDate != nil {
		s := FormatDateISO(*master.ListingDate)
		ldStr = &s
	}

	return &ScreenerRow{
		Kode:        kode,
		CompanyName: cname,
		ListingDate: ldStr,
		UW:          uwPtr,
		NetRetail:   netRetail,
		NetSemi:     pc(nets["semi"]),
		NetInst:     pc(nets["institusi"]),
		NetNonRitel: pc(nets["nonRitel"]),
		NetRepo:     pc(nets["repo"]),
		PerfIpo:     perfIpo,
		PerfAvg:     perfAvg,
		Chg:         chg,
		Ratrans:     ratrans,
		VolRatio:    volRatio,
	}
}

func screenerMetricVal(r ScreenerRow, field string) *float64 {
	switch field {
	case "uw":
		return r.UW
	case "netRetail":
		return r.NetRetail
	case "netSemi":
		return r.NetSemi
	case "netInst":
		return r.NetInst
	case "netNonRitel":
		return r.NetNonRitel
	case "netRepo":
		return r.NetRepo
	case "perfIpo":
		return r.PerfIpo
	case "perfAvg":
		return r.PerfAvg
	case "chg":
		return r.Chg
	case "ratrans":
		return r.Ratrans
	case "volRatio":
		return r.VolRatio
	default:
		return nil
	}
}

func cmpScreenerVal(op string, a, b float64) bool {
	switch op {
	case ">":
		return a > b
	case ">=":
		return a >= b
	case "<":
		return a < b
	default:
		return a <= b
	}
}

// FilterAndSortScreenerRows applies chg, custom conditions, preset mapping, and stable nulls-last sort.
func FilterAndSortScreenerRows(
	candidates []ScreenerRow,
	chg string,
	conds []ScreenerCondition,
	metric string,
	sortKey *string,
	sortDir *int,
) (ordered []ScreenerRow, activeCol, sk string, sdir int) {
	rows := make([]ScreenerRow, 0, len(candidates))
	for _, r := range candidates {
		if chg == "up" && !(r.Chg != nil && *r.Chg > 0) {
			continue
		}
		if chg == "down" && !(r.Chg != nil && *r.Chg < 0) {
			continue
		}
		passed := true
		for _, c := range conds {
			trimmed := strings.TrimSpace(c.V)
			if trimmed == "" {
				continue
			}
			targetV, err := strconv.ParseFloat(strings.ReplaceAll(trimmed, ",", "."), 64)
			if err != nil {
				continue
			}
			val := screenerMetricVal(r, c.K)
			op := c.Op
			if op == "" {
				op = ">"
			}
			if val == nil || !cmpScreenerVal(op, *val, targetV) {
				passed = false
				break
			}
		}
		if passed {
			rows = append(rows, r)
		}
	}

	activeCol = metric
	defaultDir := -1
	if preset, ok := ScrPresetDefaults[metric]; ok {
		activeCol = preset.Col
		defaultDir = preset.Dir
	}
	sk = activeCol
	if sortKey != nil && *sortKey != "" {
		sk = *sortKey
	}
	sdir = defaultDir
	if sortDir != nil {
		sdir = *sortDir
	}

	nonNull := make([]ScreenerRow, 0, len(rows))
	nulls := make([]ScreenerRow, 0)
	for _, r := range rows {
		if screenerMetricVal(r, sk) != nil {
			nonNull = append(nonNull, r)
		} else {
			nulls = append(nulls, r)
		}
	}
	desc := sdir < 0
	sort.SliceStable(nonNull, func(i, j int) bool {
		vi := *screenerMetricVal(nonNull[i], sk)
		vj := *screenerMetricVal(nonNull[j], sk)
		if desc {
			return vi > vj
		}
		return vi < vj
	})

	ordered = append(nonNull, nulls...)
	return ordered, activeCol, sk, sdir
}
