package ownership

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) groupPayload(ctx context.Context, groupID string, group shared.ConglomerateGroup) (map[string]any, error) {
	tickers := make([]string, 0, len(group.KnownTickers))
	for _, tk := range group.KnownTickers {
		tickers = append(tickers, strings.ToUpper(tk))
	}

	metaMap, err := tickerMeta(ctx, h.DB, tickers)
	if err != nil {
		return nil, err
	}
	prices, err := priceSnapshot(ctx, h.DB, tickers)
	if err != nil {
		return nil, err
	}
	holders, err := holderAggregates(ctx, h.DB, tickers)
	if err != nil {
		return nil, err
	}
	monthly, err := monthlySnapshot(ctx, h.DB, tickers)
	if err != nil {
		return nil, err
	}
	sharesMap, err := companySharesMap(ctx, h.DB, tickers)
	if err != nil {
		return nil, err
	}
	ctMetaMap := ctTickerMetaMap(ctx, h.DB, tickers)
	ctLatestPct := ctLatestPctMap(ctx, h.DB, tickers)
	ctHolderAgg := ctHolderAggregates(ctx, h.DB, tickers)
	ctSIDMap := ctLatestSIDMap(ctx, h.DB, tickers)

	constituents := make([]map[string]any, 0, len(tickers))
	for _, tk := range tickers {
		pr, hasPr := prices[tk]
		ho := holders[tk]
		mo, hasMo := monthly[tk]
		ctm := ctMetaMap[tk]
		pct, hasPct := ctLatestPct[tk]
		ctH, hasCtH := ctHolderAgg[tk]

		var pricePtr, mtdPtr, change1DPtr *float64
		var priceDatePtr *string
		if hasPr {
			cl := pr.Close
			pricePtr = &cl
			priceDatePtr = pr.Date
			change1DPtr = pr.Change1DPct
			mtdPtr = pr.MTDPct
		}

		marketCapB := ctm.MktcapB
		if marketCapB == nil && pricePtr != nil && *pricePtr != 0 {
			if ls, ok := sharesMap[tk]; ok && ls != 0 {
				v := (*pricePtr) * ls / 1_000_000_000.0
				marketCapB = &v
			}
		}

		freeFloatEst := ho.FreeFloatEst
		if freeFloatEst == 0 && ctm.FreeFloat != nil {
			freeFloatEst = roundEven(*ctm.FreeFloat*100.0, 2)
		}

		var bigPlayerPct *float64
		if hasCtH {
			v := roundEven(ctH.TotalPct, 2)
			bigPlayerPct = &v
		}

		badges := resolveBadges(ctm.Badges, ctm.Badge)
		var cong any = group.Name
		if group.Name == "" {
			cong = ctm.Konglo
		}
		cName := tk
		if m, ok := metaMap[tk]; ok && m.Name != "" {
			cName = m.Name
		} else if ctm.Issuer != nil && *ctm.Issuer != "" {
			cName = *ctm.Issuer
		}
		cSec := "-"
		if m, ok := metaMap[tk]; ok && m.Sector != "" {
			cSec = m.Sector
		}
		var monthlyObj any
		if hasMo {
			monthlyObj = MonthlySliceBrief{
				PeriodLabel:       mo.PeriodLabel,
				TotalSID:          mo.TotalSID,
				RetailPct:         mo.RetailPct,
				ForeignNamedPct:   mo.ForeignNamedPct,
				ForeignHiddenPct:  mo.ForeignHiddenPct,
				LocalNonRetailPct: mo.LocalNonRetailPct,
			}
		}
		var sidPtr *int
		if sVal, ok := ctSIDMap[tk]; ok {
			sidPtr = &sVal
		}
		var foreignTotalPct, foreignHiddenPct, retailPct *float64
		if hasPct {
			foreignTotalPct = roundEvenPtr(pct.PctAsing, 2)
			foreignHiddenPct = roundEvenPtr(pct.PctAsingLt1, 2)
			retailPct = roundEvenPtr(pct.PctRitel, 2)
		}

		constituents = append(constituents, map[string]any{
			"ticker":             tk,
			"name":               cName,
			"sector":             cSec,
			"badges":             badges,
			"badge":              ctm.Badge,
			"conglomerate":       cong,
			"konglo":             cong,
			"price":              pricePtr,
			"price_date":         priceDatePtr,
			"change_1d_pct":      change1DPtr,
			"mtd_pct":            mtdPtr,
			"holder_count":       ho.HolderCount,
			"controlling_pct":    ho.ControllingPct,
			"public_pct":         ho.PublicPct,
			"max_holder_pct":     ho.MaxHolderPct,
			"free_float_est":     freeFloatEst,
			"market_cap_b":       roundEvenPtr(marketCapB, 2),
			"big_player_pct":     bigPlayerPct,
			"foreign_total_pct":  foreignTotalPct,
			"foreign_hidden_pct": foreignHiddenPct,
			"retail_pct":         retailPct,
			"sid":                sidPtr,
			"monthly":            monthlyObj,
		})
	}

	hasMktCap := false
	for _, c := range constituents {
		if p, _ := c["market_cap_b"].(*float64); p != nil {
			hasMktCap = true
			break
		}
	}
	if hasMktCap {
		sort.SliceStable(constituents, func(i, j int) bool {
			pi, _ := constituents[i]["market_cap_b"].(*float64)
			pj, _ := constituents[j]["market_cap_b"].(*float64)
			ni, nj := pi == nil, pj == nil
			if ni != nj {
				return !ni
			}
			vi, vj := -derefFloat(pi), -derefFloat(pj)
			if vi != vj {
				return vi < vj
			}
			return constituents[i]["ticker"].(string) < constituents[j]["ticker"].(string)
		})
	} else {
		sort.SliceStable(constituents, func(i, j int) bool {
			pi, _ := constituents[i]["mtd_pct"].(*float64)
			pj, _ := constituents[j]["mtd_pct"].(*float64)
			ni, nj := pi == nil, pj == nil
			if ni != nj {
				return !ni
			}
			vi, vj := -derefFloat(pi), -derefFloat(pj)
			if vi != vj {
				return vi < vj
			}
			return constituents[i]["ticker"].(string) < constituents[j]["ticker"].(string)
		})
	}

	var mtdVals, pubVals, ctrlVals, mktCapVals []float64
	var sidSum, ctSidSum, holderRowsSum, coveredCount, monthlyCoveredCount int
	var hasMonthlySID, hasCTSID bool

	for _, c := range constituents {
		hc, _ := c["holder_count"].(int)
		bp, _ := c["big_player_pct"].(*float64)
		isCovered := hc > 0 || bp != nil
		if isCovered {
			coveredCount++
			holderRowsSum += hc
			pubVals = append(pubVals, c["public_pct"].(float64))
			ctrlVals = append(ctrlVals, c["controlling_pct"].(float64))
		}
		if mp, _ := c["mtd_pct"].(*float64); mp != nil {
			mtdVals = append(mtdVals, *mp)
		}
		if m, ok := c["monthly"].(MonthlySliceBrief); ok {
			monthlyCoveredCount++
			hasMonthlySID = true
			sidSum += m.TotalSID
		}
		if sPtr, _ := c["sid"].(*int); sPtr != nil {
			hasCTSID = true
			ctSidSum += *sPtr
		}
		if mcPtr, _ := c["market_cap_b"].(*float64); mcPtr != nil {
			mktCapVals = append(mktCapVals, *mcPtr)
		}
	}

	var avgMTD, pubAvg, ctrlAvg, mktCapTotal *float64
	var advancers, decliners int
	if len(mtdVals) > 0 {
		sum := 0.0
		for _, v := range mtdVals {
			sum += v
			if v > 0 {
				advancers++
			} else if v < 0 {
				decliners++
			}
		}
		v := roundEven(sum/float64(len(mtdVals)), 2)
		avgMTD = &v
	}
	if len(pubVals) > 0 {
		sum := 0.0
		for _, v := range pubVals {
			sum += v
		}
		v := roundEven(sum/float64(len(pubVals)), 2)
		pubAvg = &v
	}
	if len(ctrlVals) > 0 {
		sum := 0.0
		for _, v := range ctrlVals {
			sum += v
		}
		v := roundEven(sum/float64(len(ctrlVals)), 2)
		ctrlAvg = &v
	}
	if len(mktCapVals) > 0 {
		sum := 0.0
		for _, v := range mktCapVals {
			sum += v
		}
		v := roundEven(sum, 2)
		mktCapTotal = &v
	}
	var sidTotalCovered, ctSidTotal *int
	if hasMonthlySID {
		sidTotalCovered = &sidSum
	}
	if hasCTSID {
		ctSidTotal = &ctSidSum
	}

	withMTD := make([]map[string]any, 0, len(constituents))
	for _, c := range constituents {
		if p, _ := c["mtd_pct"].(*float64); p != nil {
			withMTD = append(withMTD, c)
		}
	}
	sort.SliceStable(withMTD, func(i, j int) bool {
		pi := withMTD[i]["mtd_pct"].(*float64)
		pj := withMTD[j]["mtd_pct"].(*float64)
		return *pi > *pj
	})
	if len(withMTD) > 2 {
		withMTD = withMTD[:2]
	}
	topStocks := make([]map[string]any, 0, len(withMTD))
	for _, c := range withMTD {
		topStocks = append(topStocks, map[string]any{
			"ticker":  c["ticker"],
			"name":    c["name"],
			"mtd_pct": c["mtd_pct"],
		})
	}

	hasCT := len(ctMetaMap) > 0 || len(ctLatestPct) > 0 || len(ctHolderAgg) > 0 || len(ctSIDMap) > 0
	name := group.Name
	if name == "" {
		name = groupID
	}

	return map[string]any{
		"id":                   groupID,
		"name":                 name,
		"description":          group.Description,
		"hub_entity":           group.HubEntity,
		"logo_url":             group.LogoURL,
		"catalog_code":         group.CatalogCode,
		"ticker_count":         len(tickers),
		"covered_tickers":      coveredCount,
		"avg_mtd_pct":          avgMTD,
		"advancers":            advancers,
		"decliners":            decliners,
		"public_pct_avg":       pubAvg,
		"controlling_pct_avg":  ctrlAvg,
		"holder_rows":          holderRowsSum,
		"monthly_covered":      monthlyCoveredCount,
		"sid_total_covered":    sidTotalCovered,
		"sid_total":            ctSidTotal,
		"top_stocks":           topStocks,
		"market_cap":           mktCapTotal,
		"market_cap_b":         mktCapTotal,
		"market_cap_available": len(mktCapVals) > 0,
		"market_cap_unit":      "miliar_idr",
		"ct_baseline":          hasCT,
		"constituents":         constituents,
	}, nil
}

