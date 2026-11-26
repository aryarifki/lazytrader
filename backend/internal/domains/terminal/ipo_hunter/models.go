package ipohunter

import "time"

// Underwriter describes an underwriter entry from ipo_hunter_master.underwriters JSONB.
type Underwriter struct {
	Broker string `json:"broker"`
	Role   string `json:"role"`
}

// MasterRecord holds one row from terminal_engine.ipo_hunter_master.
type MasterRecord struct {
	Ticker            string
	CompanyName       *string
	ListingDate       *time.Time
	IPOPrice          float64
	TotalListedShares *int64
	OfferingShares    *int64
	PctIPO            *string
	NilaiIPO          *float64
	Papan             *string
	LeadUnderwriter   *string
	Underwriters      []Underwriter
	Penjatahan        *float64
	JumlahSID         *int64
	DayReturns        []*float64
	Perf2W            *float64
	Perf1M            *float64
	Perf3M            *float64
	Perf6M            *float64
}

// DayRow holds one daily or D0 observation for an IPO ticker.
type DayRow struct {
	Date        time.Time
	Ticker      string
	Open        *float64
	High        *float64
	Low         *float64
	Close       *float64
	Prev        *float64
	Tavg        *float64
	Volume      *int64
	NetVolume   *int64
	BuyBrokers  map[string]float64
	SellBrokers map[string]float64
	IsD0        bool
	DayCount    int
}

// ActionLabel represents a classification label and CSS token.
type ActionLabel struct {
	Label string `json:"label"`
	Cls   string `json:"cls"`
}

// MetaResponse matches GET /api/v1/terminal/ipo-hunter/meta data payload.
type MetaResponse struct {
	TotalIPOs      int                    `json:"total_ipos"`
	D0Date         *string                `json:"d0_date"`
	EarliestDate   *string                `json:"earliest_date"`
	LatestDate     *string                `json:"latest_date"`
	AvailableDates []string               `json:"available_dates"`
	QuickPicks     []QuickPickItem        `json:"quick_picks"`
	Years          []int                  `json:"years"`
	Underwriters   []UnderwriterCountItem `json:"underwriters"`
	Tickers        []MetaTickerItem       `json:"tickers"`
}

type QuickPickItem struct {
	Ticker      string  `json:"ticker"`
	CompanyName string  `json:"company_name"`
	ListingDate string  `json:"listing_date"`
	IPOPrice    float64 `json:"ipo_price"`
}

type UnderwriterCountItem struct {
	Broker string `json:"broker"`
	Count  int    `json:"count"`
}

type MetaTickerItem struct {
	Ticker      string  `json:"ticker"`
	CompanyName string  `json:"company_name"`
	ListingDate *string `json:"listing_date"`
	IPOPrice    float64 `json:"ipo_price"`
}

// ListParams defines query parameters for GET /list.
type ListParams struct {
	UW        string
	LeadOnly  bool
	ValBucket string
	Year      *int
	Query     string
	SortBy    string
	SortDir   int
}

// ListResponse matches GET /api/v1/terminal/ipo-hunter/list data payload.
type ListResponse struct {
	Total          int            `json:"total"`
	TotalAll       int            `json:"total_all"`
	Rows           []ListRow      `json:"rows"`
	KinerjaSummary KinerjaSummary `json:"kinerja_summary"`
}

type ListRow struct {
	Ticker       string                `json:"ticker"`
	CompanyName  string                `json:"company_name"`
	ListingDate  *string               `json:"listing_date"`
	IPOPrice     float64               `json:"ipo_price"`
	NilaiIPO     *float64              `json:"nilai_ipo"`
	NilaiMiliar  *float64              `json:"nilai_miliar"`
	Papan        *string               `json:"papan"`
	Underwriters []ListUnderwriterChip `json:"underwriters"`
	JumlahSID    *int64                `json:"jumlah_sid"`
	Penjatahan   *float64              `json:"penjatahan"`
	DayReturns   []*float64            `json:"day_returns"`
	Perf2W       *float64              `json:"perf_2w"`
	Perf1M       *float64              `json:"perf_1m"`
	Perf3M       *float64              `json:"perf_3m"`
}

