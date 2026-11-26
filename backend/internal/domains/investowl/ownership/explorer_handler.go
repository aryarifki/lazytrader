package ownership

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
)

type condDef struct {
	ID     string
	Label  string
	Clause string
	Group  string
}

var stockConditionsOrdered = []condDef{
	{"all", "Semua emiten", "", "ownership"},
	{"naik_porsi", "Big Player Menambah Porsi Bulan Ini", "h.ticker IN (SELECT ticker FROM ownership_staging.v_explorer_naik_porsi)", "big_player"},
	{"baru", "Masuknya Big Player Baru Bulan Ini", "h.ticker IN (SELECT ticker FROM ownership_staging.v_explorer_baru)", "big_player"},
	{"mayoritas", "Pemegang Saham Mayoritas Menambah Porsi", "h.ticker IN (SELECT ticker FROM ownership_staging.v_explorer_mayoritas)", "big_player"},
	{"foreign_naik", "Kepemilikan Foreign Meningkat >= 0.1 pp", "h.ticker IN (SELECT ticker FROM ownership_staging.v_explorer_foreign_naik)", "big_player"},
	{"foreign_3bulan", "Foreign Akumulasi 3 Bulan Beruntun", "h.ticker IN (SELECT ticker FROM ownership_staging.v_explorer_foreign_3bulan)", "big_player"},
	{"aktif_beli", "Pembelian Aktif oleh Big Player", "h.ticker IN (SELECT ticker FROM ownership_staging.v_explorer_aktif_beli)", "big_player"},
	{"buyback", "Emiten Buyback Bulan Ini", "h.ticker IN (SELECT ticker FROM ownership_staging.v_explorer_buyback)", "big_player"},
	{"controlling_high", "Pengendali > 50%", "h.controlling_pct >= 50", "ownership"},
	{"controlling_low", "Pengendali < 20%", "h.controlling_pct < 20", "ownership"},
	{"public_high", "Publik > 30%", "h.public_pct >= 30", "ownership"},
	{"public_low", "Publik < 10%", "h.public_pct < 10", "ownership"},
	{"many_holders", "Pemegang saham >= 20", "h.holder_count >= 20", "ownership"},
	{"few_holders", "Pemegang saham <= 3", "h.holder_count <= 3", "ownership"},
	{"top_holder_dominant", "Top holder >= 50%", "h.max_holder_pct >= 50", "ownership"},
	{"treasury", "Ada saham treasury", "h.treasury_pct > 0", "ownership"},
	{"cross_holding", "Pemegang lintas emiten", "EXISTS (SELECT 1 FROM company_shareholders cs2 JOIN multi_holder mh ON mh.name = cs2.name WHERE cs2.ticker = h.ticker)", "ownership"},
	{"price_up_mtd", "Harga naik MTD > 2%", "l.close IS NOT NULL AND m.close > 0 AND (l.close / m.close - 1) > 0.02", "price"},
	{"price_down_mtd", "Harga turun MTD < -2%", "l.close IS NOT NULL AND m.close > 0 AND (l.close / m.close - 1) < -0.02", "price"},
	{"foreign_inflow_5d", "Foreign net buy 5 hari", "COALESCE(f.foreign_5d, 0) > 0", "flow"},
	{"foreign_outflow_5d", "Foreign net sell 5 hari", "COALESCE(f.foreign_5d, 0) < 0", "flow"},
}

var stockConditionMap = func() map[string]condDef {
	m := make(map[string]condDef, len(stockConditionsOrdered))
	for _, c := range stockConditionsOrdered {
		m[c.ID] = c
	}
	return m
}()

var conditionGroupLabelsOrdered = [][2]string{
	{"big_player", "Big Player (Source-Backed)"},
	{"ownership", "Kepemilikan"},
	{"price", "Harga"},
	{"flow", "Arus asing"},
}

