package ownership

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var errCTValidation = errors.New("ct validation error")

type CTExplorerParams struct {
	Mode            string
	Category        string
	Condition       string
	Ranges          *string
	MktcapTier      string
	TransaksiRange  string
	InvestorOrigin  string
	InvestorTypes   *string
	RetailRanges    *string
	PortfolioRanges *string
	Sort            string
	Order           string
	Limit           int
	Offset          int
}

type explorerSortSpec struct {
	ID     string
	Label  string
	Getter func(row map[string]any) any
}

func baseExplorerGetters() ([]string, map[string]func(map[string]any) any) {
	order := []string{
		"ticker", "name", "net", "sid", "mtd", "count",
		"market_cap", "best", "matches", "pct", "gl",
		"ticker_count", "latest", "sid_amount", "sid_pct",
		"sid_value", "sid_change", "abs_value",
	}
	getters := map[string]func(map[string]any) any{
		"ticker":     func(r map[string]any) any { return r["ticker"] },
		"name":       func(r map[string]any) any { return r["name"] },
		"net": func(r map[string]any) any {
			if v, ok := r["net_pct"]; ok {
				return v
			}
			return r["net_value_m"]
		},
		"sid":        func(r map[string]any) any { return r["net_pct"] },
		"mtd":        func(r map[string]any) any { return r["mtd_pct"] },
		"count":      func(r map[string]any) any { return r["count"] },
		"market_cap": func(r map[string]any) any { return r["market_cap_b"] },
		"best":       func(r map[string]any) any { return r["best"] },
		"matches":    func(r map[string]any) any { return r["match_count"] },
		"pct":        func(r map[string]any) any { return r["pct"] },
		"gl":         func(r map[string]any) any { return r["gl_rp_m"] },
		"ticker_count": func(r map[string]any) any {
			return r["ticker_count"]
		},
		"latest": func(r map[string]any) any { return r["latest"] },
		"sid_amount": func(r map[string]any) any {
			return math.Abs(asFloatAny(r["sid_change"]))
		},
		"sid_pct": func(r map[string]any) any {
			return math.Abs(asFloatAny(r["net_pct"]))
		},
		"sid_value": func(r map[string]any) any { return r["sid"] },
		"sid_change": func(r map[string]any) any {
			return math.Abs(asFloatAny(r["sid_change"]))
		},
		"abs_value": func(r map[string]any) any { return r["abs_value_m"] },
	}
	return order, getters
}

func asFloatAny(v any) float64 {
	if v == nil {
		return 0.0
	}
	switch x := v.(type) {
	case float64:
		return x
	case *float64:
		if x == nil {
			return 0.0
		}
		return *x
	case int:
		return float64(x)
	case *int:
		if x == nil {
			return 0.0
		}
		return float64(*x)
	case int64:
		return float64(x)
	default:
		return 0.0
	}
}

func isNilAny(v any) bool {
	if v == nil {
		return true
	}
	switch x := v.(type) {
	case *float64:
		return x == nil
	case *int:
		return x == nil
	case *string:
		return x == nil
	default:
		return false
	}
}

func derefAny(v any) any {
	switch x := v.(type) {
	case *float64:
		if x == nil {
			return nil
		}
		return *x
	case *int:
		if x == nil {
			return nil
		}
		return *x
	case *string:
		if x == nil {
			return nil
		}
		return *x
	default:
		return v
	}
}

func titleCaseSpaces(id string) string {
	words := strings.Split(strings.ReplaceAll(id, "_", " "), " ")
	for i, w := range words {
		if w == "" {
			continue
		}
		words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
	}
	return strings.Join(words, " ")
}

func containsStrSlice(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}

func parseRangeSet(raw *string, allowed []string, defaultVals []string) (map[string]bool, error) {
	out := map[string]bool{}
	if raw == nil || strings.TrimSpace(*raw) == "" {
		for _, v := range defaultVals {
			out[v] = true
		}
		return out, nil
	}
	allowedSet := map[string]bool{}
	for _, v := range allowed {
		allowedSet[v] = true
	}
	invalid := []string{}
	for _, part := range strings.Split(*raw, ",") {
		item := strings.TrimSpace(part)
		if item == "" {
			continue
		}
		if !allowedSet[item] {
			invalid = append(invalid, item)
		} else {
			out[item] = true
		}
	}
	if len(invalid) > 0 {
		sort.Strings(invalid)
		return nil, fmt.Errorf("%w: Nilai filter tidak dikenal: %s", errCTValidation, strings.Join(invalid, ", "))
	}
	return out, nil
}

func validateCTCondition(mode, category, condition string) ([]string, error) {
	var valid []string
	switch category {
	case "kepemilikan":
		if mode == "stocks" {
			valid = []string{"bertambah", "berkurang", "sidNaik", "sidTurun"}
		} else {
			valid = []string{"bertambah", "berkurang"}
		}
	case "transaksi":
		valid = []string{"bertambah", "berkurang"}
	case "advanced":
		if mode == "stocks" {
			valid = append([]string(nil), stockAdvancedOrder...)
		} else {
			valid = append([]string(nil), investorAdvancedOrder...)
		}
	default:
		return nil, fmt.Errorf("%w: Kategori explorer tidak dikenal", errCTValidation)
	}
	if !containsStrSlice(valid, condition) {
		return nil, fmt.Errorf("%w: Kondisi '%s' tidak valid untuk mode/kategori ini", errCTValidation, condition)
	}
	return valid, nil
}

