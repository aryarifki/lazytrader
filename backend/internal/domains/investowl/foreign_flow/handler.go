package foreign_flow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler serves /api/foreign-flow/* endpoints.
type Handler struct {
	DB      *pgxpool.Pool
	Cache   *cache.Store
	Sectors *SectorService
}

// GetSectorsSummary handles GET /api/foreign-flow/sectors/summary.
func (h *Handler) GetSectorsSummary(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	res, status, err := h.Sectors.GetSummary(r.Context())
	if err != nil {
		httpx.WriteError(w, status, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}

// GetSectorsHistory handles GET /api/foreign-flow/sectors/history.
func (h *Handler) GetSectorsHistory(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	q := r.URL.Query()
	sector := strings.TrimSpace(q.Get("sector"))
	if sector == "" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "sector query parameter is required")
		return
	}
	days := 30
	if raw := strings.TrimSpace(q.Get("days")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 5 || n > 100 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "days must be between 5 and 100")
			return
		}
		days = n
	}

	res, status, err := h.Sectors.GetHistory(r.Context(), sector, days)
	if err != nil {
		httpx.WriteError(w, status, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}

// GetTickerForensic handles GET /api/foreign-flow/{ticker}.
func (h *Handler) GetTickerForensic(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	ticker := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	lookbackDays := 20
	if raw := strings.TrimSpace(r.URL.Query().Get("lookback_days")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 20 || n > 250 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "lookback_days must be between 20 and 250")
			return
		}
		lookbackDays = n
	}

	var latestDate *time.Time
	err := h.DB.QueryRow(r.Context(), `SELECT MAX(date) FROM public.broker_flow WHERE ticker = $1`, ticker).Scan(&latestDate)
	if err != nil && err != pgx.ErrNoRows {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if latestDate == nil {
		httpx.WriteError(w, http.StatusNotFound, fmt.Sprintf("Emiten '%s' belum memiliki data transaksi asing di database.", ticker))
		return
	}

	cacheKey := fmt.Sprintf("ff_forensic_v5:%s:%d:%s", ticker, lookbackDays, latestDate.Format("2006-01-02"))
	if h.Cache != nil {
		if raw, ok := h.Cache.Get(r.Context(), cacheKey); ok {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(raw)
			return
		}
	}

	payload, err := ComputeForeignForensic(r.Context(), h.DB, ticker, lookbackDays)
	if err != nil {
		httpx.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if h.Cache != nil {
		if raw, err := json.Marshal(payload); err == nil {
			h.Cache.Set(r.Context(), cacheKey, raw, 21600*time.Second)
		}
	}
	httpx.WriteJSON(w, http.StatusOK, payload)
}