var investorConditionsOrdered = []condDef{
	{"all", "Semua investor", "", ""},
	{"investor_aktif", "Investor Aktif Bertransaksi Sebulan Terakhir", "name IN (SELECT investor_name FROM ownership_staging.v_explorer_investor_aktif)", ""},
	{"investor_nambah", "Investor Menambah Porsi >= Rp5 Miliar", "name IN (SELECT investor_name FROM ownership_staging.v_explorer_investor_nambah)", ""},
	{"multi_holding", "Memegang >= 5 emiten", "ticker_count >= 5", ""},
	{"controlling", "Punya saham pengendali", "controlling_count >= 1", ""},
	{"substantial", "Punya porsi >= 5%", "max_pct >= 5", ""},
	{"top_holder", "Porsi >= 25%", "max_pct >= 25", ""},
	{"cross_sector", "Lintas sektor", "sector_count >= 2", ""},
	{"single_holding", "Hanya 1 emiten", "ticker_count = 1", ""},
}

var investorConditionMap = func() map[string]condDef {
	m := make(map[string]condDef, len(investorConditionsOrdered))
	for _, c := range investorConditionsOrdered {
		m[c.ID] = c
	}
	return m
}()

var stockSortsOrdered = [][2]string{
	{"default", "Default (top holder)"},
	{"ticker", "Ticker"},
	{"name", "Nama emiten"},
	{"holders", "Jumlah pemegang"},
	{"controlling", "Porsi pengendali"},
	{"public", "Porsi publik"},
	{"top_holder", "Porsi top holder"},
	{"substantial", "Porsi substantial"},
	{"treasury", "Porsi treasury"},
	{"market_cap", "Market cap"},
	{"mtd", "Perubahan MTD"},
	{"foreign_5d", "Foreign flow 5D"},
}

var stockSortSQL = map[string]string{
	"ticker":      "h.ticker",
	"name":        "name",
	"holders":     "h.holder_count",
	"controlling": "h.controlling_pct",
	"public":      "h.public_pct",
	"top_holder":  "h.max_holder_pct",
	"substantial": "h.substantial_pct",
	"treasury":    "h.treasury_pct",
	"market_cap":  "(l.close * sh.listed_shares / 1000000000.0)",
	"mtd":         "mtd_pct",
	"foreign_5d":  "f.foreign_5d",
}

var investorSortsOrdered = [][2]string{
	{"default", "Default (jumlah emiten)"},
	{"name", "Nama investor"},
	{"tickers", "Jumlah emiten"},
	{"total_pct", "Total porsi"},
	{"max_pct", "Porsi maksimum"},
	{"controlling", "Jumlah pengendali"},
	{"substantial", "Jumlah substantial"},
	{"sectors", "Jumlah sektor"},
}

var investorSortSQL = map[string]string{
	"name":        "name",
	"tickers":     "ticker_count",
	"total_pct":   "total_pct",
	"max_pct":     "max_pct",
	"controlling": "controlling_count",
	"substantial": "substantial_count",
	"sectors":     "sector_count",
}

var unsupportedConditions = []map[string]any{
	{
		"id":     "sid_change",
		"label":  "Perubahan SID",
		"reason": "Data SID historis KSEI tidak disediakan oleh feed broker (hanya snapshot bulanan untuk 15 emiten).",
	},
	{
		"id":     "gl",
		"label":  "G/L portofolio",
		"reason": "Memerlukan histori modal transaksi lengkap untuk seluruh investor.",
	},
}

