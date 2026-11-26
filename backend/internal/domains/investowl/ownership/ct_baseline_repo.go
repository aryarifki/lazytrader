package ownership

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ctSnapshotMeta(ctx context.Context, db *pgxpool.Pool) map[string]any {
	out := map[string]any{}
	if !shared.CTAvailable(ctx, db) {
		return out
	}
	rows, err := db.Query(ctx, `SELECT key, value FROM ct_baseline.snapshot_meta`)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var key string
		var raw []byte
		if scanErr := rows.Scan(&key, &raw); scanErr != nil {
			continue
		}
		var val any
		if len(raw) > 0 {
			_ = json.Unmarshal(raw, &val)
		}
		out[key] = val
	}
	return out
}

func ctMonthLabels(ctx context.Context, db *pgxpool.Pool) map[string]string {
	out := map[string]string{}
	snap := ctSnapshotMeta(ctx, db)
	rawMonths, ok := snap["pct_months"].([]any)
	if !ok {
		return out
	}
	for _, item := range rawMonths {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		k, _ := m["key"].(string)
		if k == "" {
			continue
		}
		lbl, _ := m["label"].(string)
		if lbl == "" {
			lbl = k
		}
		out[k] = lbl
	}
	return out
}

func ctTickerMetaMap(ctx context.Context, db *pgxpool.Pool, tickers []string) map[string]CTTickerMeta {
	clean := cleanTickers(tickers)
	res := make(map[string]CTTickerMeta, len(clean))
	if len(clean) == 0 || !shared.CTAvailable(ctx, db) {
		return res
	}
	rows, err := db.Query(ctx, `
		SELECT ticker, issuer, konglo, badge, free_float::double precision,
		       mktcap_b::double precision, sec_num::double precision, price::double precision,
		       sl_ribu_lot::double precision, pct_sl::double precision, pie,
		       COALESCE(badges, ARRAY[]::text[]), (badges IS NOT NULL) AS has_badges
		FROM ct_baseline.ticker_meta
		WHERE ticker = ANY($1)
	`, clean)
	if err != nil {
		return res
	}
	defer rows.Close()

	missingKonglo := []string{}
	for rows.Next() {
		var t string
		var issuer, konglo, badge *string
		var ff, mc, secNum, price, sl, pctSL *float64
		var pieRaw []byte
		var badges []string
		var hasBadges bool
		if scanErr := rows.Scan(
			&t, &issuer, &konglo, &badge, &ff, &mc, &secNum, &price, &sl, &pctSL,
			&pieRaw, &badges, &hasBadges,
		); scanErr != nil {
			continue
		}
		var kongloPtr *string
		if konglo != nil && *konglo != "" && *konglo != "-" {
			k := *konglo
			kongloPtr = &k
		}
		var finalBadges []string
		if hasBadges {
			if badges == nil {
				finalBadges = []string{}
			} else {
				finalBadges = badges
			}
		} else if badge != nil && *badge != "" {
			finalBadges = []string{*badge}
		} else {
			finalBadges = []string{}
		}
		pie := map[string]any{}
		if len(pieRaw) > 0 {
			var parsed map[string]any
			if jsonErr := json.Unmarshal(pieRaw, &parsed); jsonErr == nil && parsed != nil {
				pie = parsed
			}
		}
		res[t] = CTTickerMeta{
			Ticker:    t,
			Issuer:    issuer,
			Konglo:    kongloPtr,
			Badge:     badge,
			Badges:    finalBadges,
			FreeFloat: ff,
			MktcapB:   mc,
			SecNum:    secNum,
			Price:     price,
			SLRibuLot: sl,
			PctSL:     pctSL,
			Pie:       pie,
		}
		if kongloPtr == nil {
			missingKonglo = append(missingKonglo, t)
		}
	}

	if len(missingKonglo) > 0 {
		cmRows, cmErr := db.Query(ctx, `
			SELECT ticker, group_name
			FROM ownership_staging.v_conglomerate_members
			WHERE ticker = ANY($1)
		`, missingKonglo)
		if cmErr == nil {
			defer cmRows.Close()
			for cmRows.Next() {
				var tk string
				var gn *string
				if scanErr := cmRows.Scan(&tk, &gn); scanErr != nil {
					continue
				}
				if meta, ok := res[tk]; ok && gn != nil && *gn != "" {
					g := *gn
					meta.Konglo = &g
					res[tk] = meta
				}
			}
		}
	}
	return res
}

