package signal

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository executes SQL queries for /api/signal/* endpoints.
type Repository struct {
	DB *pgxpool.Pool
}

func parseJSONB(raw []byte) map[string]any {
	if len(raw) == 0 {
		return map[string]any{}
	}
	var m map[string]any
	if json.Unmarshal(raw, &m) != nil || m == nil {
		return map[string]any{}
	}
	return m
}

func asFloatPtr(v any) *float64 {
	if v == nil {
		return nil
	}
	switch t := v.(type) {
	case float64:
		if math.IsNaN(t) || math.IsInf(t, 0) {
			return nil
		}
		return &t
	case int64:
		f := float64(t)
		return &f
	case int:
		f := float64(t)
		return &f
	}
	return nil
}

func asFloat(v any, fallback float64) float64 {
	if p := asFloatPtr(v); p != nil {
		return *p
	}
	return fallback
}

func asStr(v any, fallback string) string {
	if v == nil {
		return fallback
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fallback
}

func normalizeProb(v *float64) *float64 {
	if v == nil || math.IsNaN(*v) || math.IsInf(*v, 0) {
		return nil
	}
	p := *v
	if p > 1.0 {
		p = shared.Round(p/100.0, 4)
	}
	return &p
}

// GetDailySignals queries analytics_daily_signals for the latest scan date.
func (r *Repository) GetDailySignals(ctx context.Context, limit int, minScore *float64, topOnly bool) (map[string]any, int, error) {
	var latestDate *time.Time
	err := r.DB.QueryRow(ctx, `SELECT MAX(date) FROM analytics_daily_signals`).Scan(&latestDate)
	if err != nil && err != pgx.ErrNoRows {
		return nil, 500, err
	}
	if latestDate == nil {
		return nil, 404, fmt.Errorf("Belum ada sinyal yang digenerate oleh scanner algoritma.")
	}

	scoreFloor := 0.0
	if topOnly {
		scoreFloor = 65.0
	} else if minScore != nil {
		scoreFloor = *minScore
	}

	rows, err := r.DB.Query(ctx, `
		SELECT a.ticker, a.date, a.ml_win_prob::double precision, a.composite_score::double precision,
		       a.technical_score::double precision, a.smart_money_score::double precision,
		       a.sector_score::double precision, a.fundamental_score::double precision,
		       a.ml_label, a.features_snapshot, a.gate_notes
		FROM analytics_daily_signals a
		WHERE a.date = $1
		  AND a.composite_score >= $2
		ORDER BY a.composite_score DESC
		LIMIT $3
	`, *latestDate, scoreFloor, limit)
	if err != nil {
		return nil, 500, err
	}
	defer rows.Close()

	signals := make([]DailySignalItem, 0, limit)
	for rows.Next() {
		var tk string
		var dt time.Time
		var probPtr, compPtr, techPtr, smPtr, secPtr, fundPtr *float64
		var mlLabel, gateNotes *string
		var featRaw []byte
		if err := rows.Scan(&tk, &dt, &probPtr, &compPtr, &techPtr, &smPtr, &secPtr, &fundPtr, &mlLabel, &featRaw, &gateNotes); err != nil {
			return nil, 500, err
		}
		feat := parseJSONB(featRaw)
		prob := normalizeProb(probPtr)
		labelPtr := mlLabel
		if prob == nil && labelPtr != nil && (*labelPtr == "WIN" || *labelPtr == "HOLD" || *labelPtr == "LOSS") {
			na := "N/A"
			labelPtr = &na
		}

		wp := asStr(feat["wyckoff_phase"], "")
		if wp == "" {
			wp = asStr(feat["wp"], "-")
		}
		notes := ""
		if gateNotes != nil {
			notes = *gateNotes
		}

		signals = append(signals, DailySignalItem{
			Ticker:         tk,
			Date:           dt.Format("2006-01-02"),
			MLLabel:        labelPtr,
			MLWinProb:      prob,
			CompositeScore: shared.Finite(compPtr, 0),
			Scores: DailySignalScores{
				Technical:   shared.Finite(techPtr, 0),
				SmartMoney:  shared.Finite(smPtr, 0),
				Sector:      shared.Finite(secPtr, 0),
				Fundamental: shared.Finite(fundPtr, 0),
			},
			Features:        feat,
			WyckoffPhase:    wp,
			WyckoffName:     asStr(feat["wyckoff_name"], ""),
			SmartMoneyNotes: asStr(feat["smart_money_notes"], ""),
			SectorNotes:     asStr(feat["sector_notes"], ""),
			Persistence5D:   asFloat(feat["persistence_5d"], 0.5),
			RRRatio:         asFloatPtr(feat["rr_ratio"]),
			GateNotes:       notes,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, 500, err
	}

	return map[string]any{
		"latest_date": latestDate.Format("2006-01-02"),
		"count":       len(signals),
		"data":        signals,
	}, 200, nil
}

// GetBacktestStats mirrors GET /api/signal/backtest/stats.
func (r *Repository) GetBacktestStats(ctx context.Context) (*BacktestStatsResponse, error) {
	var totalSig, totalClosed, totalOpen, winCnt, lossCnt, expCnt int64
	var avgWin, avgLoss, avgHeld, sumWin, sumLoss *float64
	var latestExit, latestSig *time.Time

	err := r.DB.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COUNT(CASE WHEN status IN ('WIN', 'LOSS', 'EXPIRED') THEN 1 END),
			COUNT(CASE WHEN status = 'OPEN' THEN 1 END),
			COUNT(CASE WHEN status = 'WIN' THEN 1 END),
			COUNT(CASE WHEN status = 'LOSS' THEN 1 END),
			COUNT(CASE WHEN status = 'EXPIRED' THEN 1 END),
			COALESCE(AVG(CASE WHEN status = 'WIN' AND entry_price > 0 THEN ((exit_price - entry_price) / NULLIF(entry_price, 0)) * 100.0 END), 0.0)::double precision,
			COALESCE(AVG(CASE WHEN status = 'LOSS' AND entry_price > 0 THEN ((exit_price - entry_price) / NULLIF(entry_price, 0)) * 100.0 END), 0.0)::double precision,
			COALESCE(AVG(CASE WHEN status IN ('WIN', 'LOSS', 'EXPIRED') THEN days_held END), 0.0)::double precision,
			COALESCE(SUM(CASE WHEN status = 'WIN' AND entry_price > 0 THEN (exit_price - entry_price) / NULLIF(entry_price, 0) END), 0.0)::double precision,
			COALESCE(SUM(CASE WHEN status = 'LOSS' AND entry_price > 0 THEN ABS((exit_price - entry_price) / NULLIF(entry_price, 0)) END), 0.0)::double precision,
			MAX(exit_date),
			MAX(signal_date)
		FROM signals
	`).Scan(&totalSig, &totalClosed, &totalOpen, &winCnt, &lossCnt, &expCnt,
		&avgWin, &avgLoss, &avgHeld, &sumWin, &sumLoss, &latestExit, &latestSig)
	if err != nil {
		return nil, err
	}

	decisive := winCnt + lossCnt
	winRate := 0.0
	if totalClosed > 0 {
		winRate = shared.Round((float64(winCnt)/float64(totalClosed))*100.0, 2)
	}
	decWinRate := 0.0
	if decisive > 0 {
		decWinRate = shared.Round((float64(winCnt)/float64(decisive))*100.0, 2)
	}
	sw := shared.Finite(sumWin, 0)
	sl := shared.Finite(sumLoss, 0)
	pf := 0.0
	if sl > 0 {
		pf = shared.Round(sw/sl, 2)
	} else if sw > 0 {
		pf = 99.0
	}

	var latestAudit *string
	if latestExit != nil {
		s := latestExit.Format("2006-01-02")
		latestAudit = &s
	} else if latestSig != nil {
		s := latestSig.Format("2006-01-02")
		latestAudit = &s
	}

	return &BacktestStatsResponse{
		TotalSignals:    totalSig,
		TotalClosed:     totalClosed,
		TotalOpen:       totalOpen,
		WinCount:        winCnt,
		LossCount:       lossCnt,
		ExpiredCount:    expCnt,
		WinRate:         winRate,
		DecisiveWinRate: decWinRate,
		AvgWinPct:       shared.Round(shared.Finite(avgWin, 0), 2),
		AvgLossPct:      shared.Round(shared.Finite(avgLoss, 0), 2),
		AvgDaysHeld:     shared.Round(shared.Finite(avgHeld, 0), 1),
		ProfitFactor:    pf,
		LatestAuditDate: latestAudit,
	}, nil
}

// GetBacktestHistory mirrors GET /api/signal/backtest/history.
func (r *Repository) GetBacktestHistory(ctx context.Context, status string, limit, offset int) (map[string]any, error) {
	statusFilter := strings.ToUpper(strings.TrimSpace(status))
	if statusFilter != "ALL" && statusFilter != "WIN" && statusFilter != "LOSS" && statusFilter != "OPEN" {
		statusFilter = "ALL"
	}

	q := `
		SELECT
			s.id, s.ticker, s.signal_date,
			s.entry_price::double precision, s.tp_price::double precision, s.sl_price::double precision,
			s.exit_price::double precision, s.exit_date, s.status, s.days_held, s.features_snapshot,
			CASE
				WHEN s.exit_price IS NOT NULL AND s.entry_price > 0
				THEN ROUND(((s.exit_price - s.entry_price) / s.entry_price) * 100.0, 2)::double precision
				ELSE NULL
			END AS pnl_pct
		FROM signals s
	`
	args := []any{limit, offset}
	if statusFilter != "ALL" {
		q += ` WHERE s.status = $3`
		args = append(args, statusFilter)
	}
	q += ` ORDER BY s.signal_date DESC, (s.features_snapshot->>'composite_score')::numeric DESC NULLS LAST, s.ticker ASC LIMIT $1 OFFSET $2`

	rows, err := r.DB.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := make([]BacktestHistoryItem, 0, limit)
	for rows.Next() {
		var id, tk string
		var sigDate, exitDate *time.Time
		var entryP, tpP, slP, exitP, pnlP *float64
		var stPtr *string
		var daysHeldPtr *int
		var featRaw []byte
		if err := rows.Scan(&id, &tk, &sigDate, &entryP, &tpP, &slP, &exitP, &exitDate, &stPtr, &daysHeldPtr, &featRaw, &pnlP); err != nil {
			return nil, err
		}
		feat := parseJSONB(featRaw)
		entry := shared.Finite(entryP, 0)
		tp := shared.Finite(tpP, 0)
		sl := shared.Finite(slP, 0)

		var calcRR *float64
		if entry > sl && sl > 0 {
			v := shared.Round((tp-entry)/(entry-sl), 2)
			calcRR = &v
		}
		rrVal := asFloatPtr(feat["rr_ratio"])
		if rrVal == nil {
			rrVal = calcRR
		}

		var sdStr, edStr *string
		if sigDate != nil {
			s := sigDate.Format("2006-01-02")
			sdStr = &s
		}
		if exitDate != nil {
			s := exitDate.Format("2006-01-02")
			edStr = &s
		}
		st := "UNKNOWN"
		if stPtr != nil && *stPtr != "" {
			st = *stPtr
		}
		dh := 0
		if daysHeldPtr != nil {
			dh = *daysHeldPtr
		}

		wp := feat["wyckoff_phase"]
		if wp == nil || wp == "" {
			if v, ok := feat["wp"]; ok && v != nil {
				wp = v
			} else {
				wp = "-"
			}
		}
		compScore := feat["composite_score"]
		if compScore == nil {
			compScore = 0
		}
		smScore := feat["sm_score"]
		if smScore == nil {
			smScore = 0
		}
		secScore := feat["sec_score"]
		if secScore == nil {
			secScore = 0
		}
		wn := feat["wyckoff_name"]
		if wn == nil {
			wn = ""
		}

		history = append(history, BacktestHistoryItem{
			ID:             id,
			Ticker:         tk,
			SignalDate:     sdStr,
			EntryPrice:     entry,
			TPPrice:        tp,
			SLPrice:        sl,
			ExitPrice:      shared.FinitePtr(exitP),
			ExitDate:       edStr,
			Status:         st,
			DaysHeld:       dh,
			PnLPct:         shared.FinitePtr(pnlP),
			CompositeScore: compScore,
			SMScore:        smScore,
			SecScore:       secScore,
			WyckoffPhase:   wp,
			WyckoffName:    wn,
			RRRatio:        rrVal,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return map[string]any{
		"status_filter": statusFilter,
		"count":         len(history),
		"data":          history,
	}, nil
}

// GetOutcomesSummary mirrors GET /api/signal/outcomes/summary.
func (r *Repository) GetOutcomesSummary(ctx context.Context) (map[string]any, error) {
	var totEval, cWin, cLoss, cAct, cExp, win5D, tot5D int64
	var r1, r3, r5, r10, r20, maxRun, maxDD, avgPnL *float64

	err := r.DB.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COUNT(CASE WHEN status = 'WIN' THEN 1 END),
			COUNT(CASE WHEN status = 'LOSS' THEN 1 END),
			COUNT(CASE WHEN status = 'ACTIVE' THEN 1 END),
			COUNT(CASE WHEN status = 'EXPIRED' THEN 1 END),
			ROUND(AVG(return_1d)::numeric, 2)::double precision,
			ROUND(AVG(return_3d)::numeric, 2)::double precision,
			ROUND(AVG(return_5d)::numeric, 2)::double precision,
			ROUND(AVG(return_10d)::numeric, 2)::double precision,
			ROUND(AVG(return_20d)::numeric, 2)::double precision,
			ROUND(AVG(max_runup_20d)::numeric, 2)::double precision,
			ROUND(AVG(max_drawdown_20d)::numeric, 2)::double precision,
			ROUND(AVG(realized_pnl_pct)::numeric, 2)::double precision,
			COUNT(CASE WHEN return_5d > 0 THEN 1 END),
			COUNT(CASE WHEN return_5d IS NOT NULL THEN 1 END)
		FROM signal_outcomes
	`).Scan(&totEval, &cWin, &cLoss, &cAct, &cExp, &r1, &r3, &r5, &r10, &r20, &maxRun, &maxDD, &avgPnL, &win5D, &tot5D)
	if err != nil {
		return nil, err
	}

	winRate5D := 0.0
	if tot5D > 0 {
		winRate5D = shared.Round((float64(win5D)/float64(tot5D))*100.0, 2)
	}
	return map[string]any{
		"total_evaluated": totEval,
		"counts": map[string]any{
			"win":     cWin,
			"loss":    cLoss,
			"active":  cAct,
			"expired": cExp,
		},
		"returns": map[string]any{
			"t1_pct":      shared.Finite(r1, 0),
			"t3_pct":      shared.Finite(r3, 0),
			"t5_pct":      shared.Finite(r5, 0),
			"t10_pct":     shared.Finite(r10, 0),
			"t20_pct":     shared.Finite(r20, 0),
			"win_rate_5d": winRate5D,
		},
		"excursion_20d": map[string]any{
			"max_runup_pct":    shared.Finite(maxRun, 0),
			"max_drawdown_pct": shared.Finite(maxDD, 0),
		},
		"realized": map[string]any{
			"avg_pnl_pct": shared.Finite(avgPnL, 0),
		},
	}, nil
}

// ExportOutcomes mirrors GET /api/signal/outcomes/export.
func (r *Repository) ExportOutcomes(ctx context.Context, limit int, status string) (map[string]any, error) {
	q := `
		SELECT
			id, ticker, signal_date, entry_price::double precision,
			target_rr::double precision,
			return_1d::double precision, return_3d::double precision, return_5d::double precision,
			return_10d::double precision, return_20d::double precision,
			max_runup_20d::double precision, max_drawdown_20d::double precision,
			status, outcome_label, realized_pnl_pct::double precision,
			days_held, features_snapshot
		FROM signal_outcomes
	`
	args := []any{limit}
	if st := strings.ToUpper(strings.TrimSpace(status)); st != "" && st != "ALL" {
		q += ` WHERE status = $2`
		args = append(args, st)
	}
	q += ` ORDER BY signal_date DESC, id ASC LIMIT $1`

	rows, err := r.DB.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dataset := make([]map[string]any, 0, min(limit, 256))
	for rows.Next() {
		var id, tk string
		var sigDate time.Time
		var entryP float64
		var tRR, r1, r3, r5, r10, r20, maxRun, maxDD, realPnL *float64
		var st, outLbl *string
		var daysHeld *int
		var featRaw []byte
		if err := rows.Scan(&id, &tk, &sigDate, &entryP, &tRR, &r1, &r3, &r5, &r10, &r20,
			&maxRun, &maxDD, &st, &outLbl, &realPnL, &daysHeld, &featRaw); err != nil {
			return nil, err
		}
		dataset = append(dataset, map[string]any{
			"id":               id,
			"ticker":           tk,
			"signal_date":      sigDate.Format("2006-01-02"),
			"entry_price":      entryP,
			"target_rr":        shared.FinitePtr(tRR),
			"return_1d":        shared.FinitePtr(r1),
			"return_3d":        shared.FinitePtr(r3),
			"return_5d":        shared.FinitePtr(r5),
			"return_10d":       shared.FinitePtr(r10),
			"return_20d":       shared.FinitePtr(r20),
			"max_runup_20d":    shared.FinitePtr(maxRun),
			"max_drawdown_20d": shared.FinitePtr(maxDD),
			"status":           st,
			"outcome_label":    outLbl,
			"realized_pnl_pct": shared.FinitePtr(realPnL),
			"days_held":        daysHeld,
			"features":         parseJSONB(featRaw),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return map[string]any{
		"count": len(dataset),
		"data":  dataset,
	}, nil
}

const activeEpisodesSQL = `
WITH ranked AS (
	SELECT
		s.id, s.ticker, s.signal_date, s.entry_price, s.tp_price, s.sl_price,
		s.status, s.days_held, s.features_snapshot, s.episode_id,
		ROW_NUMBER() OVER (
			PARTITION BY COALESCE(s.episode_id, s.ticker)
			ORDER BY s.signal_date ASC, s.id ASC
		) AS episode_first_rank,
		ROW_NUMBER() OVER (
			PARTITION BY COALESCE(s.episode_id, s.ticker)
			ORDER BY s.signal_date DESC, s.id DESC
		) AS episode_last_rank,
		COUNT(*) OVER (PARTITION BY COALESCE(s.episode_id, s.ticker)) AS confirmation_count
	FROM signals s
	WHERE s.status = 'OPEN'
),
episodes AS (
	SELECT
		f.id, f.ticker, f.signal_date, f.entry_price, f.tp_price, f.sl_price,
		f.status, f.days_held, f.features_snapshot, f.confirmation_count, f.episode_id,
		l.signal_date AS latest_signal_date,
		l.features_snapshot AS latest_features_snapshot
	FROM ranked f
	JOIN ranked l
	  ON COALESCE(l.episode_id, l.ticker) = COALESCE(f.episode_id, f.ticker)
	 AND l.episode_last_rank = 1
	WHERE f.episode_first_rank = 1
)
SELECT
	e.id, e.ticker, e.signal_date,
	e.entry_price::double precision, e.tp_price::double precision, e.sl_price::double precision,
	e.status, e.days_held, e.features_snapshot,
	f.company_name, f.group_name,
	COALESCE(
		a.composite_score,
		(e.latest_features_snapshot->>'composite_score')::numeric,
		(e.features_snapshot->>'composite_score')::numeric,
		0
	)::double precision AS real_composite_score,
	COALESCE(a.features_snapshot, e.latest_features_snapshot) AS daily_features,
	e.episode_id, e.confirmation_count::int, e.latest_signal_date
FROM episodes e
LEFT JOIN (
	SELECT t.ticker, t.name AS company_name, COALESCE(cg.groups, t.sector) AS group_name
	FROM tickers t
	LEFT JOIN (
		SELECT ticker, STRING_AGG(DISTINCT group_name, ' / ' ORDER BY group_name) AS groups
		FROM ownership_staging.v_conglomerate_members
		GROUP BY ticker
	) cg ON cg.ticker = t.ticker
) f ON e.ticker = f.ticker
LEFT JOIN analytics_daily_signals a ON a.ticker = e.ticker AND a.date = e.latest_signal_date
ORDER BY e.latest_signal_date DESC, real_composite_score DESC, e.ticker ASC
LIMIT $1
`

const legacyActiveSQL = `
SELECT
	s.id, s.ticker, s.signal_date,
	s.entry_price::double precision, s.tp_price::double precision, s.sl_price::double precision,
	s.status, s.days_held, s.features_snapshot,
	f.company_name, f.group_name,
	COALESCE(a.composite_score, (s.features_snapshot->>'composite_score')::numeric, 0)::double precision AS real_composite_score,
	a.features_snapshot AS daily_features,
	NULL::varchar AS episode_id, 1::int AS confirmation_count, NULL::date AS latest_signal_date
FROM signals s
LEFT JOIN (
	SELECT t.ticker, t.name AS company_name, COALESCE(cg.groups, t.sector) AS group_name
	FROM tickers t
	LEFT JOIN (
		SELECT ticker, STRING_AGG(DISTINCT group_name, ' / ' ORDER BY group_name) AS groups
		FROM ownership_staging.v_conglomerate_members
		GROUP BY ticker
	) cg ON cg.ticker = t.ticker
) f ON s.ticker = f.ticker
LEFT JOIN analytics_daily_signals a ON (s.ticker = a.ticker AND s.signal_date = a.date)
WHERE s.status = 'OPEN'
ORDER BY s.signal_date DESC, real_composite_score DESC, s.ticker ASC
LIMIT $1
`

// GetActiveSignals mirrors GET /api/signal/active.
func (r *Repository) GetActiveSignals(ctx context.Context, limit int, episodeView bool) (map[string]any, error) {
	q := legacyActiveSQL
	if episodeView {
		q = activeEpisodesSQL
	}
	rows, err := r.DB.Query(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]ActiveSignalItem, 0, limit)
	for rows.Next() {
		var id, tk string
		var sigDate, lastConfDate *time.Time
		var entryP, tpP, slP, realComp *float64
		var stPtr, compNamePtr, grpNamePtr, epIDPtr *string
		var daysHeldPtr *int
		var confCount int
		var featRaw, dailyRaw []byte

		if err := rows.Scan(&id, &tk, &sigDate, &entryP, &tpP, &slP, &stPtr, &daysHeldPtr,
			&featRaw, &compNamePtr, &grpNamePtr, &realComp, &dailyRaw, &epIDPtr, &confCount, &lastConfDate); err != nil {
			return nil, err
		}

		feat := parseJSONB(featRaw)
		dailyFeat := parseJSONB(dailyRaw)
		entry := shared.Finite(entryP, 0)
		tp := shared.Finite(tpP, 0)
		sl := shared.Finite(slP, 0)

		targetGain, stopRisk, calcRR := 0.0, 0.0, 0.0
		if entry > 0 {
			targetGain = shared.Round(((tp-entry)/entry)*100.0, 2)
			stopRisk = shared.Round(((entry-sl)/entry)*100.0, 2)
		}
		if stopRisk > 0 {
			calcRR = shared.Round(targetGain/stopRisk, 2)
		}
		rr := calcRR
		if v := asFloatPtr(dailyFeat["rr_ratio"]); v != nil {
			rr = *v
		} else if v := asFloatPtr(feat["rr_ratio"]); v != nil {
			rr = *v
		}

		compScore := asFloat(feat["composite_score"], 0)
		if realComp != nil {
			compScore = *realComp
		}

		prob := normalizeProb(asFloatPtr(dailyFeat["ml_win_prob"]))
		if prob == nil {
			prob = normalizeProb(asFloatPtr(feat["ml_win_prob"]))
		}

		pickAny := func(k string, fallbackKey string, def any) any {
			if v, ok := dailyFeat[k]; ok && v != nil && v != "" && v != false {
				return v
			}
			if fallbackKey != "" {
				if v, ok := feat[fallbackKey]; ok && v != nil && v != "" {
					return v
				}
			} else if v, ok := feat[k]; ok && v != nil && v != "" {
				return v
			}
			return def
		}
		pickNullable := func(k, fallbackKey string) any {
			if v, ok := dailyFeat[k]; ok && v != nil {
				return v
			}
			return feat[fallbackKey]
		}

		compName := tk
		if compNamePtr != nil && *compNamePtr != "" {
			compName = *compNamePtr
		}
		grpName := "-"
		if grpNamePtr != nil && *grpNamePtr != "" {
			grpName = *grpNamePtr
		}
		var sdStr, lcdStr *string
		if sigDate != nil {
			s := sigDate.Format("2006-01-02")
			sdStr = &s
		}
		if lastConfDate != nil {
			s := lastConfDate.Format("2006-01-02")
			lcdStr = &s
		}
		st := "OPEN"
		if stPtr != nil && *stPtr != "" {
			st = *stPtr
		}
		dh := 0
		if daysHeldPtr != nil {
			dh = *daysHeldPtr
		}
		if confCount < 1 {
			confCount = 1
		}

		data = append(data, ActiveSignalItem{
			ID:                id,
			EpisodeID:         epIDPtr,
			Ticker:            tk,
			CompanyName:       compName,
			GroupName:         grpName,
			SignalDate:        sdStr,
			EntryPrice:        entry,
			TPPrice:           tp,
			SLPrice:           sl,
			TargetGainPct:     targetGain,
			StopRiskPct:       stopRisk,
			RiskRewardRatio:   rr,
			Status:            st,
			DaysHeld:          dh,
			ConfirmationCount: confCount,
			LastConfirmedDate: lcdStr,
			CompositeScore:    int(math.RoundToEven(compScore)),
			WyckoffPhase:      pickAny("wyckoff_phase", "wp", "-"),
			WyckoffName:       pickAny("wyckoff_name", "", ""),
			SmartMoneyNotes:   pickAny("smart_money_notes", "", ""),
			SectorNotes:       pickAny("sector_notes", "", ""),
			Persistence5D:     asFloat(pickNullable("persistence_5d", "persistence_5d"), 0.5),
			TurnoverIDR:       asFloat(pickNullable("turnover_idr", "turnover_idr"), 0),
			MLWinProb:         prob,
			CMF:               pickNullable("cmf", "cmf_v"),
			MFI:               pickNullable("mfi", "mfi_v"),
			RSI:               pickNullable("rsi", "rsi_v"),
			VCPGrade:          pickAny("vcp_grade", "", "NONE"),
			GateNotes:         pickAny("gate_notes", "", ""),
			Features:          feat,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return map[string]any{
		"count": len(data),
		"data":  data,
	}, nil
}