func sortCTRows(rows []map[string]any, sortKey, order, defKey string, getters map[string]func(map[string]any) any) ([]map[string]any, error) {
	key := strings.TrimSpace(sortKey)
	if key == "" {
		key = defKey
	}
	if key == "default" {
		key = defKey
	}
	getter, ok := getters[key]
	if !ok {
		return nil, fmt.Errorf("%w: Urutan '%s' tidak dikenal", errCTValidation, sortKey)
	}
	reverse := strings.ToLower(strings.TrimSpace(order)) != "asc"

	nonNull := make([]map[string]any, 0, len(rows))
	nulls := make([]map[string]any, 0)
	for _, row := range rows {
		if isNilAny(getter(row)) {
			nulls = append(nulls, row)
		} else {
			nonNull = append(nonNull, row)
		}
	}

	sort.SliceStable(nonNull, func(i, j int) bool {
		va := derefAny(getter(nonNull[i]))
		vb := derefAny(getter(nonNull[j]))
		var less bool
		switch a := va.(type) {
		case string:
			b, _ := vb.(string)
			less = strings.ToLower(a) < strings.ToLower(b)
		case int:
			switch b := vb.(type) {
			case int:
				less = a < b
			case float64:
				less = float64(a) < b
			default:
				less = fmt.Sprint(va) < fmt.Sprint(vb)
			}
		case float64:
			switch b := vb.(type) {
			case float64:
				less = a < b
			case int:
				less = a < float64(b)
			default:
				less = fmt.Sprint(va) < fmt.Sprint(vb)
			}
		default:
			less = fmt.Sprint(va) < fmt.Sprint(vb)
		}
		if reverse {
			bStr := false
			switch a := va.(type) {
			case string:
				b, _ := vb.(string)
				bStr = strings.ToLower(a) > strings.ToLower(b)
			case int:
				switch b := vb.(type) {
				case int:
					bStr = a > b
				case float64:
					bStr = float64(a) > b
				default:
					bStr = fmt.Sprint(va) > fmt.Sprint(vb)
				}
			case float64:
				switch b := vb.(type) {
				case float64:
					bStr = a > b
				case int:
					bStr = a > float64(b)
				default:
					bStr = fmt.Sprint(va) > fmt.Sprint(vb)
				}
			default:
				bStr = fmt.Sprint(va) > fmt.Sprint(vb)
			}
			return bStr
		}
		return less
	})

	return append(nonNull, nulls...), nil
}

func ctStockCommon(entry *ctContextTicker, c *ctContext) map[string]any {
	tk := entry.Ticker
	latest := c.OHLC[tk]
	price := entry.Price
	cl := latest.Close
	var mtd *float64
	if price != nil && *price != 0 && cl != nil {
		v := (*cl - *price) / *price * 100.0
		mtd = &v
	}
	issuer := entry.Issuer
	if issuer == "" {
		issuer = tk
	}
	badges := entry.Badges
	if badges == nil {
		badges = []string{}
	}
	return map[string]any{
		"ticker":       tk,
		"name":         issuer,
		"issuer":       issuer,
		"market_cap_b": roundEvenPtr(entry.MktcapB, 2),
		"price":        roundEvenPtr(price, 4),
		"close":        roundEvenPtr(cl, 4),
		"mtd_pct":      roundEvenPtr(mtd, 4),
		"badges":       badges,
		"badge":        entry.Badge,
		"conglomerate": entry.Konglo,
		"konglo":       entry.Konglo,
	}
}

func ctStockSIDRow(entry *ctContextTicker, c *ctContext, condition string) map[string]any {
	if len(c.SIDPeriods) < 2 {
		return nil
	}
	prevKey, lastKey := c.SIDPeriods[len(c.SIDPeriods)-2], c.SIDPeriods[len(c.SIDPeriods)-1]
	prev, okPrev := entry.SID[prevKey]
	last, okLast := entry.SID[lastKey]
	if !okPrev || !okLast || prev <= 0 {
		return nil
	}
	net := (float64(last)/float64(prev) - 1.0) * 100.0
	if net == 0 || (net > 0) != (condition == "sidNaik") {
		return nil
	}
	row := ctStockCommon(entry, c)
	row["net_pct"] = roundEven(net, 4)
	row["sid_change"] = last - prev
	row["sid"] = last
	row["sid_prev"] = prev
	row["period_prev"] = prevKey
	row["period_last"] = lastKey
	return row
}

func ctStockOwnershipRow(entry *ctContextTicker, c *ctContext, condition string) map[string]any {
	net := 0.0
	for _, r := range entry.Rows {
		net += rowDeltaPct(r.PctMei, r.PctJun, r.Status)
	}
	if net == 0 || (net > 0) != (condition == "bertambah") {
		return nil
	}
	row := ctStockCommon(entry, c)
	row["net_pct"] = roundEven(net, 4)
	row["holder_count"] = len(entry.Rows)
	return row
}