func resolveBadges(badges []string, badge *string) []string {
	if len(badges) > 0 {
		return badges
	}
	if badge != nil && *badge != "" {
		return []string{*badge}
	}
	return []string{}
}

func ctHolders(ctx context.Context, db *pgxpool.Pool, ticker string) []CTHolderRow {
	t := strings.ToUpper(strings.TrimSpace(ticker))
	out := []CTHolderRow{}
	if t == "" || !shared.CTAvailable(ctx, db) {
		return out
	}
	rows, err := db.Query(ctx, `
		SELECT name, lf, sl_mei::double precision, pct_mei::double precision,
		       sl_jun::double precision, pct_jun::double precision, status
		FROM ct_baseline.holders
		WHERE ticker = $1
		ORDER BY COALESCE(pct_jun, pct_mei) DESC NULLS LAST, name ASC
	`, t)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var r CTHolderRow
		if scanErr := rows.Scan(&r.Name, &r.LF, &r.SLMei, &r.PctMei, &r.SLJun, &r.PctJun, &r.Status); scanErr != nil {
			continue
		}
		out = append(out, r)
	}
	return out
}

func ctPctHistory(ctx context.Context, db *pgxpool.Pool, ticker string) []CTPctPoint {
	t := strings.ToUpper(strings.TrimSpace(ticker))
	out := []CTPctPoint{}
	if t == "" || !shared.CTAvailable(ctx, db) {
		return out
	}
	rows, err := db.Query(ctx, `
		SELECT month_key, month_index, pct_asing::double precision,
		       pct_asing_lt1::double precision, pct_ritel::double precision
		FROM ct_baseline.pct_history
		WHERE ticker = $1
		ORDER BY month_index ASC NULLS LAST, month_key ASC
	`, t)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var p CTPctPoint
		if scanErr := rows.Scan(&p.MonthKey, &p.MonthIndex, &p.PctAsing, &p.PctAsingLt1, &p.PctRitel); scanErr != nil {
			continue
		}
		out = append(out, p)
	}
	return out
}

func ctLatestPctMap(ctx context.Context, db *pgxpool.Pool, tickers []string) map[string]CTPctPoint {
	clean := cleanTickers(tickers)
	out := make(map[string]CTPctPoint, len(clean))
	if len(clean) == 0 || !shared.CTAvailable(ctx, db) {
		return out
	}
	rows, err := db.Query(ctx, `
		SELECT DISTINCT ON (ticker)
		       ticker, month_key, pct_asing::double precision,
		       pct_asing_lt1::double precision, pct_ritel::double precision
		FROM ct_baseline.pct_history
		WHERE ticker = ANY($1)
		ORDER BY ticker, month_index DESC NULLS LAST, month_key DESC
	`, clean)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var tk string
		var p CTPctPoint
		if scanErr := rows.Scan(&tk, &p.MonthKey, &p.PctAsing, &p.PctAsingLt1, &p.PctRitel); scanErr != nil {
			continue
		}
		out[tk] = p
	}
	return out
}

