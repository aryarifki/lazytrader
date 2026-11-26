package mmd

// TickerMetaSummary matches GET /api/v1/terminal/mmd/tickers items.
type TickerMetaSummary struct {
	Ticker        string   `json:"ticker"`
	EarliestDate  *string  `json:"earliest_date"`
	LatestDate    *string  `json:"latest_date"`
	TotalDays     *int     `json:"total_days"`
	LastClose     *float64 `json:"last_close"`
	LastVolume    *int64   `json:"last_volume"`
	LastTavg      *float64 `json:"last_tavg"`
	LastNetVolume *int64   `json:"last_net_volume"`
}

// TickerMetaFull represents all columns of terminal_engine.mmd_tickers_meta.
type TickerMetaFull struct {
	Ticker        string   `json:"ticker"`
	CompanyName   *string  `json:"company_name"`
	LastDate      *string  `json:"last_date"`
	LastClose     *float64 `json:"last_close"`
	LastNetVolume *int64   `json:"last_net_volume"`
	IsActive      *bool    `json:"is_active"`
	UpdatedAt     *string  `json:"updated_at"`
	EarliestDate  *string  `json:"earliest_date"`
	LatestDate    *string  `json:"latest_date"`
	TotalDays     *int     `json:"total_days"`
	LastVolume    *int64   `json:"last_volume"`
	LastTavg      *float64 `json:"last_tavg"`
	IPODate       *string  `json:"-"`
}

// DailyCandle matches one row of terminal_engine.mmd_daily_candles.
type DailyCandle struct {
	DateISO     string             `json:"date_iso"`
	DateStr     *string            `json:"date_str"`
	FileName    *string            `json:"file_name"`
	Open        *float64           `json:"open"`
	High        *float64           `json:"high"`
	Low         *float64           `json:"low"`
	Close       *float64           `json:"close"`
	Volume      *int64             `json:"volume"`
	Prev        *float64           `json:"prev"`
	Tavg        *float64           `json:"tavg"`
	NetVolume   *int64             `json:"net_volume"`
	BuyBrokers  map[string]float64 `json:"buy_brokers"`
	SellBrokers map[string]float64 `json:"sell_brokers"`
}

// BandarSeriesPoint represents an item in bandar_accumulation_series.
type BandarSeriesPoint struct {
	Date             string   `json:"date"`
	NetVolume        int64    `json:"net_volume"`
	CumulativeBandar int64    `json:"cumulative_bandar"`
	Close            *float64 `json:"close"`
}

// Badge represents a classification badge with CSS class and hex color.
type Badge struct {
	Label string `json:"label"`
	Cls   string `json:"cls"`
	Color string `json:"color"`
}

// BandarTableRow represents one row in the Top 3 / Top 5 / Average T10 / Aggregate Score table.
type BandarTableRow struct {
	Label    string  `json:"label"`
	Buy      int64   `json:"buy"`
	Sell     int64   `json:"sell"`
	Nilai    int64   `json:"nilai"`
	Ratio    float64 `json:"ratio"`
	Badge    Badge   `json:"badge"`
	IsAction bool    `json:"is_action,omitempty"`
}

// TopBrokerRow represents one row in buyers or sellers lists.
type TopBrokerRow struct {
	Broker    string  `json:"broker"`
	Net       int64   `json:"net"`
	VolumeFmt string  `json:"volume_fmt"`
	ValueJt   float64 `json:"value_jt"`
	Kategori  string  `json:"kategori"`
}

// Top10CombinedRow represents one rank row in top10_combined.
type Top10CombinedRow struct {
	Rank           int     `json:"rank"`
	BuyerBroker    string  `json:"buyer_broker"`
	BuyerNet       int64   `json:"buyer_net"`
	BuyerNetFmt    string  `json:"buyer_net_fmt"`
	BuyerValJt     float64 `json:"buyer_val_jt"`
	BuyerKategori  string  `json:"buyer_kategori"`
	SellerBroker   string  `json:"seller_broker"`
	SellerNet      int64   `json:"seller_net"`
	SellerNetFmt   string  `json:"seller_net_fmt"`
	SellerValJt    float64 `json:"seller_val_jt"`
	SellerKategori string  `json:"seller_kategori"`
	AccDist        int64   `json:"acc_dist"`
	AccDistFmt     string  `json:"acc_dist_fmt"`
}

