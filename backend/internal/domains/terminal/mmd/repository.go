package mmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/domains/terminal/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository encapsulates read-only SQL queries for MM Detector Pro.
type Repository struct {
	DB *pgxpool.Pool
}

// ExplorerCandleRow holds the latest candle + company_name per ticker for /explorer/scan.
type ExplorerCandleRow struct {
	Ticker      string
	CompanyName *string
	Close       *float64
	Prev        *float64
	Open        *float64
	Volume      *int64
	Tavg        *float64
	BuyBrokers  map[string]float64
	SellBrokers map[string]float64
}

// BigPlayerRow holds aggregated T5 metrics joined with mmd_tickers_meta.
type BigPlayerRow struct {
	Ticker     string
	LastClose  *float64
	LastVolume *int64
	LastTavg   *float64
	T5Count    int64
	LastT5Date string
}

// BrokerActivityRow represents one row from public.broker_activity.
type BrokerActivityRow struct {
	Date    string
	BuyLot  *float64
	SellLot *float64
	BuyVal  *float64
	SellVal *float64
}

// ListTickersMeta returns all master tickers ordered alphabetically.
func (r *Repository) ListTickersMeta(ctx context.Context) ([]TickerMetaSummary, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT ticker, earliest_date, latest_date, total_days,
		       last_close, last_volume, last_tavg, last_net_volume
		FROM terminal_engine.mmd_tickers_meta
		ORDER BY ticker ASC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]TickerMetaSummary, 0, 900)
	for rows.Next() {
		var item TickerMetaSummary
		var earliest, latest *time.Time
		if err := rows.Scan(
			&item.Ticker, &earliest, &latest, &item.TotalDays,
			&item.LastClose, &item.LastVolume, &item.LastTavg, &item.LastNetVolume,
		); err != nil {
			return nil, err
		}
		item.EarliestDate = formatDatePtr(earliest)
		item.LatestDate = formatDatePtr(latest)
		out = append(out, item)
	}
	return out, rows.Err()
}