func ctHolderAggregates(ctx context.Context, db *pgxpool.Pool, tickers []string) map[string]CTHolderAgg {
	clean := cleanTickers(tickers)
	out := make(map[string]CTHolderAgg, len(clean))
	if len(clean) == 0 || !shared.CTAvailable(ctx, db) {
		return out
	}
	rows, err := db.Query(ctx, `
		SELECT ticker,
		       COUNT(*)::int AS holder_count,
		       COALESCE(SUM(COALESCE(pct_jun, pct_mei)), 0)::double precision AS total_pct,
		       COALESCE(MAX(COALESCE(pct_jun, pct_mei)), 0)::double precision AS max_pct
		FROM ct_baseline.holders
		WHERE ticker = ANY($1)
		GROUP BY ticker
	`, clean)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var tk string
		var agg CTHolderAgg
		if scanErr := rows.Scan(&tk, &agg.HolderCount, &agg.TotalPct, &agg.MaxPct); scanErr != nil {
			continue
		}
		out[tk] = agg
	}
	return out
}

func ctSIDHistory(ctx context.Context, db *pgxpool.Pool, ticker string) []CTSIDPoint {
	t := strings.ToUpper(strings.TrimSpace(ticker))
	out := []CTSIDPoint{}
	if t == "" || !shared.CTAvailable(ctx, db) {
		return out
	}
	rows, err := db.Query(ctx, `
		SELECT period, sid
		FROM ct_baseline.sid_history
		WHERE ticker = $1
		ORDER BY period ASC
	`, t)
	if err != nil {
		return out
	}
	defer rows.Close()
	var prev *int
	for rows.Next() {
		var period string
		var sid int
		if scanErr := rows.Scan(&period, &sid); scanErr != nil {
			continue
		}
		var chg *int
		var chgPct *float64
		if prev != nil {
			d := sid - *prev
			chg = &d
			if *prev != 0 {
				v := roundEven(float64(d)/float64(*prev)*100.0, 2)
				chgPct = &v
			}
		}
		out = append(out, CTSIDPoint{
			Period:    period,
			SID:       sid,
			Change:    chg,
			ChangePct: chgPct,
		})
		pCopy := sid
		prev = &pCopy
	}
	return out
}

func ctLatestSIDMap(ctx context.Context, db *pgxpool.Pool, tickers []string) map[string]int {
	clean := cleanTickers(tickers)
	out := make(map[string]int, len(clean))
	if len(clean) == 0 || !shared.CTAvailable(ctx, db) {
		return out
	}
	rows, err := db.Query(ctx, `
		SELECT DISTINCT ON (ticker) ticker, sid
		FROM ct_baseline.sid_history
		WHERE ticker = ANY($1)
		ORDER BY ticker, period DESC
	`, clean)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var tk string
		var sid int
		if scanErr := rows.Scan(&tk, &sid); scanErr != nil {
			continue
		}
		out[tk] = sid
	}
	return out
}

func ctT5Rows(ctx context.Context, db *pgxpool.Pool, tickers []string, limit int) []CTT5Row {
	out := []CTT5Row{}
	if !shared.CTAvailable(ctx, db) {
		return out
	}
	if limit < 1 {
		limit = 1
	} else if limit > 1000 {
		limit = 1000
	}
	clean := cleanTickers(tickers)

	var (
		sqlStr string
		args   []any
	)
	if len(clean) > 0 {
		sqlStr = `
			SELECT ticker, date_label, broker, name, action,
			       change_lot::double precision, value_m::double precision,
			       lot_after::double precision, pct_juli::double precision
			FROM ct_baseline.t5_transactions
			WHERE ticker = ANY($1)
			ORDER BY month_num DESC NULLS LAST, day_num DESC NULLS LAST, ticker ASC, sort_order ASC
			LIMIT $2
		`
		args = []any{clean, limit}
	} else {
		sqlStr = `
			SELECT ticker, date_label, broker, name, action,
			       change_lot::double precision, value_m::double precision,
			       lot_after::double precision, pct_juli::double precision
			FROM ct_baseline.t5_transactions
			ORDER BY month_num DESC NULLS LAST, day_num DESC NULLS LAST, ticker ASC, sort_order ASC
			LIMIT $1
		`
		args = []any{limit}
	}
	rows, err := db.Query(ctx, sqlStr, args...)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var r CTT5Row
		if scanErr := rows.Scan(&r.Ticker, &r.Date, &r.Broker, &r.Name, &r.Action, &r.ChangeLot, &r.ValueM, &r.LotAfter, &r.PctJuli); scanErr != nil {
			continue
		}
		out = append(out, r)
	}
	return out
}

