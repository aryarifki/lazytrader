package mmd

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/domains/terminal/shared"
)

type prevResultContext struct {
	wtdBadge   Badge
	w1Badge    Badge
	w2Badge    Badge
	mtdBadge   Badge
	m1Badge    Badge
	mtdCandles int
}

func parseCandleDate(iso string) time.Time {
	t, err := time.Parse("2006-01-02", iso)
	if err != nil {
		return time.Time{}
	}
	return t.UTC()
}

func filterDateRange(candles []DailyCandle, start, end time.Time) []DailyCandle {
	sISO := start.Format("2006-01-02")
	eISO := end.Format("2006-01-02")
	out := make([]DailyCandle, 0, len(candles))
	for _, c := range candles {
		if c.DateISO >= sISO && c.DateISO <= eISO {
			out = append(out, c)
		}
	}
	return out
}

func computePreviousResults(sortedAll []DailyCandle) (PreviousResults, prevResultContext) {
	lastDate := parseCandleDate(sortedAll[len(sortedAll)-1].DateISO)
	pyWeekday := (int(lastDate.Weekday()) + 6) % 7 // Monday=0 .. Sunday=6
	thisMonday := lastDate.AddDate(0, 0, -pyWeekday)

	wtdCandles := filterDateRange(sortedAll, thisMonday, lastDate)
	wtdRes := computeSliceResume(wtdCandles)
	wtdRange := ""
	if len(wtdCandles) > 0 {
		wtdRange = fmt.Sprintf("%s ~ %s", wtdCandles[0].DateISO, wtdCandles[len(wtdCandles)-1].DateISO)
	}

	var w1Start, w1End time.Time
	if pyWeekday == 4 { // Friday
		w1Start = thisMonday
		w1End = thisMonday.AddDate(0, 0, 4)
	} else {
		w1End = thisMonday.AddDate(0, 0, -3)
		w1Start = thisMonday.AddDate(0, 0, -7)
	}
	w1Candles := filterDateRange(sortedAll, w1Start, w1End)
	w1Res := computeSliceResume(w1Candles)

	w2Start := w1Start.AddDate(0, 0, -7)
	w2End := w1End.AddDate(0, 0, -7)
	w2Candles := filterDateRange(sortedAll, w2Start, w2End)
	w2Res := computeSliceResume(w2Candles)

	mtdStart := time.Date(lastDate.Year(), lastDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	mtdCandles := filterDateRange(sortedAll, mtdStart, lastDate)
	mtdRes := computeSliceResume(mtdCandles)

	m1End := mtdStart.AddDate(0, 0, -1)
	m1Start := time.Date(m1End.Year(), m1End.Month(), 1, 0, 0, 0, 0, time.UTC)
	m1Candles := filterDateRange(sortedAll, m1Start, m1End)
	m1Res := computeSliceResume(m1Candles)

	m2End := m1Start.AddDate(0, 0, -1)
	m2Start := time.Date(m2End.Year(), m2End.Month(), 1, 0, 0, 0, 0, time.UTC)
	m2Candles := filterDateRange(sortedAll, m2Start, m2End)
	m2Res := computeSliceResume(m2Candles)

	prev := PreviousResults{
		Weekly: []PreviousResultRow{
			{Label: "WTD", Date: "Minggu Ini", Badge: wtdRes.Action, ChangePct: wtdRes.PriceChange, PengendaliAction: wtdRes.PengendaliAction, RetailAction: wtdRes.RetailAction, Range: wtdRange},
			{Label: "1W", Date: "1 Minggu Lalu", Badge: w1Res.Action, ChangePct: w1Res.PriceChange, PengendaliAction: w1Res.PengendaliAction, RetailAction: w1Res.RetailAction, Range: fmt.Sprintf("%s ~ %s", w1Start.Format("2006-01-02"), w1End.Format("2006-01-02"))},
			{Label: "2W", Date: "2 Minggu Lalu", Badge: w2Res.Action, ChangePct: w2Res.PriceChange, PengendaliAction: w2Res.PengendaliAction, RetailAction: w2Res.RetailAction, Range: fmt.Sprintf("%s ~ %s", w2Start.Format("2006-01-02"), w2End.Format("2006-01-02"))},
		},
		Monthly: []PreviousResultRow{
			{Label: "MTD", Date: "Bulan Ini", Badge: mtdRes.Action, ChangePct: mtdRes.PriceChange, PengendaliAction: mtdRes.PengendaliAction, RetailAction: mtdRes.RetailAction, Range: fmt.Sprintf("%s ~ %s", mtdStart.Format("2006-01-02"), lastDate.Format("2006-01-02"))},
			{Label: "1M", Date: "Bulan Lalu", Badge: m1Res.Action, ChangePct: m1Res.PriceChange, PengendaliAction: m1Res.PengendaliAction, RetailAction: m1Res.RetailAction, Range: fmt.Sprintf("%s ~ %s", m1Start.Format("2006-01-02"), m1End.Format("2006-01-02"))},
			{Label: "2M", Date: "2 Bulan Lalu", Badge: m2Res.Action, ChangePct: m2Res.PriceChange, PengendaliAction: m2Res.PengendaliAction, RetailAction: m2Res.RetailAction, Range: fmt.Sprintf("%s ~ %s", m2Start.Format("2006-01-02"), m2End.Format("2006-01-02"))},
		},
	}
	ctx := prevResultContext{
		wtdBadge: wtdRes.Action, w1Badge: w1Res.Action, w2Badge: w2Res.Action,
		mtdBadge: mtdRes.Action, m1Badge: m1Res.Action, mtdCandles: len(mtdCandles),
	}
	return prev, ctx
}

func computeAIScore(
	ticker string, sortedAll []DailyCandle, priceChange, avgDailyVal float64,
	actionCls Badge, groupPerf *GroupPerformance, ctx prevResultContext,
) (int, float64, AIScoreComponents, []sliceResume) {
	n := len(sortedAll)
	tailCount := minInt(n, 30)
	dailyResumes := make([]sliceResume, n)
	for i := n - tailCount; i < n; i++ {
		dailyResumes[i] = computeSliceResume(sortedAll[i : i+1])
	}

	// C1: 20-day transaction value
	start20 := maxInt(0, n-20)
	var sumVal20 float64
	for _, c := range sortedAll[start20:] {
		var vol int64
		if c.Volume != nil {
			vol = *c.Volume
		}
		sumVal20 += float64(vol) * candleTavgOrClose(c)
	}
	avgVal20 := avgDailyVal
	if n-start20 > 0 {
		avgVal20 = sumVal20 / float64(n-start20)
	}
	miliar20 := avgVal20 / 1e9
	c1 := 0
	switch {
	case miliar20 < 0.7:
		c1 = -50
	case miliar20 < 3.0:
		c1 = -25
	case miliar20 < 7.0:
		c1 = -10
	case miliar20 > 150.0:
		c1 = 20
	}

	// C2: 10-day price change clamped [-15, 15]
	c2 := math.Min(15.0, math.Max(-15.0, priceChange))
	if n >= 11 {
		cNow := closeVal(sortedAll[n-1])
		c10 := closeVal(sortedAll[n-11])
		pct10 := 0.0
		if c10 > 0 {
			pct10 = ((cNow - c10) / c10) * 100.0
		}
		c2 = math.Min(15.0, math.Max(-15.0, pct10))
	}

	// C3: Weekly resume combination
	c3 := resumeWeekWTD1WWeight[ctx.wtdBadge.Cls] +
		resumeWeekWTD1WWeight[ctx.w1Badge.Cls] +
		resumeWeek2WWeight[ctx.w2Badge.Cls]

	// C4: Peer group average return clamped [-20, 20]
	c4 := computeC4GroupReturn(ticker, groupPerf)

	// C5: 5-day daily resume sum
	c5 := 0
	for i := maxInt(0, n-5); i < n; i++ {
		c5 += resumeDailyWeight[dailyResumes[i].Action.Cls]
	}

	// C6: 5-day period resume minus 5-day price change
	res5 := computeSliceResume(sortedAll[maxInt(0, n-5):])
	c6 := float64(resumePeriodWeight[res5.Action.Cls]) - res5.PriceChange

	// C7: 2-day price change drop penalty
	c7 := 0.0
	if n >= 3 {
		cNow := closeVal(sortedAll[n-1])
		c2d := closeVal(sortedAll[n-3])
		pct2 := 0.0
		if c2d > 0 {
			pct2 = ((cNow - c2d) / c2d) * 100.0
		}
		if pct2 < -15.0 {
			c7 = -20.0
		} else if pct2 <= -5.0 {
			c7 = pct2
		}
	}

	// C8: Candle direction streaks (15, 10, 5 days)
	validC := make([]DailyCandle, 0, n)
	for _, c := range sortedAll {
		if c.Close != nil && *c.Close != 0 {
			validC = append(validC, c)
		}
	}
	c8 := 0
	for _, rule := range [][3]int{{15, 11, 10}, {10, 7, 6}, {5, 4, 4}} {
		nDays, th, pts := rule[0], rule[1], rule[2]
		if len(validC) >= nDays+1 {
			sl := validC[len(validC)-nDays:]
			ups, downs := 0, 0
			for idx := 1; idx < len(sl); idx++ {
				if closeVal(sl[idx]) > closeVal(sl[idx-1]) {
					ups++
				} else if closeVal(sl[idx]) < closeVal(sl[idx-1]) {
					downs++
				}
			}
			if ups >= th {
				c8 += pts
			} else if downs >= th {
				c8 -= pts
			}
		}
	}

	// C9: Monthly resume (MTD requires >= 12 candles + 1M)
	currMonthScore := 0
	if ctx.mtdCandles >= 12 {
		currMonthScore = resumeMonthCurrWeight[ctx.mtdBadge.Cls]
	}
	c9 := currMonthScore + resumeMonthPrevWeight[ctx.m1Badge.Cls]

	// C10: 10-day high volatility whipsaw penalty
	c10 := 0
	if len(validC) >= 11 {
		sl10 := validC[len(validC)-11:]
		countNaik, countTurun := 0, 0
		for i := 1; i < len(sl10); i++ {
			currP := closeVal(sl10[i])
			prevP := closeVal(sl10[i-1])
			pctD := 0.0
			if prevP > 0 {
				pctD = ((currP - prevP) / prevP) * 100.0
			}
			if pctD > 5.0 {
				countNaik++
			} else if pctD < -5.0 {
				countTurun++
			}
		}
		if (countNaik+countTurun) >= 4 && countNaik >= 1 && countTurun >= 1 {
			c10 = -20
		}
	}

	// 20-Day Correlation Indicator
	var corrSum float64
	corrCount := 0
	for idx := start20; idx < n; idx++ {
		if idx > 0 {
			pPrev := closeVal(sortedAll[idx-1])
			pCurr := closeVal(sortedAll[idx])
			pctDay := 0.0
			if pPrev > 0 {
				pctDay = ((pCurr - pPrev) / pPrev) * 100.0
			}
			corrSum += AIScoreCorrDayPoint(dailyResumes[idx].Action.Cls, pctDay)
			corrCount++
		}
	}
	var korelasiRaw float64
	if corrCount > 0 {
		korelasiRaw = shared.PyRound(corrSum/10.0, 1)
	} else if strings.Contains(actionCls.Cls, "acc") {
		korelasiRaw = 6.8
	} else if strings.Contains(actionCls.Cls, "dist") {
		korelasiRaw = -4.5
	} else {
		korelasiRaw = 1.2
	}
	cCorr := shared.PyRoundInt(float64(korelasiRaw * AIScoreCorrMult))

	// Mark Down / Mark Up
	cMark := 0
	if n >= 3 {
		currP := closeVal(sortedAll[n-1])
		c1P := closeVal(sortedAll[n-2])
		c2P := closeVal(sortedAll[n-3])
		roc1, roc2 := 0.0, 0.0
		if c1P > 0 {
			roc1 = ((currP - c1P) / c1P) * 100.0
		}
		if c2P > 0 {
			roc2 = ((currP - c2P) / c2P) * 100.0
		}
		cls2B := computeSliceResume(sortedAll[n-2:])
		cls1B := dailyResumes[n-1]
		if (roc2 < -5.0 && isMarkDownOK(cls2B.Action.Cls)) || (roc1 < -4.0 && isMarkDownOK(cls1B.Action.Cls)) {
			cMark = 10
		} else if (roc2 > 5.0 && isMarkUpOK(cls2B.Action.Cls)) || (roc1 > 4.0 && isMarkUpOK(cls1B.Action.Cls)) {
			cMark = -10
		}
	}

	// 30-day daily resume sum
	cDaily30 := 0.0
	for i := maxInt(0, n-30); i < n; i++ {
		cDaily30 += resumeDaily30Weight[dailyResumes[i].Action.Cls]
	}

	// Sideways Acc / Dist
	accBest, distWorst := 0, 0
	for _, sp := range sidewaysPeriods {
		if n >= sp.days {
			sb := computeSliceResume(sortedAll[n-sp.days:])
			if math.Abs(sb.PriceChange) <= (sp.band + 1e-9) {
				if pts := sp.acc[sb.Action.Cls]; pts > accBest {
					accBest = pts
				}
				if pts := sp.dist[sb.Action.Cls]; pts < distWorst {
					distWorst = pts
				}
			}
		}
	}
	cSideways := accBest + distWorst

	totalRaw := float64(c1) + c2 + float64(c3) + c4 + float64(c5) + c6 + c7 +
		float64(c8) + float64(c9) + float64(c10) + float64(cCorr) + float64(cMark) + cDaily30 + float64(cSideways)
	totalScaled := float64(totalRaw * AIScoreTotalScale)

	var scoreVal float64
	switch {
	case totalScaled > 300:
		scoreVal = 99.0
	case totalScaled > 161.6:
		scoreVal = 95.0 + float64(float64(totalScaled-161.6)/float64(300.0-161.6))*4.0
	case totalScaled > 0:
		scoreVal = 47.5 + float64(totalScaled*0.294)
	case totalScaled > -250:
		scoreVal = float64(float64(totalScaled+250.0)/500.0) * 95.0
	default:
		scoreVal = 0.0
	}

	aiScore := int(shared.PyRoundInt(scoreVal))
	if aiScore < AIScoreMin {
		aiScore = AIScoreMin
	}
	if aiScore > AIScoreMax {
		aiScore = AIScoreMax
	}

	comp := AIScoreComponents{
		C1: c1, C2: c2, C3: c3, C4: c4, C5: c5, C6: c6, C7: c7, C8: c8, C9: c9, C10: c10,
		CCorr: cCorr, CMark: cMark, CDaily30: cDaily30, CSideways: cSideways,
		TotalRaw: totalRaw, TotalScaled: totalScaled,
	}
	return aiScore, korelasiRaw, comp, dailyResumes
}

func computeC4GroupReturn(ticker string, gp *GroupPerformance) float64 {
	if gp == nil {
		return 0.0
	}
	for _, g := range []*GroupPerfBlock{gp.Conglomerate, gp.Sector} {
		if g == nil || len(g.Members) == 0 {
			continue
		}
		var sum float64
		count := 0
		for _, m := range g.Members {
			if m.Ticker == ticker {
				continue
			}
			if m.W1Pct != nil {
				sum += *m.W1Pct
				count++
			} else {
				sum += m.TodayPct
				count++
			}
		}
		if count > 0 {
			avgG := sum / float64(count)
			if avgG < -16.0 {
				return -20.0
			}
			if avgG > 16.0 {
				return 20.0
			}
			return avgG
		}
	}
	return 0.0
}

func closeVal(c DailyCandle) float64 {
	if c.Close != nil {
		return *c.Close
	}
	return 0.0
}

func isMarkDownOK(cls string) bool {
	return cls == "small-dist" || cls == "undetected" || cls == "small-acc" ||
		cls == "normal-acc" || cls == "big-acc" || cls == "massive-acc"
}

func isMarkUpOK(cls string) bool {
	return cls == "small-acc" || cls == "undetected" || cls == "small-dist" ||
		cls == "normal-dist" || cls == "big-dist" || cls == "massive-dist"
}
