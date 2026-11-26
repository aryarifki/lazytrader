package foreign_flow_pro

import (
	"fmt"
	"math"

	"github.com/aryarifki/lazytrader/backend/internal/domains/terminal/shared"
)

// CalculateSMA calculates Simple Moving Average over period, returning nil before period is filled.
func CalculateSMA(series []float64, period int) []*float64 {
	out := make([]*float64, len(series))
	if period <= 0 {
		return out
	}
	for i := range series {
		if i+1 < period {
			continue
		}
		var sum float64
		for j := i + 1 - period; j <= i; j++ {
			sum += series[j]
		}
		v := sum / float64(period)
		out[i] = &v
	}
	return out
}

// CalculateSMAPointers mirrors calculate_sma when intermediate elements may be nil.
func CalculateSMAPointers(series []*float64, period int) []*float64 {
	out := make([]*float64, len(series))
	if period <= 0 {
		return out
	}
	for i := range series {
		if i+1 < period {
			continue
		}
		var sum float64
		cnt := 0
		for j := i + 1 - period; j <= i; j++ {
			if series[j] != nil {
				sum += *series[j]
				cnt++
			}
		}
		if cnt > 0 {
			v := sum / float64(cnt)
			out[i] = &v
		}
	}
	return out
}

// CalculateRollingCorrelation computes rolling Pearson correlation matching calculate_rolling_correlation.
func CalculateRollingCorrelation(x, y []float64, period int) []*float64 {
	n := len(x)
	out := make([]*float64, n)
	if period <= 0 || len(y) != n {
		return out
	}
	fPeriod := float64(period)
	for i := range x {
		if i+1 < period {
			continue
		}
		start := i + 1 - period
		var sumX, sumY float64
		for j := start; j <= i; j++ {
			sumX += x[j]
			sumY += y[j]
		}
		meanX := sumX / fPeriod
		meanY := sumY / fPeriod

		var num, denX, denY float64
		for j := start; j <= i; j++ {
			dx := x[j] - meanX
			dy := y[j] - meanY
			num += dx * dy
			denX += dx * dx
			denY += dy * dy
		}
		den := math.Sqrt(denX * denY)
		var val float64
		if den != 0 {
			val = shared.PyRound(num/den, 4)
		}
		out[i] = &val
	}
	return out
}

// GetFPColorAndBadge computes Foreign Power badge and RGB/hex lerp color verbatim.
func GetFPColorAndBadge(fpScore float64) FPPoint {
	badge := "LFP"
	if fpScore >= 35.0 {
		badge = "HFP"
	} else if fpScore >= 18.0 {
		badge = "MFP"
	}

	var color string
	if fpScore >= 40.0 {
		color = "#FF6820"
	} else if fpScore <= 0.0 {
		color = "#FFB6C1"
	} else if fpScore >= 20.0 {
		t := (fpScore - 20.0) / 20.0
		r := int(211 + t*(255-211))
		g := int(211 + t*(104-211))
		b := int(211 + t*(32-211))
		color = fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
	} else {
		t := fpScore / 20.0
		r := int(255 + t*(211-255))
		g := int(182 + t*(211-182))
		b := int(193 + t*(211-193))
		color = fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
	}

	return FPPoint{
		Score: shared.PyRound(fpScore, 1),
		Badge: badge,
		Color: color,
	}
}