func ctStockTransactionRow(entry *ctContextTicker, c *ctContext, condition string) map[string]any {
	txs := c.T5ByTicker[entry.Ticker]
	net := 0.0
	for _, tx := range txs {
		act := strings.ToUpper(strings.TrimSpace(tx.Action))
		if act == "BUY" {
			net += derefFloat(tx.ValueM)
		} else if act == "SELL" {
			net -= derefFloat(tx.ValueM)
		}
	}
	if net == 0 || (net > 0) != (condition == "bertambah") {
		return nil
	}
	row := ctStockCommon(entry, c)
	row["net_value_m"] = roundEven(net, 4)
	row["count"] = len(txs)
	return row
}

func ctPctHistoryDelta(entry *ctContextTicker, metric string, monthKeys []string) (delta, latest float64, ok bool) {
	if len(monthKeys) < 2 {
		return 0, 0, false
	}
	lastKey, prevKey := monthKeys[len(monthKeys)-1], monthKeys[len(monthKeys)-2]
	lp, okL := entry.PctHistory[lastKey]
	pp, okP := entry.PctHistory[prevKey]
	if !okL || !okP {
		return 0, 0, false
	}
	var lPtr, pPtr *float64
	switch metric {
	case "pct_ritel":
		lPtr, pPtr = lp.PctRitel, pp.PctRitel
	case "pct_asing_lt1":
		lPtr, pPtr = lp.PctAsingLt1, pp.PctAsingLt1
	case "pct_asing":
		lPtr, pPtr = lp.PctAsing, pp.PctAsing
	}
	if lPtr == nil || pPtr == nil {
		return 0, 0, false
	}
	return *lPtr - *pPtr, *lPtr, true
}

func ctPctHistoryTrend3(entry *ctContextTicker, metric string, monthKeys []string) (rising, falling bool, delta, latest float64, ok bool) {
	if len(monthKeys) < 4 {
		return false, false, 0, 0, false
	}
	keys := monthKeys[len(monthKeys)-4:]
	vals := [4]float64{}
	for i, k := range keys {
		p, exists := entry.PctHistory[k]
		if !exists {
			return false, false, 0, 0, false
		}
		var ptr *float64
		if metric == "pct_asing" {
			ptr = p.PctAsing
		} else {
			ptr = p.PctRitel
		}
		if ptr == nil {
			return false, false, 0, 0, false
		}
		vals[i] = *ptr
	}
	return vals[0] < vals[1] && vals[1] < vals[2] && vals[2] < vals[3],
		vals[0] > vals[1] && vals[1] > vals[2] && vals[2] > vals[3],
		vals[3] - vals[0], vals[3], true
}

func ctSIDTrend3(sid map[string]int, periods []string) (rising, falling bool, deltaPct float64, diff, latest int, ok bool) {
	if len(periods) < 4 {
		return false, false, 0, 0, 0, false
	}
	keys := periods[len(periods)-4:]
	vals := [4]int{}
	for i, k := range keys {
		v, exists := sid[k]
		if !exists || v <= 0 {
			return false, false, 0, 0, 0, false
		}
		vals[i] = v
	}
	return vals[0] < vals[1] && vals[1] < vals[2] && vals[2] < vals[3],
		vals[0] > vals[1] && vals[1] > vals[2] && vals[2] > vals[3],
		(float64(vals[3])/float64(vals[0]) - 1.0) * 100.0,
		vals[3] - vals[0], vals[3], true
}

