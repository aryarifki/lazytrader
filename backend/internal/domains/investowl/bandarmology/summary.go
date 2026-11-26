package bandarmology

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
	"github.com/go-chi/chi/v5"
)

// Handler holds shared dependencies for /api/bandar endpoints.
type Handler struct {
	Repo      *Repository
	Cache     *cache.Store
	PythonBin string
	InvDir    string
	CFZoneID  string
}

func bandarCacheKey(prefix, fnName, raw string) string {
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

// GetUniverse handles GET /api/bandar/universe?mode=watchlist.
func (h *Handler) GetUniverse(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "watchlist"
	}
	ck := bandarCacheKey("bandar", "get_universe", mode)
	if h.tryCache(w, r, ck, 30, 3600) {
		return
	}
	tickers, err := shared.GetUniverse(r.Context(), h.Repo.DB, mode, nil)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.storeCache(w, r, ck, map[string]any{"tickers": tickers}, 30, 3600)
}

// GetUniverseByMode handles GET /api/bandar/universe/{mode}.
func (h *Handler) GetUniverseByMode(w http.ResponseWriter, r *http.Request) {
	mode := chi.URLParam(r, "mode")
	ck := bandarCacheKey("bandar", "universe_tickers", mode)
	if h.tryCache(w, r, ck, 30, 3600) {
		return
	}
	tickers, err := shared.GetUniverse(r.Context(), h.Repo.DB, mode, nil)
	if err != nil {
		h.storeCache(w, r, ck, map[string]any{"mode": mode, "tickers": []string{}, "count": 0, "error": err.Error()}, 30, 3600)
		return
	}
	if avail, err := shared.GetMasterTickers(r.Context(), h.Repo.DB, true); err == nil && len(avail) > 0 {
		aSet := make(map[string]struct{}, len(avail))
		for _, a := range avail {
			aSet[a] = struct{}{}
		}
		kept := make([]string, 0, len(tickers))
		for _, t := range tickers {
			if _, ok := aSet[t]; ok {
				kept = append(kept, t)
			}
		}
		tickers = kept
	}
	h.storeCache(w, r, ck, map[string]any{"mode": mode, "tickers": tickers, "count": len(tickers)}, 30, 3600)
}

// GetMasterTickers handles GET /api/bandar/tickers.
func (h *Handler) GetMasterTickers(w http.ResponseWriter, r *http.Request) {
	ck := bandarCacheKey("bandar", "master_tickers", "")
	if h.tryCache(w, r, ck, 30, 3600) {
		return
	}
	tickers, err := shared.GetMasterTickers(r.Context(), h.Repo.DB, true)
	if err != nil {
		h.storeCache(w, r, ck, map[string]any{"tickers": []string{}, "error": err.Error()}, 30, 3600)
		return
	}
	h.storeCache(w, r, ck, map[string]any{"tickers": tickers}, 30, 3600)
}

// GetSystemStatus handles GET /api/bandar/system-status.
func (h *Handler) GetSystemStatus(w http.ResponseWriter, r *http.Request) {
	ck := bandarCacheKey("bandar", "get_system_status", "")
	if h.tryCache(w, r, ck, 30, 300) {
		return
	}
	cnt, latest, err := h.Repo.GetSystemStatus(r.Context())
	if err != nil {
		h.storeCache(w, r, ck, map[string]any{"status": "error", "active_tickers": 0, "latest_date": "-", "error": err.Error()}, 30, 300)
		return
	}
	h.storeCache(w, r, ck, map[string]any{"status": "ok", "active_tickers": cnt, "latest_date": latest}, 30, 300)
}

// GetTickerDates handles GET /api/bandar/dates/{ticker}.
func (h *Handler) GetTickerDates(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	ck := bandarCacheKey("bandar", "ticker_dates", ticker)
	if h.tryCache(w, r, ck, 30, 3600) {
		return
	}
	dates, err := h.Repo.GetDistinctActivityDates(r.Context(), ticker)
	if err != nil {
		h.storeCache(w, r, ck, map[string]any{"dates": []string{}, "error": err.Error()}, 30, 3600)
		return
	}
	h.storeCache(w, r, ck, map[string]any{"dates": dates}, 30, 3600)
}

