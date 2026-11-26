package shared

import (
	"math"
	"sort"
	"strings"
)

// ParticipantActivityInput holds minimal activity data for participant/broker Granger causality.
type ParticipantActivityInput struct {
	Date            string
	ParticipantType string
	BrokerCode      string
	NetValue        float64
}

// ParticipantCausalityRow mirrors analysis.causality_by_participant output.
type ParticipantCausalityRow struct {
	ParticipantType string  `json:"participant_type"`
	BestLag         int     `json:"best_lag"`
	PValue          float64 `json:"p_value"`
	RawPValue       float64 `json:"raw_p_value"`
	Significant     bool    `json:"significant"`
}

// BrokerCausalityRow mirrors analysis.causality_by_broker output.
type BrokerCausalityRow struct {
	BrokerCode  string   `json:"broker_code"`
	BestLag     int      `json:"best_lag"`
	PValue      float64  `json:"p_value"`
	RawPValue   float64  `json:"raw_p_value"`
	PValueFDR   *float64 `json:"p_value_fdr"`
	Significant bool     `json:"significant"`
}

func prepCausalityPrices(prices []PriceInput, lookbackDays int) []PriceInput {
	px := make([]PriceInput, 0, len(prices))
	for _, p := range prices {
		if p.Close > 0 {
			px = append(px, p)
		}
	}
	sort.Slice(px, func(i, j int) bool { return px[i].Date < px[j].Date })
	if lookbackDays > 0 && len(px) > lookbackDays {
		px = px[len(px)-lookbackDays:]
	}
	return px
}

func runSeriesGranger(px []PriceInput, flowByDate map[string]float64, maxLags int) (int, float64, float64, bool) {
	rets := make([]float64, len(px))
	retValid := make([]bool, len(px))
	series := make([]float64, len(px))
	for i := range px {
		if i > 0 {
			rets[i] = px[i].Close/px[i-1].Close - 1.0
			retValid[i] = true
		}
		series[i] = flowByDate[px[i].Date]
	}
	statFlow, flowValid := MakeStationary(series)
	y := make([]float64, 0, len(px))
	x := make([]float64, 0, len(px))
	for i := range px {
		if retValid[i] && flowValid[i] {
			y = append(y, rets[i])
			x = append(x, statFlow[i])
		}
	}
	_, stdY := MeanStd(y)
	_, stdX := MeanStd(x)
	if len(y) < maxLags*3+2 || stdY == 0 || stdX == 0 {
		return 0, 0, 0, false
	}
	pvals, ok := GrangerSSRFTest(y, x, maxLags)
	if !ok {
		return 0, 0, 0, false
	}
	return HolmBonferroniLags(pvals)
}

// CausalityByParticipant replicates analysis.causality_by_participant.
func CausalityByParticipant(prices []PriceInput, activities []ParticipantActivityInput, maxLags, lookbackDays int) []ParticipantCausalityRow {
	px := prepCausalityPrices(prices, lookbackDays)
	if len(px) == 0 || len(activities) == 0 {
		return []ParticipantCausalityRow{}
	}
	minDate := px[0].Date

	byPartDate := make(map[string]map[string]float64)
	for _, a := range activities {
		if a.Date < minDate || strings.TrimSpace(a.ParticipantType) == "" {
			continue
		}
		m, ok := byPartDate[a.ParticipantType]
		if !ok {
			m = make(map[string]float64)
			byPartDate[a.ParticipantType] = m
		}
		m[a.Date] += a.NetValue
	}

	out := make([]ParticipantCausalityRow, 0, len(byPartDate))
	for ptype, flowMap := range byPartDate {
		bestLag, adjP, rawP, ok := runSeriesGranger(px, flowMap, maxLags)
		if !ok {
			continue
		}
		out = append(out, ParticipantCausalityRow{
			ParticipantType: ptype,
			BestLag:         bestLag,
			PValue:          adjP,
			RawPValue:       rawP,
			Significant:     adjP < 0.05,
		})
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].PValue < out[j].PValue })
	return out
}

// CausalityByBroker replicates analysis.causality_by_broker.
func CausalityByBroker(prices []PriceInput, activities []ParticipantActivityInput, topN, maxLags, lookbackDays int) []BrokerCausalityRow {
	px := prepCausalityPrices(prices, lookbackDays)
	if len(px) == 0 || len(activities) == 0 {
		return []BrokerCausalityRow{}
	}
	minDate := px[0].Date

	absByBroker := make(map[string]float64)
	byBrokerDate := make(map[string]map[string]float64)
	for _, a := range activities {
		if a.Date < minDate || a.BrokerCode == "" {
			continue
		}
		absByBroker[a.BrokerCode] += math.Abs(a.NetValue)
		m, ok := byBrokerDate[a.BrokerCode]
		if !ok {
			m = make(map[string]float64)
			byBrokerDate[a.BrokerCode] = m
		}
		m[a.Date] += a.NetValue
	}

	type bAbs struct {
		code string
		val  float64
	}
	ranked := make([]bAbs, 0, len(absByBroker))
	for c, v := range absByBroker {
		ranked = append(ranked, bAbs{code: c, val: v})
	}
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].val != ranked[j].val {
			return ranked[i].val > ranked[j].val
		}
		return ranked[i].code < ranked[j].code
	})
	if len(ranked) > topN {
		ranked = ranked[:topN]
	}

	out := make([]BrokerCausalityRow, 0, len(ranked))
	pvals := make([]float64, 0, len(ranked))
	for _, b := range ranked {
		bestLag, adjP, rawP, ok := runSeriesGranger(px, byBrokerDate[b.code], maxLags)
		if !ok {
			continue
		}
		out = append(out, BrokerCausalityRow{
			BrokerCode: b.code,
			BestLag:    bestLag,
			PValue:     adjP,
			RawPValue:  rawP,
		})
		pvals = append(pvals, adjP)
	}
	if len(out) == 0 {
		return out
	}

	fdr := BenjaminiHochberg(pvals)
	for i := range out {
		q := fdr[i]
		out[i].PValueFDR = FinitePtr(&q)
		out[i].Significant = !math.IsNaN(q) && q < 0.05
	}
	sort.SliceStable(out, func(i, j int) bool {
		qa := math.Inf(1)
		if out[i].PValueFDR != nil {
			qa = *out[i].PValueFDR
		}
		qb := math.Inf(1)
		if out[j].PValueFDR != nil {
			qb = *out[j].PValueFDR
		}
		if qa != qb {
			return qa < qb
		}
		if out[i].PValue != out[j].PValue {
			return out[i].PValue < out[j].PValue
		}
		return out[i].BrokerCode < out[j].BrokerCode
	})
	return out
}
