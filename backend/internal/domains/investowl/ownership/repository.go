package ownership

import (
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func roundFromFormat(v float64, digits int) float64 {
	r, err := strconv.ParseFloat(strconv.FormatFloat(v, 'f', digits, 64), 64)
	if err != nil || r == 0 {
		return 0.0
	}
	return r
}

func cleanTickers(tickers []string) []string {
	seen := make(map[string]struct{}, len(tickers))
	out := make([]string, 0, len(tickers))
	for _, raw := range tickers {
		t := strings.ToUpper(strings.TrimSpace(raw))
		if t == "" {
			continue
		}
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}
	return out
}

func tickerMeta(ctx context.Context, db *pgxpool.Pool, tickers []string) (map[string]TickerMeta, error) {
	clean := cleanTickers(tickers)
	out := make(map[string]TickerMeta, len(clean))
	if len(clean) == 0 || db == nil {
		return out, nil
	}
	rows, err := db.Query(ctx, `
		SELECT ticker, name, sector
		FROM tickers
		WHERE ticker = ANY($1)
	`, clean)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t string
		var name, sector *string
		if scanErr := rows.Scan(&t, &name, &sector); scanErr != nil {
			return nil, scanErr
		}
		nm := t
		if name != nil && *name != "" {
			nm = *name
		}
		sec := "-"
		if sector != nil && *sector != "" {
			sec = *sector
		}
		out[t] = TickerMeta{Ticker: t, Name: nm, Sector: sec}
	}
	return out, rows.Err()
}

func priceSnapshot(ctx context.Context, db *pgxpool.Pool, tickers []string) (map[string]PriceSnapshot, error) {
	clean := cleanTickers(tickers)
	out := make(map[string]PriceSnapshot, len(clean))
	if len(clean) == 0 || db == nil {
		return out, nil
	}
	rows, err := db.Query(ctx, `
		WITH latest AS (
			SELECT DISTINCT ON (ticker) ticker, date, close
			FROM prices
			WHERE ticker = ANY($1)
			ORDER BY ticker, date DESC
		),
		previous AS (
			SELECT DISTINCT ON (ticker) ticker, date, close
			FROM prices
			WHERE ticker = ANY($1)
			  AND date < (SELECT MAX(date) FROM prices)
			ORDER BY ticker, date DESC
		),
		month_start AS (
			SELECT DISTINCT ON (ticker) ticker, date, close
			FROM prices
			WHERE ticker = ANY($1)
			  AND date >= date_trunc('month', (SELECT MAX(date) FROM prices))
			ORDER BY ticker, date ASC
		)
		SELECT l.ticker, l.date, l.close::double precision, p.date, p.close::double precision, m.date, m.close::double precision
		FROM latest l
		LEFT JOIN previous p ON p.ticker = l.ticker
		LEFT JOIN month_start m ON m.ticker = l.ticker
	`, clean)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t string
		var lDate, pDate, mDate *time.Time
		var lClose, pClose, mClose *float64
		if scanErr := rows.Scan(&t, &lDate, &lClose, &pDate, &pClose, &mDate, &mClose); scanErr != nil {
			return nil, scanErr
		}
		lc := 0.0
		if lClose != nil {
			lc = *lClose
		}
		out[t] = PriceSnapshot{
			Date:            isoDatePtr(lDate),
			Close:           roundEven(lc, 4),
			PreviousDate:    isoDatePtr(pDate),
			PreviousClose:   roundEvenPtr(pClose, 4),
			Change1DPct:     pctChange(&lc, pClose),
			MonthStartDate:  isoDatePtr(mDate),
			MonthStartClose: roundEvenPtr(mClose, 4),
			MTDPct:          pctChange(&lc, mClose),
			Source:          "prices",
		}
	}
	return out, rows.Err()
}

