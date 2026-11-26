package ownership

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
)

// SearchInvestors handles GET /api/v1/ownership/investors.
func (h *Handler) SearchInvestors(w http.ResponseWriter, r *http.Request) {
	qMap := r.URL.Query()
	queryStr := qMap.Get("q")
	limit := 20
	if raw := strings.TrimSpace(qMap.Get("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 100 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "limit must be between 1 and 100")
			return
		}
		limit = n
	}
	includePublic := false
	if raw := strings.TrimSpace(qMap.Get("include_public")); raw != "" {
		b, err := strconv.ParseBool(raw)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "include_public must be boolean")
			return
		}
		includePublic = b
	}

	ck := ownershipCacheKey("ownership_investors", fmt.Sprintf("%s|%d|%t", queryStr, limit, includePublic))
	if h.tryCache(w, r, ck, 1800) {
		return
	}

	filters := []string{}
	args := []any{}
	if !includePublic {
		filters = append(filters, "NOT ("+publicFilterSQL+")")
		filters = append(filters, "NOT ("+treasuryFilterSQL+")")
		filters = append(filters, "name NOT IN ('-', 'Pengendali Saham', 'Non Pengendali Saham', 'Others', 'Lain-lain')")
	}
	if strings.TrimSpace(queryStr) != "" {
		args = append(args, "%"+strings.TrimSpace(queryStr)+"%")
		filters = append(filters, fmt.Sprintf("name ILIKE $%d", len(args)))
	}
	where := "TRUE"
	if len(filters) > 0 {
		where = strings.Join(filters, " AND ")
	}
	args = append(args, limit)
	limitPlaceholder := fmt.Sprintf("$%d", len(args))

	sqlStr := `
		SELECT name,
		       COUNT(DISTINCT ticker)::int AS ticker_count,
		       COALESCE(SUM(pct), 0)::double precision AS total_pct,
		       COALESCE(MAX(pct), 0)::double precision AS max_pct,
		       SUM(CASE WHEN is_controlling THEN 1 ELSE 0 END)::int AS controlling_count
		FROM company_shareholders
		WHERE ` + where + `
		GROUP BY name
		ORDER BY ticker_count DESC, total_pct DESC, name ASC
		LIMIT ` + limitPlaceholder

	rows, err := h.DB.Query(r.Context(), sqlStr, args...)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	results := []map[string]any{}
	for rows.Next() {
		var nm string
		var tkCount, ctrlCount int
		var totalPct, maxPct float64
		if scanErr := rows.Scan(&nm, &tkCount, &totalPct, &maxPct, &ctrlCount); scanErr != nil {
			httpx.WriteError(w, http.StatusInternalServerError, scanErr.Error())
			return
		}
		results = append(results, map[string]any{
			"name":              nm,
			"ticker_count":      tkCount,
			"total_pct":         roundEven(totalPct, 3),
			"max_pct":           roundEven(maxPct, 3),
			"controlling_count": ctrlCount,
		})
	}
	if err := rows.Err(); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := map[string]any{
		"query":   queryStr,
		"count":   len(results),
		"results": results,
	}
	h.storeCache(w, r, ck, resp, 1800)
}