// RetailRow represents a core retail broker row in retail_rows.
type RetailRow struct {
	Broker         string  `json:"broker"`
	Trans          int64   `json:"trans"`
	PctTrans       float64 `json:"pct_trans"`
	Net            int64   `json:"net"`
	PctNet         float64 `json:"pct_net"`
	IsRetailBandar bool    `json:"is_retail_bandar"`
}

// PengendaliGroupStat holds main participation % and net directional % for a broker category.
type PengendaliGroupStat struct {
	Main float64 `json:"main"`
	Sub  float64 `json:"sub"`
}

// ValStat holds a single scalar percentage metric.
type ValStat struct {
	Val float64 `json:"val"`
}

// PengendaliStats holds the 5-box Broker Pengendali summary.
type PengendaliStats struct {
	Institusi PengendaliGroupStat `json:"institusi"`
	SemiRitel PengendaliGroupStat `json:"semi_ritel"`
	Ritel     PengendaliGroupStat `json:"ritel"`
	NonRitel  PengendaliGroupStat `json:"non_ritel"`
	Repo      ValStat             `json:"repo"`
	Lainnya   ValStat             `json:"lainnya"`
	Ratio     float64             `json:"ratio"`
	Action    Badge               `json:"action"`
}

// PreviousResultRow represents a WTD/1W/2W or MTD/1M/2M matrix row.
type PreviousResultRow struct {
	Label            string  `json:"label"`
	Date             string  `json:"date"`
	Badge            Badge   `json:"badge"`
	ChangePct        float64 `json:"change_pct"`
	PengendaliAction Badge   `json:"pengendali_action"`
	RetailAction     Badge   `json:"retail_action"`
	Range            string  `json:"range"`
}

// PreviousResults holds weekly and monthly historical resume matrices.
type PreviousResults struct {
	Weekly  []PreviousResultRow `json:"weekly"`
	Monthly []PreviousResultRow `json:"monthly"`
}

// AIScoreComponents exposes the detailed C1-C10 + ancillary components of AI Score.
type AIScoreComponents struct {
	C1          int     `json:"c1"`
	C2          float64 `json:"c2"`
	C3          int     `json:"c3"`
	C4          float64 `json:"c4"`
	C5          int     `json:"c5"`
	C6          float64 `json:"c6"`
	C7          float64 `json:"c7"`
	C8          int     `json:"c8"`
	C9          int     `json:"c9"`
	C10         int     `json:"c10"`
	CCorr       int64   `json:"c_corr"`
	CMark       int     `json:"c_mark"`
	CDaily30    float64 `json:"c_daily30"`
	CSideways   int     `json:"c_sideways"`
	TotalRaw    float64 `json:"total_raw"`
	TotalScaled float64 `json:"total_scaled"`
}

// T5Record matches one KSEI T+5 substantial shareholder transaction.
type T5Record struct {
	Date     string  `json:"date"`
	Investor *string `json:"investor"`
	Broker   *string `json:"broker"`
	Action   *string `json:"action"`
	Lot      float64 `json:"lot"`
	ValueM   float64 `json:"value_m"`
	Price    float64 `json:"price"`
	Pct      float64 `json:"pct"`
}

// GroupPerfMember represents one peer stock inside conglomerate or sector performance.
type GroupPerfMember struct {
	Ticker    string   `json:"ticker"`
	LastClose float64  `json:"last_close"`
	TodayPct  float64  `json:"today_pct"`
	W1Pct     *float64 `json:"w1_pct"`
	M1Pct     float64  `json:"m1_pct"`
	AvgVal    float64  `json:"avg_val"`
}

// GroupPerfBlock represents either the Conglomerate or Sector performance block.
type GroupPerfBlock struct {
	GroupName string            `json:"group_name"`
	Title     string            `json:"title"`
	Type      string            `json:"type"`
	Members   []GroupPerfMember `json:"members"`
}

// GroupPerformance groups conglomerate and sector blocks for a ticker.
type GroupPerformance struct {
	Conglomerate *GroupPerfBlock `json:"conglomerate"`
	Sector       *GroupPerfBlock `json:"sector"`
}

// ExplorerBigPlayerItem matches category=bigPlayer rows from /explorer/scan.
type ExplorerBigPlayerItem struct {
	Ticker      string  `json:"ticker"`
	LastClose   float64 `json:"last_close"`
	LastVolume  int64   `json:"last_volume"`
	PriceChange float64 `json:"price_change"`
	Badge       Badge   `json:"badge"`
	Detail      string  `json:"detail"`
}