// GetDailySummary handles GET /api/bandar/daily-summary.
func (h *Handler) GetDailySummary(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	universeMode := q.Get("universe_mode")
	if universeMode == "" {
		universeMode = "all"
	}
	refresh := 0
	if raw := strings.TrimSpace(q.Get("refresh")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			refresh = n
		}
	}
	ck := bandarCacheKey("bandar", "daily_summary", fmt.Sprintf("%s|%d", universeMode, refresh))
	if refresh == 0 && h.tryCache(w, r, ck, 30, 1800) {
		return
	}

	tickers := shared.GetDynamicUniverse(r.Context(), h.Repo.DB, universeMode)
	cutoff := toDateOnly(time.Now().AddDate(0, 0, -90))
	prices, err := h.Repo.ReadPrices(r.Context(), tickers, &cutoff, nil)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	flows, err := h.Repo.ReadBrokerFlow(r.Context(), tickers, &cutoff, nil)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pxByTicker := make(map[string][]PriceRow)
	for _, p := range prices {
		pxByTicker[p.Ticker] = append(pxByTicker[p.Ticker], p)
	}
	flByTicker := make(map[string][]FlowRow)
	for _, f := range flows {
		flByTicker[f.Ticker] = append(flByTicker[f.Ticker], f)
	}

	sortedTickers := make([]string, 0, len(flByTicker))
	for tk := range flByTicker {
		sortedTickers = append(sortedTickers, tk)
	}
	sort.Strings(sortedTickers)

	type summaryItem struct {
		Ticker       string    `json:"ticker"`
		Date         string    `json:"date"`
		Signal       *string   `json:"signal"`
		SignalScore  *int      `json:"signal_score"`
		Close        *float64  `json:"close"`
		Ret5D        *float64  `json:"ret_5d"`
		ForeignNet5D *float64  `json:"foreign_net_5d"`
		Spark        []float64 `json:"spark"`
		rawScoreSort int
	}

	items := make([]summaryItem, 0, len(sortedTickers))
	counts := make(map[string]int)
	var maxAsOf string

	for _, tk := range sortedTickers {
		fsub := flByTicker[tk]
		sort.Slice(fsub, func(i, j int) bool { return fsub[i].Date.Before(fsub[j].Date) })
		last := fsub[len(fsub)-1]

		var psub []PriceRow
		for _, p := range pxByTicker[tk] {
			if !p.Date.After(last.Date) {
				psub = append(psub, p)
			}
		}
		sort.Slice(psub, func(i, j int) bool { return psub[i].Date.Before(psub[j].Date) })

		var ret5d *float64
		if len(psub) > 5 && psub[len(psub)-1].Close != nil && psub[len(psub)-6].Close != nil {
			base := *psub[len(psub)-6].Close
			if base != 0 {
				v := (*psub[len(psub)-1].Close)/base - 1.0
				ret5d = &v
			}
		}

		f5Start := max(0, len(fsub)-5)
		f5Sum := 0.0
		for _, fl := range fsub[f5Start:] {
			f5Sum += shared.Finite(fl.ForeignNetBroker, 0)
		}

		var sigScorePtr *int
		sortScore := 0
		if last.BandarSignalScore != nil {
			sVal := int(*last.BandarSignalScore)
			sigScorePtr = &sVal
			sortScore = sVal
		}

		var closePtr *float64
		if len(psub) > 0 && psub[len(psub)-1].Close != nil {
			v := *psub[len(psub)-1].Close
			closePtr = &v
		}

		spStart := max(0, len(psub)-30)
		spark := make([]float64, 0, len(psub)-spStart)
		for _, p := range psub[spStart:] {
			if p.Close != nil {
				spark = append(spark, *p.Close)
			} else {
				spark = append(spark, math.NaN())
			}
		}
		// Ensure valid finite spark values
		cleanSpark := make([]any, 0, len(spark))
		for _, sv := range spark {
			if math.IsNaN(sv) || math.IsInf(sv, 0) {
				cleanSpark = append(cleanSpark, nil)
			} else {
				cleanSpark = append(cleanSpark, sv)
			}
		}
		_ = cleanSpark

		dStr := fmtDateMidnightSpace(last.Date)
		if dStr > maxAsOf {
			maxAsOf = dStr
		}

		sigKey := "NETRAL"
		if last.BandarSignal != nil && *last.BandarSignal != "" {
			sigKey = strings.ToUpper(*last.BandarSignal)
		}
		counts[sigKey]++

		items = append(items, summaryItem{
			Ticker:       tk,
			Date:         dStr,
			Signal:       last.BandarSignal,
			SignalScore:  sigScorePtr,
			Close:        closePtr,
			Ret5D:        ret5d,
			ForeignNet5D: &f5Sum,
			Spark:        spark,
			rawScoreSort: sortScore,
		})
	}

	topCandidates := append([]summaryItem(nil), items...)
	sort.SliceStable(topCandidates, func(i, j int) bool {
		return topCandidates[i].rawScoreSort > topCandidates[j].rawScoreSort
	})
	if len(topCandidates) > 5 {
		topCandidates = topCandidates[:5]
	}

	var asOfPtr *string
	if maxAsOf != "" {
		asOfPtr = &maxAsOf
	}
	res := map[string]any{
		"as_of":          asOfPtr,
		"signal_counts":  counts,
		"items":          items,
		"top_conviction": topCandidates,
	}
	h.storeCache(w, r, ck, res, 30, 1800)
}
