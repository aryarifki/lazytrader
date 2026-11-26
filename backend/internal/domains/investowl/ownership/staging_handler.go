package ownership

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/go-chi/chi/v5"
)

func floatToInt64Ptr(v *float64) *int64 {
	if v == nil {
		return nil
	}
	n := int64(*v)
	return &n
}

// GetStagingOverview handles GET /api/v1/ownership/staging/overview.
func (h *Handler) GetStagingOverview(w http.ResponseWriter, r *http.Request) {
	ck := ownershipCacheKey("ownership_staging_overview_v1", "default")
	if h.tryCache(w, r, ck, 300) {
		return
	}

	var kseiPeriods, kseiItems, statisItems, idxReports, idxEmitens, t5Movements int
	var latestKsei, latestIdx *time.Time
	err := h.DB.QueryRow(r.Context(), `
		SELECT
			(SELECT COUNT(*)::int FROM ownership_staging.ksei_balancepos_periods) AS ksei_periods,
			(SELECT COUNT(*)::int FROM ownership_staging.ksei_balancepos_items) AS ksei_items,
			(SELECT COUNT(*)::int FROM ownership_staging.ksei_statis_items) AS statis_items,
			(SELECT COUNT(*)::int FROM ownership_staging.idx_shareholder_reports) AS idx_reports,
			(SELECT COUNT(DISTINCT ticker)::int FROM ownership_staging.idx_shareholder_reports) AS idx_emitens,
			(SELECT COUNT(*)::int FROM ownership_staging.majorholder_movements WHERE broker_code IS NOT NULL AND broker_code != '') AS t5_movements,
			(SELECT MAX(period_date) FROM ownership_staging.ksei_balancepos_periods) AS latest_ksei_period,
			(SELECT MAX(period_date) FROM ownership_staging.idx_shareholder_reports) AS latest_idx_period
	`).Scan(&kseiPeriods, &kseiItems, &statisItems, &idxReports, &idxEmitens, &t5Movements, &latestKsei, &latestIdx)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := map[string]any{
		"status": "ready",
		"schema": "ownership_staging",
		"ksei": map[string]any{
			"periods_count":          kseiPeriods,
			"balancepos_items_count": kseiItems,
			"statis_items_count":     statisItems,
			"latest_period":          isoDatePtr(latestKsei),
		},
		"idx_bae": map[string]any{
			"reports_count":  idxReports,
			"unique_emitens": idxEmitens,
			"latest_period":  isoDatePtr(latestIdx),
		},
		"t5_transactions_count": t5Movements,
	}
	h.storeCache(w, r, ck, resp, 300)
}