// GetGroups handles GET /api/v1/ownership/groups.
func (h *Handler) GetGroups(w http.ResponseWriter, r *http.Request) {
	ck := ownershipCacheKey("ownership_groups_v2", "default")
	if h.tryCache(w, r, ck, 7200) {
		return
	}
	ctx := r.Context()

	ordered, _ := shared.LoadConglomerateGroups(ctx, h.DB)
	usingCT := len(ordered) > 0
	if !usingCT {
		ordered, _ = shared.StaticConglomerateGroups()
	}

	groups := make([]map[string]any, 0, len(ordered))
	allTickers := map[string]bool{}
	for _, g := range ordered {
		payload, err := h.groupPayload(ctx, g.ID, g)
		if err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		groups = append(groups, payload)
		for _, tk := range g.KnownTickers {
			allTickers[strings.ToUpper(tk)] = true
		}
	}

	sort.SliceStable(groups, func(i, j int) bool {
		pi, _ := groups[i]["avg_mtd_pct"].(*float64)
		pj, _ := groups[j]["avg_mtd_pct"].(*float64)
		ni, nj := pi == nil, pj == nil
		if ni != nj {
			return !ni
		}
		vi, vj := -derefFloat(pi), -derefFloat(pj)
		if vi != vj {
			return vi < vj
		}
		return groups[i]["name"].(string) < groups[j]["name"].(string)
	})

	var latestPriceDate *time.Time
	if err := h.DB.QueryRow(ctx, `SELECT MAX(date) FROM prices`).Scan(&latestPriceDate); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var totalTickers int
	if err := h.DB.QueryRow(ctx, `SELECT COUNT(*)::int FROM tickers`).Scan(&totalTickers); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	mktCapAvail := false
	for _, g := range groups {
		if b, _ := g["market_cap_available"].(bool); b {
			mktCapAvail = true
			break
		}
	}
	t5Avail := len(ctT5Rows(ctx, h.DB, nil, 1)) > 0

	covNote := "Grup yang tersedia adalah preset yang sudah dipetakan di backend. Market cap, SID, dan transaksi T5 belum tersedia."
	noteCount := fmt.Sprintf("Jumlah grup mengikuti preset CONGLOMERATE_GROUPS (%d grup).", len(groups))
	if usingCT {
		covNote = "Grup, market cap, SID, dan transaksi T5 berasal dari baseline CT yang diisolasi di schema ct_baseline. Baseline ini tidak dipakai di jalur sinyal."
		noteCount = fmt.Sprintf("Jumlah grup mengikuti baseline CT (%d grup).", len(groups))
	}

	topEnd := 10
	if len(groups) < topEnd {
		topEnd = len(groups)
	}

	resp := map[string]any{
		"as_of": isoDatePtr(latestPriceDate),
		"coverage": map[string]any{
			"mapped_groups":        len(groups),
			"mapped_tickers":       len(allTickers),
			"total_tickers":        totalTickers,
			"market_cap_available": mktCapAvail,
			"t5_available":         t5Avail,
			"ct_baseline":          usingCT,
			"note":                 covNote,
		},
		"top_groups": groups[:topEnd],
		"groups":     groups,
		"notes": []string{
			"Peringkat menggunakan rata-rata sederhana perubahan harga MTD, bukan weighted market cap.",
			"Emiten tanpa harga bulan berjalan tidak masuk rata-rata.",
			noteCount,
		},
	}
	h.storeCache(w, r, ck, resp, 7200)
}