func holderAggregates(ctx context.Context, db *pgxpool.Pool, tickers []string) (map[string]HolderAggregate, error) {
	clean := cleanTickers(tickers)
	out := make(map[string]HolderAggregate, len(clean))
	if len(clean) == 0 || db == nil {
		return out, nil
	}
	q := `
		SELECT
			ticker,
			COUNT(*)::int AS holder_count,
			COALESCE(SUM(pct), 0)::double precision AS total_pct,
			COALESCE(SUM(CASE WHEN is_controlling THEN pct ELSE 0 END), 0)::double precision AS controlling_pct,
			COALESCE(SUM(CASE WHEN ` + publicFilterSQL + ` THEN pct ELSE 0 END), 0)::double precision AS public_pct,
			COALESCE(SUM(CASE WHEN ` + treasuryFilterSQL + ` THEN pct ELSE 0 END), 0)::double precision AS treasury_pct,
			COALESCE(SUM(CASE
				WHEN NOT is_controlling
				 AND NOT (` + publicFilterSQL + `)
				 AND NOT (` + treasuryFilterSQL + `)
				 AND pct >= 5 THEN pct ELSE 0 END), 0)::double precision AS substantial_pct,
			COALESCE(MAX(pct), 0)::double precision AS max_holder_pct,
			SUM(CASE WHEN is_controlling THEN 1 ELSE 0 END)::int AS controlling_count,
			SUM(CASE WHEN ` + publicFilterSQL + ` THEN 1 ELSE 0 END)::int AS public_count,
			SUM(CASE WHEN ` + treasuryFilterSQL + ` THEN 1 ELSE 0 END)::int AS treasury_count,
			SUM(CASE
				WHEN NOT is_controlling
				 AND NOT (` + publicFilterSQL + `)
				 AND NOT (` + treasuryFilterSQL + `)
				 AND pct >= 5 THEN 1 ELSE 0 END)::int AS substantial_count
		FROM company_shareholders
		WHERE ticker = ANY($1)
		GROUP BY ticker
	`
	rows, err := db.Query(ctx, q, clean)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t string
		var holderCount, ctrlCount, pubCount, treasCount, substCount int
		var total, ctrl, pub, treas, subst, maxHolder float64
		if scanErr := rows.Scan(
			&t, &holderCount, &total, &ctrl, &pub, &treas, &subst, &maxHolder,
			&ctrlCount, &pubCount, &treasCount, &substCount,
		); scanErr != nil {
			return nil, scanErr
		}
		management := math.Max(0.0, total-ctrl-pub-treas-subst)
		effectiveCtrl := ctrl
		if maxHolder >= 50 && maxHolder > effectiveCtrl {
			effectiveCtrl = maxHolder
		}
		out[t] = HolderAggregate{
			HolderCount:             holderCount,
			TotalPct:                roundEven(total, 3),
			ControllingPct:          roundEven(ctrl, 3),
			EffectiveControllingPct: roundEven(effectiveCtrl, 3),
			PublicPct:               roundEven(pub, 3),
			TreasuryPct:             roundEven(treas, 3),
			SubstantialPct:          roundEven(subst, 3),
			ManagementPct:           roundEven(management, 3),
			MaxHolderPct:            roundEven(maxHolder, 3),
			ControllingCount:        ctrlCount,
			PublicCount:             pubCount,
			TreasuryCount:           treasCount,
			SubstantialCount:        substCount,
			FreeFloatEst:            roundEven(math.Max(0.0, 100.0-effectiveCtrl-treas), 3),
		}
	}
	return out, rows.Err()
}