type ListUnderwriterChip struct {
	Broker   string `json:"broker"`
	Role     string `json:"role"`
	Category string `json:"category"`
}

type KinerjaSummary struct {
	RataRata        map[string]*int64 `json:"rata_rata"`
	KemungkinanNaik map[string]*int64 `json:"kemungkinan_naik"`
}

// AnalysisNoDataResponse matches GET /analysis/{ticker} when analysis_days is empty.
type AnalysisNoDataResponse struct {
	Ticker  string             `json:"ticker"`
	Mode    string             `json:"mode"`
	HasData bool               `json:"has_data"`
	Master  AnalysisNoDataMeta `json:"master"`
	Candles []CandleItem       `json:"candles"`
}

type AnalysisNoDataMeta struct {
	Ticker      string  `json:"ticker"`
	CompanyName string  `json:"company_name"`
	ListingDate *string `json:"listing_date"`
	IPOPrice    float64 `json:"ipo_price"`
}

// AnalysisResponse matches GET /analysis/{ticker} when data exists.
type AnalysisResponse struct {
	Ticker               string         `json:"ticker"`
	Mode                 string         `json:"mode"`
	HasData              bool           `json:"has_data"`
	Period               AnalysisPeriod `json:"period"`
	Candles              []CandleItem   `json:"candles"`
	AvailableTickerDates []string       `json:"available_ticker_dates"`
	Cards                AnalysisCards  `json:"cards"`
}

type AnalysisPeriod struct {
	Start         *string  `json:"start"`
	End           *string  `json:"end"`
	DayCount      int      `json:"day_count"`
	PriceChange   *float64 `json:"price_change"`
	AvgDailyValue *float64 `json:"avg_daily_value"`
	TotalVolume   float64  `json:"total_volume"`
	TotNetVolume  float64  `json:"tot_net_volume"`
	TotBuy        float64  `json:"tot_buy"`
	TotSell       float64  `json:"tot_sell"`
}

type CandleItem struct {
	Date      string   `json:"date"`
	Open      *float64 `json:"open"`
	High      *float64 `json:"high"`
	Low       *float64 `json:"low"`
	Close     *float64 `json:"close"`
	Prev      *float64 `json:"prev"`
	Tavg      *float64 `json:"tavg"`
	Volume    *int64   `json:"volume"`
	NetVolume *int64   `json:"net_volume"`
}

type AnalysisCards struct {
	TopBuyerSeller       TopBuyerSellerCard `json:"top_buyer_seller"`
	UnderwriterCard      UnderwriterCard    `json:"underwriter_card"`
	CategoryCard         CategoryCard       `json:"category_card"`
	AfterIPODistribution AfterIPODistCard   `json:"after_ipo_distribution"`
	BuyerSellerPie       BuyerSellerPieCard `json:"buyer_seller_pie"`
	RetailActionCard     RetailActionCard   `json:"retail_action_card"`
	InfoCard             InfoCard           `json:"info_card"`
}

type TopBuyerSellerCard struct {
	Buyers     []TBSRow  `json:"buyers"`
	Sellers    []TBSRow  `json:"sellers"`
	AccDist    []float64 `json:"acc_dist"`
	ValuePrice *float64  `json:"value_price"`
}

type TBSRow struct {
	Broker        string   `json:"broker"`
	Net           float64  `json:"net"`
	ValueJt       *float64 `json:"value_jt"`
	CategoryLabel string   `json:"category_label"`
	StyleKey      string   `json:"style_key"`
	IsUW          bool     `json:"is_uw"`
}

