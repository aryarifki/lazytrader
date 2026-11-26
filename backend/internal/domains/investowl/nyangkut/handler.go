package nyangkut

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/go-chi/chi/v5"
)

// Handler serves Dashboard Nyangkut endpoints and aliases.
type Handler struct {
	Svc   *Service
	Cache *cache.Store
}

func nyangkutCacheKey(prefix, raw string) string {
	h := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%s:%s", prefix, hex.EncodeToString(h[:])[:12])
}

func (h *Handler) tryCache(w http.ResponseWriter, r *http.Request, key string, clientMaxAge, sMaxAge int) bool {
	if h.Cache == nil {
		return false
	}
	if raw, ok := h.Cache.Get(r.Context(), key); ok {
		httpx.WriteCachedJSON(w, raw, clientMaxAge, sMaxAge, true)
		return true
	}
	return false
}

func (h *Handler) storeCache(w http.ResponseWriter, r *http.Request, key string, v any, clientMaxAge, sMaxAge int) {
	payload, err := json.Marshal(v)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if h.Cache != nil && key != "" {
		h.Cache.Set(r.Context(), key, payload, time.Duration(sMaxAge)*time.Second)
	}
	httpx.WriteCachedJSON(w, payload, clientMaxAge, sMaxAge, false)
}

// GetStockTracker handles GET /api/v1/stock/{ticker}/nyangkut-tracker and /api/nyangkut/stock/{ticker}.
func (h *Handler) GetStockTracker(w http.ResponseWriter, r *http.Request) {
	tickerClean := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	if tickerClean == "" {
		httpx.WriteError(w, http.StatusBadRequest, "Ticker saham tidak valid")
		return
	}
	ck := nyangkutCacheKey("nyangkut_tracker", tickerClean)
	if h.tryCache(w, r, ck, 30, 3600) {
		return
	}
	data, err := h.Svc.GetOrCalculateNyangkutData(r.Context(), tickerClean)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal memproses data nyangkut %s: %v", tickerClean, err))
		return
	}
	h.storeCache(w, r, ck, data, 30, 3600)
}

// GetLeaderboard handles GET /api/v1/nyangkut/leaderboard and /api/nyangkut/leaderboard.
func (h *Handler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	limit := 20
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 100 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "limit must be between 1 and 100")
			return
		}
		limit = n
	}
	ck := nyangkutCacheKey("nyangkut_leaderboard", strconv.Itoa(limit))
	if h.tryCache(w, r, ck, 30, 3600) {
		return
	}
	items, err := h.Svc.GetLeaderboard(r.Context(), limit)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal memuat leaderboard: %v", err))
		return
	}
	h.storeCache(w, r, ck, items, 30, 3600)
}

// GetSummary handles GET /api/v1/nyangkut/summary and /api/nyangkut/summary.
func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	ck := nyangkutCacheKey("nyangkut_summary", "summary")
	if h.tryCache(w, r, ck, 30, 3600) {
		return
	}
	res, err := h.Svc.GetNyangkutSummaryStats(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal memuat ringkasan nyangkut: %v", err))
		return
	}
	h.storeCache(w, r, ck, res, 30, 3600)
}

// SearchStocks handles GET /api/v1/stock/search and /api/nyangkut/search.
func (h *Handler) SearchStocks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if len(q) < 1 {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "q is required and must be at least 1 character")
		return
	}
	limit := 8
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 20 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "limit must be between 1 and 20")
			return
		}
		limit = n
	}
	ck := nyangkutCacheKey("ticker_search", fmt.Sprintf("%s|%d", q, limit))
	if h.tryCache(w, r, ck, 30, 300) {
		return
	}
	items, err := h.Svc.SearchTickersPrediction(r.Context(), q, limit)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal melakukan pencarian ticker: %v", err))
		return
	}
	h.storeCache(w, r, ck, items, 30, 300)
}

// RefreshCache handles POST /api/nyangkut/refresh using safe SCAN+UNLINK purging.
func (h *Handler) RefreshCache(w http.ResponseWriter, r *http.Request) {
	var cleared int
	if h.Cache != nil {
		cleared = h.Cache.SafePurge(r.Context(), "nyangkut*") + h.Cache.SafePurge(r.Context(), "ticker_search:*")
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":       "ok",
		"cleared_keys": cleared,
	})
}