func monthlySnapshot(ctx context.Context, db *pgxpool.Pool, tickers []string) (map[string]MonthlySnapshot, error) {
	clean := cleanTickers(tickers)
	out := make(map[string]MonthlySnapshot, len(clean))
	if len(clean) == 0 || db == nil {
		return out, nil
	}
	rows, err := db.Query(ctx, `
		SELECT DISTINCT ON (ticker)
			ticker, period_date, total_sid, sid_change_mom, sid_change_pct::double precision,
			foreign_named_pct::double precision, foreign_hidden_pct::double precision,
			retail_pct::double precision, local_non_retail_pct::double precision,
			current_price::double precision, avg_accumulation_price::double precision,
			underwater_pct::double precision, estimated_nyangkut_pct::double precision,
			rar_score::double precision, pbi_score::double precision, is_hero_badge
		FROM monthly_ownership_metrics
		WHERE ticker = ANY($1)
		ORDER BY ticker, period_date DESC
	`, clean)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t string
		var pDate *time.Time
		var totalSID, sidMoM *int
		var sidPct, fnPct, fhPct, retPct, lnrPct *float64
		var curPrice, avgAccum, underwater, nyangkut, rar, pbi *float64
		var isHero *bool
		if scanErr := rows.Scan(
			&t, &pDate, &totalSID, &sidMoM, &sidPct, &fnPct, &fhPct, &retPct, &lnrPct,
			&curPrice, &avgAccum, &underwater, &nyangkut, &rar, &pbi, &isHero,
		); scanErr != nil {
			return nil, scanErr
		}
		out[t] = MonthlySnapshot{
			PeriodDate:           isoDatePtr(pDate),
			PeriodLabel:          periodLabel(pDate),
			TotalSID:             derefInt(totalSID),
			SIDChangeMoM:         derefInt(sidMoM),
			SIDChangePct:         roundEven(derefFloat(sidPct), 3),
			ForeignNamedPct:      roundEven(derefFloat(fnPct), 3),
			ForeignHiddenPct:     roundEven(derefFloat(fhPct), 3),
			RetailPct:            roundEven(derefFloat(retPct), 3),
			LocalNonRetailPct:    roundEven(derefFloat(lnrPct), 3),
			CurrentPrice:         roundEvenPtr(curPrice, 4),
			AvgAccumulationPrice: roundEvenPtr(avgAccum, 4),
			UnderwaterPct:        roundEvenPtr(underwater, 3),
			EstimatedNyangkutPct: roundEvenPtr(nyangkut, 3),
			RARScore:             roundEvenPtr(rar, 3),
			PBIScore:             roundEvenPtr(pbi, 3),
			IsHeroBadge:          isHero != nil && *isHero,
		}
	}
	return out, rows.Err()
}

