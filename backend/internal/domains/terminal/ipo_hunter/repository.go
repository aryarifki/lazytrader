package ipohunter

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository performs indexed PostgreSQL reads for IPO Hunter without permanent daily RAM caches.
type Repository struct {
	DB *pgxpool.Pool
}

func decodeFloatMap(raw []byte) map[string]float64 {
	if len(raw) == 0 || string(raw) == "null" {
		return map[string]float64{}
	}
	var out map[string]float64
	if err := json.Unmarshal(raw, &out); err != nil || out == nil {
		return map[string]float64{}
	}
	return out
}

func decodeUnderwriters(raw []byte) []Underwriter {
	if len(raw) == 0 || string(raw) == "null" {
		return []Underwriter{}
	}
	var out []Underwriter
	if err := json.Unmarshal(raw, &out); err != nil || out == nil {
		return []Underwriter{}
	}
	return out
}

func decodeDayReturns(raw []byte) []*float64 {
	if len(raw) == 0 || string(raw) == "null" {
		return []*float64{}
	}
	var out []*float64
	if err := json.Unmarshal(raw, &out); err != nil || out == nil {
		return []*float64{}
	}
	return out
}

const masterSelectCols = `
	ticker, company_name, listing_date, ipo_price::double precision,
	total_listed_shares, offering_shares, pct_ipo, nilai_ipo::double precision,
	papan, lead_underwriter, underwriters, penjatahan::double precision,
	jumlah_sid, day_returns, perf_2w::double precision, perf_1m::double precision,
	perf_3m::double precision, perf_6m::double precision
`

func scanMaster(s interface{ Scan(...any) error }) (MasterRecord, error) {
	var m MasterRecord
	var uwRaw, drRaw []byte
	err := s.Scan(
		&m.Ticker, &m.CompanyName, &m.ListingDate, &m.IPOPrice,
		&m.TotalListedShares, &m.OfferingShares, &m.PctIPO, &m.NilaiIPO,
		&m.Papan, &m.LeadUnderwriter, &uwRaw, &m.Penjatahan,
		&m.JumlahSID, &drRaw, &m.Perf2W, &m.Perf1M,
		&m.Perf3M, &m.Perf6M,
	)
	if err != nil {
		return m, err
	}
	m.Underwriters = decodeUnderwriters(uwRaw)
	m.DayReturns = decodeDayReturns(drRaw)
	return m, nil
}

// ListMasters loads all master IPO rows ordered identically to legacy Python.
func (r *Repository) ListMasters(ctx context.Context) ([]MasterRecord, error) {
	rows, err := r.DB.Query(ctx, `SELECT `+masterSelectCols+` FROM terminal_engine.ipo_hunter_master ORDER BY listing_date DESC NULLS LAST, ticker ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]MasterRecord, 0, 320)
	for rows.Next() {
		m, err := scanMaster(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// GetMasterByTicker loads a single master IPO row, returning (nil, nil) if absent.
func (r *Repository) GetMasterByTicker(ctx context.Context, ticker string) (*MasterRecord, error) {
	row := r.DB.QueryRow(ctx, `SELECT `+masterSelectCols+` FROM terminal_engine.ipo_hunter_master WHERE ticker = $1`, ticker)
	m, err := scanMaster(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func scanDayRow(s interface{ Scan(...any) error }, isD0 bool) (DayRow, error) {
	var d DayRow
	var buyRaw, sellRaw []byte
	err := s.Scan(
		&d.Date, &d.Ticker,
		&d.Open, &d.High, &d.Low, &d.Close, &d.Prev, &d.Tavg,
		&d.Volume, &d.NetVolume, &buyRaw, &sellRaw,
	)
	if err != nil {
		return d, err
	}
	d.BuyBrokers = decodeFloatMap(buyRaw)
	d.SellBrokers = decodeFloatMap(sellRaw)
	d.IsD0 = isD0
	d.DayCount = 1
	return d, nil
}

// GetD0ByTicker loads the D0 baseline row for one ticker, returning (nil, nil) if absent.
func (r *Repository) GetD0ByTicker(ctx context.Context, ticker string) (*DayRow, error) {
	row := r.DB.QueryRow(ctx, `
		SELECT d0_date, ticker, open::double precision, high::double precision,
		       low::double precision, close::double precision, prev::double precision,
		       tavg::double precision, volume, net_volume, buy_brokers, sell_brokers
		FROM terminal_engine.ipo_hunter_d0
		WHERE ticker = $1
	`, ticker)
	d, err := scanDayRow(row, true)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// ListAllD0Map loads all 317 D0 rows keyed by ticker (only used by after=true screener).
func (r *Repository) ListAllD0Map(ctx context.Context) (map[string]*DayRow, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT d0_date, ticker, open::double precision, high::double precision,
		       low::double precision, close::double precision, prev::double precision,
		       tavg::double precision, volume, net_volume, buy_brokers, sell_brokers
		FROM terminal_engine.ipo_hunter_d0
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]*DayRow, 320)
	for rows.Next() {
		d, err := scanDayRow(rows, true)
		if err != nil {
			return nil, err
		}
		copyRow := d
		out[d.Ticker] = &copyRow
	}
	return out, rows.Err()
}

// GetAnyD0Date fetches the representative D0 date from terminal_engine.ipo_hunter_d0.
func (r *Repository) GetAnyD0Date(ctx context.Context) (*time.Time, error) {
	var dt *time.Time
	err := r.DB.QueryRow(ctx, `SELECT d0_date FROM terminal_engine.ipo_hunter_d0 WHERE d0_date IS NOT NULL LIMIT 1`).Scan(&dt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return dt, err
}

// ListDistinctDailyDates returns sorted trading dates from ipo_hunter_daily.
func (r *Repository) ListDistinctDailyDates(ctx context.Context) ([]time.Time, error) {
	rows, err := r.DB.Query(ctx, `SELECT DISTINCT date_iso FROM terminal_engine.ipo_hunter_daily ORDER BY date_iso ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dates := make([]time.Time, 0, 128)
	for rows.Next() {
		var dt time.Time
		if err := rows.Scan(&dt); err != nil {
			return nil, err
		}
		dates = append(dates, dt.UTC())
	}
	return dates, rows.Err()
}

// ListDistinctDailyTickers returns the set of tickers that have >= 1 daily row.
func (r *Repository) ListDistinctDailyTickers(ctx context.Context) (map[string]struct{}, error) {
	rows, err := r.DB.Query(ctx, `SELECT DISTINCT ticker FROM terminal_engine.ipo_hunter_daily`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]struct{}, 320)
	for rows.Next() {
		var tk string
		if err := rows.Scan(&tk); err != nil {
			return nil, err
		}
		out[tk] = struct{}{}
	}
	return out, rows.Err()
}

