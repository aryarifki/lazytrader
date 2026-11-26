package bandarmology

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository queries prices, broker_flow, broker_activity, tickers, and company_board.
type Repository struct {
	DB *pgxpool.Pool
}

// ReadPrices mirrors storage.read_prices with optional tickers and date bounds.
func (r *Repository) ReadPrices(ctx context.Context, tickers []string, startDate, endDate *time.Time) ([]PriceRow, error) {
	clauses, args := buildFilterClauses(tickers, startDate, endDate)
	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}
	q := `SELECT date, ticker, open::double precision, high::double precision, low::double precision, close::double precision, volume FROM prices` + where + ` ORDER BY ticker, date ASC`
	rows, err := r.DB.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]PriceRow, 0, 256)
	for rows.Next() {
		var p PriceRow
		if err := rows.Scan(&p.Date, &p.Ticker, &p.Open, &p.High, &p.Low, &p.Close, &p.Volume); err != nil {
			return nil, err
		}
		p.Open = shared.FinitePtr(p.Open)
		p.High = shared.FinitePtr(p.High)
		p.Low = shared.FinitePtr(p.Low)
		p.Close = shared.FinitePtr(p.Close)
		out = append(out, p)
	}
	return out, rows.Err()
}

// ReadBrokerFlow mirrors storage.read_broker_flow.
func (r *Repository) ReadBrokerFlow(ctx context.Context, tickers []string, startDate, endDate *time.Time) ([]FlowRow, error) {
	clauses, args := buildFilterClauses(tickers, startDate, endDate)
	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}
	q := `
		SELECT date, ticker, bandar_signal, bandar_signal_score::double precision,
		       foreign_net_broker::double precision, local_net_broker::double precision, gov_net_broker::double precision,
		       foreign_net_flow::double precision, domestic_net_flow::double precision, total_value::double precision,
		       foreign_signal, conclusion_broker, conclusion_flow, fetched_at
		FROM broker_flow` + where + ` ORDER BY ticker, date ASC`
	rows, err := r.DB.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]FlowRow, 0, 256)
	for rows.Next() {
		var f FlowRow
		if err := rows.Scan(
			&f.Date, &f.Ticker, &f.BandarSignal, &f.BandarSignalScore,
			&f.ForeignNetBroker, &f.LocalNetBroker, &f.GovNetBroker,
			&f.ForeignNetFlow, &f.DomesticNetFlow, &f.TotalValue,
			&f.ForeignSignal, &f.ConclusionBroker, &f.ConclusionFlow, &f.FetchedAt,
		); err != nil {
			return nil, err
		}
		f.BandarSignalScore = shared.FinitePtr(f.BandarSignalScore)
		f.ForeignNetBroker = shared.FinitePtr(f.ForeignNetBroker)
		f.LocalNetBroker = shared.FinitePtr(f.LocalNetBroker)
		f.GovNetBroker = shared.FinitePtr(f.GovNetBroker)
		f.ForeignNetFlow = shared.FinitePtr(f.ForeignNetFlow)
		f.DomesticNetFlow = shared.FinitePtr(f.DomesticNetFlow)
		f.TotalValue = shared.FinitePtr(f.TotalValue)
		out = append(out, f)
	}
	return out, rows.Err()
}

// ReadBrokerActivity mirrors storage.read_broker_activity.
func (r *Repository) ReadBrokerActivity(ctx context.Context, tickers []string, startDate, endDate *time.Time) ([]ActivityRow, error) {
	clauses, args := buildFilterClauses(tickers, startDate, endDate)
	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}
	q := `
		SELECT date, ticker, broker_code, COALESCE(participant_type, ''),
		       COALESCE(buy_value, 0)::double precision, COALESCE(sell_value, 0)::double precision,
		       COALESCE(net_value, 0)::double precision,
		       buy_lot::double precision, sell_lot::double precision, frequency::double precision,
		       buy_avg_price::double precision, sell_avg_price::double precision, fetched_at
		FROM broker_activity` + where + ` ORDER BY ticker, date ASC, net_value DESC`
	rows, err := r.DB.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]ActivityRow, 0, 512)
	for rows.Next() {
		var a ActivityRow
		if err := rows.Scan(
			&a.Date, &a.Ticker, &a.BrokerCode, &a.ParticipantType,
			&a.BuyValue, &a.SellValue, &a.NetValue,
			&a.BuyLot, &a.SellLot, &a.Frequency,
			&a.BuyAvgPrice, &a.SellAvgPrice, &a.FetchedAt,
		); err != nil {
			return nil, err
		}
		a.BuyLot = shared.FinitePtr(a.BuyLot)
		a.SellLot = shared.FinitePtr(a.SellLot)
		a.Frequency = shared.FinitePtr(a.Frequency)
		a.BuyAvgPrice = shared.FinitePtr(a.BuyAvgPrice)
		a.SellAvgPrice = shared.FinitePtr(a.SellAvgPrice)
		out = append(out, a)
	}
	return out, rows.Err()
}

