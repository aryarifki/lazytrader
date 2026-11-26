package system

// JobConfig defines metadata for a runnable SMtracker pipeline job.
type JobConfig struct {
	Key              string
	Label            string
	Script           string
	LogPrefix        string
	Description      string
	Category         string
	IsShell          bool
	SupportsUniverse bool
	SupportsWindow   bool
	Params           []string
}

// JobRunRequest mirrors routers/system.py JobRunRequest.
type JobRunRequest struct {
	Universe   *string  `json:"universe"`
	Days       *int     `json:"days"`
	Months     *string  `json:"months"`
	Windows    []int    `json:"windows"`
	Workers    *int     `json:"workers"`
	Force      bool     `json:"force"`
	RateLimit  *float64 `json:"rate_limit"`
	Details    bool     `json:"details"`
	OnlyMissing bool    `json:"only_missing"`
	Ratios     bool     `json:"ratios"`
	Actions    bool     `json:"actions"`
	Year       *int     `json:"year"`
	Quarter    *int     `json:"quarter"`
	Autofix    bool     `json:"autofix"`
	Vacuum     bool     `json:"vacuum"`
	Notify     bool     `json:"notify"`
	AlertOnly  bool     `json:"alert_only"`
	AuditDate  *string  `json:"audit_date"`
	DryRun     bool     `json:"dry_run"`
	StatusOnly bool     `json:"status_only"`
}

// DefaultJobRunRequest initializes Pydantic-compatible defaults for JobRunRequest.
func DefaultJobRunRequest() JobRunRequest {
	d := 1
	q := 4
	return JobRunRequest{
		Days:    &d,
		Quarter: &q,
		Autofix: true,
		Vacuum:  true,
		Notify:  true,
	}
}

// JobStatusItem represents one entry returned by GET /api/logs/jobs/status.
type JobStatusItem struct {
	Key              string   `json:"key"`
	Label            string   `json:"label"`
	Description      string   `json:"description"`
	Category         string   `json:"category"`
	Status           string   `json:"status"`
	PID              *int     `json:"pid"`
	StartedAt        string   `json:"started_at,omitempty"`
	ElapsedSeconds   *int     `json:"elapsed_seconds,omitempty"`
	LogFile          string   `json:"log_file,omitempty"`
	IsExternal       bool     `json:"is_external,omitempty"`
	LastFinishedAt   *string  `json:"last_finished_at,omitempty"`
	LastExitCode     *int     `json:"last_exit_code,omitempty"`
	LastLogFile      *string  `json:"last_log_file,omitempty"`
	SupportsUniverse bool     `json:"supports_universe"`
	SupportsWindow   bool     `json:"supports_window"`
	Params           []string `json:"params"`
}

// LogHistoryFile represents one log file summary entry.
type LogHistoryFile struct {
	Filename   string  `json:"filename"`
	SizeBytes  int64   `json:"size_bytes"`
	SizeKB     float64 `json:"size_kb"`
	ModifiedAt string  `json:"modified_at"`
}
