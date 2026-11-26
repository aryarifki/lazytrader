package bandarmology

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

func (h *Handler) loadCausalityInputs(r *http.Request, ticker string, startDate time.Time, endDate *time.Time) ([]shared.PriceInput, []shared.FlowInput, []shared.ParticipantActivityInput) {
	prices, errPx := h.Repo.ReadPrices(r.Context(), []string{ticker}, &startDate, endDate)
	flows, errFl := h.Repo.ReadBrokerFlow(r.Context(), []string{ticker}, &startDate, endDate)
	acts, errAct := h.Repo.ReadBrokerActivity(r.Context(), []string{ticker}, &startDate, endDate)

	var pxIn []shared.PriceInput
	if errPx == nil {
		for _, p := range prices {
			if p.Close != nil && *p.Close > 0 {
				pxIn = append(pxIn, shared.PriceInput{Date: fmtDateISO(p.Date), Ticker: p.Ticker, Close: *p.Close})
			}
		}
	}
	var flIn []shared.FlowInput
	if errFl == nil {
		for _, f := range flows {
			flIn = append(flIn, shared.FlowInput{Date: fmtDateISO(f.Date), ForeignNetBroker: shared.Finite(f.ForeignNetBroker, 0)})
		}
	}
	var actIn []shared.ParticipantActivityInput
	if errAct == nil {
		for _, a := range acts {
			actIn = append(actIn, shared.ParticipantActivityInput{
				Date:            fmtDateISO(a.Date),
				ParticipantType: a.ParticipantType,
				BrokerCode:      a.BrokerCode,
				NetValue:        a.NetValue,
			})
		}
	}
	return pxIn, flIn, actIn
}

// GetCausalityInsight handles GET /api/bandar/causality/{ticker}.
func (h *Handler) GetCausalityInsight(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	q := r.URL.Query()
	rawDate := strings.TrimSpace(q.Get("analysis_date"))
	rawWin := strings.TrimSpace(q.Get("window_days"))

	ck := bandarCacheKey("bandar_causality", "causality_insight", fmt.Sprintf("%s|%s|%s", ticker, rawDate, rawWin))
	if h.tryCache(w, r, ck, 30, 1800) {
		return
	}

	lookback := 250
	if rawWin != "" {
		if n, err := strconv.Atoi(rawWin); err == nil && n > 20 {
			lookback = n
		}
	}
	calDays := int(float64(lookback) * 1.6)

	var startDate time.Time
	var endDate *time.Time
	if rawDate != "" {
		endD, err := time.Parse("2006-01-02", rawDate)
		if err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid analysis_date")
			return
		}
		endDate = &endD
		startDate = endD.AddDate(0, 0, -calDays)
	} else {
		startDate = toDateOnly(time.Now()).AddDate(0, 0, -calDays)
	}

	pxIn, flIn, actIn := h.loadCausalityInputs(r, ticker, startDate, endDate)

	foreignCausality := shared.CausalityForeignVsPrice(ticker, pxIn, flIn, 5, 250)
	partCausality := shared.CausalityByParticipant(pxIn, actIn, 5, 250)
	brokerCausality := shared.CausalityByBroker(pxIn, actIn, 15, 5, 250)

	englishMap := map[string]string{"Asing": "Foreign", "Lokal": "Local", "Pemerintah": "Government"}
	partList := make([]map[string]any, 0, len(partCausality))
	for _, p := range partCausality {
		lbl, ok := englishMap[p.ParticipantType]
		if !ok {
			lbl = p.ParticipantType
		}
		partList = append(partList, map[string]any{
			"participant":    lbl,
			"lag":            p.BestLag,
			"p_value":        p.PValue,
			"is_significant": p.Significant,
		})
	}

	brokerList := make([]map[string]any, 0, len(brokerCausality))
	for _, b := range brokerCausality {
		brokerList = append(brokerList, map[string]any{
			"code":           b.BrokerCode,
			"lag":            b.BestLag,
			"p_value":        b.PValue,
			"is_significant": b.Significant,
		})
	}

	var grangerTest any
	if foreignCausality != nil {
		grangerTest = map[string]any{
			"is_significant": foreignCausality.IsSignificant,
			"min_p_value":    foreignCausality.MinPValue,
			"best_lag":       foreignCausality.BestLag,
		}
	}

	h.storeCache(w, r, ck, map[string]any{
		"granger_test":          grangerTest,
		"participant_causality": partList,
		"top_brokers":           brokerList,
	}, 30, 1800)
}

// GetStockCausality handles GET /api/bandar/stocks/{ticker}/causality.
func (h *Handler) GetStockCausality(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	ck := bandarCacheKey("bandar_causality_stock", "causality", ticker)
	if h.tryCache(w, r, ck, 30, 1800) {
		return
	}

	startDate := toDateOnly(time.Now()).AddDate(0, 0, -400)
	pxIn, flIn, actIn := h.loadCausalityInputs(r, ticker, startDate, nil)

	foreign := shared.CausalityForeignVsPrice(ticker, pxIn, flIn, 5, 250)
	parts := shared.CausalityByParticipant(pxIn, actIn, 5, 250)
	brokers := shared.CausalityByBroker(pxIn, actIn, 15, 5, 250)

	h.storeCache(w, r, ck, map[string]any{
		"foreign":      foreign,
		"participants": parts,
		"brokers":      brokers,
	}, 30, 1800)
}