// ListDailyByTicker loads sorted daily rows for a single ticker via primary key index.
func (r *Repository) ListDailyByTicker(ctx context.Context, ticker string) ([]DayRow, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT date_iso, ticker, open::double precision, high::double precision,
		       low::double precision, close::double precision, prev::double precision,
		       tavg::double precision, volume, net_volume, buy_brokers, sell_brokers
		FROM terminal_engine.ipo_hunter_daily
		WHERE ticker = $1
		ORDER BY date_iso ASC
	`, ticker)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]DayRow, 0, 100)
	for rows.Next() {
		d, err := scanDayRow(rows, false)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// ListHolidays returns IDX holiday ordinals and max holiday year (default 2027).
func (r *Repository) ListHolidays(ctx context.Context) (map[int]struct{}, int, error) {
	rows, err := r.DB.Query(ctx, `SELECT holiday_date FROM terminal_engine.ipo_hunter_holidays WHERE holiday_date IS NOT NULL`)
	if err != nil {
		return nil, 2027, err
	}
	defer rows.Close()

	holidays := make(map[int]struct{}, 192)
	maxYear := 0
	for rows.Next() {
		var dt time.Time
		if err := rows.Scan(&dt); err != nil {
			return nil, 2027, err
		}
		holidays[DateOrdinal(dt)] = struct{}{}
		if y := dt.UTC().Year(); y > maxYear {
			maxYear = y
		}
	}
	if maxYear == 0 {
		maxYear = 2027
	}
	return holidays, maxYear, rows.Err()
}

// ListBaseVolumesBefore loads up to 20 positive volumes strictly before loDate per ticker without JSONB.
func (r *Repository) ListBaseVolumesBefore(ctx context.Context, loDate time.Time, tickers []string) (map[string][]float64, error) {
	const sqlWithFilter = `
		WITH base AS (
			SELECT ticker, date_iso, volume,
			       ROW_NUMBER() OVER (PARTITION BY ticker ORDER BY date_iso DESC) AS rn
			FROM terminal_engine.ipo_hunter_daily
			WHERE date_iso < $1 AND volume IS NOT NULL AND volume > 0 AND ticker = ANY($2)
		)
		SELECT ticker, volume FROM base WHERE rn <= 20 ORDER BY ticker ASC, date_iso ASC
	`
	const sqlAll = `
		WITH base AS (
			SELECT ticker, date_iso, volume,
			       ROW_NUMBER() OVER (PARTITION BY ticker ORDER BY date_iso DESC) AS rn
			FROM terminal_engine.ipo_hunter_daily
			WHERE date_iso < $1 AND volume IS NOT NULL AND volume > 0
		)
		SELECT ticker, volume FROM base WHERE rn <= 20 ORDER BY ticker ASC, date_iso ASC
	`
	var rows pgx.Rows
	var err error
	if len(tickers) > 0 {
		rows, err = r.DB.Query(ctx, sqlWithFilter, loDate, tickers)
	} else {
		rows, err = r.DB.Query(ctx, sqlAll, loDate)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string][]float64, 320)
	for rows.Next() {
		var tk string
		var vol int64
		if err := rows.Scan(&tk, &vol); err != nil {
			return nil, err
		}
		out[tk] = append(out[tk], float64(vol))
	}
	return out, rows.Err()
}

// StreamDailyPeriodic streams daily rows in [loDate, hiDate] grouped by ticker in ascending ticker+date order.
func (r *Repository) StreamDailyPeriodic(
	ctx context.Context,
	loDate, hiDate time.Time,
	tickers []string,
	onTicker func(ticker string, days []DayRow),
) error {
	const sqlWithFilter = `
		SELECT date_iso, ticker, open::double precision, high::double precision,
		       low::double precision, close::double precision, prev::double precision,
		       tavg::double precision, volume, net_volume, buy_brokers, sell_brokers
		FROM terminal_engine.ipo_hunter_daily
		WHERE date_iso >= $1 AND date_iso <= $2 AND ticker = ANY($3)
		ORDER BY ticker ASC, date_iso ASC
	`
	const sqlAll = `
		SELECT date_iso, ticker, open::double precision, high::double precision,
		       low::double precision, close::double precision, prev::double precision,
		       tavg::double precision, volume, net_volume, buy_brokers, sell_brokers
		FROM terminal_engine.ipo_hunter_daily
		WHERE date_iso >= $1 AND date_iso <= $2
		ORDER BY ticker ASC, date_iso ASC
	`
	var rows pgx.Rows
	var err error
	if len(tickers) > 0 {
		rows, err = r.DB.Query(ctx, sqlWithFilter, loDate, hiDate, tickers)
	} else {
		rows, err = r.DB.Query(ctx, sqlAll, loDate, hiDate)
	}
	if err != nil {
		return err
	}
	return streamGroupedRows(rows, onTicker)
}

// StreamDailyAfterIPO streams daily rows <= hiDate strictly after each ticker's D0 date (or all rows if no D0).
func (r *Repository) StreamDailyAfterIPO(
	ctx context.Context,
	hiDate time.Time,
	tickers []string,
	onTicker func(ticker string, days []DayRow),
) error {
	const sqlWithFilter = `
		SELECT dy.date_iso, dy.ticker, dy.open::double precision, dy.high::double precision,
		       dy.low::double precision, dy.close::double precision, dy.prev::double precision,
		       dy.tavg::double precision, dy.volume, dy.net_volume, dy.buy_brokers, dy.sell_brokers
		FROM terminal_engine.ipo_hunter_daily dy
		LEFT JOIN terminal_engine.ipo_hunter_d0 d0 ON d0.ticker = dy.ticker
		WHERE dy.date_iso <= $1
		  AND (d0.d0_date IS NULL OR dy.date_iso > d0.d0_date)
		  AND dy.ticker = ANY($2)
		ORDER BY dy.ticker ASC, dy.date_iso ASC
	`
	const sqlAll = `
		SELECT dy.date_iso, dy.ticker, dy.open::double precision, dy.high::double precision,
		       dy.low::double precision, dy.close::double precision, dy.prev::double precision,
		       dy.tavg::double precision, dy.volume, dy.net_volume, dy.buy_brokers, dy.sell_brokers
		FROM terminal_engine.ipo_hunter_daily dy
		LEFT JOIN terminal_engine.ipo_hunter_d0 d0 ON d0.ticker = dy.ticker
		WHERE dy.date_iso <= $1
		  AND (d0.d0_date IS NULL OR dy.date_iso > d0.d0_date)
		ORDER BY dy.ticker ASC, dy.date_iso ASC
	`
	var rows pgx.Rows
	var err error
	if len(tickers) > 0 {
		rows, err = r.DB.Query(ctx, sqlWithFilter, hiDate, tickers)
	} else {
		rows, err = r.DB.Query(ctx, sqlAll, hiDate)
	}
	if err != nil {
		return err
	}
	return streamGroupedRows(rows, onTicker)
}

func streamGroupedRows(rows pgx.Rows, onTicker func(ticker string, days []DayRow)) error {
	defer rows.Close()
	curTicker := ""
	buf := make([]DayRow, 0, 64)
	for rows.Next() {
		d, err := scanDayRow(rows, false)
		if err != nil {
			return err
		}
		if curTicker != "" && d.Ticker != curTicker {
			onTicker(curTicker, buf)
			buf = make([]DayRow, 0, 64)
		}
		curTicker = d.Ticker
		buf = append(buf, d)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if curTicker != "" && len(buf) > 0 {
		onTicker(curTicker, buf)
	}
	return nil
}
