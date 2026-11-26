package signal

// DailySignalScores mirrors scores inside GET /api/signal/daily items.
type DailySignalScores struct {
	Technical   float64 `json:"technical"`
	SmartMoney  float64 `json:"smart_money"`
	Sector      float64 `json:"sector"`
	Fundamental float64 `json:"fundamental"`
}

// DailySignalItem mirrors one item in GET /api/signal/daily.
type DailySignalItem struct {
	Ticker          string            `json:"ticker"`
	Date            string            `json:"date"`
	MLLabel         *string           `json:"ml_label"`
	MLWinProb       *float64          `json:"ml_win_prob"`
	CompositeScore  float64           `json:"composite_score"`
	Scores          DailySignalScores `json:"scores"`
	Features        map[string]any    `json:"features"`
	WyckoffPhase    string            `json:"wyckoff_phase"`
	WyckoffName     string            `json:"wyckoff_name"`
	SmartMoneyNotes string            `json:"smart_money_notes"`
	SectorNotes     string            `json:"sector_notes"`
	Persistence5D   float64           `json:"persistence_5d"`
	RRRatio         *float64          `json:"rr_ratio"`
	GateNotes       string            `json:"gate_notes"`
}

// BacktestStatsResponse mirrors GET /api/signal/backtest/stats.
type BacktestStatsResponse struct {
	TotalSignals    int64   `json:"total_signals"`
	TotalClosed     int64   `json:"total_closed"`
	TotalOpen       int64   `json:"total_open"`
	WinCount        int64   `json:"win_count"`
	LossCount       int64   `json:"loss_count"`
	ExpiredCount    int64   `json:"expired_count"`
	WinRate         float64 `json:"win_rate"`
	DecisiveWinRate float64 `json:"decisive_win_rate"`
	AvgWinPct       float64 `json:"avg_win_pct"`
	AvgLossPct      float64 `json:"avg_loss_pct"`
	AvgDaysHeld     float64 `json:"avg_days_held"`
	ProfitFactor    float64 `json:"profit_factor"`
	LatestAuditDate *string `json:"latest_audit_date"`
}

// BacktestHistoryItem mirrors one item in GET /api/signal/backtest/history.
type BacktestHistoryItem struct {
	ID             string   `json:"id"`
	Ticker         string   `json:"ticker"`
	SignalDate     *string  `json:"signal_date"`
	EntryPrice     float64  `json:"entry_price"`
	TPPrice        float64  `json:"tp_price"`
	SLPrice        float64  `json:"sl_price"`
	ExitPrice      *float64 `json:"exit_price"`
	ExitDate       *string  `json:"exit_date"`
	Status         string   `json:"status"`
	DaysHeld       int      `json:"days_held"`
	PnLPct         *float64 `json:"pnl_pct"`
	CompositeScore any      `json:"composite_score"`
	SMScore        any      `json:"sm_score"`
	SecScore       any      `json:"sec_score"`
	WyckoffPhase   any      `json:"wyckoff_phase"`
	WyckoffName    any      `json:"wyckoff_name"`
	RRRatio        *float64 `json:"rr_ratio"`
}

// ActiveSignalItem mirrors one item in GET /api/signal/active.
type ActiveSignalItem struct {
	ID                string         `json:"id"`
	EpisodeID         *string        `json:"episode_id"`
	Ticker            string         `json:"ticker"`
	CompanyName       string         `json:"company_name"`
	GroupName         string         `json:"group_name"`
	SignalDate        *string        `json:"signal_date"`
	EntryPrice        float64        `json:"entry_price"`
	TPPrice           float64        `json:"tp_price"`
	SLPrice           float64        `json:"sl_price"`
	TargetGainPct     float64        `json:"target_gain_pct"`
	StopRiskPct       float64        `json:"stop_risk_pct"`
	RiskRewardRatio   float64        `json:"risk_reward_ratio"`
	Status            string         `json:"status"`
	DaysHeld          int            `json:"days_held"`
	ConfirmationCount int            `json:"confirmation_count"`
	LastConfirmedDate *string        `json:"last_confirmed_date"`
	CompositeScore    int            `json:"composite_score"`
	WyckoffPhase      any            `json:"wyckoff_phase"`
	WyckoffName       any            `json:"wyckoff_name"`
	SmartMoneyNotes   any            `json:"smart_money_notes"`
	SectorNotes       any            `json:"sector_notes"`
	Persistence5D     float64        `json:"persistence_5d"`
	TurnoverIDR       float64        `json:"turnover_idr"`
	MLWinProb         *float64       `json:"ml_win_prob"`
	CMF               any            `json:"cmf"`
	MFI               any            `json:"mfi"`
	RSI               any            `json:"rsi"`
	VCPGrade          any            `json:"vcp_grade"`
	GateNotes         any            `json:"gate_notes"`
	Features          map[string]any `json:"features"`
}

// ConfluenceItem mirrors one item in GET /api/signal/confluence-screener.
type ConfluenceItem struct {
	Ticker                 string   `json:"ticker"`
	CompanyName            string   `json:"company_name"`
	GroupName              string   `json:"group_name"`
	CurrentPrice           float64  `json:"current_price"`
	TotalValue             float64  `json:"total_value"`
	Signal                 string   `json:"signal"`
	ConvictionScore        float64  `json:"conviction_score"`
	BandarAvgPrice         float64  `json:"bandar_avg_price"`
	BandarPriceDistancePct float64  `json:"bandar_price_distance_pct"`
	Foreign5D              float64  `json:"foreign_5d"`
	SmartForeignNet        float64  `json:"smart_foreign_net"`
	GorenganNet            float64  `json:"gorengan_net"`
	ForeignHMMState        string   `json:"foreign_hmm_state"`
	FreeFloatEst           float64  `json:"free_float_est"`
	TopBuyer               string   `json:"top_buyer"`
	TopSeller              string   `json:"top_seller"`
	QuickTags              []string `json:"quick_tags"`
}
