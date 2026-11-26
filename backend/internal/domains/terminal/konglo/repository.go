package konglo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/domains/terminal/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database reads for Konglo Screener.
type Repository struct {
	DB *pgxpool.Pool
}

// ListReportStocks fetches all stock summary rows for /report computation.
func (r *Repository) ListReportStocks(ctx context.Context) ([]StockRow, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT ticker, name, group_name,
		       prev, open, high, low, last, chg,
		       vol, val, freq, shares_listed, market_cap
		FROM terminal_engine.konglo_stock_summary;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]StockRow, 0, 1600)
	for rows.Next() {
		var s StockRow
		if err := rows.Scan(
			&s.Ticker, &s.Name, &s.GroupName,
			&s.Prev, &s.Open, &s.High, &s.Low, &s.Last, &s.Chg,
			&s.Vol, &s.Val, &s.Freq, &s.SharesListed, &s.MarketCap,
		); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// ListGroupsUnordered preserves SQL insertion order for /report group ordering parity.
func (r *Repository) ListGroupsUnordered(ctx context.Context) ([]string, map[string][]string, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT group_name, COALESCE(tickers_arr, ARRAY(SELECT jsonb_array_elements_text(tickers))) as tickers
		FROM terminal_engine.konglo_groups;
	`)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	order := make([]string, 0, 48)
	m := make(map[string][]string, 48)
	for rows.Next() {
		var gname string
		var tickers []string
		if err := rows.Scan(&gname, &tickers); err != nil {
			return nil, nil, err
		}
		if tickers == nil {
			tickers = []string{}
		}
		order = append(order, gname)
		m[gname] = tickers
	}
	return order, m, rows.Err()
}

// GetReportTimestamps fetches latest public.prices trade date, summary updated_at, and optional IHSG summary.
func (r *Repository) GetReportTimestamps(ctx context.Context) (ReportTimestamps, error) {
	var ts ReportTimestamps
	_ = r.DB.QueryRow(ctx, `SELECT MAX(date) FROM public.prices;`).Scan(&ts.MaxTradeDate)
	_ = r.DB.QueryRow(ctx, `SELECT MAX(updated_at) FROM terminal_engine.konglo_stock_summary;`).Scan(&ts.MaxUpdated)

	var dataRaw []byte
	if err := r.DB.QueryRow(ctx, `SELECT data FROM terminal_engine.foreign_flow_ihsg_summary WHERE key = 'ihsg_summary';`).Scan(&dataRaw); err == nil && len(dataRaw) > 0 {
		var payload map[string]any
		if json.Unmarshal(dataRaw, &payload) == nil {
			if s, ok := payload["summary"].(map[string]any); ok {
				ts.IHSGSummary = s
			}
		}
	}
	return ts, nil
}

// ListGroupsOrdered returns all conglomerate groups ordered by group_name ASC.
func (r *Repository) ListGroupsOrdered(ctx context.Context) ([]GroupRow, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT group_name, COALESCE(tickers_arr, ARRAY(SELECT jsonb_array_elements_text(tickers))) as tickers, total_tickers
		FROM terminal_engine.konglo_groups
		ORDER BY group_name ASC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]GroupRow, 0, 48)
	for rows.Next() {
		var g GroupRow
		if err := rows.Scan(&g.GroupName, &g.Tickers, &g.TotalTickers); err != nil {
			return nil, err
		}
		if g.Tickers == nil {
			g.Tickers = []string{}
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

// ListStockSummary returns stock summaries optionally filtered by exact group_name.
func (r *Repository) ListStockSummary(ctx context.Context, group string, limit int) ([]SummaryStockRow, error) {
	query := `
		SELECT ticker, name, group_name,
		       prev, open, high, low, last, chg,
		       vol, val, freq, bb_price, bb_vol, bo_price, bo_vol,
		       shares_listed, market_cap, updated_at
		FROM terminal_engine.konglo_stock_summary`
	var rows pgx.Rows
	var err error
	if group != "" {
		rows, err = r.DB.Query(ctx, query+` WHERE group_name = $1 ORDER BY market_cap DESC NULLS LAST LIMIT $2;`, group, limit)
	} else {
		rows, err = r.DB.Query(ctx, query+` ORDER BY market_cap DESC NULLS LAST LIMIT $1;`, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]SummaryStockRow, 0, limit)
	for rows.Next() {
		var s SummaryStockRow
		var updatedAt *time.Time
		if err := rows.Scan(
			&s.Ticker, &s.Name, &s.GroupName,
			&s.Prev, &s.Open, &s.High, &s.Low, &s.Last, &s.Chg,
			&s.Vol, &s.Val, &s.Freq, &s.BBPrice, &s.BBVol, &s.BOPrice, &s.BOVol,
			&s.SharesListed, &s.MarketCap, &updatedAt,
		); err != nil {
			return nil, err
		}
		if updatedAt != nil {
			str := shared.FormatISO8601UTC(*updatedAt)
			s.UpdatedAt = &str
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// GetGroupByPattern finds the first matching group via ILIKE %group_name%.
func (r *Repository) GetGroupByPattern(ctx context.Context, groupName string) (*GroupRow, error) {
	var g GroupRow
	err := r.DB.QueryRow(ctx, `
		SELECT group_name, COALESCE(tickers_arr, ARRAY(SELECT jsonb_array_elements_text(tickers))) as tickers, total_tickers
		FROM terminal_engine.konglo_groups
		WHERE group_name ILIKE $1;
	`, "%"+groupName+"%").Scan(&g.GroupName, &g.Tickers, &g.TotalTickers)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if g.Tickers == nil {
		g.Tickers = []string{}
	}
	return &g, nil
}

// ListGroupMembers fetches summary rows for tickers ordered by market_cap DESC NULLS LAST.
func (r *Repository) ListGroupMembers(ctx context.Context, tickers []string) ([]GroupMemberRow, error) {
	if len(tickers) == 0 {
		return []GroupMemberRow{}, nil
	}
	rows, err := r.DB.Query(ctx, `
		SELECT ticker, name, prev, open, high, low, last, chg,
		       vol, val, freq, market_cap
		FROM terminal_engine.konglo_stock_summary
		WHERE ticker = ANY($1)
		ORDER BY market_cap DESC NULLS LAST;
	`, tickers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]GroupMemberRow, 0, len(tickers))
	for rows.Next() {
		var m GroupMemberRow
		if err := rows.Scan(
			&m.Ticker, &m.Name, &m.Prev, &m.Open, &m.High, &m.Low,
			&m.Last, &m.Chg, &m.Vol, &m.Val, &m.Freq, &m.MarketCap,
		); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
