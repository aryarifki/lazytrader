package stocks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository reads historical OHLCV bars from public.prices.
type Repository struct {
	DB *pgxpool.Pool
}

// GetHistory queries OHLCV bars descending by date and reverses to ascending chronological order.
func (r *Repository) GetHistory(ctx context.Context, ticker string, startDate, endDate *time.Time, limit int) ([]PriceBar, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database pool is nil")
	}

	clauses := []string{"ticker = $1"}
	args := []any{ticker}
	argIdx := 2

	if startDate != nil {
		clauses = append(clauses, fmt.Sprintf("date >= $%d", argIdx))
		args = append(args, *startDate)
		argIdx++
	}
	if endDate != nil {
		clauses = append(clauses, fmt.Sprintf("date <= $%d", argIdx))
		args = append(args, *endDate)
		argIdx++
	}
	args = append(args, limit)

	query := fmt.Sprintf(`
		SELECT date, ticker,
		       open::double precision, high::double precision,
		       low::double precision, close::double precision, volume
		FROM prices
		WHERE %s
		ORDER BY date DESC
		LIMIT $%d
	`, strings.Join(clauses, " AND "), argIdx)

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	desc := make([]PriceBar, 0, min(limit, 512))
	for rows.Next() {
		var d time.Time
		var b PriceBar
		if err := rows.Scan(&d, &b.Ticker, &b.Open, &b.High, &b.Low, &b.Close, &b.Volume); err != nil {
			return nil, err
		}
		b.Date = d.Format("2006-01-02")
		b.Open = shared.FinitePtr(b.Open)
		b.High = shared.FinitePtr(b.High)
		b.Low = shared.FinitePtr(b.Low)
		b.Close = shared.FinitePtr(b.Close)
		desc = append(desc, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i, j := 0, len(desc)-1; i < j; i, j = i+1, j-1 {
		desc[i], desc[j] = desc[j], desc[i]
	}
	return desc, nil
}
