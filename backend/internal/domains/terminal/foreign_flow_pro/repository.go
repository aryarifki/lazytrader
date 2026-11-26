package foreign_flow_pro

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository encapsulates read queries against terminal_engine foreign flow tables.
type Repository struct {
	DB *pgxpool.Pool
}

// ListScreener returns all rows ordered by ticker ASC matching FastAPI screener query.
func (r *Repository) ListScreener(ctx context.Context) ([]ScreenerDBRow, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT ticker, name, badges, screener_summary, updated_at
		FROM terminal_engine.foreign_flow_screener
		ORDER BY ticker ASC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]ScreenerDBRow, 0, 1100)
	for rows.Next() {
		var row ScreenerDBRow
		var badgesRaw, summaryRaw []byte
		if err := rows.Scan(&row.Ticker, &row.Name, &badgesRaw, &summaryRaw, &row.UpdatedAt); err != nil {
			return nil, err
		}
		if len(badgesRaw) > 0 {
			_ = json.Unmarshal(badgesRaw, &row.Badges)
		}
		if row.Badges == nil {
			row.Badges = []string{}
		}
		if len(summaryRaw) > 0 {
			_ = json.Unmarshal(summaryRaw, &row.ScreenerSummary)
		}
		if row.ScreenerSummary == nil {
			row.ScreenerSummary = map[string]any{}
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

// GetTickerMeta fetches optional screener metadata for a ticker.
func (r *Repository) GetTickerMeta(ctx context.Context, ticker string) (*ScreenerDBRow, error) {
	var row ScreenerDBRow
	var badgesRaw, summaryRaw []byte
	err := r.DB.QueryRow(ctx, `
		SELECT ticker, name, badges, screener_summary
		FROM terminal_engine.foreign_flow_screener
		WHERE ticker = $1;
	`, ticker).Scan(&row.Ticker, &row.Name, &badgesRaw, &summaryRaw)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if len(badgesRaw) > 0 {
		_ = json.Unmarshal(badgesRaw, &row.Badges)
	}
	if len(summaryRaw) > 0 {
		_ = json.Unmarshal(summaryRaw, &row.ScreenerSummary)
	}
	return &row, nil
}

// GetTickerDaily returns the latest limit days in chronological order (oldest -> newest).
func (r *Repository) GetTickerDaily(ctx context.Context, ticker string, limit int) ([]DailyPoint, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT date_iso, open, high, low, close,
		       foreign_net_val, volume, foreign_buy_vol, foreign_sell_vol,
		       partisipasi_asing, kepemilikan_foreign_total, kepemilikan_foreign_hidden,
		       raw_array
		FROM terminal_engine.foreign_flow_daily
		WHERE ticker = $1
		ORDER BY date_iso DESC
		LIMIT $2;
	`, ticker, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	desc := make([]DailyPoint, 0, limit)
	for rows.Next() {
		var dateVal time.Time
		var o, h, l, c, fnv, vol, fbv, fsv, part, kft, kfh *float64
		var rawBytes []byte
		if err := rows.Scan(&dateVal, &o, &h, &l, &c, &fnv, &vol, &fbv, &fsv, &part, &kft, &kfh, &rawBytes); err != nil {
			return nil, err
		}

		// Preserve Python truthiness: float(x) if x else None for OHLCV; is not None for fnv/part.
		pt := DailyPoint{
			Date:             dateVal.Format("2006-01-02"),
			Open:             nonZeroFloatPtr(o),
			High:             nonZeroFloatPtr(h),
			Low:              nonZeroFloatPtr(l),
			Close:            nonZeroFloatPtr(c),
			Volume:           nonZeroFloatPtr(vol),
			PartisipasiAsing: part,
		}
		if fnv != nil {
			pt.ForeignNetVal = *fnv
		}
		if len(rawBytes) > 0 {
			_ = json.Unmarshal(rawBytes, &pt.RawSeries)
		}
		desc = append(desc, pt)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Reverse into oldest -> newest chronological order
	for i, j := 0, len(desc)-1; i < j; i, j = i+1, j-1 {
		desc[i], desc[j] = desc[j], desc[i]
	}
	return desc, nil
}

// GetIHSGSummary returns the aggregate IHSG JSON payload and updated_at timestamp.
func (r *Repository) GetIHSGSummary(ctx context.Context) (map[string]any, *time.Time, error) {
	var dataRaw []byte
	var updatedAt *time.Time
	err := r.DB.QueryRow(ctx, `
		SELECT data, updated_at
		FROM terminal_engine.foreign_flow_ihsg_summary
		WHERE key = 'ihsg_summary';
	`).Scan(&dataRaw, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return map[string]any{}, nil, nil
		}
		return nil, nil, err
	}
	var payload map[string]any
	if len(dataRaw) > 0 {
		_ = json.Unmarshal(dataRaw, &payload)
	}
	if payload == nil {
		payload = map[string]any{}
	}
	return payload, updatedAt, nil
}

func nonZeroFloatPtr(p *float64) *float64 {
	if p == nil || *p == 0 {
		return nil
	}
	v := *p
	return &v
}
