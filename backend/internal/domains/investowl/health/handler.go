package health

import (
	"net/http"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler serves InvestOwl and Terminal healthcheck endpoints.
type Handler struct {
	DB *pgxpool.Pool
}

// InvestowlHealth mirrors GET /api/health -> {"status":"healthy","latest_date":...}.
func (h *Handler) InvestowlHealth(w http.ResponseWriter, r *http.Request) {
	var latestDate *time.Time
	err := h.DB.QueryRow(r.Context(), `SELECT MAX(date) FROM public.broker_flow`).Scan(&latestDate)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "Database health check failed")
		return
	}
	var formatted *string
	if latestDate != nil {
		s := latestDate.Format("2006-01-02")
		formatted = &s
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":      "healthy",
		"latest_date": formatted,
	})
}

// TerminalHealth mirrors GET /api/v1/terminal/health.
func (h *Handler) TerminalHealth(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":    "healthy",
		"service":   "lazytrader-fusion",
		"version":   "1.0.0",
		"timestamp": float64(time.Now().UnixNano()) / 1e9,
	})
}
