package broker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository queries public.broker_flow and public.broker_activity.
type Repository struct {
	DB *pgxpool.Pool
}

func formatUTC(t *time.Time) *string {
	if t == nil {
		return nil
	}
	u := t.UTC()
	var s string
	micros := u.Nanosecond() / 1000
	if micros == 0 {
		s = u.Format("2006-01-02T15:04:05Z")
	} else {
		s = fmt.Sprintf("%s.%06dZ", u.Format("2006-01-02T15:04:05"), micros)
	}
	return &s
}

const flowSelectCols = `
	date, ticker, bandar_signal, bandar_signal_score::double precision, foreign_signal,
	foreign_net_broker::double precision, local_net_broker::double precision, gov_net_broker::double precision,
	foreign_net_flow::double precision, domestic_net_flow::double precision, total_value::double precision,
	conclusion_broker, conclusion_flow, fetched_at
`

func scanFlowRow(row pgx.Row) (BrokerFlowRow, error) {
	var d time.Time
	var fetched *time.Time
	var out BrokerFlowRow
	err := row.Scan(
		&d, &out.Ticker, &out.BandarSignal, &out.BandarSignalScore, &out.ForeignSignal,
		&out.ForeignNetBroker, &out.LocalNetBroker, &out.GovNetBroker,
		&out.ForeignNetFlow, &out.DomesticNetFlow, &out.TotalValue,
		&out.ConclusionBroker, &out.ConclusionFlow, &fetched,
	)
	if err != nil {
		return out, err
	}
	out.Date = d.Format("2006-01-02")
	out.FetchedAt = formatUTC(fetched)
	out.BandarSignalScore = shared.FinitePtr(out.BandarSignalScore)
	out.ForeignNetBroker = shared.FinitePtr(out.ForeignNetBroker)
	out.LocalNetBroker = shared.FinitePtr(out.LocalNetBroker)
	out.GovNetBroker = shared.FinitePtr(out.GovNetBroker)
	out.ForeignNetFlow = shared.FinitePtr(out.ForeignNetFlow)
	out.DomesticNetFlow = shared.FinitePtr(out.DomesticNetFlow)
	out.TotalValue = shared.FinitePtr(out.TotalValue)
	return out, nil
}

// GetLatest returns the most recent broker_flow row for ticker.
func (r *Repository) GetLatest(ctx context.Context, ticker string) (*BrokerFlowRow, error) {
	q := `SELECT ` + flowSelectCols + ` FROM broker_flow WHERE ticker = $1 ORDER BY date DESC LIMIT 1`
	row, err := scanFlowRow(r.DB.QueryRow(ctx, q, ticker))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// GetHistory returns chronological ascending broker_flow rows up to limit.
func (r *Repository) GetHistory(ctx context.Context, ticker string, startDate, endDate *time.Time, limit int) ([]BrokerFlowRow, error) {
	clauses := []string{"ticker = $1"}
	args := []any{ticker}
	idx := 2
	if startDate != nil {
		clauses = append(clauses, fmt.Sprintf("date >= $%d", idx))
		args = append(args, *startDate)
		idx++
	}
	if endDate != nil {
		clauses = append(clauses, fmt.Sprintf("date <= $%d", idx))
		args = append(args, *endDate)
		idx++
	}
	args = append(args, limit)

	q := fmt.Sprintf(`SELECT %s FROM broker_flow WHERE %s ORDER BY date DESC LIMIT $%d`,
		flowSelectCols, strings.Join(clauses, " AND "), idx)
	rows, err := r.DB.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	desc := make([]BrokerFlowRow, 0, min(limit, 256))
	for rows.Next() {
		item, err := scanFlowRow(rows)
		if err != nil {
			return nil, err
		}
		desc = append(desc, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i, j := 0, len(desc)-1; i < j; i, j = i+1, j-1 {
		desc[i], desc[j] = desc[j], desc[i]
	}
	return desc, nil
}

// GetSummaryRows returns ascending broker_flow rows since cutoff inclusive.
func (r *Repository) GetSummaryRows(ctx context.Context, ticker string, cutoff time.Time) ([]BrokerFlowRow, error) {
	q := `SELECT ` + flowSelectCols + ` FROM broker_flow WHERE ticker = $1 AND date >= $2 ORDER BY date ASC`
	rows, err := r.DB.Query(ctx, q, ticker, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]BrokerFlowRow, 0, 64)
	for rows.Next() {
		item, err := scanFlowRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

// GetDistinctBrokers returns distinct broker codes ordered alphabetically.
func (r *Repository) GetDistinctBrokers(ctx context.Context) ([]string, error) {
	rows, err := r.DB.Query(ctx, `SELECT DISTINCT broker_code FROM broker_activity ORDER BY broker_code ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]string, 0, 128)
	for rows.Next() {
		var code *string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		if code != nil {
			if c := strings.ToUpper(strings.TrimSpace(*code)); c != "" {
				out = append(out, c)
			}
		}
	}
	return out, rows.Err()
}

// GetWhaleActivity aggregates a single broker's activity over the latest `days` active sessions.
func (r *Repository) GetWhaleActivity(ctx context.Context, code string, days int) ([]time.Time, []whaleAggRow, error) {
	dRows, err := r.DB.Query(ctx, `
		SELECT DISTINCT date
		FROM broker_activity
		WHERE broker_code = $1
		ORDER BY date DESC
		LIMIT $2
	`, code, days)
	if err != nil {
		return nil, nil, err
	}
	defer dRows.Close()

	dates := make([]time.Time, 0, days)
	for dRows.Next() {
		var dt time.Time
		if err := dRows.Scan(&dt); err != nil {
			return nil, nil, err
		}
		dates = append(dates, dt)
	}
	if err := dRows.Err(); err != nil {
		return nil, nil, err
	}
	if len(dates) == 0 {
		return nil, nil, nil
	}

	minDate, maxDate := dates[0], dates[0]
	for _, dt := range dates[1:] {
		if dt.Before(minDate) {
			minDate = dt
		}
		if dt.After(maxDate) {
			maxDate = dt
		}
	}

	aggRows, err := r.DB.Query(ctx, `
		SELECT
			ticker,
			COALESCE(SUM(buy_value), 0)::double precision AS total_buy,
			COALESCE(SUM(sell_value), 0)::double precision AS total_sell,
			COALESCE(SUM(net_value), 0)::double precision AS total_net,
			COALESCE(SUM(buy_lot - sell_lot), 0)::double precision AS net_lot,
			CASE WHEN SUM(buy_lot) > 0 THEN (SUM(buy_value) / (SUM(buy_lot) * 100))::double precision ELSE NULL END AS buy_avg_price,
			CASE WHEN SUM(sell_lot) > 0 THEN (SUM(sell_value) / (SUM(sell_lot) * 100))::double precision ELSE NULL END AS sell_avg_price
		FROM broker_activity
		WHERE broker_code = $1 AND date >= $2 AND date <= $3
		GROUP BY ticker
	`, code, minDate, maxDate)
	if err != nil {
		return nil, nil, err
	}
	defer aggRows.Close()

	out := make([]whaleAggRow, 0, 64)
	for aggRows.Next() {
		var row whaleAggRow
		if err := aggRows.Scan(
			&row.Ticker, &row.TotalBuy, &row.TotalSell, &row.TotalNet,
			&row.NetLot, &row.BuyAvgPrice, &row.SellAvgPrice,
		); err != nil {
			return nil, nil, err
		}
		out = append(out, row)
	}
	return dates, out, aggRows.Err()
}
