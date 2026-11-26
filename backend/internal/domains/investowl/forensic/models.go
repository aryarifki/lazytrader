package forensic

// TimelinePoint mirrors one element of Foreign Exodus Radar timeline.
type TimelinePoint struct {
	Date           string  `json:"date"`
	IHSGClose      float64 `json:"ihsg_close"`
	DailyFlow      float64 `json:"daily_flow"`
	CumulativeFlow float64 `json:"cumulative_flow"`
	MarketState    string  `json:"market_state,omitempty"`
}

// CulpritBroker mirrors top net-selling foreign broker entry.
type CulpritBroker struct {
	BrokerCode string  `json:"broker_code"`
	BrokerName string  `json:"broker_name"`
	NetSellIDR float64 `json:"net_sell_idr"`
}

// ExodusRadarResponse mirrors GET /api/forensic/foreign-exodus response.
type ExodusRadarResponse struct {
	MarketState    string          `json:"market_state"`
	IHSGCurrent    float64         `json:"ihsg_current"`
	IHSGChangePct  float64         `json:"ihsg_change_pct"`
	Timeline       []TimelinePoint `json:"timeline"`
	CulpritBrokers []CulpritBroker `json:"culprit_brokers"`
}
