package broker

// BrokerFlowRow mirrors schemas.BrokerFlowRow.
type BrokerFlowRow struct {
	Date              string   `json:"date"`
	Ticker            string   `json:"ticker"`
	BandarSignal      *string  `json:"bandar_signal"`
	BandarSignalScore *float64 `json:"bandar_signal_score"`
	ForeignSignal     *string  `json:"foreign_signal"`
	ForeignNetBroker  *float64 `json:"foreign_net_broker"`
	LocalNetBroker    *float64 `json:"local_net_broker"`
	GovNetBroker      *float64 `json:"gov_net_broker"`
	ForeignNetFlow    *float64 `json:"foreign_net_flow"`
	DomesticNetFlow   *float64 `json:"domestic_net_flow"`
	TotalValue        *float64 `json:"total_value"`
	ConclusionBroker  *string  `json:"conclusion_broker"`
	ConclusionFlow    *string  `json:"conclusion_flow"`
	FetchedAt         *string  `json:"fetched_at"`
}

// BrokerFlowHistoryResponse mirrors schemas.BrokerFlowHistoryResponse.
type BrokerFlowHistoryResponse struct {
	Ticker    string          `json:"ticker"`
	Count     int             `json:"count"`
	Data      []BrokerFlowRow `json:"data"`
	SortOrder string          `json:"sort_order"`
}

// BrokerFlowSummary mirrors schemas.BrokerFlowSummary.
type BrokerFlowSummary struct {
	Ticker                  string   `json:"ticker"`
	PeriodDays              int      `json:"period_days"`
	TradingDays             int      `json:"trading_days"`
	ForeignNetBrokerSum     float64  `json:"foreign_net_broker_sum"`
	LocalNetBrokerSum       float64  `json:"local_net_broker_sum"`
	GovNetBrokerSum         float64  `json:"gov_net_broker_sum"`
	TotalValueSum           float64  `json:"total_value_sum"`
	ForeignDominancePct     float64  `json:"foreign_dominance_pct"`
	LatestBandarSignal      *string  `json:"latest_bandar_signal"`
	LatestBandarSignalScore *float64 `json:"latest_bandar_signal_score"`
	LatestForeignSignal     *string  `json:"latest_foreign_signal"`
	AccumulationDays        int      `json:"accumulation_days"`
	DistributionDays        int      `json:"distribution_days"`
	LatestDate              *string  `json:"latest_date"`
}

// BrokerWhaleListItem mirrors schemas.BrokerWhaleListItem.
type BrokerWhaleListItem struct {
	BrokerCode   string `json:"broker_code"`
	Profile      string `json:"profile"`
	ProfileLabel string `json:"profile_label"`
	ProfileDesc  string `json:"profile_desc"`
	IsSmart      bool   `json:"is_smart"`
}

// WhaleStockItem mirrors schemas.WhaleStockItem.
type WhaleStockItem struct {
	Ticker       string   `json:"ticker"`
	NetValue     float64  `json:"net_value"`
	BuyValue     float64  `json:"buy_value"`
	SellValue    float64  `json:"sell_value"`
	NetLot       float64  `json:"net_lot"`
	BuyAvgPrice  *float64 `json:"buy_avg_price"`
	SellAvgPrice *float64 `json:"sell_avg_price"`
}

// BrokerWhalesDetailResponse mirrors schemas.BrokerWhalesDetailResponse.
type BrokerWhalesDetailResponse struct {
	BrokerCode       string           `json:"broker_code"`
	Profile          string           `json:"profile"`
	ProfileLabel     string           `json:"profile_label"`
	ProfileDesc      string           `json:"profile_desc"`
	IsSmart          bool             `json:"is_smart"`
	PeriodDays       int              `json:"period_days"`
	StartDate        string           `json:"start_date"`
	EndDate          string           `json:"end_date"`
	TotalBuyValue    float64          `json:"total_buy_value"`
	TotalSellValue   float64          `json:"total_sell_value"`
	TotalNetValue    float64          `json:"total_net_value"`
	TopAccumulations []WhaleStockItem `json:"top_accumulations"`
	TopDistributions []WhaleStockItem `json:"top_distributions"`
}

type whaleAggRow struct {
	Ticker       string
	TotalBuy     float64
	TotalSell    float64
	TotalNet     float64
	NetLot       float64
	BuyAvgPrice  *float64
	SellAvgPrice *float64
}
