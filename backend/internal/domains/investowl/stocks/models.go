package stocks

// PriceBar mirrors schemas.PriceBar (TradingView Lightweight Charts compatible).
type PriceBar struct {
	Date   string   `json:"date"`
	Ticker string   `json:"ticker"`
	Open   *float64 `json:"open"`
	High   *float64 `json:"high"`
	Low    *float64 `json:"low"`
	Close  *float64 `json:"close"`
	Volume *int64   `json:"volume"`
}

// PriceHistoryResponse mirrors schemas.PriceHistoryResponse.
type PriceHistoryResponse struct {
	Ticker    string     `json:"ticker"`
	Count     int        `json:"count"`
	Data      []PriceBar `json:"data"`
	SortOrder string     `json:"sort_order"`
}
