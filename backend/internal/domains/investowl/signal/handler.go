package signal

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
	"github.com/aryarifki/lazytrader/backend/internal/core/flags"
	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
)

// Handler serves /api/signal/* endpoints.
type Handler struct {
	Repo  *Repository
	Cache *cache.Store
	Flags *flags.Store
}

func queryCacheKey(prefix, name, rawQuery string) string {
	h := sha256.Sum256([]byte(rawQuery))
	return fmt.Sprintf("%s:%s:%s", prefix, name, hex.EncodeToString(h[:])[:12])
}

func (h *Handler) tryServeCache(w http.ResponseWriter, r *http.Request, key string, clientMaxAge, sMaxAge int) bool {
	if h.Cache == nil {
		return false
	}
	if raw, ok := h.Cache.Get(r.Context(), key); ok {
		httpx.WriteCachedJSON(w, raw, clientMaxAge, sMaxAge, true)
		return true
	}
	return false
}

func (h *Handler) storeAndWrite(w http.ResponseWriter, r *http.Request, key string, v any, clientMaxAge, sMaxAge int) {
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

// Refresh handles POST /api/signal/refresh.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cleared := 0
	if h.Cache != nil {
		cleared += h.Cache.SafePurge(r.Context(), "signal:*")
		cleared += h.Cache.SafePurge(r.Context(), "confluence_v2:*")
		cleared += h.Cache.SafePurge(r.Context(), "confluence_raw*")
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":       "ok",
		"cleared_keys": cleared,
	})
}

// GetDaily handles GET /api/signal/daily.
func (h *Handler) GetDaily(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := 100
	if raw := strings.TrimSpace(q.Get("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid limit")
			return
		}
		limit = n
	}
	var minScore *float64
	if raw := strings.TrimSpace(q.Get("min_score")); raw != "" {
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil || v < 0 || v > 100 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "min_score must be between 0 and 100")
			return
		}
		minScore = &v
	}
	topOnly := false
	if raw := strings.TrimSpace(q.Get("top_only")); raw != "" {
		b, err := strconv.ParseBool(raw)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid top_only")
			return
		}
		topOnly = b
	}

	ck := queryCacheKey("signal", "daily", fmt.Sprintf("%d|%v|%v", limit, minScore, topOnly))
	if h.tryServeCache(w, r, ck, 30, 300) {
		return
	}
	res, status, err := h.Repo.GetDailySignals(r.Context(), limit, minScore, topOnly)
	if err != nil {
		httpx.WriteError(w, status, err.Error())
		return
	}
	h.storeAndWrite(w, r, ck, res, 30, 300)
}

// GetBacktestStats handles GET /api/signal/backtest/stats.
func (h *Handler) GetBacktestStats(w http.ResponseWriter, r *http.Request) {
	ck := queryCacheKey("signal", "backtest_stats", "")
	if h.tryServeCache(w, r, ck, 30, 3600) {
		return
	}
	res, err := h.Repo.GetBacktestStats(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.storeAndWrite(w, r, ck, res, 30, 3600)
}

// GetBacktestHistory handles GET /api/signal/backtest/history.
func (h *Handler) GetBacktestHistory(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	status := q.Get("status")
	if status == "" {
		status = "ALL"
	}
	limit := 50
	if raw := strings.TrimSpace(q.Get("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid limit")
			return
		}
		limit = n
	}
	offset := 0
	if raw := strings.TrimSpace(q.Get("offset")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid offset")
			return
		}
		offset = n
	}

	ck := queryCacheKey("signal", "backtest_history", fmt.Sprintf("%s|%d|%d", status, limit, offset))
	if h.tryServeCache(w, r, ck, 30, 300) {
		return
	}
	res, err := h.Repo.GetBacktestHistory(r.Context(), status, limit, offset)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.storeAndWrite(w, r, ck, res, 30, 300)
}

// GetOutcomesSummary handles GET /api/signal/outcomes/summary.
func (h *Handler) GetOutcomesSummary(w http.ResponseWriter, r *http.Request) {
	ck := queryCacheKey("signal", "outcomes_summary", "")
	if h.tryServeCache(w, r, ck, 30, 300) {
		return
	}
	res, err := h.Repo.GetOutcomesSummary(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.storeAndWrite(w, r, ck, res, 30, 300)
}

// ExportOutcomes handles GET /api/signal/outcomes/export.
func (h *Handler) ExportOutcomes(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := 1000
	if raw := strings.TrimSpace(q.Get("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid limit")
			return
		}
		limit = n
	}
	res, err := h.Repo.ExportOutcomes(r.Context(), limit, q.Get("status"))
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}

// GetActive handles GET /api/signal/active.
func (h *Handler) GetActive(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid limit")
			return
		}
		limit = n
	}
	epView := true
	if h.Flags != nil {
		epView = h.Flags.GetBool("SIGNAL_EPISODE_VIEW", true)
	}
	ck := queryCacheKey("signal", "active", fmt.Sprintf("%d|%v", limit, epView))
	if h.tryServeCache(w, r, ck, 30, 300) {
		return
	}
	res, err := h.Repo.GetActiveSignals(r.Context(), limit, epView)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.storeAndWrite(w, r, ck, res, 30, 300)
}

// GetConfluenceScreener handles GET /api/signal/confluence-screener.
func (h *Handler) GetConfluenceScreener(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	universeMode := q.Get("universe_mode")
	if universeMode == "" {
		universeMode = "lq45"
	}
	var analysisDate *time.Time
	rawDate := strings.TrimSpace(q.Get("analysis_date"))
	if rawDate != "" {
		t, err := time.Parse("2006-01-02", rawDate)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid analysis_date")
			return
		}
		analysisDate = &t
	}
	windowDays := 20
	if raw := strings.TrimSpace(q.Get("window_days")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid window_days")
			return
		}
		windowDays = n
	}
	tickers := q.Get("tickers")
	minConv := 0.0
	if raw := strings.TrimSpace(q.Get("min_conviction")); raw != "" {
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid min_conviction")
			return
		}
		minConv = v
	}
	minLiq := 0.0
	if raw := strings.TrimSpace(q.Get("min_liquidity")); raw != "" {
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid min_liquidity")
			return
		}
		minLiq = v
	}
	tagFilter := q.Get("tag_filter")
	sigFilter := q.Get("signal_filter")

	ck := queryCacheKey("confluence_v2", "screener", fmt.Sprintf("%s|%s|%d|%s|%v|%v|%s|%s",
		universeMode, rawDate, windowDays, tickers, minConv, minLiq, tagFilter, sigFilter))
	if h.tryServeCache(w, r, ck, 30, 1800) {
		return
	}

	res, err := h.Repo.GetConfluenceScreener(
		r.Context(), h.Cache, universeMode, analysisDate, windowDays,
		tickers, minConv, minLiq, tagFilter, sigFilter,
	)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.storeAndWrite(w, r, ck, res, 30, 1800)
}
