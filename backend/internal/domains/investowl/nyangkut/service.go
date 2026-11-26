package nyangkut

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service implements read-only Dashboard Nyangkut calculations backed by PostgreSQL.
type Service struct {
	DB *pgxpool.Pool
}

// ComputePBIRARScores computes canonical RAR & PBI scores matching nyangkut_service.compute_pbi_rar_scores.
func ComputePBIRARScores(retailPct, retailDeltaMom, sidChangePct, underwaterPct, estimatedNyangkutPct float64) (float64, float64) {
	rarScore := roundEven(
		math.Min(200.0, (retailPct*1.5)+(math.Max(0.0, retailDeltaMom)*12.0)+(math.Max(0.0, sidChangePct)*4.0)),
		1,
	)
	if underwaterPct >= 0.0 {
		return rarScore, 0.0
	}
	pbiScore := roundEven(
		math.Min(
			100.0,
			(math.Min(rarScore, 100.0)*0.45)+(math.Min(-underwaterPct, 60.0)*0.90)+(estimatedNyangkutPct*0.15),
		),
		1,
	)
	return rarScore, pbiScore
}

func formatDateID(t time.Time) string {
	return fmt.Sprintf("%d %s %d", t.Day(), monthNamesID[t.Month()-1], t.Year())
}

// GetLatestMarketAndKSEIDates returns market price date, KSEI cutoff date, formatted labels, and lag banner.
func (s *Service) GetLatestMarketAndKSEIDates(ctx context.Context) (time.Time, time.Time, string, string, string) {
	var latestPricePtr *time.Time
	_ = s.DB.QueryRow(ctx, "SELECT MAX(date)::timestamp FROM prices WHERE close > 0 AND volume > 0").Scan(&latestPricePtr)
	latestPriceDate := time.Now().UTC()
	if latestPricePtr != nil {
		latestPriceDate = *latestPricePtr
	}

	var latestKSEIPtr *time.Time
	_ = s.DB.QueryRow(ctx, "SELECT MAX(period_date)::timestamp FROM monthly_ownership_metrics").Scan(&latestKSEIPtr)
	var kseiSnapshotDate time.Time
	if latestKSEIPtr != nil {
		kseiSnapshotDate = *latestKSEIPtr
	} else {
		firstOfMonth := time.Date(latestPriceDate.Year(), latestPriceDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		kseiSnapshotDate = firstOfMonth.AddDate(0, 0, -1)
	}

	kseiLabel := formatDateID(kseiSnapshotDate)
	priceLabel := formatDateID(latestPriceDate)
	updateLagLabel := fmt.Sprintf("Data Resmi KSEI/BAE per %s • Harga Pasar %s", kseiLabel, priceLabel)
	return latestPriceDate, kseiSnapshotDate, kseiLabel, priceLabel, updateLagLabel
}

func (s *Service) getTickerName(ctx context.Context, ticker string) string {
	var name *string
	if err := s.DB.QueryRow(ctx, "SELECT name FROM tickers WHERE ticker = $1", ticker).Scan(&name); err == nil && name != nil && *name != "" {
		return *name
	}
	var profName *string
	if err := s.DB.QueryRow(ctx, "SELECT name FROM ownership_staging.idx_company_profiles WHERE ticker = $1", ticker).Scan(&profName); err == nil && profName != nil && *profName != "" {
		return *profName
	}
	return ticker + " Tbk"
}

type priceVolRow struct {
	close float64
	vol   int64
	val   float64
}

func (s *Service) computeVolumeMetrics(ctx context.Context, ticker string, lookbackDays int) (float64, float64, []VolumeNode, float64, bool) {
	var latestPtr *time.Time
	_ = s.DB.QueryRow(ctx, "SELECT MAX(date)::timestamp FROM prices WHERE close > 0 AND volume > 0").Scan(&latestPtr)
	latestDate := time.Now().UTC()
	if latestPtr != nil {
		latestDate = *latestPtr
	}
	cutoffDate := latestDate.AddDate(0, 0, -lookbackDays)

	var data []priceVolRow
	rows, err := s.DB.Query(ctx, `
		SELECT close::float8, volume::bigint, COALESCE(NULLIF(value, 0), close * volume)::float8 AS val
		FROM prices
		WHERE ticker = $1 AND date >= $2 AND volume > 0 AND close > 0
		ORDER BY date ASC
	`, ticker, cutoffDate.Format("2006-01-02"))
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var r priceVolRow
			if scanErr := rows.Scan(&r.close, &r.vol, &r.val); scanErr == nil {
				data = append(data, r)
			}
		}
	}

	if len(data) == 0 {
		fbRows, fbErr := s.DB.Query(ctx, `
			SELECT close::float8, volume::bigint, COALESCE(NULLIF(value, 0), close * volume)::float8 AS val
			FROM prices
			WHERE ticker = $1 AND volume > 0 AND close > 0
			ORDER BY date DESC
			LIMIT 100
		`, ticker)
		if fbErr == nil {
			defer fbRows.Close()
			for fbRows.Next() {
				var r priceVolRow
				if scanErr := fbRows.Scan(&r.close, &r.vol, &r.val); scanErr == nil {
					data = append(data, r)
				}
			}
			for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
				data[i], data[j] = data[j], data[i]
			}
		}
	}

	if len(data) == 0 {
		return 0.0, 0.0, []VolumeNode{{Price: 0.0, VolumePct: 100.0, IsNyangkutZone: false}}, 0.0, true
	}

	currentPrice := roundEven(data[len(data)-1].close, 2)
	var totalVolume int64
	var totalVal float64
	var aboveVol int64
	minP, maxP := data[0].close, data[0].close
	for _, r := range data {
		totalVolume += r.vol
		totalVal += r.val
		if r.close > currentPrice {
			aboveVol += r.vol
		}
		if r.close < minP {
			minP = r.close
		}
		if r.close > maxP {
			maxP = r.close
		}
	}

	if totalVolume <= 0 {
		return currentPrice, currentPrice, []VolumeNode{{Price: currentPrice, VolumePct: 100.0, IsNyangkutZone: false}}, 0.0, true
	}

	vwap := roundEven(totalVal/float64(totalVolume), 2)
	exactNyangkutPct := roundEven(float64(aboveVol)/float64(totalVolume)*100.0, 1)

	if minP == maxP {
		return currentPrice, vwap, []VolumeNode{{Price: currentPrice, VolumePct: 100.0, IsNyangkutZone: false}}, exactNyangkutPct, false
	}

	const numBins = 5
	binSize := (maxP - minP) / float64(numBins)
	var binVolumes [numBins]int64
	var binMidPrices [numBins]float64
	for i := 0; i < numBins; i++ {
		binMidPrices[i] = roundEven(minP+(float64(i)+0.5)*binSize, 2)
	}
	for _, r := range data {
		idx := int((r.close - minP) / binSize)
		if idx >= numBins {
			idx = numBins - 1
		}
		if idx < 0 {
			idx = 0
		}
		binVolumes[idx] += r.vol
	}

	nodes := make([]VolumeNode, 0, numBins)
	for i := numBins - 1; i >= 0; i-- {
		pct := roundEven((float64(binVolumes[i])/float64(totalVolume))*100.0, 1)
		priceLvl := binMidPrices[i]
		nodes = append(nodes, VolumeNode{
			Price:          priceLvl,
			VolumePct:      pct,
			IsNyangkutZone: priceLvl > currentPrice,
		})
	}
	return currentPrice, vwap, nodes, exactNyangkutPct, false
}

