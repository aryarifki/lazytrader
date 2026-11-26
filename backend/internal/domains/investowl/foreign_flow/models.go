package foreign_flow

// SectorSummaryItem mirrors one sector entry in GET /api/foreign-flow/sectors/summary.
type SectorSummaryItem struct {
	Code              string   `json:"code"`
	Name              string   `json:"name"`
	DBSector          string   `json:"db_sector"`
	Icon              string   `json:"icon"`
	Color             string   `json:"color"`
	Close             float64  `json:"close"`
	Ret1D             float64  `json:"ret_1d"`
	Ret5D             float64  `json:"ret_5d"`
	Ret20D            float64  `json:"ret_20d"`
	Value             float64  `json:"value"`
	ValueSharePct     float64  `json:"value_share_pct"`
	ForeignNet1D      float64  `json:"foreign_net_1d"`
	ForeignNet5D      float64  `json:"foreign_net_5d"`
	ForeignNet20D     float64  `json:"foreign_net_20d"`
	ForeignState      string   `json:"foreign_state"`
	ForeignStateLabel string   `json:"foreign_state_label"`
	RSRatio           float64  `json:"rs_ratio"`
	RSMomentum        float64  `json:"rs_momentum"`
	Quadrant          string   `json:"quadrant"`
	TopStocks         []string `json:"top_stocks"`
}

// SectorBenchmark mirrors the benchmark object in /sectors/summary.
type SectorBenchmark struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Close       float64 `json:"close"`
	Ret1D       float64 `json:"ret_1d"`
	Value       float64 `json:"value"`
	MarketState string  `json:"market_state"`
}

// SectorRotationSummaryResponse mirrors GET /api/foreign-flow/sectors/summary.
type SectorRotationSummaryResponse struct {
	Date              string              `json:"date"`
	Benchmark         SectorBenchmark     `json:"benchmark"`
	MacroMarketState  string              `json:"macro_market_state"`
	TotalSectorValue  float64             `json:"total_sector_value"`
	TotalForeignNet1D float64             `json:"total_foreign_net_1d"`
	Sectors           []SectorSummaryItem `json:"sectors"`
}

// SectorHistoryPoint mirrors one time-series bar in /sectors/history.
type SectorHistoryPoint struct {
	Date                  string  `json:"date"`
	Close                 float64 `json:"close"`
	Value                 float64 `json:"value"`
	ForeignNet            float64 `json:"foreign_net"`
	CumulativeForeignFlow float64 `json:"cumulative_foreign_flow"`
	RSRatio               float64 `json:"rs_ratio"`
}

// SectorBrokerItem mirrors savior_brokers / culprit_brokers entries.
type SectorBrokerItem struct {
	BrokerCode string  `json:"broker_code"`
	BrokerName string  `json:"broker_name"`
	NetValue   float64 `json:"net_value"`
}

// SectorHistoryResponse mirrors GET /api/foreign-flow/sectors/history.
type SectorHistoryResponse struct {
	Sector         string               `json:"sector"`
	Name           string               `json:"name"`
	Count          int                  `json:"count"`
	Data           []SectorHistoryPoint `json:"data"`
	SaviorBrokers  []SectorBrokerItem   `json:"savior_brokers"`
	CulpritBrokers []SectorBrokerItem   `json:"culprit_brokers"`
}

// MicroRow represents one foreign broker activity row used in heatmap generation.
type MicroRow struct {
	Date       string
	BrokerCode string
	NetValue   float64
}

// BrokerHeatmap mirrors compute_broker_heatmap output.
type BrokerHeatmap struct {
	XDates     []string    `json:"x_dates"`
	YBrokers   []string    `json:"y_brokers"`
	MatrixData [][]any     `json:"matrix_data"`
}

// ForensicBrokerEntry mirrors top_holders / top_sellers entries.
type ForensicBrokerEntry struct {
	BrokerCode     string   `json:"broker_code"`
	BrokerName     string   `json:"broker_name"`
	NetValue       float64  `json:"net_value"`
	NetLot         float64  `json:"net_lot"`
	BuyValue       float64  `json:"buy_value"`
	SellValue      float64  `json:"sell_value"`
	AvgBuyPrice    float64  `json:"avg_buy_price"`
	AvgSellPrice   float64  `json:"avg_sell_price"`
	FloatingPnLPct *float64 `json:"floating_pnl_pct"`
	ActiveDays     int      `json:"active_days"`
}