func ctStockAdvancedRow(key string, entry *ctContextTicker, c *ctContext) map[string]any {
	tk := entry.Ticker
	rows := entry.Rows
	txs := c.T5ByTicker[tk]

	t5InDays := func(maxDays int, names map[string]bool) (*float64, int) {
		net := 0.0
		cnt := 0
		found := false
		for _, tx := range txs {
			if names != nil && !names[strings.ToLower(strings.TrimSpace(tx.Name))] {
				continue
			}
			d, ok := t5Date(tx, c.Year)
			if !ok {
				continue
			}
			diffDays := int(c.T5AnchorDate.Sub(d).Hours() / 24)
			if diffDays < 0 || diffDays > maxDays {
				continue
			}
			found = true
			cnt++
			if tx.Action == "BUY" {
				net += derefFloat(tx.ValueM)
			} else if tx.Action == "SELL" {
				net -= derefFloat(tx.ValueM)
			}
		}
		if !found {
			return nil, cnt
		}
		return &net, cnt
	}

	t5CurrentMonth := func() (*float64, int) {
		net := 0.0
		cnt := 0
		found := false
		for _, tx := range txs {
			d, ok := t5Date(tx, c.Year)
			if !ok || int(d.Month()) != c.GLMonthNum || d.Year() != c.GLYear {
				continue
			}
			found = true
			cnt++
			if tx.Action == "BUY" {
				net += derefFloat(tx.ValueM)
			} else if tx.Action == "SELL" {
				net -= derefFloat(tx.ValueM)
			}
		}
		if !found {
			return nil, cnt
		}
		return &net, cnt
	}

	switch key {
	case "naikPorsi":
		net := 0.0
		for _, r := range rows {
			net += rowDeltaPct(r.PctMei, r.PctJun, r.Status)
		}
		if net <= 0.1 {
			return nil
		}
		row := ctStockCommon(entry, c)
		row["pct_change"] = roundEven(net, 4)
		return row

	case "aktifBeli":
		type invGroup struct {
			name    string
			count   int
			net     float64
			last    time.Time
			hasLast bool
		}
		byInv := map[string]*invGroup{}
		invOrder := []string{}
		for _, tx := range txs {
			d, ok := t5Date(tx, c.Year)
			if !ok || int(d.Month()) != c.GLMonthNum || d.Year() != c.GLYear {
				continue
			}
			nm := strings.TrimSpace(tx.Name)
			if nm == "" {
				continue
			}
			g, ok := byInv[nm]
			if !ok {
				g = &invGroup{name: nm}
				byInv[nm] = g
				invOrder = append(invOrder, nm)
			}
			g.count++
			if tx.Action == "BUY" {
				g.net += derefFloat(tx.ValueM)
			} else if tx.Action == "SELL" {
				g.net -= derefFloat(tx.ValueM)
			}
			if !g.hasLast || d.After(g.last) {
				g.last = d
				g.hasLast = true
			}
		}
		var best *invGroup
		for _, nm := range invOrder {
			g := byInv[nm]
			if g.count <= 5 || !g.hasLast {
				continue
			}
			diffDays := int(c.T5AnchorDate.Sub(g.last).Hours() / 24)
			if diffDays < 0 || diffDays > 5 || g.net <= 0 {
				continue
			}
			if best == nil || g.net > best.net {
				best = g
			}
		}
		if best == nil {
			return nil
		}
		row := ctStockCommon(entry, c)
		row["count"] = best.count
		row["net_value_m"] = roundEven(best.net, 4)
		row["holder_name"] = best.name
		return row

	case "buyback":
		normIssuer := normalizeCompany(entry.Issuer)
		treasRows := []ctContextHolder{}
		for _, r := range rows {
			if normalizeCompany(r.Name) == normIssuer {
				treasRows = append(treasRows, r)
			}
		}
		if len(treasRows) == 0 {
			return nil
		}
		pctChg := 0.0
		names := map[string]bool{}
		for _, r := range treasRows {
			pctChg += rowDeltaPct(r.PctMei, r.PctJun, r.Status)
			names[strings.ToLower(strings.TrimSpace(r.Name))] = true
		}
		net30, _ := t5InDays(30, names)
		if !(pctChg > 0) && !(net30 != nil && *net30 > 0) {
			return nil
		}
		row := ctStockCommon(entry, c)
		row["pct_change"] = roundEven(pctChg, 4)
		row["net_value_m"] = roundEven(derefFloat(net30), 4)
		return row

	case "mayoritas":
		candidates := []ctContextHolder{}
		for _, r := range rows {
			if strings.ToLower(strings.TrimSpace(r.Status)) != "new" {
				candidates = append(candidates, r)
			}
		}
		if len(candidates) == 0 {
			return nil
		}
		largest := candidates[0]
		for _, r := range candidates[1:] {
			if derefFloat(r.PctJun) > derefFloat(largest.PctJun) {
				largest = r
			}
		}
		if derefFloat(largest.PctJun) < 10 {
			return nil
		}
		pctChg := rowDeltaPct(largest.PctMei, largest.PctJun, largest.Status)
		if pctChg <= 0 {
			return nil
		}
		net30, _ := t5InDays(30, map[string]bool{strings.ToLower(strings.TrimSpace(largest.Name)): true})
		row := ctStockCommon(entry, c)
		row["pct_change"] = roundEven(pctChg, 4)
		row["net_value_m"] = roundEven(derefFloat(net30), 4)
		row["holder_name"] = strings.TrimSpace(largest.Name)
		return row

	case "baru":
		newRows := []ctContextHolder{}
		for _, r := range rows {
			if strings.ToLower(strings.TrimSpace(r.Status)) == "new" {
				newRows = append(newRows, r)
			}
		}
		if len(newRows) == 0 {
			return nil
		}
		pctChg := 0.0
		for _, r := range newRows {
			pctChg += rowDeltaPct(r.PctMei, r.PctJun, r.Status)
		}
		row := ctStockCommon(entry, c)
		row["pct_change"] = roundEven(pctChg, 4)
		return row

	case "retailTurun", "hiddenNaik", "foreignNaik":
		metric := "pct_asing"
		if key == "retailTurun" {
			metric = "pct_ritel"
		} else if key == "hiddenNaik" {
			metric = "pct_asing_lt1"
		}
		delta, latest, ok := ctPctHistoryDelta(entry, metric, c.MonthKeys)
		if !ok {
			return nil
		}
		if key == "retailTurun" && delta > -0.1 {
			return nil
		}
		if (key == "hiddenNaik" || key == "foreignNaik") && delta < 0.1 {
			return nil
		}
		row := ctStockCommon(entry, c)
		row["pct_change"] = roundEven(delta, 4)
		row["latest"] = roundEven(latest, 4)
		return row

	case "masukFlat":
		pctChg := 0.0
		for _, r := range rows {
			pctChg += rowDeltaPct(r.PctMei, r.PctJun, r.Status)
		}
		if pctChg <= 0.1 {
			return nil
		}
		netMonth, _ := t5CurrentMonth()
		if netMonth == nil || *netMonth <= 0 {
			return nil
		}
		row := ctStockCommon(entry, c)
		if mtdPtr, _ := row["mtd_pct"].(*float64); mtdPtr != nil && *mtdPtr >= 2 {
			return nil
		}
		row["pct_change"] = roundEven(pctChg, 4)
		row["net_value_m"] = roundEven(*netMonth, 4)
		return row

	case "foreign3Bulan", "retail3Bulan":
		metric := "pct_ritel"
		if key == "foreign3Bulan" {
			metric = "pct_asing"
		}
		rising, falling, delta, latest, ok := ctPctHistoryTrend3(entry, metric, c.MonthKeys)
		if !ok {
			return nil
		}
		if key == "foreign3Bulan" && !rising {
			return nil
		}
		if key == "retail3Bulan" && !falling {
			return nil
		}
		row := ctStockCommon(entry, c)
		row["pct_change"] = roundEven(delta, 4)
		row["latest"] = roundEven(latest, 4)
		return row

	case "sid3Turun", "sid3Naik":
		rising, falling, deltaPct, diff, latest, ok := ctSIDTrend3(entry.SID, c.SIDPeriods)
		if !ok {
			return nil
		}
		if key == "sid3Turun" && !falling {
			return nil
		}
		if key == "sid3Naik" && !rising {
			return nil
		}
		row := ctStockCommon(entry, c)
		row["pct_change"] = roundEven(deltaPct, 4)
		row["sid_change"] = diff
		row["latest"] = latest
		return row

	case "retailRendah":
		var latestPtr *float64
		if len(c.MonthKeys) > 0 {
			if p, exists := entry.PctHistory[c.MonthKeys[len(c.MonthKeys)-1]]; exists {
				latestPtr = p.PctRitel
			}
		}
		if latestPtr == nil || *latestPtr >= 5 {
			return nil
		}
		if derefFloat(entry.MktcapB) <= 0 {
			return nil
		}
		row := ctStockCommon(entry, c)
		row["latest"] = roundEvenPtr(latestPtr, 4)
		return row
	}
	return nil
}