type momRow struct {
	periodDate        time.Time
	totalSid          int64
	sidChangeMoM      int64
	sidChangePct      float64
	foreignNamedPct   float64
	foreignHiddenPct  float64
	retailPct         float64
	localNonRetailPct float64
}

func ppChange(pts []TrendPoint) string {
	if len(pts) < 2 {
		return "Snapshot tunggal"
	}
	delta := pts[len(pts)-1].Val - pts[len(pts)-2].Val
	if delta == 0 {
		delta = 0.0
	}
	return fmt.Sprintf("%+.2fpp vs %s", delta, pts[len(pts)-2].Month)
}

func sidChangeLabel(pts []TrendPoint) string {
	if len(pts) < 2 || pts[len(pts)-2].Val <= 0 {
		return "Snapshot tunggal"
	}
	pct := ((pts[len(pts)-1].Val - pts[len(pts)-2].Val) / pts[len(pts)-2].Val) * 100.0
	if pct == 0 {
		pct = 0.0
	}
	return fmt.Sprintf("%+.2f%% vs %s", pct, pts[len(pts)-2].Month)
}

func (s *Service) buildOverviewAndTrends(
	ctx context.Context,
	ticker string,
	currentPrice float64,
	momHistory []momRow,
	totalSid int64,
	retailPct float64,
	foreignNamedPct float64,
	foreignHiddenPct float64,
) (OverviewStats, []TrendCard) {
	var kongloPtr *string
	_ = s.DB.QueryRow(ctx, `
		SELECT STRING_AGG(DISTINCT group_name, ' / ' ORDER BY group_name)
		FROM ownership_staging.v_conglomerate_members
		WHERE ticker = $1
	`, ticker).Scan(&kongloPtr)
	if kongloPtr == nil || *kongloPtr == "" {
		_ = s.DB.QueryRow(ctx, "SELECT conglomerate_group FROM master_emiten WHERE ticker = $1", ticker).Scan(&kongloPtr)
	}
	kongloStr := "-"
	if kongloPtr != nil && *kongloPtr != "" {
		kongloStr = *kongloPtr
	}

	var totSharesPtr, totScriptlessPtr, pctScriptlessPtr *float64
	_ = s.DB.QueryRow(ctx, `
		SELECT total_shares::float8, total_scriptless::float8, pct_scriptless::float8
		FROM ownership_staging.v_composition_ksei_normalized
		WHERE ticker = $1
		ORDER BY period_date DESC
		LIMIT 1
	`, ticker).Scan(&totSharesPtr, &totScriptlessPtr, &pctScriptlessPtr)

	var totalShares, totalScriptless, pctScriptless float64
	if totSharesPtr != nil {
		totalShares = *totSharesPtr
	}
	if totScriptlessPtr != nil {
		totalScriptless = *totScriptlessPtr
	}
	if pctScriptlessPtr != nil {
		pctScriptless = *pctScriptlessPtr
	}

	if totalShares <= 0 {
		var csShares *float64
		_ = s.DB.QueryRow(ctx, "SELECT listed_shares::float8 FROM company_shares WHERE ticker = $1 ORDER BY as_of DESC LIMIT 1", ticker).Scan(&csShares)
		if csShares != nil {
			totalShares = *csShares
		}
	}

	marketCapStr := "—"
	if totalShares > 0 && currentPrice > 0 {
		mcapVal := (currentPrice * totalShares) / 1e9
		marketCapStr = formatIntDot(int64(math.RoundToEven(mcapVal))) + " M"
	}

	var scriplessLots int64
	if totalScriptless > 0 {
		scriplessLots = int64(totalScriptless / 100.0)
	}
	scriplessStr := "—"
	if scriplessLots > 0 {
		scriplessStr = formatIntDot(scriplessLots)
	}
	scriplessPctText := "Data scripless tidak tersedia"
	if pctScriptless > 0 {
		scriplessPctText = fmt.Sprintf("%.2f%% dari total saham", pctScriptless)
	}

	var pubPctPtr *float64
	_ = s.DB.QueryRow(ctx, `
		SELECT COALESCE(SUM(pct), 0)::float8
		FROM company_shareholders
		WHERE ticker = $1 AND (name ILIKE 'Masyarakat%' OR name ILIKE 'Publik%')
	`, ticker).Scan(&pubPctPtr)
	freeFloatVal := 0.0
	if pubPctPtr != nil {
		freeFloatVal = *pubPctPtr
	}
	if freeFloatVal <= 0 {
		freeFloatVal = math.Min(100.0, math.Max(0.0, retailPct+foreignHiddenPct))
	}

	avgRetailStr := "—"
	if totalSid > 0 && totalScriptless > 0 && currentPrice > 0 {
		avgLots := int64((totalScriptless * (retailPct / 100.0) / float64(totalSid)) / 100.0)
		avgRp := roundEven((float64(avgLots)*100.0*currentPrice)/1e6, 1)
		avgRetailStr = fmt.Sprintf("%s lot / Rp%.1f jt", formatIntDot(avgLots), avgRp)
	}

	sidLabel := "0 SID"
	if totalSid > 0 {
		sidLabel = formatIntDot(totalSid) + " SID"
	}

	overviewStats := OverviewStats{
		MarketCap:           marketCapStr,
		TotalSahamScripless: scriplessStr,
		ScriplessPctText:    scriplessPctText,
		FreeFloat:           fmt.Sprintf("%.2f%%", freeFloatVal),
		JumlahPemegangSaham: sidLabel,
		AvgSahamRetail:      avgRetailStr,
		Konglo:              kongloStr,
	}

	foreignTotal := roundEven(foreignNamedPct+foreignHiddenPct, 2)
	var ftPts, fhPts, rtPts, sidPts []TrendPoint
	if len(momHistory) > 0 {
		for _, r := range momHistory {
			mLabel := monthNamesShort[r.periodDate.Month()-1]
			ftPts = append(ftPts, TrendPoint{Month: mLabel, Val: roundEven(r.foreignNamedPct+r.foreignHiddenPct, 2)})
			fhPts = append(fhPts, TrendPoint{Month: mLabel, Val: roundEven(r.foreignHiddenPct, 2)})
			rtPts = append(rtPts, TrendPoint{Month: mLabel, Val: roundEven(r.retailPct, 2)})
			sidPts = append(sidPts, TrendPoint{Month: mLabel, Val: float64(r.totalSid)})
		}
	} else {
		ftPts = []TrendPoint{{Month: "Kini", Val: foreignTotal}}
		fhPts = []TrendPoint{{Month: "Kini", Val: roundEven(foreignHiddenPct, 2)}}
		rtPts = []TrendPoint{{Month: "Kini", Val: roundEven(retailPct, 2)}}
		sidPts = []TrendPoint{{Month: "Kini", Val: float64(totalSid)}}
	}

	trendCards := []TrendCard{
		{
			Label:  "% FOREIGN TOTAL",
			Value:  fmt.Sprintf("%.2f%%", foreignTotal),
			Change: ppChange(ftPts),
			Color:  "#3b82f6",
			Points: ftPts,
		},
		{
			Label:  "% FOREIGN HIDDEN PLAYER",
			Value:  fmt.Sprintf("%.2f%%", foreignHiddenPct),
			Change: ppChange(fhPts),
			Color:  "#f59e0b",
			Points: fhPts,
		},
		{
			Label:  "% RETAIL (PAHLAWAN BURSA)",
			Value:  fmt.Sprintf("%.2f%%", retailPct),
			Change: ppChange(rtPts),
			Color:  "#10b981",
			Points: rtPts,
		},
		{
			Label:  "JUMLAH SID (INVESTOR)",
			Value:  formatIntDot(totalSid),
			Change: sidChangeLabel(sidPts),
			Color:  "#8b5cf6",
			Points: sidPts,
		},
	}
	return overviewStats, trendCards
}