func buildFilterClauses(tickers []string, startDate, endDate *time.Time) ([]string, []any) {
	var clauses []string
	var args []any
	idx := 1
	if len(tickers) > 0 {
		up := make([]string, len(tickers))
		for i, t := range tickers {
			up[i] = strings.ToUpper(strings.TrimSpace(t))
		}
		clauses = append(clauses, fmt.Sprintf("ticker = ANY($%d)", idx))
		args = append(args, up)
		idx++
	}
	if startDate != nil {
		clauses = append(clauses, fmt.Sprintf("date >= $%d", idx))
		args = append(args, *startDate)
		idx++
	}
	if endDate != nil {
		clauses = append(clauses, fmt.Sprintf("date <= $%d", idx))
		args = append(args, *endDate)
	}
	return clauses, args
}

// GetMaxActivityDate returns MAX(date) from broker_activity for ticker.
func (r *Repository) GetMaxActivityDate(ctx context.Context, ticker string) *time.Time {
	var d *time.Time
	_ = r.DB.QueryRow(ctx, `SELECT MAX(date) FROM broker_activity WHERE ticker = $1`, ticker).Scan(&d)
	return d
}

// GetDistinctActivityDates returns sorted ascending ISO dates for ticker.
func (r *Repository) GetDistinctActivityDates(ctx context.Context, ticker string) ([]string, error) {
	rows, err := r.DB.Query(ctx, `SELECT DISTINCT date FROM broker_activity WHERE ticker = $1 ORDER BY date ASC`, ticker)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]string, 0, 128)
	for rows.Next() {
		var d time.Time
		if err := rows.Scan(&d); err == nil {
			out = append(out, fmtDateISO(d))
		}
	}
	return out, rows.Err()
}

// GetSystemStatus returns active ticker count and latest broker_flow date.
func (r *Repository) GetSystemStatus(ctx context.Context) (int64, string, error) {
	var cnt int64
	if err := r.DB.QueryRow(ctx, `SELECT COUNT(*) FROM tickers WHERE is_active = TRUE`).Scan(&cnt); err != nil {
		return 0, "-", err
	}
	var latest *time.Time
	if err := r.DB.QueryRow(ctx, `SELECT MAX(date) FROM broker_flow`).Scan(&latest); err != nil {
		return cnt, "-", err
	}
	if latest == nil {
		return cnt, "-", nil
	}
	return cnt, fmtDateISO(*latest), nil
}

// GetCompanyMeta fetches company name and conglomerate/sector group for ticker.
func (r *Repository) GetCompanyMeta(ctx context.Context, ticker string) (name, group string) {
	name = "Bursa Efek Indonesia"
	group = "Sektor Perusahaan"
	var n, g *string
	err := r.DB.QueryRow(ctx, `
		SELECT t.name, COALESCE(cg.groups, t.sector) AS group_name
		FROM tickers t
		LEFT JOIN (
			SELECT ticker, STRING_AGG(DISTINCT group_name, ' / ' ORDER BY group_name) AS groups
			FROM ownership_staging.v_conglomerate_members
			WHERE ticker = $1
			GROUP BY ticker
		) cg ON cg.ticker = t.ticker
		WHERE t.ticker = $1
		LIMIT 1
	`, ticker).Scan(&n, &g)
	if err != nil && err != pgx.ErrNoRows {
		return name, group
	}
	if n != nil && *n != "" {
		name = *n
	}
	if g != nil && *g != "" {
		group = *g
	}
	return name, group
}

// CheckDynamicBoardMatch queries company_board for Jusuf/Danatama interlocks when broker == "II".
func (r *Repository) CheckDynamicBoardMatch(ctx context.Context, ticker, brokerCode string) map[string]string {
	code := strings.ToUpper(strings.TrimSpace(brokerCode))
	if code != "II" || r == nil || r.DB == nil {
		return nil
	}
	bInfo := GetBrokerInfo(code)
	rows, err := r.DB.Query(ctx, `
		SELECT name, role, title FROM company_board
		WHERE ticker = $1 AND (name ILIKE '%Jusuf%' OR name ILIKE '%Danatama%')
	`, ticker)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		var role, title *string
		if rows.Scan(&name, &role, &title) == nil {
			lbl := ""
			if title != nil && *title != "" {
				lbl = *title
			} else if role != nil {
				lbl = *role
			}
			names = append(names, fmt.Sprintf("%s (%s)", name, lbl))
		}
	}
	if len(names) == 0 {
		return nil
	}
	return map[string]string{
		"badge":         "INTERLOCKING BOARD",
		"relation_type": "BOARD_INTERLOCK",
		"entity":        bInfo["name"],
		"relation":      fmt.Sprintf("Terdeteksi kesamaan nama pengurus emiten dengan pendiri/komisaris Danatama: %s.", strings.Join(names, ", ")),
		"confidence":    "HIGH",
	}
}