func ctInvestorOwnershipRows(c *ctContext, condition string, ranges map[string]bool, origin string, selectedTypes map[string]bool) []map[string]any {
	results := []map[string]any{}
	for _, name := range c.InvestorOrder {
		holdings := c.InvestorHoldings[name]
		lf := investorLFHoldings(holdings)
		if !passInvestorType(lf, origin, selectedTypes) {
			continue
		}
		matches := []map[string]any{}
		for _, r := range holdings {
			delta := rowDeltaPct(r.PctMei, r.PctJun, r.Status)
			if delta == 0 || (delta > 0) != (condition == "bertambah") {
				continue
			}
			if !ranges[bucketPct(math.Abs(delta))] {
				continue
			}
			var iss any
			if r.RawIssuer != nil {
				iss = *r.RawIssuer
			}
			matches = append(matches, map[string]any{
				"ticker": strings.ToUpper(r.Ticker),
				"issuer": iss,
				"delta":  roundEven(delta, 4),
			})
		}
		if len(matches) == 0 {
			continue
		}
		sort.SliceStable(matches, func(i, j int) bool {
			return math.Abs(asFloatAny(matches[i]["delta"])) > math.Abs(asFloatAny(matches[j]["delta"]))
		})
		results = append(results, map[string]any{
			"name":         name,
			"lf":           lf,
			"lf_label":     lfLabel(lf),
			"ticker_count": len(holdings),
			"matches":      matches,
			"match_count":  len(matches),
			"best":         matches[0]["delta"],
		})
	}
	sort.SliceStable(results, func(i, j int) bool {
		ai, bi := math.Abs(asFloatAny(results[i]["best"])), math.Abs(asFloatAny(results[j]["best"]))
		if ai != bi {
			return ai > bi
		}
		ci, _ := results[i]["match_count"].(int)
		cj, _ := results[j]["match_count"].(int)
		return ci > cj
	})
	return results
}