type UnderwriterCard struct {
	HasUW         bool                 `json:"has_uw"`
	TotalNet      float64              `json:"total_net"`
	TotalPctNet   *float64             `json:"total_pct_net"`
	TotalTrans    float64              `json:"total_trans"`
	TotalPctTrans *float64             `json:"total_pct_trans"`
	LotIPO        *float64             `json:"lot_ipo"`
	PctOfIPO      *float64             `json:"pct_of_ipo"`
	Status        string               `json:"status"`
	Rows          []UnderwriterCardRow `json:"rows"`
}

type UnderwriterCardRow struct {
	Broker   string   `json:"broker"`
	Role     string   `json:"role"`
	Net      float64  `json:"net"`
	PctNet   *float64 `json:"pct_net"`
	Trans    float64  `json:"trans"`
	PctTrans *float64 `json:"pct_trans"`
}

type CategoryCard struct {
	Rows  []CategoryCardRow `json:"rows"`
	Total CategoryCardTotal `json:"total"`
}

type CategoryCardRow struct {
	Key           string   `json:"key"`
	Label         string   `json:"label"`
	BandarBrokers []string `json:"bandar_brokers"`
	Net           float64  `json:"net"`
	PctNet        *float64 `json:"pct_net"`
	Trans         float64  `json:"trans"`
	PctTrans      *float64 `json:"pct_trans"`
}

type CategoryCardTotal struct {
	Net      float64  `json:"net"`
	PctNet   *float64 `json:"pct_net"`
	Trans    float64  `json:"trans"`
	PctTrans *float64 `json:"pct_trans"`
}

type AfterIPODistCard struct {
	Enabled bool     `json:"enabled"`
	Sell    DistSide `json:"sell"`
	Buy     DistSide `json:"buy"`
}

type DistSide struct {
	Available bool        `json:"available"`
	Empty     *bool       `json:"empty,omitempty"`
	IPOLot    *float64    `json:"ipo_lot,omitempty"`
	Slices    []DistSlice `json:"slices"`
	Total     float64     `json:"total"`
}

type DistSlice struct {
	Key        string  `json:"key"`
	Label      string  `json:"label"`
	Color      string  `json:"color"`
	Value      float64 `json:"value"`
	PctIPO     float64 `json:"pct_ipo"`
	SubBrokers string  `json:"sub_brokers"`
	PctRing    float64 `json:"pct_ring"`
}

type BuyerSellerPieCard struct {
	Net   PieModePair `json:"net"`
	Gross PieModePair `json:"gross"`
}

type PieModePair struct {
	Buyer  DonutSide `json:"buyer"`
	Seller DonutSide `json:"seller"`
}

type DonutSide struct {
	Total  float64      `json:"total"`
	Slices []DonutSlice `json:"slices"`
}

type DonutSlice struct {
	Key        string  `json:"key"`
	Label      string  `json:"label"`
	Color      string  `json:"color"`
	Value      float64 `json:"value"`
	Pct        float64 `json:"pct"`
	TopBrokers string  `json:"top_brokers"`
}

type RetailActionCard struct {
	RetailRatio            *float64            `json:"retail_ratio"`
	RetailAction           ActionLabel         `json:"retail_action"`
	RetailBandarAutoResult *string             `json:"retail_bandar_auto_result"`
	RetailBandarBroker     *string             `json:"retail_bandar_broker"`
	RetailBandarBrokers    []string            `json:"retail_bandar_brokers"`
	RetailBandarReasons    map[string][]string `json:"retail_bandar_reasons"`
	EffectiveRetailList    []string            `json:"effective_retail_list"`
	Rows                   []RetailActionRow   `json:"rows"`
	TotalExclAuto          RetailTotalExclAuto `json:"total_excl_auto"`
}

type RetailActionRow struct {
	Broker         string   `json:"broker"`
	Net            float64  `json:"net"`
	PctNet         *float64 `json:"pct_net"`
	Trans          float64  `json:"trans"`
	PctTrans       *float64 `json:"pct_trans"`
	IsAutoBandar   bool     `json:"is_auto_bandar"`
	IsActiveBandar bool     `json:"is_active_bandar"`
}