func parseOptInt(q url.Values, key string, minVal int) (*int, bool) {
	raw := strings.TrimSpace(q.Get(key))
	if raw == "" {
		return nil, true
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v < minVal {
		return nil, false
	}
	return &v, true
}

func parseOptFloat(q url.Values, key string, minVal, maxVal float64, hasMax bool) (*float64, bool) {
	raw := strings.TrimSpace(q.Get(key))
	if raw == "" {
		return nil, true
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil || v < minVal || (hasMax && v > maxVal) {
		return nil, false
	}
	return &v, true
}

// GetExplorer handles GET /api/v1/ownership/explorer.
func (h *Handler) GetExplorer(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	mode := q.Get("mode")
	if mode == "" {
		mode = "stocks"
	}
	if mode != "stocks" && mode != "investors" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "mode must be stocks or investors")
		return
	}
	condition := q.Get("condition")
	if condition == "" {
		condition = "all"
	}
	searchQ := q.Get("q")
	sector := q.Get("sector")

	minHolders, ok := parseOptInt(q, "min_holders", 0)
	if !ok {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid min_holders")
		return
	}
	maxHolders, ok := parseOptInt(q, "max_holders", 0)
	if !ok {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid max_holders")
		return
	}
	minCtrl, ok := parseOptFloat(q, "min_controlling_pct", 0, 100, true)
	if !ok {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid min_controlling_pct")
		return
	}
	maxCtrl, ok := parseOptFloat(q, "max_controlling_pct", 0, 100, true)
	if !ok {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid max_controlling_pct")
		return
	}
	minPub, ok := parseOptFloat(q, "min_public_pct", 0, 100, true)
	if !ok {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid min_public_pct")
		return
	}
	maxPub, ok := parseOptFloat(q, "max_public_pct", 0, 100, true)
	if !ok {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid max_public_pct")
		return
	}
	minTop, ok := parseOptFloat(q, "min_top_holder_pct", 0, 100, true)
	if !ok {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid min_top_holder_pct")
		return
	}
	maxTop, ok := parseOptFloat(q, "max_top_holder_pct", 0, 100, true)
	if !ok {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid max_top_holder_pct")
		return
	}
	minMktCap, ok := parseOptFloat(q, "min_market_cap_b", 0, 0, false)
	if !ok {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid min_market_cap_b")
		return
	}
	maxMktCap, ok := parseOptFloat(q, "max_market_cap_b", 0, 0, false)
	if !ok {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid max_market_cap_b")
		return
	}

	mtd := q.Get("mtd")
	if mtd == "" {
		mtd = "all"
	}
	if mtd != "all" && mtd != "up" && mtd != "down" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "mtd must be all, up, or down")
		return
	}
	fFlow := q.Get("foreign_flow")
	if fFlow == "" {
		fFlow = "all"
	}
	if fFlow != "all" && fFlow != "inflow" && fFlow != "outflow" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "foreign_flow must be all, inflow, or outflow")
		return
	}
	hasTreasury := false
	if raw := strings.TrimSpace(q.Get("has_treasury")); raw != "" {
		b, bErr := strconv.ParseBool(raw)
		if bErr != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid has_treasury")
			return
		}
		hasTreasury = b
	}

	minTickers, ok := parseOptInt(q, "min_tickers", 0)
	if !ok {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid min_tickers")
		return
	}
	maxTickers, ok := parseOptInt(q, "max_tickers", 0)
	if !ok {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid max_tickers")
		return
	}
	minMaxPct, ok := parseOptFloat(q, "min_max_pct", 0, 100, true)
	if !ok {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid min_max_pct")
		return
	}
	maxMaxPct, ok := parseOptFloat(q, "max_max_pct", 0, 100, true)
	if !ok {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid max_max_pct")
		return
	}
	minSecCount, ok := parseOptInt(q, "min_sector_count", 0)
	if !ok {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid min_sector_count")
		return
	}

	sortKey := q.Get("sort")
	if sortKey == "" {
		sortKey = "default"
	}
	order := q.Get("order")
	if order == "" {
		order = "desc"
	}
	if order != "asc" && order != "desc" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "order must be asc or desc")
		return
	}
	limit := 50
	if raw := strings.TrimSpace(q.Get("limit")); raw != "" {
		n, nErr := strconv.Atoi(raw)
		if nErr != nil || n < 1 || n > 200 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "limit must be between 1 and 200")
			return
		}
		limit = n
	}
	offset := 0
	if raw := strings.TrimSpace(q.Get("offset")); raw != "" {
		n, nErr := strconv.Atoi(raw)
		if nErr != nil || n < 0 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "offset must be >= 0")
			return
		}
		offset = n
	}

	ck := ownershipCacheKey("ownership_explorer_v3", canonicalQueryString(q))
	if h.tryCache(w, r, ck, 1800) {
		return
	}

	ctx := r.Context()
	sectors, err := explorerSectors(ctx, h.DB)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	termTrim := strings.TrimSpace(searchQ)
	hasTerm := termTrim != ""

	if mode == "investors" {
		cond, ok := investorConditionMap[condition]
		if !ok {
			httpx.WriteError(w, http.StatusBadRequest, "Kondisi investor tidak dikenal")
			return
		}
		sortFound := false
		for _, pair := range investorSortsOrdered {
			if pair[0] == sortKey {
				sortFound = true
				break
			}
		}
		if !sortFound {
			httpx.WriteError(w, http.StatusBadRequest, "Urutan investor tidak dikenal")
			return
		}

		args := []any{limit, offset}
		addArg := func(v any) string {
			args = append(args, v)
			return fmt.Sprintf("$%d", len(args))
		}

		whereParts := []string{}
		if cond.Clause != "" {
			whereParts = append(whereParts, cond.Clause)
		}
		if hasTerm {
			whereParts = append(whereParts, "name ILIKE "+addArg("%"+termTrim+"%"))
		}
		if minTickers != nil {
			whereParts = append(whereParts, "ticker_count >= "+addArg(*minTickers))
		}
		if maxTickers != nil {
			whereParts = append(whereParts, "ticker_count <= "+addArg(*maxTickers))
		}
		if minMaxPct != nil {
			whereParts = append(whereParts, "max_pct >= "+addArg(*minMaxPct))
		}
		if maxMaxPct != nil {
			whereParts = append(whereParts, "max_pct <= "+addArg(*maxMaxPct))
		}
		if minSecCount != nil {
			whereParts = append(whereParts, "sector_count >= "+addArg(*minSecCount))
		}
		where := "TRUE"
		if len(whereParts) > 0 {
			where = strings.Join(whereParts, " AND ")
		}

		orderClause := "ticker_count DESC, total_pct DESC, name ASC"
		if sortKey != "default" {
			orderClause = investorSortSQL[sortKey] + " " + strings.ToUpper(order) + " NULLS LAST, name ASC"
		}

		sqlStr := `
			WITH inv AS (
				SELECT cs.name,
				       COUNT(DISTINCT cs.ticker)::int AS ticker_count,
				       COALESCE(SUM(cs.pct), 0)::double precision AS total_pct,
				       COALESCE(MAX(cs.pct), 0)::double precision AS max_pct,
				       SUM(CASE WHEN cs.is_controlling THEN 1 ELSE 0 END)::int AS controlling_count,
				       SUM(CASE WHEN NOT cs.is_controlling
				                 AND NOT (cs.name ILIKE 'Masyarakat%' OR cs.name ILIKE 'Publik%')
				                 AND NOT (cs.name ILIKE '%Treasury%')
				                 AND cs.pct >= 5 THEN 1 ELSE 0 END)::int AS substantial_count,
				       COUNT(DISTINCT COALESCE(t.sector, '-'))::int AS sector_count,
				       (array_agg(cs.ticker ORDER BY cs.pct DESC))[1:5] AS sample_tickers
				FROM company_shareholders cs
				LEFT JOIN tickers t ON t.ticker = cs.ticker
				WHERE NOT (cs.name ILIKE 'Masyarakat%' OR cs.name ILIKE 'Publik%')
				  AND NOT (cs.name ILIKE '%Treasury%')
				  AND cs.name NOT IN ('-', 'Pengendali Saham', 'Non Pengendali Saham', 'Others', 'Lain-lain')
				GROUP BY cs.name
			)
			SELECT name, ticker_count, total_pct, max_pct, controlling_count,
			       substantial_count, sector_count, sample_tickers,
			       COUNT(*) OVER()::int AS total_count
			FROM inv
			WHERE ` + where + `
			ORDER BY ` + orderClause + `
			LIMIT $1 OFFSET $2
		`

		rows, qErr := h.DB.Query(ctx, sqlStr, args...)
		if qErr != nil {
			httpx.WriteError(w, http.StatusInternalServerError, qErr.Error())
			return
		}
		defer rows.Close()

		results := []map[string]any{}
		total := 0
		for rows.Next() {
			var nm string
			var tkCnt, ctrlCnt, substCnt, secCnt, totCnt int
			var totPct, maxPct float64
			var samples []string
			if scanErr := rows.Scan(&nm, &tkCnt, &totPct, &maxPct, &ctrlCnt, &substCnt, &secCnt, &samples, &totCnt); scanErr != nil {
				httpx.WriteError(w, http.StatusInternalServerError, scanErr.Error())
				return
			}
			total = totCnt
			if samples == nil {
				samples = []string{}
			}
			results = append(results, map[string]any{
				"name":              nm,
				"ticker_count":      tkCnt,
				"total_pct":         roundEven(totPct, 3),
				"max_pct":           roundEven(maxPct, 3),
				"controlling_count": ctrlCnt,
				"substantial_count": substCnt,
				"sector_count":      secCnt,
				"sample_tickers":    samples,
			})
		}
		if err := rows.Err(); err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		condMeta := make([]map[string]any, 0, len(investorConditionsOrdered))
		for _, c := range investorConditionsOrdered {
			condMeta = append(condMeta, map[string]any{"id": c.ID, "label": c.Label, "available": true})
		}
		sortMeta := make([]map[string]any, 0, len(investorSortsOrdered))
		for _, s := range investorSortsOrdered {
			sortMeta = append(sortMeta, map[string]any{"id": s[0], "label": s[1]})
		}

		resp := map[string]any{
			"mode":                   "investors",
			"condition":              condition,
			"condition_label":        cond.Label,
			"count":                  len(results),
			"total":                  total,
			"limit":                  limit,
			"offset":                 offset,
			"results":                results,
			"conditions":             condMeta,
			"condition_groups":       []any{},
			"sorts":                  sortMeta,
			"sectors":                sectors,
			"unsupported_conditions": unsupportedConditions,
			"data_quality": map[string]any{
				"source":              "company_shareholders + ownership_staging.v_explorer_*",
				"periodized_holdings": true,
				"t5_transactions":     true,
				"notes": []string{
					"Kondisi Big Player dihitung langsung dari view terindeks ownership_staging.",
					"Data mutasi majorholder dan akumulasi aktif tersambung.",
					"Filter sektor dan pasar berlaku untuk mode saham.",
				},
			},
		}
		h.storeCache(w, r, ck, resp, 1800)
		return
	}

	// mode == "stocks"
	cond, ok := stockConditionMap[condition]
	if !ok {
		httpx.WriteError(w, http.StatusBadRequest, "Kondisi saham tidak dikenal")
		return
	}
	sortFound := false
	for _, pair := range stockSortsOrdered {
		if pair[0] == sortKey {
			sortFound = true
			break
		}
	}
	if !sortFound {
		httpx.WriteError(w, http.StatusBadRequest, "Urutan saham tidak dikenal")
		return
	}

	args := []any{limit, offset}
	addArg := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf("$%d", len(args))
	}

	whereParts := []string{}
	if cond.Clause != "" {
		whereParts = append(whereParts, cond.Clause)
	}
	if hasTerm {
		p := addArg("%" + termTrim + "%")
		whereParts = append(whereParts, "(h.ticker ILIKE "+p+" OR t.name ILIKE "+p+")")
	}
	if strings.TrimSpace(sector) != "" {
		whereParts = append(whereParts, "LOWER(COALESCE(t.sector, '')) = LOWER("+addArg(strings.TrimSpace(sector))+")")
	}
	if minHolders != nil {
		whereParts = append(whereParts, "h.holder_count >= "+addArg(*minHolders))
	}
	if maxHolders != nil {
		whereParts = append(whereParts, "h.holder_count <= "+addArg(*maxHolders))
	}
	if minCtrl != nil {
		whereParts = append(whereParts, "h.controlling_pct >= "+addArg(*minCtrl))
	}
	if maxCtrl != nil {
		whereParts = append(whereParts, "h.controlling_pct <= "+addArg(*maxCtrl))
	}
	if minPub != nil {
		whereParts = append(whereParts, "h.public_pct >= "+addArg(*minPub))
	}
	if maxPub != nil {
		whereParts = append(whereParts, "h.public_pct <= "+addArg(*maxPub))
	}
	if minTop != nil {
		whereParts = append(whereParts, "h.max_holder_pct >= "+addArg(*minTop))
	}
	if maxTop != nil {
		whereParts = append(whereParts, "h.max_holder_pct <= "+addArg(*maxTop))
	}
	mktExpr := "(l.close * sh.listed_shares / 1000000000.0)"
	if minMktCap != nil {
		whereParts = append(whereParts, "l.close IS NOT NULL AND sh.listed_shares IS NOT NULL AND "+mktExpr+" >= "+addArg(*minMktCap))
	}
	if maxMktCap != nil {
		whereParts = append(whereParts, "l.close IS NOT NULL AND sh.listed_shares IS NOT NULL AND "+mktExpr+" <= "+addArg(*maxMktCap))
	}
	if mtd == "up" {
		whereParts = append(whereParts, "l.close IS NOT NULL AND m.close > 0 AND (l.close / m.close - 1) > 0")
	} else if mtd == "down" {
		whereParts = append(whereParts, "l.close IS NOT NULL AND m.close > 0 AND (l.close / m.close - 1) < 0")
	}
	if fFlow == "inflow" {
		whereParts = append(whereParts, "COALESCE(f.foreign_5d, 0) > 0")
	} else if fFlow == "outflow" {
		whereParts = append(whereParts, "COALESCE(f.foreign_5d, 0) < 0")
	}
	if hasTreasury {
		whereParts = append(whereParts, "h.treasury_pct > 0")
	}

	where := "TRUE"
	if len(whereParts) > 0 {
		where = strings.Join(whereParts, " AND ")
	}
	orderClause := "h.max_holder_pct DESC NULLS LAST, h.ticker ASC"
	if sortKey != "default" {
		orderClause = stockSortSQL[sortKey] + " " + strings.ToUpper(order) + " NULLS LAST, h.ticker ASC"
	}

	sqlStr := `
		WITH multi_holder AS (
			SELECT name
			FROM company_shareholders
			WHERE NOT (` + publicFilterSQL + `)
			  AND NOT (` + treasuryFilterSQL + `)
			GROUP BY name
			HAVING COUNT(DISTINCT ticker) >= 2
		),
		holder_agg AS (
			SELECT ticker,
			       COUNT(*)::int AS holder_count,
			       COALESCE(SUM(pct), 0)::double precision AS total_pct,
			       COALESCE(SUM(CASE WHEN is_controlling THEN pct ELSE 0 END), 0)::double precision AS controlling_pct,
			       COALESCE(SUM(CASE WHEN ` + publicFilterSQL + ` THEN pct ELSE 0 END), 0)::double precision AS public_pct,
			       COALESCE(SUM(CASE WHEN ` + treasuryFilterSQL + ` THEN pct ELSE 0 END), 0)::double precision AS treasury_pct,
			       COALESCE(SUM(CASE WHEN NOT is_controlling
			                         AND NOT (` + publicFilterSQL + `)
			                         AND NOT (` + treasuryFilterSQL + `)
			                         AND pct >= 5 THEN pct ELSE 0 END), 0)::double precision AS substantial_pct,
			       COALESCE(MAX(pct), 0)::double precision AS max_holder_pct,
			       SUM(CASE WHEN is_controlling THEN 1 ELSE 0 END)::int AS controlling_count,
			       SUM(CASE WHEN ` + publicFilterSQL + ` THEN 1 ELSE 0 END)::int AS public_count
			FROM company_shareholders
			GROUP BY ticker
		),
		latest_prices AS (
			SELECT DISTINCT ON (ticker) ticker, date, close
			FROM prices
			ORDER BY ticker, date DESC
		),
		month_prices AS (
			SELECT DISTINCT ON (ticker) ticker, date, close
			FROM prices
			WHERE date >= date_trunc('month', (SELECT MAX(date) FROM prices))
			ORDER BY ticker, date ASC
		),
		flow_latest AS (
			SELECT MAX(date) AS max_date FROM broker_flow
		),
		flow_agg AS (
			SELECT ticker, SUM(COALESCE(foreign_net_broker, 0))::double precision AS foreign_5d
			FROM broker_flow, flow_latest
			WHERE date >= flow_latest.max_date - INTERVAL '7 days'
			GROUP BY ticker
		)
		SELECT h.ticker,
		       COALESCE(t.name, h.ticker) AS name,
		       COALESCE(t.sector, '-') AS sector,
		       h.holder_count, h.total_pct, h.controlling_pct, h.public_pct,
		       h.treasury_pct, h.substantial_pct, h.max_holder_pct,
		       h.controlling_count, h.public_count,
		       l.close::double precision AS latest_close, l.date AS latest_date,
		       m.close::double precision AS month_close,
		       CASE WHEN l.close IS NOT NULL AND m.close > 0
		            THEN (((l.close / m.close) - 1) * 100)::double precision ELSE NULL END AS mtd_pct,
		       f.foreign_5d,
		       CASE WHEN l.close IS NOT NULL AND sh.listed_shares IS NOT NULL
		            THEN (l.close * sh.listed_shares / 1000000000.0)::double precision
		            ELSE NULL END AS market_cap_b,
		       COUNT(*) OVER()::int AS total_count
		FROM holder_agg h
		LEFT JOIN tickers t ON t.ticker = h.ticker
		LEFT JOIN latest_prices l ON l.ticker = h.ticker
		LEFT JOIN month_prices m ON m.ticker = h.ticker
		LEFT JOIN flow_agg f ON f.ticker = h.ticker
		LEFT JOIN company_shares sh ON sh.ticker = h.ticker
		WHERE ` + where + `
		ORDER BY ` + orderClause + `
		LIMIT $1 OFFSET $2
	`

	rows, qErr := h.DB.Query(ctx, sqlStr, args...)
	if qErr != nil {
		httpx.WriteError(w, http.StatusInternalServerError, qErr.Error())
		return
	}
	defer rows.Close()

	type stockRowTmp struct {
		ticker, name, sector                                         string
		hCnt, ctrlCnt, pubCnt                                        int
		totPct, ctrlPct, pubPct, treasPct, substPct, maxHolderPct    float64
		latestClose, monthClose, mtdPct, foreign5d, mktCapB          *float64
		latestDate                                                   *time.Time
	}
	rawRows := []stockRowTmp{}
	pageTickers := []string{}
	total := 0
	for rows.Next() {
		var rItem stockRowTmp
		var totCnt int
		if scanErr := rows.Scan(
			&rItem.ticker, &rItem.name, &rItem.sector,
			&rItem.hCnt, &rItem.totPct, &rItem.ctrlPct, &rItem.pubPct,
			&rItem.treasPct, &rItem.substPct, &rItem.maxHolderPct,
			&rItem.ctrlCnt, &rItem.pubCnt,
			&rItem.latestClose, &rItem.latestDate, &rItem.monthClose,
			&rItem.mtdPct, &rItem.foreign5d, &rItem.mktCapB, &totCnt,
		); scanErr != nil {
			httpx.WriteError(w, http.StatusInternalServerError, scanErr.Error())
			return
		}
		total = totCnt
		rawRows = append(rawRows, rItem)
		pageTickers = append(pageTickers, rItem.ticker)
	}
	if err := rows.Err(); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ctMetaMap := ctTickerMetaMap(ctx, h.DB, pageTickers)
	results := make([]map[string]any, 0, len(rawRows))
	for _, rItem := range rawRows {
		tm := ctMetaMap[rItem.ticker]
		badges := resolveBadges(tm.Badges, tm.Badge)
		results = append(results, map[string]any{
			"ticker":            rItem.ticker,
			"name":              rItem.name,
			"sector":            rItem.sector,
			"holder_count":      rItem.hCnt,
			"total_pct":         roundEven(rItem.totPct, 3),
			"controlling_pct":   roundEven(rItem.ctrlPct, 3),
			"public_pct":        roundEven(rItem.pubPct, 3),
			"treasury_pct":      roundEven(rItem.treasPct, 3),
			"substantial_pct":   roundEven(rItem.substPct, 3),
			"max_holder_pct":    roundEven(rItem.maxHolderPct, 3),
			"controlling_count": rItem.ctrlCnt,
			"public_count":      rItem.pubCnt,
			"price":             roundEvenPtr(rItem.latestClose, 4),
			"price_date":        isoDatePtr(rItem.latestDate),
			"mtd_pct":           roundEvenPtr(rItem.mtdPct, 2),
			"foreign_5d":        roundEvenPtr(rItem.foreign5d, 2),
			"market_cap_b":      roundEvenPtr(rItem.mktCapB, 2),
			"badges":            badges,
			"badge":             tm.Badge,
			"conglomerate":      tm.Konglo,
			"konglo":            tm.Konglo,
		})
	}

	var sharesAvail bool
	_ = h.DB.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM company_shares)`).Scan(&sharesAvail)

	condMeta := make([]map[string]any, 0, len(stockConditionsOrdered))
	for _, c := range stockConditionsOrdered {
		condMeta = append(condMeta, map[string]any{
			"id":        c.ID,
			"label":     c.Label,
			"available": true,
			"category":  c.Group,
		})
	}
	condGrpMeta := make([]map[string]any, 0, len(conditionGroupLabelsOrdered))
	for _, g := range conditionGroupLabelsOrdered {
		condGrpMeta = append(condGrpMeta, map[string]any{"id": g[0], "label": g[1]})
	}
	sortMeta := make([]map[string]any, 0, len(stockSortsOrdered))
	for _, s := range stockSortsOrdered {
		sortMeta = append(sortMeta, map[string]any{"id": s[0], "label": s[1]})
	}

	resp := map[string]any{
		"mode":                   "stocks",
		"condition":              condition,
		"condition_label":        cond.Label,
		"count":                  len(results),
		"total":                  total,
		"limit":                  limit,
		"offset":                 offset,
		"results":                results,
		"conditions":             condMeta,
		"condition_groups":       condGrpMeta,
		"sorts":                  sortMeta,
		"sectors":                sectors,
		"unsupported_conditions": unsupportedConditions,
		"data_quality": map[string]any{
			"source":              "company_shareholders + ownership_staging + prices + broker_flow + company_shares",
			"periodized_holdings": true,
			"t5_transactions":     true,
			"market_cap":          sharesAvail,
			"notes": []string{
				"Kondisi Big Player dihitung langsung dari view SQL terindeks ownership_staging.",
				"Data mutasi majorholder dan akumulasi pemegang saham aktif tersambung.",
				"Market cap = harga terakhir x listed shares IDX (miliar Rupiah); null jika shares belum tersedia.",
				"Filter dan sorting dijalankan di server; hasil mengikuti snapshot terbaru.",
			},
		},
	}
	h.storeCache(w, r, ck, resp, 1800)
}

func canonicalQueryString(q url.Values) string {
	if len(q) == 0 {
		return ""
	}
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		vals := append([]string(nil), q[k]...)
		sort.Strings(vals)
		parts = append(parts, k+"="+strings.Join(vals, ","))
	}
	return strings.Join(parts, "&")
}