func ctInvestorTransactionRows(c *ctContext, condition, valueRange, origin string, selectedTypes map[string]bool) []map[string]any {
	results := []map[string]any{}
	for _, name := range c.T5NameOrder {
		txs := c.T5ByName[name]
		holdings := c.InvestorHoldings[name]
		lf := investorLFHoldings(holdings)
		if !passInvestorType(lf, origin, selectedTypes) {
			continue
		}
		net := 0.0
		tSet := map[string]bool{}
		for _, tx := range txs {
			if tx.Action == "BUY" {
				net += derefFloat(tx.ValueM)
			} else if tx.Action == "SELL" {
				net -= derefFloat(tx.ValueM)
			}
			tSet[strings.ToUpper(tx.Ticker)] = true
		}
		if net == 0 || (net > 0) != (condition == "bertambah") {
			continue
		}
		if bucketValue(math.Abs(net)) != valueRange {
			continue
		}
		tickers := make([]string, 0, len(tSet))
		for tk := range tSet {
			tickers = append(tickers, tk)
		}
		sort.Strings(tickers)
		results = append(results, map[string]any{
			"name":         name,
			"lf":           lf,
			"lf_label":     lfLabel(lf),
			"net_value_m":  roundEven(net, 4),
			"count":        len(txs),
			"ticker_count": len(tickers),
			"tickers":      tickers,
		})
	}
	sort.SliceStable(results, func(i, j int) bool {
		return math.Abs(asFloatAny(results[i]["net_value_m"])) > math.Abs(asFloatAny(results[j]["net_value_m"]))
	})
	return results
}

func ctInvestorAdvancedRows(key string, c *ctContext, ranges map[string]bool, origin string, selectedTypes map[string]bool) []map[string]any {
	results := []map[string]any{}
	if key == "investorAktif" || key == "investorNambah" {
		for _, name := range c.T5NameOrder {
			txs := c.T5ByName[name]
			holdings := c.InvestorHoldings[name]
			lf := investorLFHoldings(holdings)
			if !passInvestorType(lf, origin, selectedTypes) {
				continue
			}
			cnt := 0
			net := 0.0
			absVal := 0.0
			byTicker := map[string]float64{}
			for _, tx := range txs {
				d, ok := t5Date(tx, c.Year)
				if !ok {
					continue
				}
				diffDays := int(c.T5AnchorDate.Sub(d).Hours() / 24)
				if diffDays < 0 || diffDays > 30 {
					continue
				}
				cnt++
				val := derefFloat(tx.ValueM)
				signed := 0.0
				if tx.Action == "BUY" {
					signed = val
				} else if tx.Action == "SELL" {
					signed = -val
				}
				net += signed
				absVal += math.Abs(val)
				tk := strings.ToUpper(tx.Ticker)
				byTicker[tk] += signed
			}
			if key == "investorAktif" {
				if cnt < 5 {
					continue
				}
				tks := make([]string, 0, len(byTicker))
				for tk := range byTicker {
					tks = append(tks, tk)
				}
				sort.Strings(tks)
				results = append(results, map[string]any{
					"name":        name,
					"lf":          lf,
					"lf_label":    lfLabel(lf),
					"count":       cnt,
					"net_value_m": roundEven(net, 4),
					"abs_value_m": roundEven(absVal, 4),
					"tickers":     tks,
				})
			} else {
				if net < 5 {
					continue
				}
				type tkVal struct {
					tk string
					v  float64
				}
				pairs := make([]tkVal, 0, len(byTicker))
				for tk, v := range byTicker {
					pairs = append(pairs, tkVal{tk: tk, v: v})
				}
				sort.SliceStable(pairs, func(i, j int) bool {
					return pairs[i].v > pairs[j].v
				})
				tks := make([]map[string]any, 0, len(pairs))
				for _, p := range pairs {
					tks = append(tks, map[string]any{
						"ticker":      p.tk,
						"net_value_m": roundEven(p.v, 4),
					})
				}
				results = append(results, map[string]any{
					"name":        name,
					"lf":          lf,
					"lf_label":    lfLabel(lf),
					"net_value_m": roundEven(net, 4),
					"count":       cnt,
					"tickers":     tks,
				})
			}
		}
		if key == "investorAktif" {
			sort.SliceStable(results, func(i, j int) bool {
				ci, _ := results[i]["count"].(int)
				cj, _ := results[j]["count"].(int)
				return ci > cj
			})
		} else {
			sort.SliceStable(results, func(i, j int) bool {
				return asFloatAny(results[i]["net_value_m"]) > asFloatAny(results[j]["net_value_m"])
			})
		}
		return results
	}

	dirUp := key == "portofolioNaik"
	for _, name := range c.InvestorOrder {
		holdings := c.InvestorHoldings[name]
		lf := investorLFHoldings(holdings)
		if !passInvestorType(lf, origin, selectedTypes) {
			continue
		}
		gl := investorGL(c, name, holdings)
		if gl == nil || gl.GLPct == nil {
			continue
		}
		pct := *gl.GLPct
		if dirUp && pct <= 0 {
			continue
		}
		if !dirUp && pct >= 0 {
			continue
		}
		if !ranges[bucketPortfolio(math.Abs(pct))] {
			continue
		}
		results = append(results, map[string]any{
			"name":         name,
			"lf":           lf,
			"lf_label":     lfLabel(lf),
			"gl_rp_m":      roundEven(gl.GL/1_000_000_000.0, 4),
			"pct":          roundEven(pct, 4),
			"ticker_count": len(holdings),
		})
	}
	sort.SliceStable(results, func(i, j int) bool {
		return math.Abs(asFloatAny(results[i]["pct"])) > math.Abs(asFloatAny(results[j]["pct"]))
	})
	return results
}