// GetInvestorDetail handles GET /api/v1/ownership/investors/detail.
func (h *Handler) GetInvestorDetail(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) < 2 {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "name is required and must be at least 2 characters")
		return
	}
	trimmedName := strings.TrimSpace(name)
	ck := ownershipCacheKey("ownership_investor_detail_v2", name)
	if h.tryCache(w, r, ck, 1800) {
		return
	}

	ctx := r.Context()
	rows, err := h.DB.Query(ctx, `
		SELECT cs.ticker, COALESCE(t.name, cs.ticker), COALESCE(t.sector, '-'),
		       cs.pct::double precision, cs.is_controlling
		FROM company_shareholders cs
		LEFT JOIN tickers t ON t.ticker = cs.ticker
		WHERE LOWER(cs.name) = LOWER($1)
		ORDER BY cs.pct DESC, cs.ticker ASC
	`, trimmedName)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	type invRawRow struct {
		ticker        string
		companyName   string
		sector        string
		pct           float64
		isControlling bool
	}
	rawRows := []invRawRow{}
	tickers := []string{}
	for rows.Next() {
		var rItem invRawRow
		if scanErr := rows.Scan(&rItem.ticker, &rItem.companyName, &rItem.sector, &rItem.pct, &rItem.isControlling); scanErr != nil {
			httpx.WriteError(w, http.StatusInternalServerError, scanErr.Error())
			return
		}
		rawRows = append(rawRows, rItem)
		tickers = append(tickers, rItem.ticker)
	}
	if err := rows.Err(); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(rawRows) == 0 {
		httpx.WriteError(w, http.StatusNotFound, "Investor tidak ditemukan")
		return
	}

	prices, err := priceSnapshot(ctx, h.DB, tickers)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	monthly, err := monthlySnapshot(ctx, h.DB, tickers)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ctMeta := ctTickerMetaMap(ctx, h.DB, tickers)

	holdings := make([]map[string]any, 0, len(rawRows))
	sectorSet := map[string]bool{}
	var totalPct, maxPct float64
	var ctrlCount, substCount, pubCount int

	for _, rr := range rawRows {
		pctRounded := roundEven(rr.pct, 4)
		cat := classifyHolder(name, rr.pct, rr.isControlling)
		totalPct += pctRounded
		if pctRounded > maxPct {
			maxPct = pctRounded
		}
		if rr.isControlling {
			ctrlCount++
		}
		if cat == "substantial" {
			substCount++
		} else if cat == "public" {
			pubCount++
		}
		if rr.sector != "" && rr.sector != "-" {
			sectorSet[rr.sector] = true
		}

		var priceObj any
		if p, ok := prices[rr.ticker]; ok {
			priceObj = p
		}
		var monthlyObj any
		if m, ok := monthly[rr.ticker]; ok {
			monthlyObj = MonthlySliceBrief{
				PeriodLabel:       m.PeriodLabel,
				TotalSID:          m.TotalSID,
				RetailPct:         m.RetailPct,
				ForeignNamedPct:   m.ForeignNamedPct,
				ForeignHiddenPct:  m.ForeignHiddenPct,
				LocalNonRetailPct: m.LocalNonRetailPct,
			}
		}
		tm := ctMeta[rr.ticker]
		badges := resolveBadges(tm.Badges, tm.Badge)
		holdings = append(holdings, map[string]any{
			"ticker":         rr.ticker,
			"name":           rr.companyName,
			"sector":         rr.sector,
			"pct":            pctRounded,
			"is_controlling": rr.isControlling,
			"category":       cat,
			"price":          priceObj,
			"badges":         badges,
			"badge":          tm.Badge,
			"conglomerate":   tm.Konglo,
			"konglo":         tm.Konglo,
			"monthly":        monthlyObj,
		})
	}

	sectors := make([]string, 0, len(sectorSet))
	for s := range sectorSet {
		sectors = append(sectors, s)
	}
	sort.Strings(sectors)

	resp := map[string]any{
		"name": trimmedName,
		"summary": map[string]any{
			"ticker_count":      len(holdings),
			"total_pct":         roundEven(totalPct, 3),
			"max_pct":           roundEven(maxPct, 3),
			"controlling_count": ctrlCount,
			"substantial_count": substCount,
			"public_count":      pubCount,
			"sectors":           sectors,
		},
		"holdings":     holdings,
		"transactions": []any{},
		"data_quality": map[string]any{
			"periodized_holdings": false,
			"t5_transactions":     false,
			"gl":                  false,
			"market_cap":          false,
			"notes": []string{
				"Portofolio dihitung dari snapshot pemegang saham saat ini.",
				"Nilai rupiah, G/L, dan transaksi T5 belum tersedia.",
			},
		},
	}
	h.storeCache(w, r, ck, resp, 1800)
}