// ProcessForeignFlowAnalytics computes all Foreign Flow Pro overlays, signals, and statistics.
func ProcessForeignFlowAnalytics(points []DailyPoint) Analytics {
	n := len(points)
	closes := make([]float64, n)
	highs := make([]float64, n)
	lows := make([]float64, n)
	opens := make([]float64, n)
	ffDaily := make([]float64, n)
	dates := make([]string, n)

	valuesMiliar := make([]float64, n)
	fcoDaily := make([]float64, n)
	kfDaily := make([]float64, n)
	xlDaily := make([]float64, n)
	srDaily := make([]float64, n)
	partSeries := make([]float64, n)
	katSeries := make([]float64, n)
	kfhSeries := make([]float64, n)
	kretSeries := make([]float64, n)

	for i, p := range points {
		dates[i] = p.Date
		c := floatPtrOr(p.Close, 0.0)
		closes[i] = c
		highs[i] = floatPtrOr(p.High, c)
		lows[i] = floatPtrOr(p.Low, c)
		opens[i] = floatPtrOr(p.Open, c)
		ffDaily[i] = p.ForeignNetVal
		vol := floatPtrOr(p.Volume, 0.0)

		raw := p.RawSeries
		var valM, fcoVal, kfVal, xlVal, srVal, part, kat, kret, kfh float64
		if len(raw) >= 16 {
			if f, ok := rawFloat(raw, 5); ok {
				valM = f
			} else {
				valM = vol * c / 1e9
			}
			if f, ok := rawFloat(raw, 14); ok {
				fcoVal = f
			} else {
				fcoVal = p.ForeignNetVal
			}
			if f, ok := rawFloat(raw, 15); ok {
				kfVal = f
			}
			if f, ok := rawFloat(raw, 20); ok {
				xlVal = f
			}
			if f, ok := rawFloat(raw, 21); ok {
				srVal = f
			}
			if f, ok := rawFloat(raw, 9); ok {
				part = f
			} else {
				part = floatPtrOr(p.PartisipasiAsing, 0.0)
			}
			if f, ok := rawFloat(raw, 10); ok {
				kat = f
			} else {
				kat = 50.0
			}
			if f, ok := rawFloat(raw, 11); ok {
				kret = f
			} else {
				kret = 20.0
			}
			if f, ok := rawFloat(raw, 12); ok {
				kfh = f
			} else {
				kfh = 10.0
			}
		} else {
			valM = (vol * c) / 1e9
			fcoVal = p.ForeignNetVal
			part = floatPtrOr(p.PartisipasiAsing, 0.0)
			kat = 50.0
			kret = 20.0
			kfh = 10.0
		}

		valuesMiliar[i] = shared.PyRound(valM, 2)
		fcoDaily[i] = shared.PyRound(fcoVal, 4)
		kfDaily[i] = shared.PyRound(kfVal, 4)
		xlDaily[i] = shared.PyRound(xlVal, 4)
		srDaily[i] = shared.PyRound(srVal, 4)
		partSeries[i] = shared.PyRound(part, 2)
		katSeries[i] = shared.PyRound(kat, 2)
		kfhSeries[i] = shared.PyRound(kfh, 2)
		kretSeries[i] = shared.PyRound(kret, 2)
	}

	// 1. Moving Averages
	smaClose10 := CalculateSMA(closes, 10)
	smaClose20 := CalculateSMA(closes, 20)
	smaClose40 := CalculateSMA(closes, 40)
	smaClose60 := CalculateSMA(closes, 60)
	smaClose80 := CalculateSMA(closes, 80)
	smaPart5 := CalculateSMA(partSeries, 5)
	smaFF5 := CalculateSMA(ffDaily, 5)
	smaFF25 := CalculateSMA(ffDaily, 25)
	smaFF50 := CalculateSMA(ffDaily, 50)
	smaVal20 := CalculateSMA(valuesMiliar, 20)

	// 2. Cumulative Lines
	cumFF := make([]float64, n)
	cumFCo := make([]float64, n)
	cumKF := make([]float64, n)
	cumXL := make([]float64, n)
	cumSR := make([]float64, n)
	cumLNR := make([]float64, n)
	var runFF, runFCo, runKF, runXL, runSR, runLNR float64
	for i := 0; i < n; i++ {
		runFF += ffDaily[i]
		runFCo += fcoDaily[i]
		runKF += kfDaily[i]
		runXL += xlDaily[i]
		runSR += srDaily[i]
		runLNR += -(ffDaily[i] + kfDaily[i] + srDaily[i])

		cumFF[i] = shared.PyRound(runFF, 2)
		cumFCo[i] = shared.PyRound(runFCo, 2)
		cumKF[i] = shared.PyRound(runKF, 2)
		cumXL[i] = shared.PyRound(runXL, 2)
		cumSR[i] = shared.PyRound(runSR, 2)
		cumLNR[i] = shared.PyRound(runLNR, 2)
	}

	// 3. Rolling Correlations
	corr20 := CalculateRollingCorrelation(closes, cumFF, 20)
	corrKF20 := CalculateRollingCorrelation(closes, cumKF, 20)

	// 4. Foreign Power 2021 & Special Candle Detector
	fpResults := make([]FPPoint, n)
	specialCandles := make([]SpecialCandle, 0)
	for i := 0; i < n; i++ {
		c := closes[i]
		part5 := partSeries[i]
		if smaPart5[i] != nil && *smaPart5[i] != 0 {
			part5 = *smaPart5[i]
		}
		kat := katSeries[i]
		kfh := kfhSeries[i]
		cor := 0.5
		if corr20[i] != nil {
			cor = *corr20[i]
		}

		bobot := 3.0
		if cor < 0 {
			bobot = 6.0
		}
		// Explicit float64() prevents ARM64 FMA contraction so IEEE-754 intermediate rounding matches CPython.
		weightedKatKfh := (float64(0.4*kat) + float64(0.6*kfh)) / 2.0
		fpScore := (float64(2.0*part5) + weightedKatKfh + float64(cor*bobot*5.0)) / 3.0
		if fpScore < 0.0 {
			fpScore = 0.0
		} else if fpScore > 100.0 {
			fpScore = 100.0
		}
		fpResults[i] = GetFPColorAndBadge(fpScore)

		isOrange, isYellow, isPurple, isRed := false, false, false, false
		ffCurr := ffDaily[i]
		sff5 := smaFF5[i]
		sff25 := smaFF25[i]
		sff50 := smaFF50[i]

		// ORANGE: SIGNAL_MET (Massive Accumulation)
		if sff5 != nil && sff25 != nil && fpScore > 20 {
			condFlow := (ffCurr > *sff25 && ffCurr > *sff5) || (sff50 != nil && ffCurr > *sff50 && ffCurr > *sff5)
			condPrice := false
			for _, ma := range []*float64{smaClose10[i], smaClose20[i], smaClose40[i], smaClose60[i], smaClose80[i]} {
				if ma != nil && c > *ma {
					condPrice = true
					break
				}
			}
			if condFlow && condPrice {
				isOrange = true
			}
		}

		// YELLOW: MARK_UP_FF
		if !isOrange && i >= 3 && fpScore > 20 {
			cPrev1 := closes[i-1]
			cPrev3 := closes[i-3]
			var roc1, roc3 float64
			if cPrev1 > 0 {
				roc1 = (c - cPrev1) / cPrev1 * 100.0
			}
			if cPrev3 > 0 {
				roc3 = (c - cPrev3) / cPrev3 * 100.0
			}
			if roc1 > 1.5 && math.Abs(roc3) > 2.5 {
				smaC9 := smaClose10[i]
				if smaC9 != nil && *smaC9 != 0 && c > *smaC9 && sff5 != nil && *sff5 != 0 && ffCurr < *sff5 {
					isYellow = true
				}
			}
		}

		// PURPLE: FOREIGN_SUDDEN_SELL
		if !isOrange && !isYellow && i >= 20 && fpScore > 20 {
			val20 := floatPtrOr(smaVal20[i], 0.0)
			if val20 > 3.0 {
				start := i - 40
				if start < 0 {
					start = 0
				}
				llv40 := ffDaily[start]
				for k := start + 1; k <= i; k++ {
					if ffDaily[k] < llv40 {
						llv40 = ffDaily[k]
					}
				}
				if ffCurr <= llv40 && ffCurr < -10.0 {
					isPurple = true
				}
			}
		}

		// RED: PAHLAWAN_BURSA_NYANGKUT
		if !isOrange && !isYellow && !isPurple && i >= 20 {
			start := i - 20
			if start < 0 {
				start = 0
			}
			var sumKF float64
			for k := start; k <= i; k++ {
				sumKF += kfDaily[k]
			}
			if sumKF > 0 {
				cPrev2 := closes[i-2]
				cPrev3 := closes[i-3]
				cPrev5 := closes[i-5]
				var roc2, roc3, roc5 float64
				if cPrev2 > 0 {
					roc2 = (c - cPrev2) / cPrev2 * 100.0
				}
				if cPrev3 > 0 {
					roc3 = (c - cPrev3) / cPrev3 * 100.0
				}
				if cPrev5 > 0 {
					roc5 = (c - cPrev5) / cPrev5 * 100.0
				}
				if roc2 < -5.0 || roc3 < -5.0 || roc5 < -5.0 {
					isRed = true
				}
			}
		}

		switch {
		case isOrange:
			specialCandles = append(specialCandles, SpecialCandle{Date: points[i].Date, Type: "ORANGE", Color: "#FF6820", Label: "Massive Acc"})
		case isYellow:
			specialCandles = append(specialCandles, SpecialCandle{Date: points[i].Date, Type: "YELLOW", Color: "#FFFF00", Label: "Mark Up"})
		case isPurple:
			specialCandles = append(specialCandles, SpecialCandle{Date: points[i].Date, Type: "PURPLE", Color: "#C724F5", Label: "Sudden Sell"})
		case isRed:
			specialCandles = append(specialCandles, SpecialCandle{Date: points[i].Date, Type: "RED", Color: "#D62839", Label: "PB Nyangkut"})
		}
	}

	// 5. Highlights for Today
	streak := 0
	for i := n - 1; i >= 0; i-- {
		if ffDaily[i] > 0 {
			streak++
		} else {
			break
		}
	}
	highlights := make([]string, 0, 3)
	if streak >= 2 {
		highlights = append(highlights, fmt.Sprintf("Foreign Inflow beruntun selama %d hari perdagangan.", streak))
	}
	latestFP := fpResults[n-1]
	highlights = append(highlights, fmt.Sprintf("Foreign Power berada pada status %s (%.1f pts).", latestFP.Badge, latestFP.Score))

	wStart := n - 20
	if wStart < 0 {
		wStart = 0
	}
	var last20FF float64
	for i := wStart; i < n; i++ {
		last20FF += ffDaily[i]
	}
	if last20FF >= 0 {
		highlights = append(highlights, fmt.Sprintf("Akumulasi asing 20 hari terakhir: +Rp %.2f Miliar.", last20FF))
	} else {
		highlights = append(highlights, fmt.Sprintf("Distribusi asing 20 hari terakhir: -Rp %.2f Miliar.", math.Abs(last20FF)))
	}

	// 6. Statistics
	stats := buildStatistics(points, closes, ffDaily, cumFF, cumFCo, cumKF, cumXL, cumSR, specialCandles)

	var lastPrice, priceChange float64
	if n > 0 {
		lastPrice = closes[n-1]
		if closes[0] > 0 {
			priceChange = shared.PyRound((closes[n-1]-closes[0])/closes[0]*100.0, 2)
		}
	}

	return Analytics{
		Dates:                    dates,
		Opens:                    opens,
		Highs:                    highs,
		Lows:                     lows,
		Closes:                   closes,
		ValuesMiliar:             valuesMiliar,
		FFDailyVal:               ffDaily,
		FCoDailyVal:              fcoDaily,
		KFDailyVal:               kfDaily,
		XLDailyVal:               xlDaily,
		SemiRitelDailyVal:        srDaily,
		CumFF:                    cumFF,
		CumFCo:                   cumFCo,
		CumKF:                    cumKF,
		CumXL:                    cumXL,
		CumSR:                    cumSR,
		CumLNR:                   cumLNR,
		PartisipasiSeries:        partSeries,
		KepemilikanAsingTotal:    katSeries,
		KepemilikanForeignHidden: kfhSeries,
		KepemilikanRitel:         kretSeries,
		Correlation20:            roundFloatPtrSlice(corr20, 2),
		CorrelationKF20:          roundFloatPtrSlice(corrKF20, 2),
		MA20Close:                smaClose20,
		FPSeries:                 fpResults,
		LatestFP:                 latestFP,
		SpecialCandles:           specialCandles,
		Highlights:               highlights,
		LastPrice:                lastPrice,
		PriceChange:              priceChange,
		Statistics:               stats,
	}
}

func rawFloat(raw []any, idx int) (float64, bool) {
	if idx < 0 || idx >= len(raw) || raw[idx] == nil {
		return 0, false
	}
	switch v := raw[idx].(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}

func floatPtrOr(p *float64, fallback float64) float64 {
	if p == nil {
		return fallback
	}
	return *p
}

func roundFloatPtrSlice(in []*float64, digits int) []*float64 {
	out := make([]*float64, len(in))
	for i, p := range in {
		if p != nil {
			v := shared.PyRound(*p, digits)
			out[i] = &v
		}
	}
	return out
}