// ExplorerScanItem matches category=normal/advanced rows from /explorer/scan.
type ExplorerScanItem struct {
	Ticker      string  `json:"ticker"`
	CompanyName *string `json:"company_name"`
	LastClose   float64 `json:"last_close"`
	LastVolume  int64   `json:"last_volume"`
	PriceChange float64 `json:"price_change"`
	Ratio       float64 `json:"ratio"`
	Badge       Badge   `json:"badge"`
	Detail      string  `json:"detail"`
}

// BrokerDrilldownDay represents one day of broker activity.
type BrokerDrilldownDay struct {
	Date    string  `json:"date"`
	BuyLot  float64 `json:"buy_lot"`
	BuyVal  float64 `json:"buy_val"`
	BuyAvg  float64 `json:"buy_avg"`
	SellLot float64 `json:"sell_lot"`
	SellVal float64 `json:"sell_val"`
	SellAvg float64 `json:"sell_avg"`
	NetLot  float64 `json:"net_lot"`
	NetVal  float64 `json:"net_val"`
}

// BrokerDrilldownSummary summarizes broker buy/sell totals across the date window.
type BrokerDrilldownSummary struct {
	TotalBuyLot     float64 `json:"total_buy_lot"`
	TotalBuyValueM  float64 `json:"total_buy_value_m"`
	TotalBuyValueB  float64 `json:"total_buy_value_b"`
	AvgBuyPrice     float64 `json:"avg_buy_price"`
	TotalSellLot    float64 `json:"total_sell_lot"`
	TotalSellValueM float64 `json:"total_sell_value_m"`
	TotalSellValueB float64 `json:"total_sell_value_b"`
	AvgSellPrice    float64 `json:"avg_sell_price"`
	NetLot          float64 `json:"net_lot"`
	NetValueM       float64 `json:"net_value_m"`
	NetValueB       float64 `json:"net_value_b"`
	DaysActive      int     `json:"days_active"`
}

// AnalysisResult represents the full MM Detector Pro analysis payload.
type AnalysisResult struct {
	Ticker                     string             `json:"ticker"`
	TotalDays                  int                `json:"total_days"`
	DateFrom                   string             `json:"date_from"`
	DateTo                     string             `json:"date_to"`
	LastPrice                  float64            `json:"last_price"`
	PriceChange                float64            `json:"price_change"`
	TotalBuy                   int64              `json:"total_buy"`
	TotalSell                  int64              `json:"total_sell"`
	TotalVolume                int64              `json:"total_volume"`
	TotNetVolume               int64              `json:"tot_net_volume"`
	AvgDailyValue              float64            `json:"avg_daily_value"`
	Buyers                     []TopBrokerRow     `json:"buyers"`
	Sellers                    []TopBrokerRow     `json:"sellers"`
	AccDistRows                []int64            `json:"acc_dist_rows"`
	Top10Combined              []Top10CombinedRow `json:"top10_combined"`
	BandarTable                []BandarTableRow   `json:"bandar_table"`
	BandarAction               Badge              `json:"bandar_action"`
	RetailRows                 []RetailRow        `json:"retail_rows"`
	RetailRatio                float64            `json:"retail_ratio"`
	RetailAction               Badge              `json:"retail_action"`
	RetailBandarAutoResult     *string            `json:"retail_bandar_auto_result"`
	ActiveRetailBandar         *string            `json:"active_retail_bandar"`
	RetailBandarAutoEnabled    bool               `json:"retail_bandar_auto_enabled"`
	RetailBandarManualOverride string             `json:"retail_bandar_manual_override"`
	EffectiveRetailBrokers     []string           `json:"effective_retail_brokers"`
	PengendaliStats            PengendaliStats    `json:"pengendali_stats"`
	Resume                     Badge              `json:"resume"`
	AIScore                    int                `json:"ai_score"`
	Korelasi                   float64            `json:"korelasi"`
	AIScoreComponents          AIScoreComponents  `json:"ai_score_components"`
	PreviousResults            PreviousResults    `json:"previous_results"`
	AICommentary               []string           `json:"ai_commentary"`
	Commentary                 []string           `json:"commentary"`
	BigPlayerSummary           []string           `json:"big_player_summary"`
	T5Records                  []T5Record         `json:"t5_records"`
	GroupPerformance           *GroupPerformance  `json:"group_performance"`
}
