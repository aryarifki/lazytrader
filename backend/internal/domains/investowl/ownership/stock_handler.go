package ownership

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler serves all /api/v1/ownership/* endpoints.
type Handler struct {
	DB                *pgxpool.Pool
	Cache             *cache.Store
	SyncCatalogRunner func(r *http.Request) (map[string]any, error)
}

func (h *Handler) tryCache(w http.ResponseWriter, r *http.Request, key string, sMaxAge int) bool {
	if h.Cache == nil || key == "" {
		return false
	}
	if raw, ok := h.Cache.Get(r.Context(), key); ok {
		httpx.WriteCachedJSON(w, raw, 30, sMaxAge, true)
		return true
	}
	return false
}

func (h *Handler) storeCache(w http.ResponseWriter, r *http.Request, key string, v any, sMaxAge int) {
	payload, err := json.Marshal(v)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if h.Cache != nil && key != "" {
		h.Cache.Set(r.Context(), key, payload, time.Duration(sMaxAge)*time.Second)
	}
	httpx.WriteCachedJSON(w, payload, 30, sMaxAge, false)
}

func buildMonthlyComposition(m MonthlySnapshot) CompositionBlock {
	items := []CompositionItem{
		{
			Key:         "foreign_named",
			Label:       "Foreign Named",
			Pct:         m.ForeignNamedPct,
			Color:       monthlyColors["foreign_named"],
			Description: "Kepemilikan asing yang teridentifikasi by name.",
		},
		{
			Key:         "foreign_hidden",
			Label:       "Foreign Hidden",
			Pct:         m.ForeignHiddenPct,
			Color:       monthlyColors["foreign_hidden"],
			Description: "Perkiraan kepemilikan asing yang tidak teridentifikasi by name.",
		},
		{
			Key:         "retail",
			Label:       "Retail",
			Pct:         m.RetailPct,
			Color:       monthlyColors["retail"],
			Description: "Investor ritel / pahlawan bursa.",
		},
		{
			Key:         "local_non_retail",
			Label:       "Lokal Non-Ritel",
			Pct:         m.LocalNonRetailPct,
			Color:       monthlyColors["local_non_retail"],
			Description: "Investor lokal non-ritel dan korporasi.",
		},
	}
	total := 0.0
	for _, it := range items {
		total += it.Pct
	}
	return CompositionBlock{
		Basis:       "monthly_ownership_metrics",
		PeriodDate:  m.PeriodDate,
		PeriodLabel: m.PeriodLabel,
		TotalPct:    roundEven(total, 3),
		Items:       items,
		Note:        "Komposisi berasal dari snapshot KSEI bulanan yang tersedia di database.",
	}
}

func buildSnapshotComposition(summary HolderAggregate) CompositionBlock {
	items := []CompositionItem{
		{
			Key:         "controlling",
			Label:       "Controlling",
			Pct:         summary.ControllingPct,
			Color:       snapshotColors["controlling"],
			Description: "Pemegang saham pengendali.",
		},
		{
			Key:         "substantial",
			Label:       "Substantial >5%",
			Pct:         summary.SubstantialPct,
			Color:       snapshotColors["substantial"],
			Description: "Pemegang saham non-pengendali di atas 5%.",
		},
		{
			Key:         "public",
			Label:       "Publik / Retail",
			Pct:         summary.PublicPct,
			Color:       snapshotColors["public"],
			Description: "Masyarakat / publik.",
		},
		{
			Key:         "treasury",
			Label:       "Saham Treasury",
			Pct:         summary.TreasuryPct,
			Color:       snapshotColors["treasury"],
			Description: "Saham treasury yang dicatat di tabel kepemilikan.",
		},
		{
			Key:         "management",
			Label:       "Manajemen / Lainnya",
			Pct:         summary.ManagementPct,
			Color:       snapshotColors["management"],
			Description: "Manajemen dan pemegang saham kecil lainnya.",
		},
	}
	known := 0.0
	for _, it := range items {
		known += it.Pct
	}
	other := math.Max(0.0, 100.0-known)
	if other >= 0.5 {
		items = append(items, CompositionItem{
			Key:         "other",
			Label:       "Tidak Terklasifikasi",
			Pct:         roundEven(other, 3),
			Color:       snapshotColors["other"],
			Description: "Selisih terhadap 100% karena keterbatasan snapshot.",
		})
	}
	return CompositionBlock{
		Basis:       "company_shareholders_snapshot",
		PeriodDate:  nil,
		PeriodLabel: "Snapshot pemegang saham",
		TotalPct:    roundEven(known+other, 3),
		Items:       items,
		Note: "Klasifikasi foreign/retail belum tersedia untuk emiten ini. " +
			"Yang ditampilkan adalah klasifikasi langsung dari snapshot pemegang saham.",
	}
}

func buildCTComposition(points []CTPctPoint, label *string) *CompositionBlock {
	if len(points) == 0 {
		return nil
	}
	latest := points[len(points)-1]
	fTotal := derefFloat(latest.PctAsing)
	fHidden := derefFloat(latest.PctAsingLt1)
	fNamed := math.Max(0.0, fTotal-fHidden)
	retail := derefFloat(latest.PctRitel)
	localNR := math.Max(0.0, 100.0-fTotal-retail)
	items := []CompositionItem{
		{
			Key:         "foreign_named",
			Label:       "Foreign Named",
			Pct:         roundEven(fNamed, 4),
			Color:       ctCompositionColors["foreign_named"],
			Description: "Foreign investor yang tercatat >=1%.",
		},
		{
			Key:         "foreign_hidden",
			Label:       "Foreign Hidden",
			Pct:         roundEven(fHidden, 4),
			Color:       ctCompositionColors["foreign_hidden"],
			Description: "Foreign investor di bawah ambang 1%.",
		},
		{
			Key:         "retail",
			Label:       "Retail",
			Pct:         roundEven(retail, 4),
			Color:       ctCompositionColors["retail"],
			Description: "Pemegang ritel / masyarakat.",
		},
		{
			Key:         "local_non_retail",
			Label:       "Lokal Non-Retail",
			Pct:         roundEven(localNR, 4),
			Color:       ctCompositionColors["local_non_retail"],
			Description: "Lokal non-ritel dan sisa komposisi.",
		},
	}
	total := 0.0
	for _, it := range items {
		total += it.Pct
	}
	pLabel := "Baseline CT"
	if label != nil && *label != "" {
		pLabel = *label
	} else if latest.MonthKey != "" {
		pLabel = latest.MonthKey
	}
	return &CompositionBlock{
		Basis:       "ct_baseline_pct_history",
		PeriodLabel: pLabel,
		PeriodDate:  nil,
		TotalPct:    roundEven(total, 3),
		Items:       items,
		Note: "Komposisi dari baseline CT (pct_asing, pct_asing_lt1, pct_ritel). " +
			"Foreign Named = total foreign - foreign hidden; lokal non-ritel = sisa.",
	}
}

func buildCTHistory(points []CTPctPoint, sidPoints []CTSIDPoint, labels map[string]string) []CTHistoryItem {
	out := []CTHistoryItem{}
	for idx, pt := range points {
		if pt.PctAsing == nil {
			continue
		}
		fTotal := *pt.PctAsing
		fHidden := derefFloat(pt.PctAsingLt1)
		retail := derefFloat(pt.PctRitel)
		lbl := pt.MonthKey
		if mapped, ok := labels[pt.MonthKey]; ok && mapped != "" {
			lbl = mapped
		}
		var pDate *string
		var totSID, sidChg *int
		var sidChgPct *float64
		if idx < len(sidPoints) {
			sp := sidPoints[idx]
			pCopy := sp.Period
			sCopy := sp.SID
			pDate = &pCopy
			totSID = &sCopy
			sidChg = sp.Change
			sidChgPct = sp.ChangePct
		}
		out = append(out, CTHistoryItem{
			PeriodLabel:       lbl,
			PeriodDate:        pDate,
			TotalSID:          totSID,
			SIDChangeMoM:      sidChg,
			SIDChangePct:      sidChgPct,
			ForeignTotalPct:   roundEven(fTotal, 4),
			ForeignNamedPct:   roundEven(math.Max(0.0, fTotal-fHidden), 4),
			ForeignHiddenPct:  roundEven(fHidden, 4),
			RetailPct:         roundEven(retail, 4),
			LocalNonRetailPct: roundEven(math.Max(0.0, 100.0-fTotal-retail), 4),
		})
	}
	return out
}

func buildDataQuality(
	hasHolders, hasMonthly, hasMarketCap, hasPeriodized, hasT5, hasSID, hasCTBaseline bool,
	referenceAsOf *string,
) map[string]any {
	notes := []string{}
	if !hasMonthly {
		notes = append(notes, "Snapshot KSEI bulanan (foreign/retail/SID) belum tersedia untuk emiten ini.")
	}
	if !hasHolders {
		notes = append(notes, "Tabel pemegang saham belum memiliki baris untuk emiten ini.")
	}
	if !hasPeriodized {
		notes = append(notes, "Riwayat kepemilikan per investor (periodized holdings) belum tersedia.")
	}
	if !hasT5 {
		notes = append(notes, "Transaksi T5 belum tersedia untuk emiten ini.")
	}
	if !hasSID {
		notes = append(notes, "Riwayat SID bulanan belum tersedia untuk emiten ini.")
	}
	if hasMarketCap {
		notes = append(notes, "Market cap dihitung dari harga terakhir x listed shares IDX.")
	} else {
		notes = append(notes, "Market cap dan jumlah saham beredar belum tersedia untuk emiten ini.")
	}
	if hasCTBaseline {
		stamp := ""
		if referenceAsOf != nil && *referenceAsOf != "" {
			stamp = " per " + *referenceAsOf
		}
		notes = append(notes, fmt.Sprintf(
			"Baseline CT%s bersifat statis dan quarantined: dipakai untuk parity ownership, bukan data live, dan tidak dipakai di jalur sinyal.",
			stamp,
		))
	}
	return map[string]any{
		"ownership_snapshot":  hasHolders,
		"monthly_composition": hasMonthly,
		"periodized_holdings": hasPeriodized,
		"t5_transactions":     hasT5,
		"sid_history":         hasSID,
		"market_cap":          hasMarketCap,
		"shares_outstanding":  hasMarketCap,
		"gl":                  false,
		"ct_baseline":         hasCTBaseline,
		"reference_snapshot":  hasCTBaseline,
		"reference_as_of":     referenceAsOf,
		"notes":               notes,
	}
}

// GetStockOwnership handles GET /api/v1/ownership/stock/{ticker}.
func (h *Handler) GetStockOwnership(w http.ResponseWriter, r *http.Request) {
	t := strings.TrimSpace(strings.ReplaceAll(strings.ToUpper(chi.URLParam(r, "ticker")), ".JK", ""))
	if t == "" {
		httpx.WriteError(w, http.StatusBadRequest, "Ticker tidak valid")
		return
	}
	ck := ownershipCacheKey("ownership_stock_v11", t)
	if h.tryCache(w, r, ck, 3600) {
		return
	}

	ctx := r.Context()
	metaMap, err := tickerMeta(ctx, h.DB, []string{t})
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	priceMap, err := priceSnapshot(ctx, h.DB, []string{t})
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	monthlyMap, err := monthlySnapshot(ctx, h.DB, []string{t})
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	aggMap, err := holderAggregates(ctx, h.DB, []string{t})
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	holders, err := holderRows(ctx, h.DB, t)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	shares, err := companyShares(ctx, h.DB, t)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	meta, hasMeta := metaMap[t]
	priceSnap, hasPrice := priceMap[t]
	monthlySnap, hasMonthly := monthlyMap[t]
	aggSnap, hasAgg := aggMap[t]

	ctMetaAll := ctTickerMetaMap(ctx, h.DB, []string{t})
	ctMetaRow, hasCTMeta := ctMetaAll[t]
	ctHoldersList := ctHolders(ctx, h.DB, t)
	ctPoints := ctPctHistory(ctx, h.DB, t)
	ctSID := ctSIDHistory(ctx, h.DB, t)
	ctT5 := ctT5Rows(ctx, h.DB, []string{t}, 100)
	ctLabels := ctMonthLabels(ctx, h.DB)
	ctSnapMeta := ctSnapshotMeta(ctx, h.DB)

	var ctLastUpdate *string
	if metaObj, ok := ctSnapMeta["meta"].(map[string]any); ok {
		if lu, ok2 := metaObj["last_update_date"].(string); ok2 && lu != "" {
			ctLastUpdate = &lu
		}
	}

	var priceOut any
	var priceClose *float64
	var priceDate *string
	var priceSource *string
	if hasPrice {
		priceOut = priceSnap
		cl := priceSnap.Close
		priceClose = &cl
		priceDate = priceSnap.Date
		src := priceSnap.Source
		priceSource = &src
	} else if hasCTMeta && ctMetaRow.Price != nil {
		cl := roundEven(*ctMetaRow.Price, 4)
		priceClose = &cl
		src := "ct_baseline"
		priceSource = &src
		priceOut = CTPriceSnapshot{
			Date:   nil,
			Close:  cl,
			Source: src,
		}
	}

	var listedShares *float64
	var sharesSource, sharesAsOf *string
	if shares != nil && shares.ListedShares != 0 {
		ls := shares.ListedShares
		listedShares = &ls
		src := shares.Source
		sharesSource = &src
		sharesAsOf = shares.AsOf
	} else if hasCTMeta && ctMetaRow.SecNum != nil && *ctMetaRow.SecNum != 0 {
		ls := *ctMetaRow.SecNum
		listedShares = &ls
		src := "ct_baseline.ticker_meta"
		sharesSource = &src
		sharesAsOf = nil
	}

	capProv := marketCapProvenance(
		priceClose, priceDate, priceSource,
		listedShares, sharesAsOf, sharesSource,
		ctMetaRow.MktcapB,
	)
	var marketCap, marketCapB *float64
	if capProv.Value != nil {
		mc := roundEven(*capProv.Value, 2)
		marketCap = &mc
		mcb := roundEven(mc/1_000_000_000.0, 2)
		marketCapB = &mcb
	}

	if !hasMeta && priceOut == nil && !hasMonthly && len(holders) == 0 && !hasCTMeta && len(ctHoldersList) == 0 {
		httpx.WriteError(w, http.StatusNotFound, fmt.Sprintf("Data kepemilikan %s tidak ditemukan", t))
		return
	}

	summary := aggSnap
	if !hasAgg {
		tot := 0.0
		maxP := 0.0
		for _, hr := range holders {
			tot += hr.Pct
			if hr.Pct > maxP {
				maxP = hr.Pct
			}
		}
		summary = HolderAggregate{
			HolderCount:  len(holders),
			TotalPct:     roundEven(tot, 3),
			MaxHolderPct: maxP,
		}
	}

	var composition CompositionBlock
	var history []MonthlyHistoryItem
	if hasMonthly {
		composition = buildMonthlyComposition(monthlySnap)
		history, err = monthlyHistory(ctx, h.DB, t, 8)
		if err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		composition = buildSnapshotComposition(summary)
		history = []MonthlyHistoryItem{}
	}

	var ctCompLabel *string
	if len(ctPoints) > 0 {
		mk := ctPoints[len(ctPoints)-1].MonthKey
		if lbl, ok := ctLabels[mk]; ok {
			ctCompLabel = &lbl
		}
	}
	ctComposition := buildCTComposition(ctPoints, ctCompLabel)
	ctHistory := buildCTHistory(ctPoints, ctSID, ctLabels)

	var smtrackerTotalSID *int
	if hasMonthly {
		v := monthlySnap.TotalSID
		smtrackerTotalSID = &v
	}
	totalSID := smtrackerTotalSID
	if len(ctSID) > 0 {
		v := ctSID[len(ctSID)-1].SID
		totalSID = &v
	}

	var ctFreeFloatPct *float64
	if ctMetaRow.FreeFloat != nil {
		v := roundEven(*ctMetaRow.FreeFloat*100.0, 2)
		ctFreeFloatPct = &v
	}
	freeFloatPct := summary.FreeFloatEst
	if freeFloatPct == 0 && ctFreeFloatPct != nil {
		freeFloatPct = *ctFreeFloatPct
	}

	conglomerate := ctMetaRow.Konglo
	if conglomerate == nil {
		cmRows, cmErr := h.DB.Query(ctx, `SELECT group_name FROM ownership_staging.v_conglomerate_members WHERE ticker = $1 LIMIT 1`, t)
		if cmErr == nil {
			if cmRows.Next() {
				var gn *string
				if scanErr := cmRows.Scan(&gn); scanErr == nil && gn != nil && *gn != "" {
					conglomerate = gn
				}
			}
			cmRows.Close()
		}
	}

	badges := resolveBadges(ctMetaRow.Badges, ctMetaRow.Badge)
	hasCT := hasCTMeta || len(ctHoldersList) > 0 || len(ctPoints) > 0 || len(ctSID) > 0 || len(ctT5) > 0

	// Staging Composition
	compRows, err := h.DB.Query(ctx, `
		SELECT label, item_kind, shares::double precision, normalized_pct::double precision,
		       color_light, color_dark, report_date
		FROM ownership_staging.v_composition_normalized
		WHERE ticker = $1
		ORDER BY report_date DESC, position ASC
	`, t)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer compRows.Close()
	var stagingComposition any
	var latestCompDate *time.Time
	compItems := []map[string]any{}
	for compRows.Next() {
		var label, kind string
		var sharesVal, normPct *float64
		var cLight, cDark *string
		var rDate *time.Time
		if scanErr := compRows.Scan(&label, &kind, &sharesVal, &normPct, &cLight, &cDark, &rDate); scanErr != nil {
			httpx.WriteError(w, http.StatusInternalServerError, scanErr.Error())
			return
		}
		if latestCompDate == nil && rDate != nil {
			d := *rDate
			latestCompDate = &d
		}
		if latestCompDate != nil && rDate != nil && !rDate.Equal(*latestCompDate) {
			continue
		}
		var sharesOut *int64
		if sharesVal != nil {
			v := int64(*sharesVal)
			sharesOut = &v
		}
		var colorOut any
		if cLight != nil && *cLight != "" {
			colorOut = *cLight
		} else if cDark != nil && *cDark != "" {
			colorOut = *cDark
		}
		compItems = append(compItems, map[string]any{
			"label":  label,
			"kind":   kind,
			"shares": sharesOut,
			"pct":    normPct,
			"color":  colorOut,
		})
	}
	if err := compRows.Err(); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if latestCompDate != nil && len(compItems) > 0 {
		stagingComposition = map[string]any{
			"report_date": latestCompDate.Format("2006-01-02"),
			"items":       compItems,
			"source":      "ownership_staging.v_composition_normalized",
		}
	}

	// Staging Holders
	shRows, err := h.DB.Query(ctx, `
		SELECT report_date, name, investor_type, classification, location, nationality,
		       pct::double precision, total_shares::double precision, previous_pct::double precision,
		       previous_total_shares::double precision, change_shares::double precision, change_pct::double precision
		FROM ownership_staging.holders
		WHERE ticker = $1 AND report_date = (
			SELECT MAX(report_date) FROM ownership_staging.holders WHERE ticker = $1
		)
		ORDER BY pct DESC, name ASC
	`, t)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer shRows.Close()
	var stagingHolders any
	var latestHolderDate *time.Time
	var stagingHoldersItems []map[string]any
	for shRows.Next() {
		var rDate *time.Time
		var nm string
		var invType, cls, loc, nat *string
		var pct, totShares, prevPct, prevTotShares, chgShares, chgPct *float64
		if scanErr := shRows.Scan(
			&rDate, &nm, &invType, &cls, &loc, &nat,
			&pct, &totShares, &prevPct, &prevTotShares, &chgShares, &chgPct,
		); scanErr != nil {
			httpx.WriteError(w, http.StatusInternalServerError, scanErr.Error())
			return
		}
		if latestHolderDate == nil && rDate != nil {
			d := *rDate
			latestHolderDate = &d
		}
		cp := 0.0
		if chgPct != nil {
			cp = roundEven(*chgPct, 4)
		}
		pp := roundEvenPtr(prevPct, 4)
		status := "TETAP"
		if pp == nil || *pp == 0 {
			status = "BARU"
		} else if cp > 0 {
			status = "NAIK"
		} else if cp < 0 {
			status = "TURUN"
		}
		stagingHoldersItems = append(stagingHoldersItems, map[string]any{
			"name":                  nm,
			"investor_type":         nonEmptyPtr(invType),
			"classification":        nonEmptyPtr(cls),
			"location":              nonEmptyPtr(loc),
			"nationality":           nonEmptyPtr(nat),
			"pct":                   roundEven(derefFloat(pct), 4),
			"total_shares":          derefFloat(totShares),
			"previous_pct":          pp,
			"previous_total_shares": prevTotShares,
			"change_shares":         derefFloat(chgShares),
			"change_pct":            cp,
			"status":                status,
		})
	}
	if err := shRows.Err(); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(stagingHoldersItems) > 0 {
		stagingHolders = map[string]any{
			"report_date": isoDatePtr(latestHolderDate),
			"holders":     stagingHoldersItems,
			"source":      "ownership_staging.holders",
		}
	}

	// Staging Movements
	mvRows, err := h.DB.Query(ctx, `
		SELECT movement_date, name, action_type, change_shares::double precision, change_pct::double precision,
		       current_shares::double precision, current_pct::double precision, price::double precision,
		       value::double precision, broker_code, data_source_label
		FROM ownership_staging.majorholder_movements
		WHERE ticker = $1
		ORDER BY movement_date DESC, value DESC NULLS LAST
		LIMIT 100
	`, t)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer mvRows.Close()
	stagingMovements := []map[string]any{}
	var mvBuyCnt, mvSellCnt int
	var mvBuyValM, mvSellValM float64
	for mvRows.Next() {
		var mDate *time.Time
		var nm string
		var actType, broker, srcLbl *string
		var chgSh, chgPct, curSh, curPct, priceV, valV *float64
		if scanErr := mvRows.Scan(&mDate, &nm, &actType, &chgSh, &chgPct, &curSh, &curPct, &priceV, &valV, &broker, &srcLbl); scanErr != nil {
			httpx.WriteError(w, http.StatusInternalServerError, scanErr.Error())
			return
		}
		act := "TRADE"
		if actType != nil && *actType != "" {
			act = strings.ToUpper(*actType)
		}
		var valM *float64
		if valV != nil {
			v := roundEven(*valV/1_000_000_000.0, 2)
			valM = &v
		}
		if act == "BUY" {
			mvBuyCnt++
			mvBuyValM += derefFloat(valM)
		} else if act == "SELL" {
			mvSellCnt++
			mvSellValM += derefFloat(valM)
		}
		stagingMovements = append(stagingMovements, map[string]any{
			"date":           isoDatePtr(mDate),
			"name":           nm,
			"action":         act,
			"change_shares":  derefFloat(chgSh),
			"change_pct":     roundEvenPtr(chgPct, 4),
			"current_shares": curSh,
			"current_pct":    roundEvenPtr(curPct, 4),
			"price":          roundEvenPtr(priceV, 2),
			"value":          roundEvenPtr(valV, 2),
			"value_m":        valM,
			"broker":         nonEmptyPtr(broker),
			"source_label":   nonEmptyPtr(srcLbl),
		})
	}
	if err := mvRows.Err(); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	stagingMovementsSummary := map[string]any{
		"count":        len(stagingMovements),
		"buy_count":    mvBuyCnt,
		"sell_count":   mvSellCnt,
		"buy_value_m":  roundEven(mvBuyValM, 2),
		"sell_value_m": roundEven(mvSellValM, 2),
		"net_value_m":  roundEven(mvBuyValM-mvSellValM, 2),
	}

	// Staging UBO
	uboRows, err := h.DB.Query(ctx, `
		SELECT from_report_date, to_report_date, name, pct_to::double precision, pct_from::double precision, subholdings
		FROM ownership_staging.ubo_owners
		WHERE ticker = $1
		ORDER BY pct_to DESC NULLS LAST
	`, t)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer uboRows.Close()
	stagingUBO := []map[string]any{}
	for uboRows.Next() {
		var fromDate, toDate *time.Time
		var nm string
		var pctTo, pctFrom *float64
		var subRaw []byte
		if scanErr := uboRows.Scan(&fromDate, &toDate, &nm, &pctTo, &pctFrom, &subRaw); scanErr != nil {
			httpx.WriteError(w, http.StatusInternalServerError, scanErr.Error())
			return
		}
		subs := []any{}
		if len(subRaw) > 0 {
			var parsed []any
			if jErr := json.Unmarshal(subRaw, &parsed); jErr == nil && parsed != nil {
				subs = parsed
			}
		}
		stagingUBO = append(stagingUBO, map[string]any{
			"name":         nm,
			"pct":          roundEven(derefFloat(pctTo), 4),
			"previous_pct": roundEvenPtr(pctFrom, 4),
			"subholdings":  subs,
			"report_date":  isoDatePtr(toDate),
		})
	}
	if err := uboRows.Err(); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	displayName := t
	if hasMeta && meta.Name != "" {
		displayName = meta.Name
	} else if ctMetaRow.Issuer != nil && *ctMetaRow.Issuer != "" {
		displayName = *ctMetaRow.Issuer
	}
	sector := "-"
	if hasMeta && meta.Sector != "" {
		sector = meta.Sector
	}
	var profSrc *string
	if hasMeta {
		s := "smtracker"
		profSrc = &s
	} else if hasCTMeta {
		s := "ct_baseline"
		profSrc = &s
	}

	var monthlyDatePtr *string
	if hasMonthly {
		monthlyDatePtr = monthlySnap.PeriodDate
	}
	var stagingHolderDatePtr any
	if shMap, ok := stagingHolders.(map[string]any); ok {
		stagingHolderDatePtr = shMap["report_date"]
	}
	var stagingMovementDatePtr any
	if len(stagingMovements) > 0 {
		stagingMovementDatePtr = stagingMovements[0]["date"]
	}

	resp := map[string]any{
		"ticker":       t,
		"name":         displayName,
		"sector":       sector,
		"badges":       badges,
		"badge":        ctMetaRow.Badge,
		"conglomerate": conglomerate,
		"konglo":       conglomerate,
		"price":        priceOut,
		"profile": map[string]any{
			"market_cap":             marketCap,
			"market_cap_b":           marketCapB,
			"shares_outstanding":     listedShares,
			"free_float_pct":         freeFloatPct,
			"ct_free_float_pct":      ctFreeFloatPct,
			"total_sid":              totalSID,
			"smtracker_total_sid":    smtrackerTotalSID,
			"avg_retail_shares":      nil,
			"conglomerate":           conglomerate,
			"konglo":                 conglomerate,
			"badge":                  ctMetaRow.Badge,
			"badges":                 badges,
			"scripless_lot":          ctMetaRow.SLRibuLot,
			"pct_scrip":              ctMetaRow.PctSL,
			"source":                 profSrc,
			"sources": map[string]any{
				"profile":    profSrc,
				"price":      priceSource,
				"shares":     sharesSource,
				"market_cap": capProv.Source,
			},
			"market_cap_basis":       capProv.Basis,
			"market_cap_independent": capProv.Independent,
		},
		"summary":                   summary,
		"composition":               composition,
		"staging_composition":       stagingComposition,
		"history":                   history,
		"ct_composition":            ctComposition,
		"ct_history":                ctHistory,
		"holders":                   holders,
		"staging_holders":           stagingHolders,
		"ct_holders":                ctHoldersList,
		"ct_sid_history":            ctSID,
		"transactions":              ctT5,
		"staging_movements":         stagingMovements,
		"staging_movements_summary": stagingMovementsSummary,
		"staging_ubo":               stagingUBO,
		"t5_summary":                ctT5Summary(ctT5),
		"data_quality": buildDataQuality(
			len(holders) > 0 || len(ctHoldersList) > 0 || stagingHolders != nil,
			hasMonthly,
			marketCap != nil,
			len(ctHoldersList) > 0 || stagingHolders != nil,
			len(ctT5) > 0 || len(stagingMovements) > 0,
			len(ctSID) > 0,
			hasCT,
			ctLastUpdate,
		),
		"as_of": map[string]any{
			"price":             priceDate,
			"monthly":           monthlyDatePtr,
			"shares":            sharesAsOf,
			"market_cap":        capProv.AsOf,
			"shareholders":      "snapshot saat ini",
			"staging_holders":   stagingHolderDatePtr,
			"staging_movements": stagingMovementDatePtr,
			"ct_baseline":       ctLastUpdate,
		},
	}
	h.storeCache(w, r, ck, resp, 3600)
}

func nonEmptyPtr(p *string) *string {
	if p == nil || *p == "" {
		return nil
	}
	v := *p
	return &v
}
