package bandarmology

import (
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
	"github.com/go-chi/chi/v5"
)

func parseTSOrFallback(raw string, fallback time.Time) time.Time {
	s := strings.TrimSpace(raw)
	if s == "" || s == "undefined" || s == "null" {
		return fallback
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fallback
	}
	return t
}

var brokerFlowProfileDescriptions = map[string]string{
	"smart_foreign":       "Directional foreign institutions",
	"local_institutional": "Local institution-like accounts",
	"market_maker":        "Active on both sides; net position matters",
	"bandar_gorengan":     "Speculative operator profile",
	"retail":              "Retail-heavy platforms",
	"lainnya":             "Outside defined behavioral profiles",
}

// GetBrokerFlowDetail handles GET /api/bandar/broker-flow/{ticker}.
func (h *Handler) GetBrokerFlowDetail(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	q := r.URL.Query()
	rawDate := strings.TrimSpace(q.Get("analysis_date"))

	windowDays := 20
	if raw := strings.TrimSpace(q.Get("window_days")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			windowDays = n
		}
	}
	brokerCodes := q.Get("broker_codes")
	flowMode := q.Get("flow_mode")
	if flowMode == "" {
		flowMode = "Cumulative"
	}
	distMode := q.Get("dist_mode")
	distDate := q.Get("dist_date")
	distStart := q.Get("dist_start")
	distEnd := q.Get("dist_end")

	ck := bandarCacheKey("bandar_brokerflow", "broker_flow_detail", fmt.Sprintf("%s|%s|%d|%s|%s|%s|%s|%s|%s",
		ticker, rawDate, windowDays, brokerCodes, flowMode, distMode, distDate, distStart, distEnd))
	if h.tryCache(w, r, ck, 30, 1800) {
		return
	}

	var analysisTS time.Time
	if rawDate != "" {
		t, err := time.Parse("2006-01-02", rawDate)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid analysis_date")
			return
		}
		analysisTS = t
	} else if maxD := h.Repo.GetMaxActivityDate(r.Context(), ticker); maxD != nil {
		analysisTS = *maxD
	} else {
		analysisTS = toDateOnly(time.Now())
	}

	lookbackDays := max(windowDays+60, 180)
	fetchStart := analysisTS.AddDate(0, 0, -lookbackDays)
	fetchEnd := analysisTS
	if distStart != "" {
		if ds, err := time.Parse("2006-01-02", strings.TrimSpace(distStart)); err == nil && ds.Before(fetchStart) {
			fetchStart = ds
		}
	}

	actDF, err := h.Repo.ReadBrokerActivity(r.Context(), []string{ticker}, &fetchStart, &fetchEnd)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(actDF) == 0 {
		h.storeCache(w, r, ck, map[string]any{"error": "No broker history for " + ticker}, 30, 1800)
		return
	}

	windowStart := analysisTS.AddDate(0, 0, -windowDays)
	var actWin []ActivityRow
	for _, a := range actDF {
		if !a.Date.Before(windowStart) && !a.Date.After(analysisTS) {
			actWin = append(actWin, a)
		}
	}
	if len(actWin) == 0 {
		h.storeCache(w, r, ck, map[string]any{"error": "No activity in window"}, 30, 1800)
		return
	}

	absMap := make(map[string]float64)
	for _, a := range actWin {
		absMap[a.BrokerCode] += math.Abs(a.NetValue)
	}
	allCodes := make([]string, 0, len(absMap))
	for c := range absMap {
		allCodes = append(allCodes, c)
	}
	sort.Strings(allCodes)

	ranked := append([]string(nil), allCodes...)
	sort.SliceStable(ranked, func(i, j int) bool {
		if absMap[ranked[i]] != absMap[ranked[j]] {
			return absMap[ranked[i]] > absMap[ranked[j]]
		}
		return ranked[i] < ranked[j]
	})

	var selected []string
	if strings.TrimSpace(brokerCodes) != "" {
		for _, p := range strings.Split(brokerCodes, ",") {
			if s := strings.TrimSpace(p); s != "" {
				selected = append(selected, s)
			}
		}
	}
	if len(selected) == 0 {
		selected = append([]string(nil), ranked[:min(len(ranked), 3)]...)
	}
	defaultCodes := append([]string(nil), ranked[:min(len(ranked), 3)]...)

	compareData := BrokerCompareData(actWin, selected, flowMode)

	var distStartTS, distEndTS time.Time
	switch distMode {
	case "Single day":
		distStartTS = parseTSOrFallback(distDate, analysisTS)
		distEndTS = parseTSOrFallback(distDate, analysisTS)
	case "Date range":
		distStartTS = parseTSOrFallback(distStart, windowStart)
		distEndTS = parseTSOrFallback(distEnd, analysisTS)
	default:
		distStartTS = windowStart
		distEndTS = analysisTS
	}

	distData := BrokerDistributionDataRange(actDF, distStartTS, distEndTS)
	summary := BrokerSummaryTable(distData)

	type bpKey struct{ broker, ptype string }
	type bpAgg struct{ buy, sell, net, freq float64 }
	grouped := make(map[bpKey]*bpAgg)
	for _, a := range actWin {
		k := bpKey{a.BrokerCode, a.ParticipantType}
		g, ok := grouped[k]
		if !ok {
			g = &bpAgg{}
			grouped[k] = g
		}
		g.buy += a.BuyValue
		g.sell += a.SellValue
		g.net += a.NetValue
		g.freq += shared.Finite(a.Frequency, 0)
	}
	detailRows := make([]map[string]any, 0, len(grouped))
	for k, g := range grouped {
		detailRows = append(detailRows, map[string]any{
			"broker":  k.broker,
			"type":    shared.ParticipantLabel(k.ptype),
			"subtype": BrokerSubtype(k.ptype, g.net, g.freq),
			"buy":     g.buy,
			"sell":    g.sell,
			"net":     g.net,
			"freq":    g.freq,
		})
	}
	sort.SliceStable(detailRows, func(i, j int) bool {
		return math.Abs(detailRows[i]["net"].(float64)) > math.Abs(detailRows[j]["net"].(float64))
	})

	profEntries := ProfileFlowFromActivity(actWin)
	profRows := make([]map[string]any, 0, len(profEntries))
	for _, pe := range profEntries {
		desc := pe.Description
		if d, ok := brokerFlowProfileDescriptions[pe.Profile]; ok {
			desc = d
		}
		profRows = append(profRows, map[string]any{
			"profile":     pe.Profile,
			"label":       pe.Label,
			"description": desc,
			"net":         pe.Net,
			"top_brokers": pe.TopBrokers,
		})
	}

	res := map[string]any{
		"ticker":                ticker,
		"analysis_date":         fmtDateISO(analysisTS),
		"window_start":          fmtDateISO(windowStart),
		"all_codes":             allCodes,
		"ranked_codes":          ranked,
		"default_codes":         defaultCodes,
		"selected_codes":        selected,
		"compare_chart":         compareData,
		"distribution":          distData,
		"summary":               summary,
		"profile_flow":          profRows,
		"profile_broker_detail": ProfileBrokerDetailTable(actWin),
		"detail_rows":           detailRows,
	}
	h.storeCache(w, r, ck, res, 30, 1800)
}