// GetOrCalculateNyangkutData returns full nyangkut-tracker payload matching nyangkut_service.get_or_calculate_nyangkut_data.
func (s *Service) GetOrCalculateNyangkutData(ctx context.Context, ticker string) (StockDataResponse, error) {
	tickerUpper := strings.ToUpper(strings.TrimSpace(ticker))
	name := s.getTickerName(ctx, tickerUpper)

	_, _, kseiLabel, _, updateLagLabel := s.GetLatestMarketAndKSEIDates(ctx)
	currentPrice, vwap, volumeNodes, exactNyangkutPct, isEstimated := s.computeVolumeMetrics(ctx, tickerUpper, 180)

	rows, err := s.DB.Query(ctx, `
		SELECT period_date::timestamp,
		       COALESCE(total_sid, 0)::bigint,
		       COALESCE(sid_change_mom, 0)::bigint,
		       COALESCE(sid_change_pct, 0)::float8,
		       COALESCE(foreign_named_pct, 0)::float8,
		       COALESCE(foreign_hidden_pct, 0)::float8,
		       COALESCE(retail_pct, 0)::float8,
		       COALESCE(local_non_retail_pct, 0)::float8
		FROM monthly_ownership_metrics
		WHERE ticker = $1
		ORDER BY period_date DESC
		LIMIT 6
	`, tickerUpper)
	if err != nil {
		return StockDataResponse{}, err
	}
	defer rows.Close()

	var momDesc []momRow
	for rows.Next() {
		var r momRow
		if scanErr := rows.Scan(
			&r.periodDate, &r.totalSid, &r.sidChangeMoM, &r.sidChangePct,
			&r.foreignNamedPct, &r.foreignHiddenPct, &r.retailPct, &r.localNonRetailPct,
		); scanErr != nil {
			return StockDataResponse{}, scanErr
		}
		momDesc = append(momDesc, r)
	}
	if err := rows.Err(); err != nil {
		return StockDataResponse{}, err
	}

	momHistory := make([]momRow, len(momDesc))
	for i := range momDesc {
		momHistory[i] = momDesc[len(momDesc)-1-i]
	}

	var foreignNamedPct, foreignHiddenPct, retailPct, localNonRetailPct, sidChangePct, retailDeltaMom float64
	var totalSid, sidChangeMoM int64

	if len(momHistory) > 0 {
		latest := momHistory[len(momHistory)-1]
		foreignNamedPct = roundEven(latest.foreignNamedPct, 2)
		foreignHiddenPct = roundEven(latest.foreignHiddenPct, 2)
		retailPct = roundEven(latest.retailPct, 2)
		localNonRetailPct = roundEven(latest.localNonRetailPct, 2)
		totalSid = latest.totalSid
		sidChangeMoM = latest.sidChangeMoM
		sidChangePct = roundEven(latest.sidChangePct, 2)
		if len(momHistory) >= 2 {
			prev := momHistory[len(momHistory)-2]
			retailDeltaMom = retailPct - roundEven(prev.retailPct, 2)
		}
	} else {
		isEstimated = true
	}

	ownership := []OwnershipCategory{
		{
			Category: "Foreign Named",
			Pct:      foreignNamedPct,
			Color:    "#3b82f6",
			Desc:     "Foreign Named (Institusi Asing Terdaftar)",
		},
		{
			Category: "Foreign Hidden Player",
			Pct:      foreignHiddenPct,
			Color:    "#0ea5e9",
			Desc:     "Foreign Hidden Player (Nominee / Custody)",
		},
		{
			Category: "Retail (Pahlawan Bursa)",
			Pct:      retailPct,
			Color:    "#10b981",
			Desc:     "Retail (Pahlawan Bursa Investor Ritel)",
		},
		{
			Category: "Lokal Non-Ritel",
			Pct:      localNonRetailPct,
			Color:    "#475569",
			Desc:     "Lokal Non-Ritel (Institusi Domestik)",
		},
	}

	avgAccum := vwap
	if avgAccum <= 0 {
		avgAccum = currentPrice
	}
	underwaterPct := 0.0
	if avgAccum > 0 {
		underwaterPct = roundEven(((currentPrice-avgAccum)/avgAccum)*100.0, 2)
	}
	estimatedNyangkutPct := exactNyangkutPct
	rarScore, pbiScore := ComputePBIRARScores(retailPct, retailDeltaMom, sidChangePct, underwaterPct, estimatedNyangkutPct)
	isHeroBadge := underwaterPct < -5.0 && (retailDeltaMom > 0.0 || sidChangeMoM > 0)

	var sparklineSid []int64
	if len(momHistory) > 0 {
		sparklineSid = make([]int64, len(momHistory))
		for i, r := range momHistory {
			sparklineSid[i] = r.totalSid
		}
	} else {
		sparklineSid = []int64{totalSid}
	}
	if len(sparklineSid) == 1 {
		sparklineSid = []int64{sparklineSid[0], sparklineSid[0]}
	}

	overviewStats, trendCards := s.buildOverviewAndTrends(
		ctx, tickerUpper, currentPrice, momHistory, totalSid, retailPct, foreignNamedPct, foreignHiddenPct,
	)

	return StockDataResponse{
		Ticker:               tickerUpper,
		Name:                 name,
		CurrentPrice:         currentPrice,
		AvgAccumulationPrice: avgAccum,
		UnderwaterPct:        underwaterPct,
		EstimatedNyangkutPct: estimatedNyangkutPct,
		IsEstimated:          isEstimated,
		TotalSid:             totalSid,
		SidChangeMoM:         sidChangeMoM,
		SidChangePct:         sidChangePct,
		RarScore:             rarScore,
		PbiScore:             pbiScore,
		IsHeroBadge:          isHeroBadge,
		Ownership:            ownership,
		VolumeNodes:          volumeNodes,
		PeriodDate:           kseiLabel,
		UpdateLagLabel:       updateLagLabel,
		SparklineSid:         sparklineSid,
		OverviewStats:        overviewStats,
		TrendCards:           trendCards,
	}, nil
}

