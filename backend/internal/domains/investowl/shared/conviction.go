package shared

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

// LabelComponent mirrors services/bandarmology_service.py label_component.
func LabelComponent(signal string) float64 {
	switch strings.ToUpper(strings.TrimSpace(signal)) {
	case "AKUMULASI_KUAT", "STRONG_ACCUMULATION":
		return 100
	case "AKUMULASI", "ACCUMULATION", "NET_BUY":
		return 80
	case "NETRAL", "NEUTRAL":
		return 50
	case "DISTRIBUSI", "DISTRIBUTION", "NET_SELL":
		return 25
	case "DISTRIBUSI_KUAT", "STRONG_DISTRIBUTION":
		return 0
	default:
		return 40
	}
}

// PValueComponent mirrors services/bandarmology_service.py p_value_component.
func PValueComponent(pVal *float64) float64 {
	if pVal == nil || math.IsNaN(*pVal) || math.IsInf(*pVal, 0) {
		return 25.0
	}
	p := *pVal
	if p <= 0.01 {
		return 100.0
	}
	if p <= 0.05 {
		return 80.0
	}
	if p <= 0.10 {
		return 55.0
	}
	if p <= 0.25 {
		return 35.0
	}
	return 25.0
}

// ForeignComponent mirrors services/bandarmology_service.py foreign_component.
func ForeignComponent(val *float64, scale float64) float64 {
	if val == nil || math.IsNaN(*val) || math.IsInf(*val, 0) {
		return 50.0
	}
	denom := scale
	if denom <= 0 {
		denom = 25_000_000_000.0
	}
	if denom < 1_000_000_000.0 {
		denom = 1_000_000_000.0
	}
	return Round(50.0+50.0*math.Tanh(*val/denom), 2)
}

// BrokerWinComponent mirrors services/bandarmology_service.py broker_win_component.
func BrokerWinComponent(scan []AlphaScanRow, ticker string) (float64, string) {
	sub := make([]AlphaScanRow, 0)
	for _, r := range scan {
		if ticker == "" || r.Ticker == ticker {
			sub = append(sub, r)
		}
	}
	if len(sub) == 0 {
		return 50, "No broker validation sample"
	}
	sort.SliceStable(sub, func(i, j int) bool {
		a, b := sub[i], sub[j]
		if a.Significant != b.Significant {
			return a.Significant
		}
		oa, ob := a.PValueOneSided, b.PValueOneSided
		if math.IsNaN(oa) != math.IsNaN(ob) {
			return !math.IsNaN(oa)
		}
		if !math.IsNaN(oa) && oa != ob {
			return oa < ob
		}
		return a.MeanFwdReturn > b.MeanFwdReturn
	})
	top := sub[0]
	score := math.Max(0, math.Min(100, top.WinRate*100))
	pct := int(math.RoundToEven(top.WinRate * 100))
	note := top.BrokerCode + " win rate " + strconv.Itoa(pct) + "%"
	return score, note
}

// ConvictionOutput mirrors bandarmology_service.conviction_score return dict.
type ConvictionOutput struct {
	Score              float64
	PValue             *float64
	CausalityComponent float64
	SignalComponent    float64
	ForeignComponent   float64
	BrokerComponent    float64
	BrokerNote         string
}

// ConvictionScore mirrors services/bandarmology_service.py conviction_score.
func ConvictionScore(signal string, foreign5d *float64, scan []AlphaScanRow, ticker string, prices []PriceInput, flows []FlowInput, totalValues5d []float64) ConvictionOutput {
	// Python passes prices (sorted asc) and flow to causality_foreign_vs_price(..., lookback_days=250)
	causality := CausalityForeignVsPrice(ticker, prices, flows, 5, 250)
	var pPtr *float64
	if causality != nil && !math.IsNaN(causality.MinPValue) {
		v := causality.MinPValue
		pPtr = &v
	}
	pScore := PValueComponent(pPtr)
	sScore := LabelComponent(signal)

	scale := 25_000_000_000.0
	if len(totalValues5d) > 0 {
		var tot5d float64
		start := len(totalValues5d) - 5
		if start < 0 {
			start = 0
		}
		for _, tv := range totalValues5d[start:] {
			tot5d += tv
		}
		if tot5d > 0 {
			scale = math.Max(tot5d*0.15, 1_000_000_000.0)
		}
	}
	fScore := ForeignComponent(foreign5d, scale)
	wScore, wNote := BrokerWinComponent(scan, ticker)
	score := (pScore * 0.30) + (sScore * 0.30) + (fScore * 0.20) + (wScore * 0.20)

	return ConvictionOutput{
		Score:              Round(score, 1),
		PValue:             pPtr,
		CausalityComponent: pScore,
		SignalComponent:    sScore,
		ForeignComponent:   fScore,
		BrokerComponent:    wScore,
		BrokerNote:         wNote,
	}
}