func ctExplorer(ctx context.Context, db *pgxpool.Pool, p CTExplorerParams) (map[string]any, error) {
	c, err := loadCTContext(ctx, db)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return map[string]any{
			"mode":      p.Mode,
			"category":  p.Category,
			"condition": p.Condition,
			"count":     0,
			"total":     0,
			"results":   []any{},
			"source":    "unavailable",
			"data_quality": map[string]any{
				"ct_baseline": false,
				"notes":       []string{"Baseline CT tidak tersedia."},
			},
		}, nil
	}

	mode := p.Mode
	if mode != "investors" {
		mode = "stocks"
	}
	category := p.Category
	if category != "transaksi" && category != "advanced" {
		category = "kepemilikan"
	}
	validConditions, err := validateCTCondition(mode, category, p.Condition)
	if err != nil {
		return nil, err
	}
	if !containsStrSlice(mktcapTiers, p.MktcapTier) {
		return nil, fmt.Errorf("%w: Tier market cap tidak dikenal", errCTValidation)
	}
	if category == "transaksi" && !containsStrSlice(valueRanges, p.TransaksiRange) {
		return nil, fmt.Errorf("%w: Rentang transaksi tidak dikenal", errCTValidation)
	}
	if p.InvestorOrigin != "Local" && p.InvestorOrigin != "Foreign" {
		return nil, fmt.Errorf("%w: Asal investor tidak dikenal", errCTValidation)
	}

	selectedTypes := map[string]bool{}
	allowedOpts := investorTypeOptionsForeign
	if p.InvestorOrigin == "Local" {
		allowedOpts = investorTypeOptionsLocal
	}
	if p.InvestorTypes == nil || strings.TrimSpace(*p.InvestorTypes) == "" {
		for _, o := range allowedOpts {
			selectedTypes[o] = true
		}
	} else {
		allowedSet := map[string]bool{}
		for _, o := range allowedOpts {
			allowedSet[o] = true
		}
		for _, part := range strings.Split(*p.InvestorTypes, ",") {
			item := strings.ToUpper(strings.TrimSpace(part))
			if item == "" {
				continue
			}
			if !allowedSet[item] {
				return nil, fmt.Errorf("%w: Tipe investor tidak dikenal untuk origin ini", errCTValidation)
			}
			selectedTypes[item] = true
		}
	}

	var ownRanges, sidRangeSet map[string]bool
	if category == "kepemilikan" && p.Condition != "sidNaik" && p.Condition != "sidTurun" {
		ownRanges, err = parseRangeSet(p.Ranges, ownershipRanges, []string{">10"})
	} else {
		ownRanges, err = parseRangeSet(nil, ownershipRanges, ownershipRanges)
	}
	if err != nil {
		return nil, err
	}
	if category == "kepemilikan" && (p.Condition == "sidNaik" || p.Condition == "sidTurun") {
		sidRangeSet, err = parseRangeSet(p.Ranges, sidRanges, []string{">10"})
	} else {
		sidRangeSet, err = parseRangeSet(nil, sidRanges, sidRanges)
	}
	if err != nil {
		return nil, err
	}
	retailBuckets, err := parseRangeSet(p.RetailRanges, retailRanges, retailRanges)
	if err != nil {
		return nil, err
	}
	portfolioBuckets, err := parseRangeSet(p.PortfolioRanges, portfolioRanges, portfolioRanges)
	if err != nil {
		return nil, err
	}

	sortOrderKeys, getters := baseExplorerGetters()
	defaultSort := "ticker"
	results := []map[string]any{}

	switch {
	case category == "kepemilikan" && mode == "stocks":
		if p.Condition == "sidNaik" || p.Condition == "sidTurun" {
			for _, tk := range c.TickerOrder {
				entry := c.Tickers[tk]
				row := ctStockSIDRow(entry, c, p.Condition)
				if row == nil {
					continue
				}
				if !sidRangeSet[bucketSIDPct(math.Abs(asFloatAny(row["net_pct"])))] {
					continue
				}
				if !passMktcap(entry, p.MktcapTier) {
					continue
				}
				results = append(results, row)
			}
			defaultSort = "sid"
		} else {
			for _, tk := range c.TickerOrder {
				entry := c.Tickers[tk]
				row := ctStockOwnershipRow(entry, c, p.Condition)
				if row == nil {
					continue
				}
				if !ownRanges[bucketPct(math.Abs(asFloatAny(row["net_pct"])))] {
					continue
				}
				if !passMktcap(entry, p.MktcapTier) {
					continue
				}
				results = append(results, row)
			}
			defaultSort = "net"
		}
		sort.SliceStable(results, func(i, j int) bool {
			return math.Abs(asFloatAny(results[i]["net_pct"])) > math.Abs(asFloatAny(results[j]["net_pct"]))
		})

	case category == "kepemilikan":
		results = ctInvestorOwnershipRows(c, p.Condition, ownRanges, p.InvestorOrigin, selectedTypes)
		defaultSort = "best"

	case category == "transaksi" && mode == "stocks":
		for _, tk := range c.TickerOrder {
			entry := c.Tickers[tk]
			row := ctStockTransactionRow(entry, c, p.Condition)
			if row == nil {
				continue
			}
			if bucketValue(math.Abs(asFloatAny(row["net_value_m"]))) != p.TransaksiRange {
				continue
			}
			if !passMktcap(entry, p.MktcapTier) {
				continue
			}
			results = append(results, row)
		}
		defaultSort = "net"
		getters["net"] = func(r map[string]any) any { return r["net_value_m"] }
		sort.SliceStable(results, func(i, j int) bool {
			return math.Abs(asFloatAny(results[i]["net_value_m"])) > math.Abs(asFloatAny(results[j]["net_value_m"]))
		})

	case category == "transaksi":
		results = ctInvestorTransactionRows(c, p.Condition, p.TransaksiRange, p.InvestorOrigin, selectedTypes)
		defaultSort = "net"
		getters["net"] = func(r map[string]any) any { return r["net_value_m"] }

	case category == "advanced" && mode == "stocks":
		for _, tk := range c.TickerOrder {
			entry := c.Tickers[tk]
			if p.Condition != "retailRendah" && !passMktcap(entry, p.MktcapTier) {
				continue
			}
			row := ctStockAdvancedRow(p.Condition, entry, c)
			if row == nil {
				continue
			}
			if p.Condition == "retailRendah" {
				if !retailBuckets[bucketRetail(asFloatAny(row["latest"]))] {
					continue
				}
			}
			results = append(results, row)
		}
		defaultSort = "ticker"
		getters["pct"] = func(r map[string]any) any { return r["pct_change"] }
		getters["net"] = func(r map[string]any) any {
			if v, ok := r["net_value_m"]; ok {
				return v
			}
			return r["pct_change"]
		}
		sort.SliceStable(results, func(i, j int) bool {
			si, _ := results[i]["ticker"].(string)
			sj, _ := results[j]["ticker"].(string)
			return strings.TrimSpace(si) < strings.TrimSpace(sj)
		})

	case category == "advanced":
		results = ctInvestorAdvancedRows(p.Condition, c, portfolioBuckets, p.InvestorOrigin, selectedTypes)
		switch p.Condition {
		case "investorAktif":
			defaultSort = "count"
		case "investorNambah":
			defaultSort = "net"
		case "portofolioNaik", "portofolioTurun":
			defaultSort = "pct"
		default:
			defaultSort = "name"
		}
	}

	sortTrim := strings.TrimSpace(p.Sort)
	if sortTrim != "" && sortTrim != "default" {
		results, err = sortCTRows(results, sortTrim, p.Order, defaultSort, getters)
		if err != nil {
			return nil, err
		}
	}

	lim := p.Limit
	if lim < 1 {
		lim = 1
	} else if lim > 200 {
		lim = 200
	}
	off := p.Offset
	if off < 0 {
		off = 0
	}
	total := len(results)
	start := off
	if start > total {
		start = total
	}
	end := start + lim
	if end > total {
		end = total
	}
	page := results[start:end]

	condOpts := make([]map[string]any, 0, len(validConditions))
	for _, key := range validConditions {
		lbl := conditionLabels[key]
		if lbl == "" {
			lbl = key
		}
		condOpts = append(condOpts, map[string]any{
			"id":       key,
			"label":    lbl,
			"category": category,
		})
	}
	sortOpts := make([]map[string]any, 0, len(sortOrderKeys))
	for _, key := range sortOrderKeys {
		sortOpts = append(sortOpts, map[string]any{
			"id":    key,
			"label": titleCaseSpaces(key),
		})
	}
	tierOpts := make([]map[string]any, 0, len(mktcapTiers))
	for _, key := range mktcapTiers {
		tierOpts = append(tierOpts, map[string]any{
			"id":    key,
			"label": mktcapLabels[key],
		})
	}

	condLbl := conditionLabels[p.Condition]
	if condLbl == "" {
		condLbl = p.Condition
	}

	return map[string]any{
		"mode":            mode,
		"category":        category,
		"condition":       p.Condition,
		"condition_label": condLbl,
		"count":           len(page),
		"total":           total,
		"limit":           lim,
		"offset":          off,
		"results":         page,
		"conditions":      condOpts,
		"sorts":           sortOpts,
		"filters": map[string]any{
			"ownership_ranges": ownershipRanges,
			"sid_ranges":       sidRanges,
			"value_ranges":     valueRanges,
			"mktcap_tiers":     tierOpts,
			"retail_ranges":    retailRanges,
			"portfolio_ranges": portfolioRanges,
			"investor_origins": []string{"Local", "Foreign"},
			"investor_types": map[string]any{
				"Local":   investorTypeOptionsLocal,
				"Foreign": investorTypeOptionsForeign,
			},
		},
		"period": periodPayload(c),
		"source": "ct_baseline",
		"data_quality": map[string]any{
			"ct_baseline":         true,
			"periodized_holdings": true,
			"t5_transactions":     true,
			"gl":                  true,
			"market_cap":          true,
			"notes": []string{
				"Explorer CT dihitung server-side dari baseline quarantined.",
				"T5 dan G/L tidak masuk ke jalur sinyal.",
				"Default origin investor adalah Local; ubah ke Foreign untuk melihat investor asing.",
			},
		},
	}, nil
}
