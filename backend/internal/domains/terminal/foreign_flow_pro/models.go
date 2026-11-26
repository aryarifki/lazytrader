package foreign_flow_pro

import "time"

// ScreenerDBRow represents a row from terminal_engine.foreign_flow_screener.
type ScreenerDBRow struct {
	Ticker          string
	Name            *string
	Badges          []string
	ScreenerSummary map[string]any
	UpdatedAt       *time.Time
}

// ScreenerItem matches one entry in GET /api/v1/terminal/foreign-flow/screener.
type ScreenerItem struct {
	Ticker              string         `json:"ticker"`
	Name                string         `json:"name"`
	Close               float64        `json:"close"`
	ChangePct           float64        `json:"changePct"`
	FPBadge             *string        `json:"fpBadge"`
	FPScore             float64        `json:"fpScore"`
	AvgTransValue10d    float64        `json:"avgTransValue10d"`
	AccumVal            float64        `json:"accumVal"`
	PriceChangeWindow   float64        `json:"priceChangeWindow"`
	Badges              []string       `json:"badges"`
	DailyClose          any            `json:"dailyClose"`
	DailyValue          any            `json:"dailyValue"`
	DailyFFValue        any            `json:"dailyFFValue"`
	DailyFCoValue       any            `json:"dailyFCoValue"`
	DailyKFValue        any            `json:"dailyKFValue"`
	DailyXLValue        any            `json:"dailyXLValue"`
	DailySemiRitelValue any            `json:"dailySemiRitelValue"`
	MaxDailyFF60d       any            `json:"maxDailyFF60d"`
	MinDailyFF60d       any            `json:"minDailyFF60d"`
	Summary             map[string]any `json:"summary"`
}

// DailyPoint represents a single chronological trading day input to analytics.
type DailyPoint struct {
	Date             string
	Open             *float64
	High             *float64
	Low              *float64
	Close            *float64
	Volume           *float64
	ForeignNetVal    float64
	PartisipasiAsing *float64
	RawSeries        []any
}

// FPPoint represents Foreign Power 2021 score, badge, and color for one day.
type FPPoint struct {
	Score float64 `json:"score"`
	Badge string  `json:"badge"`
	Color string  `json:"color"`
}

// SpecialCandle represents a detected special candlestick signal.
type SpecialCandle struct {
	Date  string `json:"date"`
	Type  string `json:"type"`
	Color string `json:"color"`
	Label string `json:"label"`
}

// CorrelationStats holds smoothed indicator correlations for 5 flow series.
type CorrelationStats struct {
	FF        *float64 `json:"ff"`
	FCo       *float64 `json:"fco"`
	KF        *float64 `json:"kf"`
	XL        *float64 `json:"xl"`
	SemiRitel *float64 `json:"semiritel"`
}

// ConsecutiveStats holds average consecutive foreign inflow streak lengths.
type ConsecutiveStats struct {
	D50  *float64 `json:"50d"`
	D100 *float64 `json:"100d"`
	D200 *float64 `json:"200d"`
}

// InflowNaikWindow holds inflow & price-up count and percentage for a window.
type InflowNaikWindow struct {
	Count int      `json:"count"`
	Pct   *float64 `json:"pct"`
}

// InflowNaikStats holds 50d/100d/200d inflow-naik breakdowns.
type InflowNaikStats struct {
	D50  InflowNaikWindow `json:"50d"`
	D100 InflowNaikWindow `json:"100d"`
	D200 InflowNaikWindow `json:"200d"`
}

// MarkdownStats holds foreign-inflow markdown day counts across windows.
type MarkdownStats struct {
	D50  int `json:"50d"`
	D100 int `json:"100d"`
	D200 int `json:"200d"`
}

// SignalHorizon holds horizon win-rate statistics for 5d, 10d, or 20d.
type SignalHorizon struct {
	N      int     `json:"n"`
	Up     int     `json:"up"`
	Pct    float64 `json:"pct"`
	Median float64 `json:"median"`
}

// SignalTypeStats holds event count and horizon map for one signal type.
type SignalTypeStats struct {
	Events   int                      `json:"events"`
	Horizons map[string]SignalHorizon `json:"horizons"`
}

// SignalStats holds horizon statistics for all 4 signal triggers.
type SignalStats struct {
	StreakInflow SignalTypeStats `json:"streak_inflow"`
	MassiveAcc   SignalTypeStats `json:"massive_acc"`
	Markup       SignalTypeStats `json:"markup"`
	SuddenSell   SignalTypeStats `json:"sudden_sell"`
}

// AkumulasiStats holds directional VWAP and net foreign accumulation summary.
type AkumulasiStats struct {
	AvgPrice   *int64  `json:"avg_price"`
	TotalFlowM float64 `json:"total_flow_m"`
	Status     string  `json:"status"`
}

// Statistics holds all statistical tables inside analytics.
type Statistics struct {
	Correlations      CorrelationStats `json:"correlations"`
	ConsecutiveInflow ConsecutiveStats `json:"consecutive_inflow"`
	InflowHargaNaik   InflowNaikStats  `json:"inflow_harga_naik"`
	MarkDown          MarkdownStats    `json:"mark_down"`
	SignalStats       SignalStats      `json:"signal_stats"`
	AkumulasiDefault  AkumulasiStats   `json:"akumulasi_default"`
}

// Analytics matches the exact payload returned by process_foreign_flow_analytics.
type Analytics struct {
	Dates                    []string        `json:"dates"`
	Opens                    []float64       `json:"opens"`
	Highs                    []float64       `json:"highs"`
	Lows                     []float64       `json:"lows"`
	Closes                   []float64       `json:"closes"`
	ValuesMiliar             []float64       `json:"values_miliar"`
	FFDailyVal               []float64       `json:"ff_daily_val"`
	FCoDailyVal              []float64       `json:"fco_daily_val"`
	KFDailyVal               []float64       `json:"kf_daily_val"`
	XLDailyVal               []float64       `json:"xl_daily_val"`
	SemiRitelDailyVal        []float64       `json:"semiritel_daily_val"`
	CumFF                    []float64       `json:"cum_ff"`
	CumFCo                   []float64       `json:"cum_fco"`
	CumKF                    []float64       `json:"cum_kf"`
	CumXL                    []float64       `json:"cum_xl"`
	CumSR                    []float64       `json:"cum_sr"`
	CumLNR                   []float64       `json:"cum_lnr"`
	PartisipasiSeries        []float64       `json:"partisipasi_series"`
	KepemilikanAsingTotal    []float64       `json:"kepemilikan_asing_total"`
	KepemilikanForeignHidden []float64       `json:"kepemilikan_foreign_hidden"`
	KepemilikanRitel         []float64       `json:"kepemilikan_ritel"`
	Correlation20            []*float64      `json:"correlation_20"`
	CorrelationKF20          []*float64      `json:"correlation_kf_20"`
	MA20Close                []*float64      `json:"ma20_close"`
	FPSeries                 []FPPoint       `json:"fp_series"`
	LatestFP                 FPPoint         `json:"latest_fp"`
	SpecialCandles           []SpecialCandle `json:"special_candles"`
	Highlights               []string        `json:"highlights"`
	LastPrice                float64         `json:"last_price"`
	PriceChange              float64         `json:"price_change"`
	Statistics               Statistics      `json:"statistics"`
}