func ctT5Summary(rows []CTT5Row) CTT5Summary {
	var buyCnt, sellCnt int
	var buyVal, sellVal float64
	for _, r := range rows {
		act := ""
		if r.Action != nil {
			act = strings.ToUpper(strings.TrimSpace(*r.Action))
		}
		v := derefFloat(r.ValueM)
		if act == "BUY" {
			buyCnt++
			buyVal += v
		} else if act == "SELL" {
			sellCnt++
			sellVal += v
		}
	}
	return CTT5Summary{
		Count:      len(rows),
		BuyCount:   buyCnt,
		SellCount:  sellCnt,
		BuyValueM:  roundEven(buyVal, 4),
		SellValueM: roundEven(sellVal, 4),
		NetValueM:  roundEven(buyVal-sellVal, 4),
	}
}

func ctT5MapSummary(rows []map[string]any) CTT5Summary {
	var buyCnt, sellCnt int
	var buyVal, sellVal float64
	for _, r := range rows {
		act, _ := r["action"].(string)
		act = strings.ToUpper(strings.TrimSpace(act))
		v, _ := r["value_m"].(float64)
		if act == "BUY" {
			buyCnt++
			buyVal += v
		} else if act == "SELL" {
			sellCnt++
			sellVal += v
		}
	}
	return CTT5Summary{
		Count:      len(rows),
		BuyCount:   buyCnt,
		SellCount:  sellCnt,
		BuyValueM:  roundEven(buyVal, 4),
		SellValueM: roundEven(sellVal, 4),
		NetValueM:  roundEven(buyVal-sellVal, 4),
	}
}

