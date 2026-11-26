package foreign_flow_pro

import (
	"math"
	"strconv"

	"github.com/aryarifki/lazytrader/backend/internal/domains/terminal/shared"
)

func buildStatistics(
	points []DailyPoint,
	closes, ffDaily, cumFF, cumFCo, cumKF, cumXL, cumSR []float64,
	specialCandles []SpecialCandle,
) Statistics {
	n := len(points)

	calcIndicatorCorr := func(cumSeries []float64) *float64 {
		if len(cumSeries) < 10 {
			return nil
		}
		c5 := CalculateRollingCorrelation(cumSeries, closes, 5)
		c20 := CalculateRollingCorrelation(cumSeries, closes, 20)
		c50 := CalculateRollingCorrelation(cumSeries, closes, 50)
		c100 := CalculateRollingCorrelation(cumSeries, closes, 100)
		combined := make([]float64, len(cumSeries))
		for idx := range cumSeries {
			var sum float64
			cnt := 0
			for _, arr := range [][]*float64{c5, c20, c50, c100} {
				if arr[idx] != nil {
					sum += *arr[idx]
					cnt++
				}
			}
			if cnt > 0 {
				combined[idx] = (sum / float64(cnt)) * 10.0
			}
		}
		smoothed := CalculateSMA(combined, 10)
		if len(smoothed) == 0 || smoothed[len(smoothed)-1] == nil {
			return nil
		}
		v := shared.PyRound(*smoothed[len(smoothed)-1], 2)
		return &v
	}

	calcAvgConsecutive := func(days int) *float64 {
		w := days
		if n < w {
			w = n
		}
		if w < 1 {
			return nil
		}
		streaks := make([]int, 0)
		current := 0
		for _, v := range ffDaily[n-w:] {
			if v > 0 {
				current++
			} else {
				if current >= 2 {
					streaks = append(streaks, current)
				}
				current = 0
			}
		}
		if current >= 2 {
			streaks = append(streaks, current)
		}
		if len(streaks) == 0 {
			return nil
		}
		var sum int
		for _, s := range streaks {
			sum += s
		}
		v := shared.PyRound(float64(sum)/float64(len(streaks)), 1)
		return &v
	}

	calcInflowHargaNaik := func(days int) InflowNaikWindow {
		w := days
		if n < w {
			w = n
		}
		if w < 2 {
			return InflowNaikWindow{Count: 0, Pct: nil}
		}
		start := n - w
		if start < 1 {
			start = 1
		}
		inflowCount := 0
		inflowUpCount := 0
		for idx := start; idx < n; idx++ {
			if ffDaily[idx] > 0 {
				inflowCount++
				if closes[idx] > closes[idx-1] {
					inflowUpCount++
				}
			}
		}
		var pct *float64
		if inflowCount > 0 {
			v := shared.PyRound(float64(inflowUpCount)/float64(inflowCount)*100.0, 0)
			pct = &v
		}
		return InflowNaikWindow{Count: inflowCount, Pct: pct}
	}

	calcMarkdown := func(days int) int {
		w := days
		if n < w {
			w = n
		}
		if w < 2 {
			return 0
		}
		start := n - w
		if start < 1 {
			start = 1
		}
		count := 0
		for idx := start; idx < n; idx++ {
			if closes[idx] < closes[idx-1] && ffDaily[idx] > 0 {
				count++
			}
		}
		return count
	}

	// Build fast date+type index matching `any(sc["date"] == points[idx]["date"] and sc["type"] == ...)`
	candleByDate := make(map[string]map[string]bool, len(specialCandles))
	for _, sc := range specialCandles {
		if candleByDate[sc.Date] == nil {
			candleByDate[sc.Date] = make(map[string]bool, 2)
		}
		candleByDate[sc.Date][sc.Type] = true
	}

	calcSignalType := func(sigType string) SignalTypeStats {
		horizons := [3]int{5, 10, 20}
		res := map[string]SignalHorizon{
			"5":  {N: 0, Up: 0, Pct: 0, Median: 0.0},
			"10": {N: 0, Up: 0, Pct: 0, Median: 0.0},
			"20": {N: 0, Up: 0, Pct: 0, Median: 0.0},
		}
		events := 0
		for idx := 0; idx < n-21; idx++ {
			triggered := false
			switch sigType {
			case "streak_inflow":
				if idx >= 3 && ffDaily[idx] > 0 && ffDaily[idx-1] > 0 && ffDaily[idx-2] > 0 && ffDaily[idx-3] > 0 {
					triggered = true
				}
			case "massive_acc":
				triggered = candleByDate[points[idx].Date]["ORANGE"]
			case "markup":
				triggered = candleByDate[points[idx].Date]["YELLOW"]
			case "sudden_sell":
				triggered = candleByDate[points[idx].Date]["PURPLE"]
			}
			if triggered {
				events++
				cNow := closes[idx]
				if cNow > 0 {
					for _, h := range horizons {
						if idx+h < n {
							key := strconv.Itoa(h)
							cur := res[key]
							cur.N++
							ret := (closes[idx+h] - cNow) / cNow * 100.0
							if ret > 0 {
								cur.Up++
							}
							res[key] = cur
						}
					}
				}
			}
		}
		for _, h := range horizons {
			key := strconv.Itoa(h)
			cur := res[key]
			if cur.N > 0 {
				cur.Pct = shared.PyRound(float64(cur.Up)/float64(cur.N)*100.0, 0)
			}
			res[key] = cur
		}
		return SignalTypeStats{Events: events, Horizons: res}
	}

	calcAkumulasi := func(days int) AkumulasiStats {
		w := days
		if n < w {
			w = n
		}
		if w < 1 {
			return AkumulasiStats{AvgPrice: nil, TotalFlowM: 0.0, Status: "DISTRIBUTION"}
		}
		var buyShares, buyMoney, sellShares, sellMoney, totalM float64
		for _, p := range points[n-w:] {
			c := floatPtrOr(p.Close, 0.0)
			o := nonZeroOr(p.Open, c)
			h := nonZeroOr(p.High, c)
			l := nonZeroOr(p.Low, c)
			avgP := (o + h + l + c) / 4.0
			shares := 0.0
			if f, ok := rawFloat(p.RawSeries, 6); ok {
				shares = f
			}
			if shares > 0 {
				buyShares += shares
				buyMoney += shares * avgP
			} else if shares < 0 {
				absS := math.Abs(shares)
				sellShares += absS
				sellMoney += absS * avgP
			}
			totalM += p.ForeignNetVal
		}
		var avgPrice *int64
		status := "DISTRIBUTION"
		if totalM >= 0 {
			status = "ACCUMULATION"
			if buyShares > 0 {
				v := shared.PyRoundInt(buyMoney / buyShares)
				avgPrice = &v
			}
		} else if sellShares > 0 {
			v := shared.PyRoundInt(sellMoney / sellShares)
			avgPrice = &v
		}
		return AkumulasiStats{
			AvgPrice:   avgPrice,
			TotalFlowM: shared.PyRound(totalM, 2),
			Status:     status,
		}
	}

	return Statistics{
		Correlations: CorrelationStats{
			FF:        calcIndicatorCorr(cumFF),
			FCo:       calcIndicatorCorr(cumFCo),
			KF:        calcIndicatorCorr(cumKF),
			XL:        calcIndicatorCorr(cumXL),
			SemiRitel: calcIndicatorCorr(cumSR),
		},
		ConsecutiveInflow: ConsecutiveStats{
			D50:  calcAvgConsecutive(50),
			D100: calcAvgConsecutive(100),
			D200: calcAvgConsecutive(200),
		},
		InflowHargaNaik: InflowNaikStats{
			D50:  calcInflowHargaNaik(50),
			D100: calcInflowHargaNaik(100),
			D200: calcInflowHargaNaik(200),
		},
		MarkDown: MarkdownStats{
			D50:  calcMarkdown(50),
			D100: calcMarkdown(100),
			D200: calcMarkdown(200),
		},
		SignalStats: SignalStats{
			StreakInflow: calcSignalType("streak_inflow"),
			MassiveAcc:   calcSignalType("massive_acc"),
			Markup:       calcSignalType("markup"),
			SuddenSell:   calcSignalType("sudden_sell"),
		},
		AkumulasiDefault: calcAkumulasi(20),
	}
}

func nonZeroOr(p *float64, fallback float64) float64 {
	if p == nil || *p == 0 {
		return fallback
	}
	return *p
}