// GetLeaderboard returns Pahlawan Bursa leaderboard ordered by PBI score DESC, underwater_pct ASC.
func (s *Service) GetLeaderboard(ctx context.Context, limit int) ([]LeaderboardHeroItem, error) {
	_, kseiDate, _, _, _ := s.GetLatestMarketAndKSEIDates(ctx)
	rows, err := s.DB.Query(ctx, `
		SELECT m.ticker,
		       COALESCE(t.name, m.ticker) AS name,
		       COALESCE(m.pbi_score, 0)::float8,
		       COALESCE(m.rar_score, 0)::float8,
		       COALESCE(m.underwater_pct, 0)::float8 AS price_drop_pct,
		       COALESCE(m.sid_change_mom, 0)::bigint AS sid_growth_mom,
		       COALESCE(m.current_price, 0)::float8,
		       COALESCE(m.avg_accumulation_price, 0)::float8
		FROM monthly_ownership_metrics m
		LEFT JOIN tickers t ON m.ticker = t.ticker
		WHERE m.period_date = $1
		  AND m.pbi_score > 0
		  AND m.underwater_pct < 0
		ORDER BY m.pbi_score DESC, m.underwater_pct ASC
		LIMIT $2
	`, kseiDate.Format("2006-01-02"), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []LeaderboardHeroItem{}
	rank := 1
	for rows.Next() {
		var item LeaderboardHeroItem
		if scanErr := rows.Scan(
			&item.Ticker, &item.Name, &item.PbiScore, &item.RarScore,
			&item.PriceDropPct, &item.SidGrowthMoM, &item.CurrentPrice, &item.AvgAccumulationPrice,
		); scanErr != nil {
			return nil, scanErr
		}
		item.Rank = rank
		out = append(out, item)
		rank++
	}
	return out, rows.Err()
}

// GetNyangkutSummaryStats returns macro aggregate statistics for retail nyangkut metrics.
func (s *Service) GetNyangkutSummaryStats(ctx context.Context) (SummaryStatsResponse, error) {
	_, kseiDate, kseiLabel, _, _ := s.GetLatestMarketAndKSEIDates(ctx)
	var totalStocks, totalSid, heroCount *int64
	var avgNyangkut *float64

	err := s.DB.QueryRow(ctx, `
		SELECT COUNT(*)::bigint AS total_stocks,
		       (AVG(estimated_nyangkut_pct) FILTER (WHERE current_price > 0))::float8 AS avg_nyangkut_pct,
		       SUM(total_sid)::bigint AS total_tracked_sid,
		       (COUNT(*) FILTER (WHERE is_hero_badge = TRUE))::bigint AS hero_stocks_count
		FROM monthly_ownership_metrics
		WHERE period_date = $1
	`, kseiDate.Format("2006-01-02")).Scan(&totalStocks, &avgNyangkut, &totalSid, &heroCount)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return SummaryStatsResponse{}, err
	}

	res := SummaryStatsResponse{PeriodDate: kseiLabel}
	if totalStocks != nil {
		res.TotalStocksAnalyzed = *totalStocks
	}
	if avgNyangkut != nil {
		res.AvgRetailNyangkutPct = roundEven(*avgNyangkut, 1)
	}
	if totalSid != nil {
		res.TotalTrackedSid = *totalSid
	}
	if heroCount != nil {
		res.HeroStocksCount = *heroCount
	}
	return res, nil
}

