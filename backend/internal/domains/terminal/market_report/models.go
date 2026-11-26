package market_report

import (
	"context"
	"net/http"

	"github.com/aryarifki/lazytrader/backend/internal/core/auth"
)

// LatestReport holds metadata and text contents for the newest report file.
type LatestReport struct {
	Session      int     `json:"session"`
	Filename     string  `json:"filename"`
	Date         string  `json:"date"`
	GeneratedAt  string  `json:"generated_at"`
	MaskedText   string  `json:"masked_text"`
	UnmaskedText *string `json:"unmasked_text"`
}

// EngineState stores live status of the Market Report runner.
type EngineState struct {
	Status       string        `json:"status"`
	IsRunning    bool          `json:"is_running"`
	Session      int           `json:"session"`
	Progress     int           `json:"progress"`
	Total        int           `json:"total"`
	Pct          float64       `json:"pct"`
	Message      string        `json:"message"`
	StartedAt    *string       `json:"started_at"`
	CompletedAt  *string       `json:"completed_at"`
	Error        *string       `json:"error"`
	LatestReport *LatestReport `json:"latest_report"`
}

// StatusData matches the "data" payload of GET /api/v1/terminal/market-report/status.
type StatusData struct {
	Status         string        `json:"status"`
	IsRunning      bool          `json:"is_running"`
	Session        int           `json:"session"`
	Progress       int           `json:"progress"`
	Total          int           `json:"total"`
	Pct            float64       `json:"pct"`
	Message        string        `json:"message"`
	StartedAt      *string       `json:"started_at"`
	CompletedAt    *string       `json:"completed_at"`
	Error          *string       `json:"error"`
	LatestReport   *LatestReport `json:"latest_report"`
	TokenValid     bool          `json:"token_valid"`
	CurrentTimeWIB string        `json:"current_time_wib"`
	RecentLogs     []string      `json:"recent_logs"`
}

// HistoryItem matches one entry in GET /api/v1/terminal/market-report/history.
type HistoryItem struct {
	Filename   string `json:"filename"`
	Session    int    `json:"session"`
	Date       string `json:"date"`
	ModifiedAt string `json:"modified_at"`
	SizeBytes  int64  `json:"size_bytes"`
}

// ReportDetail matches the "report" payload of GET /api/v1/terminal/market-report/detail.
type ReportDetail struct {
	Session      int     `json:"session"`
	Filename     string  `json:"filename"`
	MaskedText   string  `json:"masked_text"`
	UnmaskedText *string `json:"unmasked_text"`
}

// GenerateReportRequest mirrors Pydantic GenerateReportRequest.
type GenerateReportRequest struct {
	Session   *int     `json:"session"`
	RateLimit *float64 `json:"rate_limit"`
	Workers   *int     `json:"workers"`
}

// UpdateTokenRequest mirrors Pydantic UpdateTokenRequest.
type UpdateTokenRequest struct {
	Token *string `json:"token"`
}

// CommandRunner runs an external process and invokes onLine for each combined stdout/stderr line.
type CommandRunner func(ctx context.Context, bin string, args []string, extraEnv map[string]string, onLine func(string)) (int, error)

// Deps supplies external dependencies and paths to Market Report handlers/engine.
type Deps struct {
	Guard             *auth.Guard
	TerminalAdminKey  string
	ReportDir         string
	SMtrackerDir      string
	ScriptPath        string
	StockbitVerifyURL string
	HTTPClient        *http.Client
	Runner            CommandRunner
}