// GetStagingStockProfile handles GET /api/v1/ownership/staging/stock/{ticker}.
func (h *Handler) GetStagingStockProfile(w http.ResponseWriter, r *http.Request) {
	tickerClean := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	if tickerClean == "" {
		httpx.WriteError(w, http.StatusBadRequest, "Ticker tidak valid")
		return
	}
	ck := ownershipCacheKey("ownership_staging_stock_v1", tickerClean)
	if h.tryCache(w, r, ck, 600) {
		return
	}
	ctx := r.Context()

	compRows, err := h.DB.Query(ctx, `
		SELECT
			period_date, ticker, total_shares::double precision, total_scriptless::double precision,
			pct_scriptless::double precision, pct_asing::double precision, pct_asing_gt1::double precision,
			pct_asing_lt1::double precision, pct_ritel_ksei_id::double precision, pct_individual_gt1::double precision,
			pct_ritel_pure::double precision,
			pct_local_corporate::double precision, pct_local_insurance::double precision,
			pct_local_pension::double precision, pct_local_mutual_fund::double precision,
			pct_local_bank::double precision, pct_local_securities::double precision,
			pct_local_foundation::double precision, pct_local_other::double precision,
			pct_local_institution::double precision,
			pct_foreign_corporate::double precision, pct_foreign_insurance::double precision,
			pct_foreign_pension::double precision, pct_foreign_bank::double precision,
			pct_foreign_mutual_fund::double precision, pct_foreign_securities::double precision,
			pct_foreign_foundation::double precision, pct_foreign_other::double precision,
			pct_foreign_institution::double precision
		FROM ownership_staging.v_composition_ksei_normalized
		WHERE ticker = $1
		ORDER BY period_date DESC
	`, tickerClean)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer compRows.Close()

	compList := []map[string]any{}
	for compRows.Next() {
		var pDate *time.Time
		var tk string
		var totSh, totSL, pctSL, pctAsing, pctAsingGT1, pctAsingLT1, pctRitelID, pctIndivGT1, pctRitelPure *float64
		var lCorp, lIns, lPen, lMF, lBank, lSec, lFound, lOther, lInst *float64
		var fCorp, fIns, fPen, fBank, fMF, fSec, fFound, fOther, fInst *float64
		if scanErr := compRows.Scan(
			&pDate, &tk, &totSh, &totSL, &pctSL, &pctAsing, &pctAsingGT1, &pctAsingLT1,
			&pctRitelID, &pctIndivGT1, &pctRitelPure,
			&lCorp, &lIns, &lPen, &lMF, &lBank, &lSec, &lFound, &lOther, &lInst,
			&fCorp, &fIns, &fPen, &fBank, &fMF, &fSec, &fFound, &fOther, &fInst,
		); scanErr != nil {
			httpx.WriteError(w, http.StatusInternalServerError, scanErr.Error())
			return
		}
		compList = append(compList, map[string]any{
			"period_date":        isoStr(pDate),
			"total_shares":       floatToInt64Ptr(totSh),
			"total_scriptless":   floatToInt64Ptr(totSL),
			"pct_scriptless":     pctSL,
			"pct_asing":          pctAsing,
			"pct_asing_gt1":      pctAsingGT1,
			"pct_asing_lt1":      pctAsingLT1,
			"pct_ritel_ksei_id":  pctRitelID,
			"pct_individual_gt1": pctIndivGT1,
			"pct_ritel_pure":     pctRitelPure,
			"local_breakdown": map[string]any{
				"corporate":         lCorp,
				"insurance":         lIns,
				"pension":           lPen,
				"mutual_fund":       lMF,
				"bank":              lBank,
				"securities":        lSec,
				"foundation":        lFound,
				"other":             lOther,
				"institution_total": lInst,
			},
			"foreign_breakdown": map[string]any{
				"corporate":         fCorp,
				"insurance":         fIns,
				"pension":           fPen,
				"bank":              fBank,
				"mutual_fund":       fMF,
				"securities":        fSec,
				"foundation":        fFound,
				"other":             fOther,
				"institution_total": fInst,
			},
		})
	}
	if err := compRows.Err(); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sidRows, err := h.DB.Query(ctx, `
		SELECT
			period_date, month_key, total_shareholders, prev_shareholders,
			change_shareholders, pct_shareholders_change::double precision, bae_name,
			ksei_pct_ritel::double precision, ksei_local_id_shares::double precision,
			avg_retail_shares::double precision, attachment_url
		FROM ownership_staging.v_sid_history
		WHERE ticker = $1
		ORDER BY period_date DESC
	`, tickerClean)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer sidRows.Close()

	sidList := []map[string]any{}
	for sidRows.Next() {
		var pDate *time.Time
		var mk string
		var totSH, prevSH, chgSH *int
		var pctChg, kseiRitel, kseiLocalID, avgRet *float64
		var bae, attach *string
		if scanErr := sidRows.Scan(
			&pDate, &mk, &totSH, &prevSH, &chgSH, &pctChg, &bae,
			&kseiRitel, &kseiLocalID, &avgRet, &attach,
		); scanErr != nil {
			httpx.WriteError(w, http.StatusInternalServerError, scanErr.Error())
			return
		}
		sidList = append(sidList, map[string]any{
			"period_date":             isoStr(pDate),
			"month_key":               mk,
			"total_shareholders":      totSH,
			"prev_shareholders":       prevSH,
			"change_shareholders":     chgSH,
			"pct_shareholders_change": pctChg,
			"bae_name":                nonEmptyPtr(bae),
			"ksei_pct_ritel":          kseiRitel,
			"ksei_local_id_shares":    floatToInt64Ptr(kseiLocalID),
			"avg_retail_shares":       floatToInt64Ptr(avgRet),
			"attachment_url":          nonEmptyPtr(attach),
		})
	}
	if err := sidRows.Err(); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	t5Rows, err := h.DB.Query(ctx, `
		SELECT
			movement_date, investor_name, broker, action,
			change_lot::double precision, value_m::double precision,
			price::double precision, change_pct::double precision
		FROM ownership_staging.v_t5_transactions
		WHERE ticker = $1
		ORDER BY movement_date DESC
		LIMIT 30
	`, tickerClean)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer t5Rows.Close()

	t5List := []map[string]any{}
	for t5Rows.Next() {
		var mDate *time.Time
		var invName, broker, action *string
		var chgLot, valM, priceV, chgPct *float64
		if scanErr := t5Rows.Scan(&mDate, &invName, &broker, &action, &chgLot, &valM, &priceV, &chgPct); scanErr != nil {
			httpx.WriteError(w, http.StatusInternalServerError, scanErr.Error())
			return
		}
		t5List = append(t5List, map[string]any{
			"date":          isoStr(mDate),
			"investor_name": cleanStrRaw(invName),
			"broker":        cleanStrRaw(broker),
			"action":        cleanStrRaw(action),
			"change_lot":    derefFloat(chgLot),
			"value_m":       derefFloat(valM),
			"price":         derefFloat(priceV),
			"change_pct":    derefFloat(chgPct),
		})
	}
	if err := t5Rows.Err(); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	metaMap, err := tickerMeta(ctx, h.DB, []string{tickerClean})
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	meta, hasMeta := metaMap[tickerClean]
	if len(compList) == 0 && len(sidList) == 0 && !hasMeta {
		httpx.WriteError(w, http.StatusNotFound, fmt.Sprintf("Emiten %s tidak ditemukan", tickerClean))
		return
	}
	companyName := tickerClean
	sector := "-"
	if hasMeta {
		companyName = meta.Name
		sector = meta.Sector
	}

	resp := map[string]any{
		"ticker":              tickerClean,
		"company_name":        companyName,
		"sector":              sector,
		"data_source":         "ownership_staging (KSEI K15/K16 + IDX BAE)",
		"composition_history": compList,
		"sid_history":         sidList,
		"t5_movements":        t5List,
	}
	h.storeCache(w, r, ck, resp, 600)
}