// SearchTickersPrediction predicts tickers and issuer names ordered by exact/prefix match relevance.
func (s *Service) SearchTickersPrediction(ctx context.Context, query string, limit int) ([]TickerSearchResult, error) {
	qClean := strings.ToUpper(strings.TrimSpace(query))
	if qClean == "" {
		return []TickerSearchResult{}, nil
	}
	_, kseiDate, _, _, _ := s.GetLatestMarketAndKSEIDates(ctx)
	term := "%" + qClean + "%"
	prefix := qClean + "%"

	rows, err := s.DB.Query(ctx, `
		SELECT t.ticker,
		       t.name,
		       COALESCE(t.sector, 'Saham BEI') AS sector,
		       COALESCE(m.is_hero_badge, FALSE) AS is_hero_badge,
		       m.underwater_pct::float8
		FROM tickers t
		LEFT JOIN monthly_ownership_metrics m ON t.ticker = m.ticker AND m.period_date = $4
		WHERE t.ticker ILIKE $1 OR t.name ILIKE $1
		ORDER BY CASE
		           WHEN t.ticker = $2 THEN 0
		           WHEN t.ticker LIKE $3 THEN 1
		           WHEN t.name ILIKE $3 THEN 2
		           ELSE 3
		         END,
		         t.ticker ASC
		LIMIT $5
	`, term, qClean, prefix, kseiDate.Format("2006-01-02"), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []TickerSearchResult{}
	for rows.Next() {
		var r TickerSearchResult
		var namePtr *string
		if scanErr := rows.Scan(&r.Ticker, &namePtr, &r.Sector, &r.IsHeroBadge, &r.UnderwaterPct); scanErr != nil {
			return nil, scanErr
		}
		if namePtr != nil {
			r.Name = *namePtr
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
