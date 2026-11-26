package foreign_flow

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/forensic"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SectorService implements sector rotation RRG and history endpoints.
type SectorService struct {
	DB       *pgxpool.Pool
	Cache    *cache.Store
	Forensic *forensic.Service
}

func fillPivotSeries(dates []string, raw map[string]float64) ([]float64, bool) {
	out := make([]float64, len(dates))
	hasAny := false
	last := math.NaN()
	for i, d := range dates {
		if v, ok := raw[d]; ok && !math.IsNaN(v) {
			last = v
			hasAny = true
		}
		out[i] = last
	}
	if !hasAny {
		return nil, false
	}
	last = math.NaN()
	for i := len(dates) - 1; i >= 0; i-- {
		if !math.IsNaN(out[i]) {
			last = out[i]
		} else {
			out[i] = last
		}
	}
	return out, true
}

func rollingRatio100(series []float64, window, minPeriods int) []float64 {
	n := len(series)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		start := i - window + 1
		if start < 0 {
			start = 0
		}
		var sum float64
		cnt := 0
		for j := start; j <= i; j++ {
			if !math.IsNaN(series[j]) && !math.IsInf(series[j], 0) {
				sum += series[j]
				cnt++
			}
		}
		if cnt < minPeriods || cnt == 0 {
			out[i] = 100.0
			continue
		}
		mean := sum / float64(cnt)
		if mean == 0 || math.IsNaN(series[i]) {
			out[i] = 100.0
			continue
		}
		val := (series[i] / mean) * 100.0
		if math.IsNaN(val) || math.IsInf(val, 0) {
			val = 100.0
		}
		out[i] = val
	}
	return out
}