// GetGroupDetail handles GET /api/v1/ownership/groups/{group_id}.
func (h *Handler) GetGroupDetail(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "group_id")
	key := strings.ToLower(strings.TrimSpace(rawID))
	ck := ownershipCacheKey("ownership_group_detail_v3", rawID)
	if h.tryCache(w, r, ck, 7200) {
		return
	}

	ctx := r.Context()
	group, found := shared.GetConglomerateGroup(ctx, h.DB, key)
	if !found {
		_, staticByID := shared.StaticConglomerateGroups()
		if st, ok := staticByID[key]; ok {
			group = &st
			found = true
		}
	}
	if !found || group == nil {
		httpx.WriteError(w, http.StatusNotFound, "Grup konglomerasi tidak ditemukan")
		return
	}

	canonicalKey := group.ID
	if canonicalKey == "" {
		canonicalKey = key
	}
	payload, err := h.groupPayload(ctx, canonicalKey, *group)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	constituents := payload["constituents"].([]map[string]any)
	tickers := make([]string, 0, len(constituents))
	for _, c := range constituents {
		tickers = append(tickers, c["ticker"].(string))
	}

	topHoldings := []map[string]any{}
	if len(tickers) > 0 {
		q := `
			SELECT name,
			       COUNT(DISTINCT ticker)::int AS ticker_count,
			       COALESCE(SUM(pct), 0)::double precision AS total_pct,
			       json_agg(
			           json_build_object(
			               'ticker', ticker,
			               'pct', pct::double precision,
			               'is_controlling', is_controlling
			           ) ORDER BY pct DESC
			       ) AS holdings
			FROM company_shareholders
			WHERE ticker = ANY($1)
			  AND NOT (` + publicFilterSQL + `)
			  AND NOT (` + treasuryFilterSQL + `)
			  AND name NOT IN ('-', 'Pengendali Saham', 'Non Pengendali Saham', 'Others', 'Lain-lain')
			GROUP BY name
			HAVING COUNT(DISTINCT ticker) >= 2
			ORDER BY ticker_count DESC, total_pct DESC, name ASC
			LIMIT 25
		`
		rows, qErr := h.DB.Query(ctx, q, tickers)
		if qErr != nil {
			httpx.WriteError(w, http.StatusInternalServerError, qErr.Error())
			return
		}
		for rows.Next() {
			var nm string
			var tkCnt int
			var totPct float64
			var hRaw []byte
			if scanErr := rows.Scan(&nm, &tkCnt, &totPct, &hRaw); scanErr != nil {
				rows.Close()
				httpx.WriteError(w, http.StatusInternalServerError, scanErr.Error())
				return
			}
			var rawItems []map[string]any
			if len(hRaw) > 0 {
				_ = json.Unmarshal(hRaw, &rawItems)
			}
			cleanItems := make([]map[string]any, 0, len(rawItems))
			for _, it := range rawItems {
				tkStr, _ := it["ticker"].(string)
				isCtrl, _ := it["is_controlling"].(bool)
				cleanItems = append(cleanItems, map[string]any{
					"ticker":         tkStr,
					"pct":            roundEven(asFloat(it["pct"]), 3),
					"is_controlling": isCtrl,
				})
			}
			topHoldings = append(topHoldings, map[string]any{
				"name":         nm,
				"ticker_count": tkCnt,
				"total_pct":    roundEven(totPct, 3),
				"holdings":     cleanItems,
			})
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	transactions := ctT5Rows(ctx, h.DB, tickers, 300)
	t5Sum := ctT5Summary(transactions)
	hasCT, _ := payload["ct_baseline"].(bool)
	mktAvail, _ := payload["market_cap_available"].(bool)
	monCovered, _ := payload["monthly_covered"].(int)

	notes := []string{
		"Tabel emiten memakai snapshot pemegang saham dan harga terbaru.",
		"Jejaring kepemilikan memakai endpoint network-graph existing.",
	}
	if mktAvail {
		notes = append(notes, "Market cap dihitung dari baseline CT atau harga terakhir x listed shares IDX.")
	} else {
		notes = append(notes, "Market cap belum tersedia untuk grup ini.")
	}
	if len(transactions) > 0 {
		notes = append(notes, "Transaksi T5 berasal dari baseline CT (quarantined), bukan data sinyal.")
	} else {
		notes = append(notes, "Transaksi T5 belum tersedia untuk emiten grup ini.")
	}

	sidTotalPtr, _ := payload["sid_total"].(*int)
	resp := map[string]any{
		"group": map[string]any{
			"id":                   payload["id"],
			"name":                 payload["name"],
			"description":          payload["description"],
			"hub_entity":           payload["hub_entity"],
			"logo_url":             payload["logo_url"],
			"catalog_code":         payload["catalog_code"],
			"ticker_count":         payload["ticker_count"],
			"covered_tickers":      payload["covered_tickers"],
			"avg_mtd_pct":          payload["avg_mtd_pct"],
			"advancers":            payload["advancers"],
			"decliners":            payload["decliners"],
			"public_pct_avg":       payload["public_pct_avg"],
			"controlling_pct_avg":  payload["controlling_pct_avg"],
			"holder_rows":          payload["holder_rows"],
			"monthly_covered":      payload["monthly_covered"],
			"sid_total_covered":    payload["sid_total_covered"],
			"sid_total":            payload["sid_total"],
			"top_stocks":           payload["top_stocks"],
			"market_cap":           payload["market_cap"],
			"market_cap_b":         payload["market_cap_b"],
			"market_cap_available": payload["market_cap_available"],
			"market_cap_unit":      payload["market_cap_unit"],
			"ct_baseline":          hasCT,
		},
		"constituents": constituents,
		"top_holdings": topHoldings,
		"transactions": transactions,
		"t5_summary":   t5Sum,
		"data_quality": map[string]any{
			"shareholder_snapshot": true,
			"price_mtd":            true,
			"monthly_composition":  monCovered > 0,
			"market_cap":           mktAvail,
			"t5_transactions":      len(transactions) > 0,
			"sid_history":          sidTotalPtr != nil,
			"ct_baseline":          hasCT,
			"notes":                notes,
		},
	}
	h.storeCache(w, r, ck, resp, 7200)
}

// SyncGroupsCatalog handles POST /api/v1/ownership/groups/sync-catalog.
func (h *Handler) SyncGroupsCatalog(w http.ResponseWriter, r *http.Request) {
	if h.SyncCatalogRunner != nil {
		res, err := h.SyncCatalogRunner(r)
		if err != nil {
			httpx.WriteJSON(w, http.StatusOK, map[string]any{
				"status":          "error",
				"message":         err.Error(),
				"catalogs_synced": 0,
				"members_synced":  0,
				"unique_tickers":  0,
			})
			return
		}
		if h.Cache != nil {
			for _, pat := range []string{"ownership_groups*", "ownership_group_detail*", "shareholders*"} {
				h.Cache.SafePurge(r.Context(), pat)
			}
		}
		httpx.WriteJSON(w, http.StatusOK, res)
		return
	}

	smDir := "/root/SMtracker"
	pyBin := "/root/SMtracker/venv/bin/python3"
	if _, err := os.Stat(pyBin); err != nil {
		pyBin = "python3"
	}
	cmd := exec.CommandContext(r.Context(), pyBin, "-c",
		"import json; from sqlalchemy.orm import Session; from idx_bandarmology import storage; from idx_bandarmology.stockbit_catalog import sync_conglomerate_catalogs; db=Session(storage.engine); print(json.dumps(sync_conglomerate_catalogs(db))); db.close()",
	)
	cmd.Dir = smDir
	cmd.Env = append(os.Environ(), "PYTHONPATH="+smDir+"/src:"+smDir)
	out, err := cmd.Output()
	if err != nil {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{
			"status":          "error",
			"message":         "Tidak dapat mengambil katalog dari Stockbit: " + err.Error(),
			"catalogs_synced": 0,
			"members_synced":  0,
			"unique_tickers":  0,
		})
		return
	}
	var parsed map[string]any
	if jErr := json.Unmarshal([]byte(strings.TrimSpace(string(out))), &parsed); jErr != nil {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{
			"status":          "error",
			"message":         "Output sinkronisasi tidak valid",
			"catalogs_synced": 0,
			"members_synced":  0,
			"unique_tickers":  0,
		})
		return
	}
	if h.Cache != nil {
		for _, pat := range []string{"ownership_groups*", "ownership_group_detail*", "shareholders*"} {
			h.Cache.SafePurge(r.Context(), pat)
		}
	}
	httpx.WriteJSON(w, http.StatusOK, parsed)
}
