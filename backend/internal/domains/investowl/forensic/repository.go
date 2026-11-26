package forensic

import (
	"context"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository queries macro foreign flow and culprit brokers from PostgreSQL.
type Repository struct {
	DB *pgxpool.Pool
}

// FetchMacroTimeline executes the exact q_macro SQL query from routers/forensic.py.
func (r *Repository) FetchMacroTimeline(ctx context.Context, lookback int) ([]TimelinePoint, error) {
	q := `
		WITH latest_dates AS (
			SELECT date FROM idx_index_summary WHERE index_code = 'COMPOSITE' ORDER BY date DESC LIMIT $1
		),
		target_dates AS (
			SELECT date FROM latest_dates
			UNION ALL
			SELECT date FROM (SELECT DISTINCT date FROM prices ORDER BY date DESC LIMIT $1) p_dates
			WHERE NOT EXISTS (SELECT 1 FROM latest_dates)
		),
		min_d AS (
			SELECT MIN(date) AS min_date FROM target_dates
		)
		SELECT p.date,
		       ROUND(s.close, 2)::double precision AS ihsg_close,
		       SUM((p.foreign_buy - p.foreign_sell) *
		           CASE
		               WHEN p.volume > 0 AND p.value > 0 THEN (p.value / p.volume)
		               ELSE COALESCE(p.close, 0)
		           END
		       )::double precision AS daily_flow
		FROM prices p
		CROSS JOIN min_d
		JOIN target_dates ld ON p.date = ld.date
		LEFT JOIN idx_index_summary s ON p.date = s.date AND s.index_code = 'COMPOSITE'
		WHERE p.date >= min_d.min_date
		GROUP BY p.date, s.close
		ORDER BY p.date ASC
	`
	rows, err := r.DB.Query(ctx, q, lookback)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	timeline := make([]TimelinePoint, 0, lookback)
	runningCum := 0.0
	for rows.Next() {
		var d time.Time
		var ihsgPtr, flowPtr *float64
		if err := rows.Scan(&d, &ihsgPtr, &flowPtr); err != nil {
			return nil, err
		}
		dFlow := shared.Finite(flowPtr, 0.0)
		runningCum += dFlow
		ihsgVal := shared.Finite(ihsgPtr, 0.0)
		timeline = append(timeline, TimelinePoint{
			Date:           d.Format("2006-01-02"),
			IHSGClose:      ihsgVal,
			DailyFlow:      dFlow,
			CumulativeFlow: runningCum,
		})
	}
	return timeline, rows.Err()
}

// FetchCulpritBrokers executes q_culprits from routers/forensic.py.
func (r *Repository) FetchCulpritBrokers(ctx context.Context) ([]CulpritBroker, error) {
	q := `
		SELECT broker_code, SUM(net_value)::double precision AS net_sell_idr
		FROM broker_activity
		WHERE (participant_type ILIKE '%asing%' OR broker_code IN ('ZP','AK','BK','RX','TP','KZ','YU','CS','CG','ML','MS','DB'))
		  AND net_value < 0
		  AND date >= (SELECT MAX(date) - INTERVAL '14 days' FROM broker_activity)
		GROUP BY broker_code
		ORDER BY net_sell_idr ASC
		LIMIT 5
	`
	rows, err := r.DB.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	culprits := make([]CulpritBroker, 0, 5)
	for rows.Next() {
		var rawCode string
		var netSellPtr *float64
		if err := rows.Scan(&rawCode, &netSellPtr); err != nil {
			return nil, err
		}
		code := strings.ToUpper(strings.TrimSpace(rawCode))
		name, ok := shared.BrokerNameMap[code]
		if !ok {
			name = "Broker " + code
		}
		culprits = append(culprits, CulpritBroker{
			BrokerCode: code,
			BrokerName: name,
			NetSellIDR: shared.Finite(netSellPtr, 0.0),
		})
	}
	return culprits, rows.Err()
}
