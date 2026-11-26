package shareholders

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/go-chi/chi/v5"
)

// Handler serves /api/bandar/shareholders endpoints.
type Handler struct {
	Svc   *Service
	Cache *cache.Store
}

func shCacheKey(prefix, fnName, raw string) string {
	h := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%s:%s:%s", prefix, fnName, hex.EncodeToString(h[:])[:12])
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

// GetTopCrossHoldings handles GET /api/bandar/shareholders/top-cross-holdings.
func (h *Handler) GetTopCrossHoldings(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := 25
	if raw := strings.TrimSpace(q.Get("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 100 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "limit must be between 1 and 100")
			return
		}
		limit = n
	}
	minHoldings := 2
	if raw := strings.TrimSpace(q.Get("min_holdings")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 2 || n > 20 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "min_holdings must be between 2 and 20")
			return
		}
		minHoldings = n
	}

	ck := shCacheKey("shareholders", "get_top_cross_holdings", fmt.Sprintf("%d|%d", limit, minHoldings))
	if h.tryCache(w, r, ck, 30, 7200) {
		return
	}
	res, err := h.Svc.GetTopCrossHoldings(r.Context(), limit, minHoldings)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.storeCache(w, r, ck, res, 30, 7200)
}

// SearchShareholders handles GET /api/bandar/shareholders/search.
func (h *Handler) SearchShareholders(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if len(q) < 2 {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "q is required and must be at least 2 characters")
		return
	}
	limit := 30
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 100 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "limit must be between 1 and 100")
			return
		}
		limit = n
	}

	ck := shCacheKey("shareholders", "search_shareholders", fmt.Sprintf("%s|%d", q, limit))
	if h.tryCache(w, r, ck, 30, 3600) {
		return
	}
	res, err := h.Svc.SearchShareholders(r.Context(), q, limit)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.storeCache(w, r, ck, res, 30, 3600)
}

// GetConglomerateGroups handles GET /api/bandar/shareholders/groups.
func (h *Handler) GetConglomerateGroups(w http.ResponseWriter, r *http.Request) {
	ck := shCacheKey("shareholders_v2", "get_conglomerate_groups", "default")
	if h.tryCache(w, r, ck, 30, 7200) {
		return
	}
	res := h.Svc.GetConglomerateGroups(r.Context())
	h.storeCache(w, r, ck, res, 30, 7200)
}

// GetOwnershipNetworkGraph handles GET /api/bandar/shareholders/network-graph.
func (h *Handler) GetOwnershipNetworkGraph(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	group := q.Get("group")
	if !q.Has("group") && group == "" {
		group = "barito"
	}
	ticker := q.Get("ticker")
	minPct := 1.0
	if raw := strings.TrimSpace(q.Get("min_pct")); raw != "" {
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil || v < 0.1 || v > 50.0 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "min_pct must be between 0.1 and 50.0")
			return
		}
		minPct = v
	}

	ck := shCacheKey("shareholders", "get_ownership_network_graph", fmt.Sprintf("%s|%s|%.4f", group, ticker, minPct))
	if h.tryCache(w, r, ck, 30, 7200) {
		return
	}
	res, err := h.Svc.GetOwnershipNetworkGraph(r.Context(), group, ticker, minPct)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.storeCache(w, r, ck, res, 30, 7200)
}

// GetStockShareholders handles GET /api/bandar/shareholders/{ticker}.
func (h *Handler) GetStockShareholders(w http.ResponseWriter, r *http.Request) {
	rawTicker := chi.URLParam(r, "ticker")
	normTicker := strings.TrimSpace(strings.ReplaceAll(strings.ToUpper(rawTicker), ".JK", ""))
	if normTicker == "" {
		httpx.WriteError(w, http.StatusBadRequest, "ticker is required")
		return
	}

	ck := shCacheKey("shareholders", "get_stock_shareholders", normTicker)
	if h.tryCache(w, r, ck, 30, 7200) {
		return
	}
	res, err := h.Svc.GetStockShareholders(r.Context(), rawTicker)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, fmt.Sprintf("Tidak ada data pemegang saham untuk ticker '%s'", normTicker))
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.storeCache(w, r, ck, res, 30, 7200)
}

// RefreshCache handles POST /api/bandar/shareholders/refresh using safe SCAN+UNLINK purging.
func (h *Handler) RefreshCache(w http.ResponseWriter, r *http.Request) {
	var cleared int
	if h.Cache != nil {
		cleared = h.Cache.SafePurge(r.Context(), "shareholders*")
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":       "ok",
		"cleared_keys": cleared,
	})
}