// GetStagingParity handles GET /api/v1/ownership/staging/parity/{ticker}.
func (h *Handler) GetStagingParity(w http.ResponseWriter, r *http.Request) {
	tickerClean := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	if tickerClean == "" {
		httpx.WriteError(w, http.StatusBadRequest, "Ticker tidak valid")
		return
	}
	ck := ownershipCacheKey("ownership_staging_parity_v1", tickerClean)
	if h.tryCache(w, r, ck, 600) {
		return
	}

	rows, err := h.DB.Query(r.Context(), `
		SELECT
			k.period_date,
			k.pct_ritel_pure::double precision AS ksei_ritel_pure,
			k.pct_ritel_ksei_id::double precision AS ksei_ritel_id,
			k.pct_asing::double precision AS ksei_foreign,
			ct.pct_ritel::double precision AS ct_retail,
			ct.pct_asing::double precision AS ct_foreign,
			ROUND(ABS(k.pct_ritel_pure - ct.pct_ritel)::numeric, 4)::double precision AS retail_diff_pp,
			ROUND(ABS(k.pct_asing - ct.pct_asing)::numeric, 4)::double precision AS foreign_diff_pp
		FROM ownership_staging.v_composition_ksei_normalized k
		LEFT JOIN ct_baseline.pct_history ct
		  ON ct.ticker = k.ticker
		 AND ct.month_key = CASE EXTRACT(MONTH FROM k.period_date)::int
		      WHEN 5 THEN 'mei'
		      WHEN 6 THEN 'jun'
		      WHEN 7 THEN 'jul'
		      WHEN 8 THEN 'agu'
		      WHEN 9 THEN 'sep'
		      WHEN 10 THEN 'okt'
		      WHEN 11 THEN 'nov'
		      WHEN 12 THEN 'des'
		      WHEN 1 THEN 'jan'
		      WHEN 2 THEN 'feb'
		      WHEN 3 THEN 'mar'
		      WHEN 4 THEN 'apr'
		    END
		WHERE k.ticker = $1
		ORDER BY k.period_date DESC
	`, tickerClean)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	records := []map[string]any{}
	maxRetailDiff := 0.0
	maxForeignDiff := 0.0
	for rows.Next() {
		var pDate *time.Time
		var rPure, rID, kFor, ctRet, ctFor, rDiff, fDiff *float64
		if scanErr := rows.Scan(&pDate, &rPure, &rID, &kFor, &ctRet, &ctFor, &rDiff, &fDiff); scanErr != nil {
			httpx.WriteError(w, http.StatusInternalServerError, scanErr.Error())
			return
		}
		if rDiff != nil && *rDiff > maxRetailDiff {
			maxRetailDiff = *rDiff
		}
		if fDiff != nil && *fDiff > maxForeignDiff {
			maxForeignDiff = *fDiff
		}
		isMatched := rDiff != nil && *rDiff <= 0.05 && fDiff != nil && *fDiff <= 0.05
		records = append(records, map[string]any{
			"period_date":       isoStr(pDate),
			"ksei_ritel_pure":   rPure,
			"ksei_ritel_id_raw": rID,
			"ksei_foreign":      kFor,
			"ct_retail":         ctRet,
			"ct_foreign":        ctFor,
			"retail_diff_pp":    rDiff,
			"foreign_diff_pp":   fDiff,
			"is_matched":        isMatched,
		})
	}
	if err := rows.Err(); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(records) == 0 {
		httpx.WriteError(w, http.StatusNotFound, fmt.Sprintf("Data paritas untuk %s tidak ditemukan di staging", tickerClean))
		return
	}

	parityStatus := "ACCEPTABLE"
	if maxForeignDiff <= 0.05 {
		parityStatus = "EXCELLENT"
	}
	resp := map[string]any{
		"ticker":              tickerClean,
		"comparison":          "KSEI Staging vs CT Baseline",
		"max_retail_diff_pp":  maxRetailDiff,
		"max_foreign_diff_pp": maxForeignDiff,
		"parity_status":       parityStatus,
		"months_compared":     len(records),
		"history":             records,
	}
	h.storeCache(w, r, ck, resp, 600)
}

func isoStr(t *time.Time) string {
	if t == nil {
		return "None"
	}
	return t.Format("2006-01-02")
}

func cleanStrRaw(p *string) string {
	if p == nil {
		return "None"
	}
	return *p
}
