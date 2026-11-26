package konglo

import "time"

// StockRow holds a row from terminal_engine.konglo_stock_summary for /report calculation.
type StockRow struct {
	Ticker       string
	Name         *string
	GroupName    *string
	Prev         *float64
	Open         *float64
	High         *float64
	Low          *float64
	Last         *float64
	Chg          *float64
	Vol          *int64
	Val          *float64
	Freq         *int64
	SharesListed *float64
	MarketCap    *float64
}

// SummaryStockRow matches GET /api/v1/terminal/konglo/summary item JSON.
type SummaryStockRow struct {
	Ticker       string   `json:"ticker"`
	Name         *string  `json:"name"`
	GroupName    *string  `json:"group_name"`
	Prev         *float64 `json:"prev"`
	Open         *float64 `json:"open"`
	High         *float64 `json:"high"`
	Low          *float64 `json:"low"`
	Last         *float64 `json:"last"`
	Chg          *float64 `json:"chg"`
	Vol          *int64   `json:"vol"`
	Val          *float64 `json:"val"`
	Freq         *int64   `json:"freq"`
	BBPrice      *float64 `json:"bb_price"`
	BBVol        *int64   `json:"bb_vol"`
	BOPrice      *float64 `json:"bo_price"`
	BOVol        *int64   `json:"bo_vol"`
	SharesListed *float64 `json:"shares_listed"`
	MarketCap    *float64 `json:"market_cap"`
	UpdatedAt    *string  `json:"updated_at"`
}

// GroupRow matches GET /api/v1/terminal/konglo/groups item JSON.
type GroupRow struct {
	GroupName    string   `json:"group_name"`
	Tickers      []string `json:"tickers"`
	TotalTickers *int     `json:"total_tickers"`
}

// GroupMemberRow matches member items in GET /api/v1/terminal/konglo/group/{group_name}.
type GroupMemberRow struct {
	Ticker    string   `json:"ticker"`
	Name      *string  `json:"name"`
	Prev      *float64 `json:"prev"`
	Open      *float64 `json:"open"`
	High      *float64 `json:"high"`
	Low       *float64 `json:"low"`
	Last      *float64 `json:"last"`
	Chg       *float64 `json:"chg"`
	Vol       *int64   `json:"vol"`
	Val       *float64 `json:"val"`
	Freq      *int64   `json:"freq"`
	MarketCap *float64 `json:"market_cap"`
}

// ReportStockItem represents a normalized stock entry inside Konglo report tabs.
type ReportStockItem struct {
	Ticker string  `json:"ticker"`
	Name   string  `json:"name"`
	Group  string  `json:"group"`
	Prev   float64 `json:"prev"`
	Close  float64 `json:"close"`
	Pct    float64 `json:"pct"`
	Val    float64 `json:"val"`
	Vol    float64 `json:"vol"`
	Freq   int64   `json:"freq"`
	Shares float64 `json:"shares"`
	MCChg  float64 `json:"mc_chg"`
}

// ReportSummary holds market breadth and optional IHSG summary fields.
type ReportSummary struct {
	StocksUp          int     `json:"stocks_up"`
	StocksDown        int     `json:"stocks_down"`
	StocksFlat        int     `json:"stocks_flat"`
	TotalActiveStocks int     `json:"total_active_stocks"`
	TotalStocks       int     `json:"total_stocks"`
	TotalTurnoverM    float64 `json:"total_turnover_m"`
	IHSGClose         string  `json:"ihsg_close,omitempty"`
	IHSGPct           string  `json:"ihsg_pct,omitempty"`
}

// OverviewCard holds overall konglo breadth metrics.
type OverviewCard struct {
	TotalKongloStocks int     `json:"total_konglo_stocks"`
	Rising            int     `json:"rising"`
	Falling           int     `json:"falling"`
	Flat              int     `json:"flat"`
	TotalMCChg        float64 `json:"total_mc_chg"`
	MissingCount      int     `json:"missing_count"`
}

// LeaderCard represents a best/worst group highlight card.
type LeaderCard struct {
	Group string  `json:"group"`
	Val   float64 `json:"val"`
	Count int     `json:"count"`
}

// KongloCards groups overview and best/worst leader cards.
type KongloCards struct {
	Overview OverviewCard `json:"overview"`
	BestMC   LeaderCard   `json:"best_mc"`
	BestPct  LeaderCard   `json:"best_pct"`
	WorstMC  LeaderCard   `json:"worst_mc"`
	WorstPct LeaderCard   `json:"worst_pct"`
}