func monthlyHistory(ctx context.Context, db *pgxpool.Pool, ticker string, limit int) ([]MonthlyHistoryItem, error) {
	if limit <= 0 {
		limit = 8
	}
	rows, err := db.Query(ctx, `
		SELECT period_date, total_sid, sid_change_mom, sid_change_pct::double precision,
		       foreign_named_pct::double precision, foreign_hidden_pct::double precision,
		       retail_pct::double precision, local_non_retail_pct::double precision,
		       current_price::double precision, underwater_pct::double precision,
		       estimated_nyangkut_pct::double precision, is_hero_badge
		FROM monthly_ownership_metrics
		WHERE ticker = $1
		ORDER BY period_date DESC
		LIMIT $2
	`, ticker, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tmp := make([]MonthlyHistoryItem, 0, limit)
	for rows.Next() {
		var pDate *time.Time
		var totalSID, sidMoM *int
		var sidPct, fnPct, fhPct, retPct, lnrPct, curPrice, underwater, nyangkut *float64
		var isHero *bool
		if scanErr := rows.Scan(
			&pDate, &totalSID, &sidMoM, &sidPct, &fnPct, &fhPct, &retPct, &lnrPct,
			&curPrice, &underwater, &nyangkut, &isHero,
		); scanErr != nil {
			return nil, scanErr
		}
		fn := derefFloat(fnPct)
		fh := derefFloat(fhPct)
		tmp = append(tmp, MonthlyHistoryItem{
			PeriodDate:           isoDatePtr(pDate),
			PeriodLabel:          periodLabel(pDate),
			TotalSID:             derefInt(totalSID),
			SIDChangeMoM:         derefInt(sidMoM),
			SIDChangePct:         roundEven(derefFloat(sidPct), 3),
			ForeignTotalPct:      roundEven(fn+fh, 3),
			ForeignNamedPct:      roundEven(fn, 3),
			ForeignHiddenPct:     roundEven(fh, 3),
			RetailPct:            roundEven(derefFloat(retPct), 3),
			LocalNonRetailPct:    roundEven(derefFloat(lnrPct), 3),
			CurrentPrice:         roundEvenPtr(curPrice, 4),
			UnderwaterPct:        roundEvenPtr(underwater, 3),
			EstimatedNyangkutPct: roundEvenPtr(nyangkut, 3),
			IsHeroBadge:          isHero != nil && *isHero,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i, j := 0, len(tmp)-1; i < j; i, j = i+1, j-1 {
		tmp[i], tmp[j] = tmp[j], tmp[i]
	}
	return tmp, nil
}

func holderRows(ctx context.Context, db *pgxpool.Pool, ticker string) ([]HolderRow, error) {
	rows, err := db.Query(ctx, `
		SELECT name, pct::double precision, is_controlling
		FROM company_shareholders
		WHERE ticker = $1
		ORDER BY pct DESC, name ASC
	`, ticker)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []HolderRow{}
	for rows.Next() {
		var name string
		var pct *float64
		var isCtrl *bool
		if scanErr := rows.Scan(&name, &pct, &isCtrl); scanErr != nil {
			return nil, scanErr
		}
		p := derefFloat(pct)
		ctrl := isCtrl != nil && *isCtrl
		out = append(out, HolderRow{
			Name:          name,
			Pct:           roundEven(p, 4),
			IsControlling: ctrl,
			Category:      classifyHolder(name, p, ctrl),
		})
	}
	return out, rows.Err()
}

func companyShares(ctx context.Context, db *pgxpool.Pool, ticker string) (*CompanyShares, error) {
	rows, err := db.Query(ctx, `
		SELECT as_of, listed_shares::double precision, tradeable_shares::double precision
		FROM company_shares
		WHERE ticker = $1
		ORDER BY as_of DESC NULLS LAST
		LIMIT 1
	`, ticker)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, rows.Err()
	}
	var asOf *time.Time
	var listed, tradeable *float64
	if scanErr := rows.Scan(&asOf, &listed, &tradeable); scanErr != nil {
		return nil, scanErr
	}
	return &CompanyShares{
		AsOf:            isoDatePtr(asOf),
		ListedShares:    derefFloat(listed),
		TradeableShares: derefFloat(tradeable),
		Source:          "company_shares",
	}, rows.Err()
}

func companySharesMap(ctx context.Context, db *pgxpool.Pool, tickers []string) (map[string]float64, error) {
	clean := cleanTickers(tickers)
	out := make(map[string]float64, len(clean))
	if len(clean) == 0 || db == nil {
		return out, nil
	}
	rows, err := db.Query(ctx, `
		SELECT DISTINCT ON (ticker) ticker, listed_shares::double precision
		FROM company_shares
		WHERE ticker = ANY($1)
		ORDER BY ticker, as_of DESC NULLS LAST
	`, clean)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t string
		var listed *float64
		if scanErr := rows.Scan(&t, &listed); scanErr != nil {
			return nil, scanErr
		}
		if listed != nil {
			out[t] = *listed
		}
	}
	return out, rows.Err()
}

func indexHistory(ctx context.Context, db *pgxpool.Pool, indexCode string, limit int) ([]IndexHistoryPoint, error) {
	if limit <= 0 {
		limit = 320
	}
	rows, err := db.Query(ctx, `
		SELECT date, close::double precision, volume, value, market_cap::double precision
		FROM idx_index_summary
		WHERE index_code = $1
		ORDER BY date DESC
		LIMIT $2
	`, indexCode, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tmp := make([]IndexHistoryPoint, 0, limit)
	for rows.Next() {
		var d *time.Time
		var cl, mc *float64
		var vol, val *int64
		if scanErr := rows.Scan(&d, &cl, &vol, &val, &mc); scanErr != nil {
			return nil, scanErr
		}
		tmp = append(tmp, IndexHistoryPoint{
			Date:      isoDatePtr(d),
			Close:     roundEven(derefFloat(cl), 4),
			Volume:    derefInt64(vol),
			Value:     derefInt64(val),
			MarketCap: roundEvenPtr(mc, 4),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i, j := 0, len(tmp)-1; i < j; i, j = i+1, j-1 {
		tmp[i], tmp[j] = tmp[j], tmp[i]
	}
	return tmp, nil
}

func explorerSectors(ctx context.Context, db *pgxpool.Pool) ([]string, error) {
	rows, err := db.Query(ctx, `
		SELECT DISTINCT sector
		FROM tickers
		WHERE sector IS NOT NULL AND btrim(sector) <> ''
		ORDER BY sector
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var s string
		if scanErr := rows.Scan(&s); scanErr != nil {
			return nil, scanErr
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func isoDatePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}

func derefInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func derefInt64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

func derefFloat(p *float64) float64 {
	if p == nil {
		return 0.0
	}
	return *p
}