func getIHSGComposition(ctx context.Context, db *pgxpool.Pool) map[string]any {
	if !shared.CTAvailable(ctx, db) {
		return nil
	}
	mRows, err := db.Query(ctx, `
		SELECT month_key, total_trn::double precision, foreign_named_trn::double precision,
		       foreign_hidden_trn::double precision, local_nonret_trn::double precision,
		       retail_trn::double precision
		FROM ct_baseline.ihsg_month
	`)
	if err != nil {
		return nil
	}
	defer mRows.Close()

	months := map[string]map[string]any{}
	monthCount := 0
	for mRows.Next() {
		var mk string
		var tot, fn, fh, ln, rt *float64
		if scanErr := mRows.Scan(&mk, &tot, &fn, &fh, &ln, &rt); scanErr != nil {
			continue
		}
		monthCount++
		tVal := derefFloat(tot)
		fnVal := derefFloat(fn)
		fhVal := derefFloat(fh)
		lnVal := derefFloat(ln)
		rtVal := derefFloat(rt)
		pctOf := func(part float64) float64 {
			if tVal == 0 {
				return 0.0
			}
			return roundEven(part/tVal*100.0, 2)
		}
		months[mk] = map[string]any{
			"total_trn":          tVal,
			"foreign_named_trn":  fnVal,
			"foreign_named_pct":  pctOf(fnVal),
			"foreign_hidden_trn": fhVal,
			"foreign_hidden_pct": pctOf(fhVal),
			"local_nonret_trn":   lnVal,
			"local_nonret_pct":   pctOf(lnVal),
			"retail_trn":         rtVal,
			"retail_pct":         pctOf(rtVal),
		}
	}

	pRows, pErr := db.Query(ctx, `SELECT key, value FROM ct_baseline.ihsg_payload`)
	if pErr != nil {
		return nil
	}
	defer pRows.Close()
	payloads := map[string]any{}
	payloadCount := 0
	for pRows.Next() {
		var k string
		var raw []byte
		if scanErr := pRows.Scan(&k, &raw); scanErr != nil {
			continue
		}
		payloadCount++
		var val any
		if len(raw) > 0 {
			_ = json.Unmarshal(raw, &val)
		}
		payloads[k] = val
	}

	if monthCount == 0 && payloadCount == 0 {
		return nil
	}

	var top20Raw []any
	if s, ok := payloads["top20_jun"].([]any); ok {
		top20Raw = s
	}
	pieExtra, _ := payloads["pie_extra"].(map[string]any)
	if pieExtra == nil {
		pieExtra = map[string]any{}
	}

	top20Tickers := make([]string, 0, len(top20Raw))
	for _, item := range top20Raw {
		if m, ok := item.(map[string]any); ok {
			if tk, ok := m["ticker"].(string); ok && tk != "" {
				top20Tickers = append(top20Tickers, tk)
			}
		}
	}
	top20Meta := ctTickerMetaMap(ctx, db, top20Tickers)

	top20 := make([]map[string]any, 0, len(top20Raw))
	for idx, item := range top20Raw {
		row, _ := item.(map[string]any)
		if row == nil {
			row = map[string]any{}
		}
		tkStr, _ := row["ticker"].(string)
		tClean := strings.ToUpper(strings.TrimSpace(tkStr))
		meta := top20Meta[tClean]
		badges := meta.Badges
		if badges == nil {
			badges = []string{}
		}
		top20 = append(top20, map[string]any{
			"rank":               idx + 1,
			"ticker":             tClean,
			"issuer":             row["issuer"],
			"badges":             badges,
			"conglomerate":       meta.Konglo,
			"foreign_named_trn":  asFloat(row["foreign_named_trn"]),
			"foreign_named_pct":  asFloat(row["foreign_named_pct"]),
			"foreign_hidden_trn": asFloat(row["foreign_hidden_trn"]),
			"foreign_hidden_pct": asFloat(row["foreign_hidden_pct"]),
			"local_nonret_trn":   asFloat(row["local_nonret_trn"]),
			"local_nonret_pct":   asFloat(row["local_nonret_pct"]),
			"retail_trn":         asFloat(row["retail_trn"]),
			"retail_pct":         asFloat(row["retail_pct"]),
		})
	}

	jun := months["jun"]
	mei := months["mei"]
	scriptlessTrn := asFloat(jun["total_trn"])
	warkatTrn := asFloat(pieExtra["warkat_trn"])
	totalMktcapTrn := asFloat(pieExtra["total_mktcap_trn"])
	if totalMktcapTrn == 0 {
		totalMktcapTrn = scriptlessTrn + warkatTrn
	}

	pctMkt := func(v float64) float64 {
		if totalMktcapTrn == 0 {
			return 0.0
		}
		return roundEven(v/totalMktcapTrn*100.0, 2)
	}
	donutScripless := []map[string]any{
		{
			"label":     "Scriptless (C-BEST)",
			"value_trn": scriptlessTrn,
			"pct":       pctMkt(scriptlessTrn),
			"color":     "#1DB87A",
		},
		{
			"label":     "Warkat (Fisik)",
			"value_trn": warkatTrn,
			"pct":       pctMkt(warkatTrn),
			"color":     "#64748B",
		},
	}

	fn := asFloat(jun["foreign_named_trn"])
	fh := asFloat(jun["foreign_hidden_trn"])
	forTot := fn + fh
	ln := asFloat(jun["local_nonret_trn"])
	rt := asFloat(jun["retail_trn"])
	pctScr := func(v float64) float64 {
		if scriptlessTrn == 0 {
			return 0.0
		}
		return roundEven(v/scriptlessTrn*100.0, 2)
	}
	donutInvestorType := []map[string]any{
		{
			"label":     "Foreign (Asing)",
			"value_trn": roundEven(forTot, 2),
			"pct":       pctScr(forTot),
			"color":     "#3B82F6",
		},
		{
			"label":     "Local Non-Retail",
			"value_trn": roundEven(ln, 2),
			"pct":       pctScr(ln),
			"color":     "#10B981",
		},
		{
			"label":     "Retail",
			"value_trn": roundEven(rt, 2),
			"pct":       pctScr(rt),
			"color":     "#F59E0B",
		},
	}

	nForGE50 := int(asFloat(pieExtra["n_foreign_ge50"]))
	if nForGE50 == 0 {
		nForGE50 = 180
	}
	nTotTickers := int(asFloat(pieExtra["n_tickers"]))
	if nTotTickers == 0 {
		nTotTickers = 1005
	}
	nForLT50 := nTotTickers - nForGE50
	if nForLT50 < 0 {
		nForLT50 = 0
	}
	pctTk := func(n int) float64 {
		if nTotTickers == 0 {
			return 0.0
		}
		return roundEven(float64(n)/float64(nTotTickers)*100.0, 2)
	}
	donutForeignDist := []map[string]any{
		{
			"label": "Asing ≥ 50%",
			"count": nForGE50,
			"pct":   pctTk(nForGE50),
			"color": "#8B5CF6",
		},
		{
			"label": "Asing < 50%",
			"count": nForLT50,
			"pct":   pctTk(nForLT50),
			"color": "#06B6D4",
		},
	}

	summaryComparison := []map[string]any{
		{
			"category": "Total Saham Scriptless",
			"prev_trn": mei["total_trn"],
			"curr_trn": jun["total_trn"],
			"diff_trn": roundEven(asFloat(jun["total_trn"])-asFloat(mei["total_trn"]), 2),
			"pct":      100.0,
		},
		{
			"category": "Foreign Named",
			"prev_trn": mei["foreign_named_trn"],
			"curr_trn": jun["foreign_named_trn"],
			"diff_trn": roundEven(asFloat(jun["foreign_named_trn"])-asFloat(mei["foreign_named_trn"]), 2),
			"pct":      jun["foreign_named_pct"],
		},
		{
			"category": "Foreign Hidden Player",
			"prev_trn": mei["foreign_hidden_trn"],
			"curr_trn": jun["foreign_hidden_trn"],
			"diff_trn": roundEven(asFloat(jun["foreign_hidden_trn"])-asFloat(mei["foreign_hidden_trn"]), 2),
			"pct":      jun["foreign_hidden_pct"],
		},
		{
			"category": "Local Non-Retail",
			"prev_trn": mei["local_nonret_trn"],
			"curr_trn": jun["local_nonret_trn"],
			"diff_trn": roundEven(asFloat(jun["local_nonret_trn"])-asFloat(mei["local_nonret_trn"]), 2),
			"pct":      jun["local_nonret_pct"],
		},
		{
			"category": "Retail",
			"prev_trn": mei["retail_trn"],
			"curr_trn": jun["retail_trn"],
			"diff_trn": roundEven(asFloat(jun["retail_trn"])-asFloat(mei["retail_trn"]), 2),
			"pct":      jun["retail_pct"],
		},
	}

	return map[string]any{
		"months":             months,
		"summary_comparison": summaryComparison,
		"donuts": map[string]any{
			"scripless_vs_warkat":  donutScripless,
			"investor_type":        donutInvestorType,
			"foreign_distribution": donutForeignDist,
		},
		"top20":            top20,
		"total_mktcap_trn": totalMktcapTrn,
		"warkat_trn":       warkatTrn,
		"n_tickers":        nTotTickers,
		"n_foreign_ge50":   nForGE50,
	}
}

func asFloat(v any) float64 {
	if v == nil {
		return 0.0
	}
	switch x := v.(type) {
	case float64:
		return x
	case float32:
		return float64(x)
	case int:
		return float64(x)
	case int64:
		return float64(x)
	default:
		return 0.0
	}
}