// GetTickerMeta fetches full metadata for one ticker or nil if absent.
func (r *Repository) GetTickerMeta(ctx context.Context, ticker string) (*TickerMetaFull, error) {
	var m TickerMetaFull
	var lastDate, earliest, latest, updatedAt *time.Time
	err := r.DB.QueryRow(ctx, `
		SELECT ticker, company_name, last_date, last_close, last_net_volume,
		       is_active, updated_at, earliest_date, latest_date,
		       total_days, last_volume, last_tavg
		FROM terminal_engine.mmd_tickers_meta
		WHERE ticker = $1;
	`, ticker).Scan(
		&m.Ticker, &m.CompanyName, &lastDate, &m.LastClose, &m.LastNetVolume,
		&m.IsActive, &updatedAt, &earliest, &latest,
		&m.TotalDays, &m.LastVolume, &m.LastTavg,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	m.LastDate = formatDatePtr(lastDate)
	m.EarliestDate = formatDatePtr(earliest)
	m.LatestDate = formatDatePtr(latest)
	if updatedAt != nil {
		s := shared.FormatISO8601UTC(*updatedAt)
		m.UpdatedAt = &s
	}
	return &m, nil
}

// GetAllCandles fetches all historical candles ordered oldest -> newest.
func (r *Repository) GetAllCandles(ctx context.Context, ticker string) ([]DailyCandle, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT date_iso, date_str, file_name,
		       open, high, low, close, volume, prev, tavg, net_volume,
		       buy_brokers, sell_brokers
		FROM terminal_engine.mmd_daily_candles
		WHERE ticker = $1
		ORDER BY date_iso ASC;
	`, ticker)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCandles(rows)
}

// GetLatestCandles fetches the latest `limit` candles ordered oldest -> newest.
func (r *Repository) GetLatestCandles(ctx context.Context, ticker string, limit int) ([]DailyCandle, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT date_iso, date_str, file_name,
		       open, high, low, close, volume, prev, tavg, net_volume,
		       buy_brokers, sell_brokers
		FROM (
			SELECT date_iso, date_str, file_name,
			       open, high, low, close, volume, prev, tavg, net_volume,
			       buy_brokers, sell_brokers
			FROM terminal_engine.mmd_daily_candles
			WHERE ticker = $1
			ORDER BY date_iso DESC
			LIMIT $2
		) sub
		ORDER BY date_iso ASC;
	`, ticker, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCandles(rows)
}

func scanCandles(rows pgx.Rows) ([]DailyCandle, error) {
	out := make([]DailyCandle, 0, 100)
	for rows.Next() {
		var c DailyCandle
		var dateVal time.Time
		var buysRaw, sellsRaw []byte
		if err := rows.Scan(
			&dateVal, &c.DateStr, &c.FileName,
			&c.Open, &c.High, &c.Low, &c.Close, &c.Volume, &c.Prev, &c.Tavg, &c.NetVolume,
			&buysRaw, &sellsRaw,
		); err != nil {
			return nil, err
		}
		c.DateISO = dateVal.Format("2006-01-02")
		if len(buysRaw) > 0 {
			_ = json.Unmarshal(buysRaw, &c.BuyBrokers)
		}
		if c.BuyBrokers == nil {
			c.BuyBrokers = map[string]float64{}
		}
		if len(sellsRaw) > 0 {
			_ = json.Unmarshal(sellsRaw, &c.SellBrokers)
		}
		if c.SellBrokers == nil {
			c.SellBrokers = map[string]float64{}
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// ListLatestCandlesForExplorer loads the latest candle per ticker for explorer/scan.
func (r *Repository) ListLatestCandlesForExplorer(ctx context.Context) ([]ExplorerCandleRow, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT DISTINCT ON (c.ticker)
		       c.ticker, m.company_name, c.close, c.prev, c.open,
		       c.volume, c.tavg, c.buy_brokers, c.sell_brokers
		FROM terminal_engine.mmd_daily_candles c
		JOIN terminal_engine.mmd_tickers_meta m ON c.ticker = m.ticker
		ORDER BY c.ticker, c.date_iso DESC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]ExplorerCandleRow, 0, 900)
	for rows.Next() {
		var row ExplorerCandleRow
		var buysRaw, sellsRaw []byte
		if err := rows.Scan(
			&row.Ticker, &row.CompanyName, &row.Close, &row.Prev, &row.Open,
			&row.Volume, &row.Tavg, &buysRaw, &sellsRaw,
		); err != nil {
			return nil, err
		}
		if len(buysRaw) > 0 {
			_ = json.Unmarshal(buysRaw, &row.BuyBrokers)
		}
		if row.BuyBrokers == nil {
			row.BuyBrokers = map[string]float64{}
		}
		if len(sellsRaw) > 0 {
			_ = json.Unmarshal(sellsRaw, &row.SellBrokers)
		}
		if row.SellBrokers == nil {
			row.SellBrokers = map[string]float64{}
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

// ListBigPlayerT5Summary returns T5-active tickers ordered by most recent movement_date.
func (r *Repository) ListBigPlayerT5Summary(ctx context.Context, limit int) ([]BigPlayerRow, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT DISTINCT t.ticker, m.last_close, m.last_volume, m.last_tavg,
		       COUNT(t.movement_date) AS t5_count,
		       MAX(t.movement_date) AS last_t5_date
		FROM ownership_staging.v_t5_transactions t
		JOIN terminal_engine.mmd_tickers_meta m ON t.ticker = m.ticker
		GROUP BY t.ticker, m.last_close, m.last_volume, m.last_tavg
		ORDER BY last_t5_date DESC
		LIMIT $1;
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]BigPlayerRow, 0, limit)
	for rows.Next() {
		var row BigPlayerRow
		var lastT5 time.Time
		if err := rows.Scan(
			&row.Ticker, &row.LastClose, &row.LastVolume, &row.LastTavg,
			&row.T5Count, &lastT5,
		); err != nil {
			return nil, err
		}
		row.LastT5Date = lastT5.Format("2006-01-02")
		out = append(out, row)
	}
	return out, rows.Err()
}

// ListT5Transactions loads T5 records for a ticker ordered by movement_date DESC.
func (r *Repository) ListT5Transactions(ctx context.Context, ticker string, limit int) ([]T5Record, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT movement_date, investor_name, broker, action,
		       change_lot, value_m, price, change_pct
		FROM ownership_staging.v_t5_transactions
		WHERE ticker = $1
		ORDER BY movement_date DESC
		LIMIT $2;
	`, ticker, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]T5Record, 0, limit)
	for rows.Next() {
		var rec T5Record
		var moveDate time.Time
		var lot, valM, price, pct *float64
		if err := rows.Scan(
			&moveDate, &rec.Investor, &rec.Broker, &rec.Action,
			&lot, &valM, &price, &pct,
		); err != nil {
			return nil, err
		}
		rec.Date = moveDate.Format("2006-01-02")
		if lot != nil {
			rec.Lot = *lot
		}
		if valM != nil {
			rec.ValueM = *valM
		}
		if price != nil {
			rec.Price = *price
		}
		if pct != nil {
			rec.Pct = *pct
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

// GetGroupPerformance fetches real conglomerate and sectoral performance matching mmd.py.
func (r *Repository) GetGroupPerformance(ctx context.Context, ticker string) (*GroupPerformance, error) {
	res := &GroupPerformance{}

	var groupName string
	var tickersArr []string
	err := r.DB.QueryRow(ctx, `
		SELECT group_name, tickers_arr
		FROM terminal_engine.konglo_groups
		WHERE $1 = ANY(tickers_arr)
		LIMIT 1;
	`, ticker).Scan(&groupName, &tickersArr)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if err == nil && len(tickersArr) > 0 {
		members, qErr := r.queryGroupMembers(ctx, tickersArr)
		if qErr != nil {
			return nil, qErr
		}
		res.Conglomerate = &GroupPerfBlock{
			GroupName: groupName,
			Title:     fmt.Sprintf("Kinerja Konglomerasi %s", groupName),
			Type:      "Konglomerasi",
			Members:   members,
		}
	}

	var sector, subSector *string
	err = r.DB.QueryRow(ctx, `
		SELECT sector, sub_sector
		FROM ownership_staging.company_meta
		WHERE ticker = $1
		LIMIT 1;
	`, ticker).Scan(&sector, &subSector)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	_ = subSector
	if err == nil && sector != nil && *sector != "" {
		secTickers, sErr := r.listSectorTickers(ctx, *sector)
		if sErr != nil {
			return nil, sErr
		}
		if len(secTickers) > 0 {
			members, qErr := r.queryGroupMembers(ctx, secTickers)
			if qErr != nil {
				return nil, qErr
			}
			res.Sector = &GroupPerfBlock{
				GroupName: *sector,
				Title:     fmt.Sprintf("Kinerja Sektor %s", *sector),
				Type:      "Sektor",
				Members:   members,
			}
		}
	}

	return res, nil
}

func (r *Repository) listSectorTickers(ctx context.Context, sectorName string) ([]string, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT ticker
		FROM ownership_staging.company_meta
		WHERE sector = $1 AND status = 'STATUS_ACTIVE'
		LIMIT 15;
	`, sectorName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]string, 0, 15)
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repository) queryGroupMembers(ctx context.Context, tickers []string) ([]GroupPerfMember, error) {
	rows, err := r.DB.Query(ctx, `
		WITH ranked AS (
			SELECT ticker, date_iso, close, prev, volume, tavg,
			       ROW_NUMBER() OVER (PARTITION BY ticker ORDER BY date_iso DESC) AS rn
			FROM terminal_engine.mmd_daily_candles
			WHERE ticker = ANY($1)
		)
		SELECT
			r1.ticker,
			r1.close AS last_close,
			ROUND(((r1.close - r1.prev) / NULLIF(r1.prev, 0) * 100)::numeric, 2)::double precision AS today_pct,
			ROUND(((r1.close - r6.close) / NULLIF(r6.close, 0) * 100)::numeric, 2)::double precision AS w1_pct,
			ROUND(((r1.close - r20.close) / NULLIF(r20.close, 0) * 100)::numeric, 2)::double precision AS m1_pct,
			(SELECT ROUND(AVG(volume * COALESCE(tavg, close))::numeric, 0)::double precision
			 FROM terminal_engine.mmd_daily_candles
			 WHERE ticker = r1.ticker AND date_iso >= r1.date_iso - INTERVAL '30 days') AS avg_val
		FROM ranked r1
		LEFT JOIN ranked r6 ON r1.ticker = r6.ticker AND r6.rn = 6
		LEFT JOIN ranked r20 ON r1.ticker = r20.ticker AND r20.rn = 20
		WHERE r1.rn = 1
		ORDER BY today_pct DESC;
	`, tickers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]GroupPerfMember, 0, len(tickers))
	for rows.Next() {
		var m GroupPerfMember
		var lastClose, todayPct, m1Pct, avgVal *float64
		if err := rows.Scan(&m.Ticker, &lastClose, &todayPct, &m.W1Pct, &m1Pct, &avgVal); err != nil {
			return nil, err
		}
		if lastClose != nil {
			m.LastClose = *lastClose
		}
		if todayPct != nil {
			m.TodayPct = *todayPct
		}
		if m1Pct != nil {
			m.M1Pct = *m1Pct
		}
		if avgVal != nil {
			m.AvgVal = *avgVal
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// GetLatestMacroRepoRatio queries v_ksei_overbooking_monthly matching mmd.py lines 380-392.
func (r *Repository) GetLatestMacroRepoRatio(ctx context.Context) float64 {
	var fopRatio *float64
	err := r.DB.QueryRow(ctx, `
		SELECT fop_ratio
		FROM ownership_staging.v_ksei_overbooking_monthly
		WHERE fop_ratio IS NOT NULL
		ORDER BY period_month DESC
		LIMIT 1;
	`).Scan(&fopRatio)
	if err != nil || fopRatio == nil {
		return 0.0
	}
	return *fopRatio
}

// GetBrokerFirmName looks up the latest firm_name in idx_broker_summary, falling back to broker code if absent.
func (r *Repository) GetBrokerFirmName(ctx context.Context, broker string) (*string, error) {
	var firmName *string
	err := r.DB.QueryRow(ctx, `
		SELECT firm_name
		FROM public.idx_broker_summary
		WHERE id_firm = $1
		ORDER BY date DESC
		LIMIT 1;
	`, broker).Scan(&firmName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			b := broker
			return &b, nil
		}
		return nil, err
	}
	return firmName, nil
}

// ListBrokerActivity loads buy/sell lot and value records for a ticker and broker.
func (r *Repository) ListBrokerActivity(ctx context.Context, ticker, broker, dateFrom, dateTo string) ([]BrokerActivityRow, error) {
	var b strings.Builder
	b.WriteString(`SELECT date, buy_lot, sell_lot, buy_value, sell_value FROM public.broker_activity WHERE ticker = $1 AND broker_code = $2`)
	args := []any{ticker, broker}
	argIdx := 3
	if dateFrom != "" {
		b.WriteString(fmt.Sprintf(" AND date >= $%d", argIdx))
		args = append(args, dateFrom)
		argIdx++
	}
	if dateTo != "" {
		b.WriteString(fmt.Sprintf(" AND date <= $%d", argIdx))
		args = append(args, dateTo)
	}
	b.WriteString(" ORDER BY date DESC;")

	rows, err := r.DB.Query(ctx, b.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]BrokerActivityRow, 0, 128)
	for rows.Next() {
		var row BrokerActivityRow
		var dateVal time.Time
		if err := rows.Scan(&dateVal, &row.BuyLot, &row.SellLot, &row.BuyVal, &row.SellVal); err != nil {
			return nil, err
		}
		row.Date = dateVal.Format("2006-01-02")
		out = append(out, row)
	}
	return out, rows.Err()
}

func formatDatePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}
