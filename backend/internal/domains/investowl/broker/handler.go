package broker

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
	"github.com/go-chi/chi/v5"
)

// Handler serves /api/broker-flow endpoints.
type Handler struct {
	Repo *Repository
	Svc  *Service
}

// GetLatest handles GET /api/broker-flow/{ticker}/latest.
func (h *Handler) GetLatest(w http.ResponseWriter, r *http.Request) {
	ticker := shared.NormalizeTicker(chi.URLParam(r, "ticker"))
	row, err := h.Repo.GetLatest(r.Context(), ticker)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if row == nil {
		httpx.WriteError(w, http.StatusNotFound, fmt.Sprintf("Tidak ada data broker_flow untuk ticker '%s'", ticker))
		return
	}
	httpx.WriteJSON(w, http.StatusOK, row)
}

// GetHistory handles GET /api/broker-flow/{ticker}/history.
func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	ticker := shared.NormalizeTicker(chi.URLParam(r, "ticker"))
	q := r.URL.Query()

	var startDate, endDate *time.Time
	if raw := strings.TrimSpace(q.Get("start_date")); raw != "" {
		t, err := time.Parse("2006-01-02", raw)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "Invalid start_date format (expected YYYY-MM-DD)")
			return
		}
		startDate = &t
	}
	if raw := strings.TrimSpace(q.Get("end_date")); raw != "" {
		t, err := time.Parse("2006-01-02", raw)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "Invalid end_date format (expected YYYY-MM-DD)")
			return
		}
		endDate = &t
	}

	limit := 250
	if raw := strings.TrimSpace(q.Get("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 2500 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "limit must be between 1 and 2500")
			return
		}
		limit = n
	}

	rows, err := h.Repo.GetHistory(r.Context(), ticker, startDate, endDate, limit)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, BrokerFlowHistoryResponse{
		Ticker:    ticker,
		Count:     len(rows),
		Data:      rows,
		SortOrder: "asc",
	})
}

// GetSummary handles GET /api/broker-flow/{ticker}/summary.
func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	ticker := shared.NormalizeTicker(chi.URLParam(r, "ticker"))
	days := 30
	if raw := strings.TrimSpace(r.URL.Query().Get("days")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 365 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "days must be between 1 and 365")
			return
		}
		days = n
	}

	summary, err := h.Svc.BuildSummary(r.Context(), ticker, days)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if summary == nil {
		httpx.WriteError(w, http.StatusNotFound, fmt.Sprintf("Tidak ada data broker_flow untuk '%s' dalam %d hari terakhir", ticker, days))
		return
	}
	httpx.WriteJSON(w, http.StatusOK, summary)
}

// GetWhalesList handles GET /api/broker-flow/whales.
func (h *Handler) GetWhalesList(w http.ResponseWriter, r *http.Request) {
	items, err := h.Svc.GetWhalesList(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, items)
}

// GetWhaleDetail handles GET /api/broker-flow/whales/{broker_code}.
func (h *Handler) GetWhaleDetail(w http.ResponseWriter, r *http.Request) {
	code := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "broker_code")))
	q := r.URL.Query()

	days := 20
	if raw := strings.TrimSpace(q.Get("days")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 100 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "days must be between 1 and 100")
			return
		}
		days = n
	}

	limit := 20
	if raw := strings.TrimSpace(q.Get("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 100 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "limit must be between 1 and 100")
			return
		}
		limit = n
	}

	detail, err := h.Svc.GetWhaleDetail(r.Context(), code, days, limit)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if detail == nil {
		httpx.WriteError(w, http.StatusNotFound, fmt.Sprintf("Tidak ada data aktivitas transaksi untuk broker '%s'", code))
		return
	}
	httpx.WriteJSON(w, http.StatusOK, detail)
}