// Top5Summary holds top-5 gainers and losers by % and MC change.
type Top5Summary struct {
	Top5PctUp   []ReportStockItem `json:"top5_pct_up"`
	Top5MCUp    []ReportStockItem `json:"top5_mc_up"`
	Top5PctDown []ReportStockItem `json:"top5_pct_down"`
	Top5MCDown  []ReportStockItem `json:"top5_mc_down"`
}

// Top15Tables holds top-15 ranking tables by % and MC change.
type Top15Tables struct {
	PctUp   []ReportStockItem `json:"pct_up"`
	MCUp    []ReportStockItem `json:"mc_up"`
	PctDown []ReportStockItem `json:"pct_down"`
	MCDown  []ReportStockItem `json:"mc_down"`
}

// Top2Entry holds a top-2 stock highlight within a conglomerate group.
type Top2Entry struct {
	Ticker string  `json:"ticker"`
	Pct    float64 `json:"pct"`
}

// GroupStat holds aggregated statistics and member list for a conglomerate group.
type GroupStat struct {
	Group        string            `json:"group"`
	TotalMembers int               `json:"total_members"`
	MCChg        float64           `json:"mc_chg"`
	TotalMC      float64           `json:"total_mc"`
	MCPrev       float64           `json:"mc_prev"`
	MCPct        float64           `json:"mc_pct"`
	TotalVal     float64           `json:"total_val"`
	TotalFreq    int64             `json:"total_freq"`
	AvgPct       float64           `json:"avg_pct"`
	Top2         []Top2Entry       `json:"top2"`
	Members      []ReportStockItem `json:"members"`
}

// FundamentalSummary holds benchmark stats across the 50 fundamental stocks.
type FundamentalSummary struct {
	TotalStocks  int               `json:"total_stocks"`
	MissingCount int               `json:"missing_count"`
	Rising       int               `json:"rising"`
	Falling      int               `json:"falling"`
	Flat         int               `json:"flat"`
	TotalMCChg   float64           `json:"total_mc_chg"`
	MCPct        float64           `json:"mc_pct"`
	Top5PctUp    []ReportStockItem `json:"top5_pct_up"`
	Top5MCUp     []ReportStockItem `json:"top5_mc_up"`
	Top5PctDown  []ReportStockItem `json:"top5_pct_down"`
	Top5MCDown   []ReportStockItem `json:"top5_mc_down"`
}

// MissingKonglo records an inactive/missing ticker within a conglomerate group.
type MissingKonglo struct {
	Ticker string `json:"ticker"`
	Group  string `json:"group"`
}

// ReportMetadata holds formatted Indonesian EOD timestamps.
type ReportMetadata struct {
	Timestamp      string `json:"timestamp"`
	TradeDate      string `json:"trade_date"`
	DateDisplay    string `json:"date_display"`
	TimeDisplay    string `json:"time_display"`
	LastUpdatedStr string `json:"last_updated_str"`
	SourceStatus   string `json:"source_status"`
	SourceTitle    string `json:"source_title"`
	SourceSub      string `json:"source_sub"`
}

// Report matches the full payload under "report" in GET /api/v1/terminal/konglo/report.
type Report struct {
	Summary            ReportSummary      `json:"summary"`
	KongloCards        KongloCards        `json:"konglo_cards"`
	Top5Summary        Top5Summary        `json:"top5_summary"`
	Top15Tables        Top15Tables        `json:"top15_tables"`
	Top5Konglo         []GroupStat        `json:"top5_konglo"`
	AllGroupStats      []GroupStat        `json:"all_group_stats"`
	FundamentalSummary FundamentalSummary `json:"fundamental_summary"`
	FundamentalReport  []ReportStockItem  `json:"fundamental_report"`
	MissingKonglo      []MissingKonglo    `json:"missing_konglo"`
	MissingFunda       []string           `json:"missing_funda"`
	Metadata           *ReportMetadata    `json:"metadata,omitempty"`
}

// ReportTimestamps holds latest trade date, stock summary updated_at, and optional IHSG summary map.
type ReportTimestamps struct {
	MaxTradeDate *time.Time
	MaxUpdated   *time.Time
	IHSGSummary  map[string]any
}