// GetSummary computes GET /api/foreign-flow/sectors/summary.
func (s *SectorService) GetSummary(ctx context.Context) (*SectorRotationSummaryResponse, int, error) {
	const cacheKey = "sector_rotation_summary_v4"
	if s.Cache != nil {
		if raw, ok := s.Cache.Get(ctx, cacheKey); ok {
			var cached SectorRotationSummaryResponse
			if json.Unmarshal(raw, &cached) == nil {
				return &cached, 200, nil
			}
		}
	}

	idxRows, err := s.DB.Query(ctx, `
		SELECT date, index_code, close::double precision, COALESCE(value, 0)::double precision
		FROM idx_index_summary
		WHERE index_code IN (
			'COMPOSITE', 'IDXENERGY', 'IDXBASIC', 'IDXINDUST',
			'IDXNONCYC', 'IDXCYCLIC', 'IDXHEALTH', 'IDXFINANCE',
			'IDXPROPERT', 'IDXTECHNO', 'IDXINFRA', 'IDXTRANS'
		)
		AND date >= (SELECT MAX(date) - INTERVAL '75 days' FROM idx_index_summary)
		ORDER BY date ASC
	`)
	if err != nil {
		return nil, 500, err
	}
	defer idxRows.Close()

	dateSet := make(map[string]struct{})
	closeRaw := make(map[string]map[string]float64)
	valRaw := make(map[string]map[string]float64)
	rowCount := 0
	for idxRows.Next() {
		var dt time.Time
		var code string
		var closePtr *float64
		var val float64
		if err := idxRows.Scan(&dt, &code, &closePtr, &val); err != nil {
			return nil, 500, err
		}
		rowCount++
		dStr := dt.Format("2006-01-02")
		dateSet[dStr] = struct{}{}
		if closePtr != nil && !math.IsNaN(*closePtr) {
			if closeRaw[code] == nil {
				closeRaw[code] = make(map[string]float64)
			}
			closeRaw[code][dStr] = *closePtr
		}
		if valRaw[code] == nil {
			valRaw[code] = make(map[string]float64)
		}
		valRaw[code][dStr] = shared.Finite(&val, 0)
	}
	if err := idxRows.Err(); err != nil {
		return nil, 500, err
	}
	if rowCount == 0 {
		return nil, 404, fmt.Errorf("Data indeks sektoral belum tersedia di database.")
	}

	dates := make([]string, 0, len(dateSet))
	for d := range dateSet {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	ffRows, err := s.DB.Query(ctx, `
		WITH trading_dates AS (
			SELECT DISTINCT date FROM broker_flow ORDER BY date DESC LIMIT 20
		),
		ranked_dates AS (
			SELECT date, ROW_NUMBER() OVER (ORDER BY date DESC) AS rn FROM trading_dates
		),
		latest_dates AS (
			SELECT
				MAX(CASE WHEN rn = 1 THEN date END) AS max_date,
				MIN(CASE WHEN rn <= 5 THEN date END) AS date_5d_start,
				MIN(date) AS date_20d_start
			FROM ranked_dates
		),
		sector_flows AS (
			SELECT t.sector,
				   SUM(CASE WHEN bf.date = ld.max_date THEN bf.foreign_net_broker ELSE 0 END)::double precision AS net_1d,
				   SUM(CASE WHEN bf.date >= ld.date_5d_start THEN bf.foreign_net_broker ELSE 0 END)::double precision AS net_5d,
				   SUM(CASE WHEN bf.date >= ld.date_20d_start THEN bf.foreign_net_broker ELSE 0 END)::double precision AS net_20d
			FROM broker_flow bf
			CROSS JOIN latest_dates ld
			JOIN tickers t ON bf.ticker = t.ticker
			WHERE bf.date >= ld.date_20d_start AND bf.date <= ld.max_date AND t.sector IS NOT NULL
			GROUP BY t.sector
		),
		ranked_stocks AS (
			SELECT t.sector, bf.ticker, bf.total_value,
				   ROW_NUMBER() OVER(PARTITION BY t.sector ORDER BY bf.total_value DESC NULLS LAST) AS rn
			FROM broker_flow bf
			CROSS JOIN latest_dates ld
			JOIN tickers t ON bf.ticker = t.ticker
			WHERE bf.date = ld.max_date AND t.sector IS NOT NULL AND t.is_active = true
		),
		top_tickers_agg AS (
			SELECT sector, array_agg(ticker ORDER BY rn) AS top_tickers
			FROM ranked_stocks
			WHERE rn <= 5
			GROUP BY sector
		)
		SELECT sf.sector, sf.net_1d, sf.net_5d, sf.net_20d, COALESCE(ta.top_tickers, '{}') AS top_tickers
		FROM sector_flows sf
		LEFT JOIN top_tickers_agg ta ON sf.sector = ta.sector
	`)
	if err != nil {
		return nil, 500, err
	}
	defer ffRows.Close()

	type ffAgg struct {
		n1, n5, n20 float64
		top         []string
	}
	ffBySec := make(map[string]ffAgg)
	for ffRows.Next() {
		var sec string
		var n1, n5, n20 *float64
		var tops []string
		if err := ffRows.Scan(&sec, &n1, &n5, &n20, &tops); err != nil {
			return nil, 500, err
		}
		ffBySec[sec] = ffAgg{
			n1:  shared.Finite(n1, 0),
			n5:  shared.Finite(n5, 0),
			n20: shared.Finite(n20, 0),
			top: tops,
		}
	}
	if err := ffRows.Err(); err != nil {
		return nil, 500, err
	}

	nDates := len(dates)
	latestIdx := nDates - 1
	prevIdx := max(0, nDates-2)
	idx5D := max(0, nDates-6)
	idx20D := max(0, nDates-21)
	latestDate := dates[latestIdx]

	compClose, okComp := fillPivotSeries(dates, closeRaw["COMPOSITE"])
	if !okComp {
		return nil, 404, fmt.Errorf("Data benchmark COMPOSITE belum tersedia.")
	}
	ihsgNow := compClose[latestIdx]
	ihsgPrev := compClose[prevIdx]
	ihsgRet1D := 0.0
	if ihsgPrev != 0 {
		ihsgRet1D = shared.Round(((ihsgNow/ihsgPrev)-1)*100, 2)
	}
	ihsgVal := valRaw["COMPOSITE"][latestDate]

	macroState := "NORMAL"
	if s.Forensic != nil {
		if ex, err := s.Forensic.GetForeignExodusRadar(ctx, 45); err == nil && ex != nil && ex.MarketState != "" {
			macroState = ex.MarketState
		}
	}

	sectors := make([]SectorSummaryItem, 0, len(SectorOrder))
	totalSectorVal := 0.0
	totalForeignNet1D := 0.0

	for _, sCode := range SectorOrder {
		meta := SectorMeta[sCode]
		secClose, ok := fillPivotSeries(dates, closeRaw[sCode])
		if !ok {
			continue
		}
		cNow := secClose[latestIdx]
		cPrev := secClose[prevIdx]
		c5 := secClose[idx5D]
		c20 := secClose[idx20D]

		ret1D, ret5D, ret20D := 0.0, 0.0, 0.0
		if cPrev != 0 {
			ret1D = shared.Round(((cNow/cPrev)-1)*100, 2)
		}
		if c5 != 0 {
			ret5D = shared.Round(((cNow/c5)-1)*100, 2)
		}
		if c20 != 0 {
			ret20D = shared.Round(((cNow/c20)-1)*100, 2)
		}

		valNow := valRaw[sCode][latestDate]
		totalSectorVal += valNow

		rsRaw := make([]float64, nDates)
		for i := 0; i < nDates; i++ {
			if compClose[i] != 0 {
				rsRaw[i] = secClose[i] / compClose[i]
			} else {
				rsRaw[i] = math.NaN()
			}
		}
		rsRatioSeries := rollingRatio100(rsRaw, 14, 3)
		rsMomSeries := rollingRatio100(rsRatioSeries, 5, 2)

		currRS := shared.Round(rsRatioSeries[latestIdx], 2)
		currMom := shared.Round(rsMomSeries[latestIdx], 2)

		quadrant := "Improving"
		switch {
		case currRS >= 100 && currMom >= 100:
			quadrant = "Leading"
		case currRS >= 100 && currMom < 100:
			quadrant = "Weakening"
		case currRS < 100 && currMom < 100:
			quadrant = "Lagging"
		}

		ff := ffBySec[meta.DBSector]
		totalForeignNet1D += ff.n1
		fState, fLabel := EvaluateSectorForeignState(ff.n1, ff.n5, ret5D, ret20D, macroState)

		topStocks := ff.top
		if len(topStocks) == 0 {
			topStocks = append([]string{}, DefaultFallbackStocks[sCode]...)
		}

		sectors = append(sectors, SectorSummaryItem{
			Code:              sCode,
			Name:              meta.Name,
			DBSector:          meta.DBSector,
			Icon:              meta.Icon,
			Color:             meta.Color,
			Close:             cNow,
			Ret1D:             ret1D,
			Ret5D:             ret5D,
			Ret20D:            ret20D,
			Value:             valNow,
			ValueSharePct:     0,
			ForeignNet1D:      ff.n1,
			ForeignNet5D:      ff.n5,
			ForeignNet20D:     ff.n20,
			ForeignState:      fState,
			ForeignStateLabel: fLabel,
			RSRatio:           currRS,
			RSMomentum:        currMom,
			Quadrant:          quadrant,
			TopStocks:         topStocks,
		})
	}

	if totalSectorVal > 0 {
		for i := range sectors {
			sectors[i].ValueSharePct = shared.Round((sectors[i].Value/totalSectorVal)*100, 2)
		}
	}
	sort.SliceStable(sectors, func(i, j int) bool {
		return sectors[i].Ret1D > sectors[j].Ret1D
	})

	payload := &SectorRotationSummaryResponse{
		Date: latestDate,
		Benchmark: SectorBenchmark{
			Code:        "COMPOSITE",
			Name:        "IHSG (Indeks Harga Saham Gabungan)",
			Close:       ihsgNow,
			Ret1D:       ihsgRet1D,
			Value:       ihsgVal,
			MarketState: macroState,
		},
		MacroMarketState:  macroState,
		TotalSectorValue:  totalSectorVal,
		TotalForeignNet1D: totalForeignNet1D,
		Sectors:           sectors,
	}

	if s.Cache != nil {
		if raw, err := json.Marshal(payload); err == nil {
			s.Cache.Set(ctx, cacheKey, raw, 300*time.Second)
		}
	}
	return payload, 200, nil
}

// GetHistory computes GET /api/foreign-flow/sectors/history.
func (s *SectorService) GetHistory(ctx context.Context, sector string, days int) (*SectorHistoryResponse, int, error) {
	sec := strings.ToUpper(strings.TrimSpace(sector))
	meta, ok := SectorMeta[sec]
	if !ok {
		return nil, 404, fmt.Errorf("Sektor '%s' tidak valid atau tidak terdaftar.", sec)
	}

	cacheKey := fmt.Sprintf("sector_history_v4:%s:%d", sec, days)
	if s.Cache != nil {
		if raw, ok := s.Cache.Get(ctx, cacheKey); ok {
			var cached SectorHistoryResponse
			if json.Unmarshal(raw, &cached) == nil {
				return &cached, 200, nil
			}
		}
	}

	idxRows, err := s.DB.Query(ctx, `
		SELECT date, index_code, close::double precision, COALESCE(value, 0)::double precision
		FROM idx_index_summary
		WHERE index_code IN ($1, 'COMPOSITE')
		  AND date >= (SELECT MAX(date) - INTERVAL '150 days' FROM idx_index_summary)
		ORDER BY date ASC
	`, sec)
	if err != nil {
		return nil, 500, err
	}
	defer idxRows.Close()

	dateSet := make(map[string]struct{})
	closeRaw := map[string]map[string]float64{sec: {}, "COMPOSITE": {}}
	valRaw := make(map[string]float64)
	rowCount := 0
	for idxRows.Next() {
		var dt time.Time
		var code string
		var closePtr *float64
		var val float64
		if err := idxRows.Scan(&dt, &code, &closePtr, &val); err != nil {
			return nil, 500, err
		}
		rowCount++
		dStr := dt.Format("2006-01-02")
		dateSet[dStr] = struct{}{}
		if closePtr != nil && !math.IsNaN(*closePtr) {
			closeRaw[code][dStr] = *closePtr
		}
		if code == sec {
			valRaw[dStr] = shared.Finite(&val, 0)
		}
	}
	if err := idxRows.Err(); err != nil {
		return nil, 500, err
	}
	if rowCount == 0 {
		return nil, 404, fmt.Errorf("Data historis untuk sektor '%s' tidak ditemukan.", sec)
	}

	dates := make([]string, 0, len(dateSet))
	for d := range dateSet {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	secClose, okSec := fillPivotSeries(dates, closeRaw[sec])
	compClose, okComp := fillPivotSeries(dates, closeRaw["COMPOSITE"])
	if !okSec || !okComp {
		return nil, 404, fmt.Errorf("Data indeks %s tidak lengkap.", sec)
	}

	flowRows, err := s.DB.Query(ctx, `
		SELECT bf.date, SUM(bf.foreign_net_broker)::double precision AS daily_flow
		FROM broker_flow bf
		JOIN tickers t ON bf.ticker = t.ticker
		WHERE t.sector = $1
		  AND bf.date >= (SELECT MAX(date) - INTERVAL '150 days' FROM broker_flow)
		GROUP BY bf.date
		ORDER BY bf.date ASC
	`, meta.DBSector)
	if err != nil {
		return nil, 500, err
	}
	defer flowRows.Close()
	flowMap := make(map[string]float64)
	for flowRows.Next() {
		var dt time.Time
		var fPtr *float64
		if err := flowRows.Scan(&dt, &fPtr); err != nil {
			return nil, 500, err
		}
		flowMap[dt.Format("2006-01-02")] = shared.Finite(fPtr, 0)
	}
	if err := flowRows.Err(); err != nil {
		return nil, 500, err
	}

	brkRows, err := s.DB.Query(ctx, `
		SELECT ba.broker_code, SUM(ba.net_value)::double precision AS net_val
		FROM broker_activity ba
		JOIN tickers t ON ba.ticker = t.ticker
		LEFT JOIN public.broker_taxonomy bt ON bt.code = ba.broker_code
		WHERE t.sector = $1
		  AND (ba.participant_type ILIKE '%asing%' OR ba.broker_code IN ('ZP','AK','BK','RX','TP','KZ','YU','CS','CG','ML','MS','DB'))
		  AND COALESCE(bt.category, '') NOT IN ('retail', 'semi_ritel')
		  AND ba.date >= (SELECT MAX(date) - INTERVAL '14 days' FROM broker_activity)
		GROUP BY ba.broker_code
		ORDER BY net_val DESC
	`, meta.DBSector)
	if err != nil {
		return nil, 500, err
	}
	defer brkRows.Close()

	type brkVal struct {
		code string
		val  float64
	}
	allBrokers := make([]brkVal, 0, 32)
	for brkRows.Next() {
		var c string
		var vPtr *float64
		if err := brkRows.Scan(&c, &vPtr); err != nil {
			return nil, 500, err
		}
		allBrokers = append(allBrokers, brkVal{
			code: strings.ToUpper(strings.TrimSpace(c)),
			val:  shared.Finite(vPtr, 0),
		})
	}
	if err := brkRows.Err(); err != nil {
		return nil, 500, err
	}

	saviors := make([]SectorBrokerItem, 0, 3)
	for _, b := range allBrokers {
		if b.val > 0 {
			name, exists := shared.BrokerNameMap[b.code]
			if !exists {
				name = "Broker " + b.code
			}
			saviors = append(saviors, SectorBrokerItem{BrokerCode: b.code, BrokerName: name, NetValue: b.val})
			if len(saviors) >= 3 {
				break
			}
		}
	}

	culprits := make([]SectorBrokerItem, 0, 3)
	for i := len(allBrokers) - 1; i >= 0; i-- {
		b := allBrokers[i]
		if b.val < 0 {
			name, exists := shared.BrokerNameMap[b.code]
			if !exists {
				name = "Broker " + b.code
			}
			culprits = append(culprits, SectorBrokerItem{BrokerCode: b.code, BrokerName: name, NetValue: b.val})
			if len(culprits) >= 3 {
				break
			}
		}
	}

	rsRaw := make([]float64, len(dates))
	for i := range dates {
		if compClose[i] != 0 {
			rsRaw[i] = secClose[i] / compClose[i]
		} else {
			rsRaw[i] = math.NaN()
		}
	}
	rsRatioSeries := rollingRatio100(rsRaw, 14, 3)

	startIdx := 0
	if len(dates) > days {
		startIdx = len(dates) - days
	}
	points := make([]SectorHistoryPoint, 0, len(dates)-startIdx)
	runCum := 0.0
	for i := startIdx; i < len(dates); i++ {
		dStr := dates[i]
		dFlow := flowMap[dStr]
		runCum += dFlow
		cVal := shared.Round(secClose[i], 2)
		if math.IsNaN(secClose[i]) || math.IsInf(secClose[i], 0) {
			cVal = 0
		}
		points = append(points, SectorHistoryPoint{
			Date:                  dStr,
			Close:                 cVal,
			Value:                 valRaw[dStr],
			ForeignNet:            dFlow,
			CumulativeForeignFlow: runCum,
			RSRatio:               shared.Round(rsRatioSeries[i], 2),
		})
	}

	payload := &SectorHistoryResponse{
		Sector:         sec,
		Name:           meta.Name,
		Count:          len(points),
		Data:           points,
		SaviorBrokers:  saviors,
		CulpritBrokers: culprits,
	}

	if s.Cache != nil {
		if raw, err := json.Marshal(payload); err == nil {
			s.Cache.Set(ctx, cacheKey, raw, 300*time.Second)
		}
	}
	return payload, 200, nil
}
