package stocks

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
	"github.com/go-chi/chi/v5"
)

// Handler serves /api/stocks endpoints.
type Handler struct {
	Repo *Repository
}

// GetHistory handles GET /api/stocks/{ticker}/history.
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

	limit := 500
	if raw := strings.TrimSpace(q.Get("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 5000 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "limit must be between 1 and 5000")
			return
		}
		limit = n
	}

	rows, err := h.Repo.GetHistory(r.Context(), ticker, startDate, endDate, limit)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.WriteJSON(w, http.StatusOK, PriceHistoryResponse{
		Ticker:    ticker,
		Count:     len(rows),
		Data:      rows,
		SortOrder: "asc",
	})
}