type RetailTotalExclAuto struct {
	SumTrans float64  `json:"sum_trans"`
	SumNet   float64  `json:"sum_net"`
	PctTrans *float64 `json:"pct_trans"`
	PctNet   *float64 `json:"pct_net"`
}

type InfoCard struct {
	Ticker              string    `json:"ticker"`
	CompanyName         string    `json:"company_name"`
	ListingDate         *string   `json:"listing_date"`
	IPOPrice            float64   `json:"ipo_price"`
	LastClose           *float64  `json:"last_close"`
	VWAP                *float64  `json:"vwap"`
	IPOVsClosePct       *int64    `json:"ipo_vs_close_pct"`
	AvgVsClosePct       *int64    `json:"avg_vs_close_pct"`
	LotOffer            *float64  `json:"lot_offer"`
	LotTotal            *float64  `json:"lot_total"`
	PctIPO              *string   `json:"pct_ipo"`
	NilaiMiliar         *float64  `json:"nilai_miliar"`
	Penjatahan          *float64  `json:"penjatahan"`
	Papan               *string   `json:"papan"`
	JumlahSID           *int64    `json:"jumlah_sid"`
	TradingDaysSinceIPO *int      `json:"trading_days_since_ipo"`
	HolidayWarn         bool      `json:"holiday_warn"`
	MaxHolidayYear      int       `json:"max_holiday_year"`
	Gap                 *GapRange `json:"gap"`
	NoD0                bool      `json:"no_d0"`
	InD0                bool      `json:"in_d0"`
	D0Date              *string   `json:"d0_date"`
	NLater              int       `json:"n_later"`
}

type GapRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// ScreenerCondition matches Pydantic ScreenerCondition.
type ScreenerCondition struct {
	K  string `json:"k"`
	Op string `json:"op"`
	V  string `json:"v"`
}

// ScreenerRequest matches Pydantic ScreenerRequest.
type ScreenerRequest struct {
	After    bool                `json:"after"`
	DateFrom *string             `json:"date_from"`
	DateTo   *string             `json:"date_to"`
	Quick    *int                `json:"quick"`
	Year     *int                `json:"year"`
	Chg      string              `json:"chg"`
	Metric   string              `json:"metric"`
	Conds    []ScreenerCondition `json:"conds"`
	SortKey  *string             `json:"sort_key"`
	Dir      *int                `json:"dir"`
}

// ScreenerRow represents one scored ticker in POST /screener.
type ScreenerRow struct {
	Kode        string   `json:"kode"`
	CompanyName string   `json:"company_name"`
	ListingDate *string  `json:"listing_date"`
	UW          *float64 `json:"uw"`
	NetRetail   *float64 `json:"netRetail"`
	NetSemi     *float64 `json:"netSemi"`
	NetInst     *float64 `json:"netInst"`
	NetNonRitel *float64 `json:"netNonRitel"`
	NetRepo     *float64 `json:"netRepo"`
	PerfIpo     *float64 `json:"perfIpo"`
	PerfAvg     *float64 `json:"perfAvg"`
	Chg         *float64 `json:"chg"`
	Ratrans     *float64 `json:"ratrans"`
	VolRatio    *float64 `json:"volRatio"`
}

// ScreenerResponse matches POST /api/v1/terminal/ipo-hunter/screener data payload.
type ScreenerResponse struct {
	Rows      []ScreenerRow `json:"rows"`
	Total     int           `json:"total"`
	NDays     int           `json:"n_days"`
	Lo        *string       `json:"lo"`
	Hi        *string       `json:"hi"`
	After     bool          `json:"after"`
	ActiveCol string        `json:"active_col"`
	SortKey   string        `json:"sort_key"`
	SortDir   int           `json:"sort_dir"`
}
