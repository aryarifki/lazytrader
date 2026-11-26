package mmd

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"
)

func makeTestCandle(d time.Time, prev, closePx float64, buys, sells map[string]float64, vol int64) DailyCandle {
	hi := math.Max(prev, closePx)
	lo := math.Min(prev, closePx)
	return DailyCandle{
		DateISO:     d.Format("2006-01-02"),
		Prev:        &prev,
		Open:        &prev,
		High:        &hi,
		Low:         &lo,
		Close:       &closePx,
		Volume:      &vol,
		Tavg:        &closePx,
		BuyBrokers:  buys,
		SellBrokers: sells,
	}
}

func TestPyFloorDivAndFormatting(t *testing.T) {
	if got := pyFloorDiv(-12266888, 3); got != -4088963 {
		t.Fatalf("pyFloorDiv(-12266888, 3) = %d, want -4088963", got)
	}
	if got := pyFloorDiv(10, 3); got != 3 {
		t.Fatalf("pyFloorDiv(10, 3) = %d, want 3", got)
	}
	if got := formatCommaInt(4454428); got != "4,454,428" {
		t.Fatalf("formatCommaInt = %q", got)
	}
	if got := formatSignedCommaInt(-1702919); got != "-1,702,919" {
		t.Fatalf("formatSignedCommaInt(-) = %q", got)
	}
	if got := formatSignedCommaInt(1250); got != "+1,250" {
		t.Fatalf("formatSignedCommaInt(+) = %q", got)
	}
	if got := formatSignedCommaInt(0); got != "0" {
		t.Fatalf("formatSignedCommaInt(0) = %q", got)
	}
	if got := formatCommaFloat(98674.12, 1); got != "98,674.1" {
		t.Fatalf("formatCommaFloat = %q", got)
	}
}

func TestSortedJSONBKeysAndCommentaryBrokerBaru(t *testing.T) {
	m := map[string]float64{"CD": 100, "AK": 50, "ABC": 10, "BK": 25}
	keys := sortedJSONBKeys(m)
	want := []string{"AK", "BK", "CD", "ABC"}
	if len(keys) != len(want) {
		t.Fatalf("len=%d", len(keys))
	}
	for i := range want {
		if keys[i] != want[i] {
			t.Fatalf("keys[%d]=%q want %q", i, keys[i], want[i])
		}
	}

	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	candles := make([]DailyCandle, 0, 21)
	for i := 0; i < 20; i++ {
		candles = append(candles, makeTestCandle(start.AddDate(0, 0, i), 1000, 1000, map[string]float64{"AK": 100}, map[string]float64{"YP": 100}, 1000))
	}
	candles = append(candles, makeTestCandle(start.AddDate(0, 0, 20), 1000, 1010, map[string]float64{"CD": 500, "ZP": 200}, map[string]float64{"YP": 700}, 1000))
	msg := commentaryBrokerBaru(candles)
	if !strings.Contains(msg, "broker baru CD ") {
		t.Fatalf("unexpected commentaryBrokerBaru: %q", msg)
	}
}

func TestCommentaryKorelasiParuhBranchesReachable(t *testing.T) {
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	accBuy := map[string]float64{"AK": 8000}
	accSell := map[string]float64{"YP": 1000, "PD": 1000, "XC": 1000, "NI": 1000, "KK": 1000, "XL": 1000, "EP": 1000, "CC": 1000}

	candlesWorsen := make([]DailyCandle, 0, 20)
	resumesWorsen := make([]sliceResume, 0, 20)
	px := 1000.0
	for i := 0; i < 20; i++ {
		d := start.AddDate(0, 0, i)
		mult := 0.97
		if i < 10 {
			mult = 1.04
		}
		nxt := px * mult
		c := makeTestCandle(d, px, nxt, accBuy, accSell, 10000)
		candlesWorsen = append(candlesWorsen, c)
		resumesWorsen = append(resumesWorsen, computeSliceResume([]DailyCandle{c}))
		px = nxt
	}
	msgWorsen := commentaryKorelasiParuh(candlesWorsen, resumesWorsen)
	if !strings.Contains(msgWorsen, "mulai memburuk") {
		t.Fatalf("expected memburuk, got %q", msgWorsen)
	}

	candlesImprove := make([]DailyCandle, 0, 20)
	resumesImprove := make([]sliceResume, 0, 20)
	px = 1000.0
	for i := 0; i < 20; i++ {
		d := start.AddDate(0, 0, i)
		mult := 1.04
		if i < 10 {
			mult = 0.97
		}
		nxt := px * mult
		c := makeTestCandle(d, px, nxt, accBuy, accSell, 10000)
		candlesImprove = append(candlesImprove, c)
		resumesImprove = append(resumesImprove, computeSliceResume([]DailyCandle{c}))
		px = nxt
	}
	msgImprove := commentaryKorelasiParuh(candlesImprove, resumesImprove)
	if !strings.Contains(msgImprove, "mulai menguat") {
		t.Fatalf("expected menguat, got %q", msgImprove)
	}
}

func TestC4PrefersW1PctAndFallbacksToTodayPct(t *testing.T) {
	d := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	c := makeTestCandle(d, 1000, 1010, map[string]float64{"AK": 5000}, map[string]float64{"YP": 5000}, 10000)
	w12, w8, w1 := 12.0, 8.0, 1.0

	gpW1 := &GroupPerformance{
		Conglomerate: &GroupPerfBlock{
			GroupName: "TEST", Type: "Konglomerasi",
			Members: []GroupPerfMember{
				{Ticker: "AAA", TodayPct: -5.0, W1Pct: &w12, M1Pct: 20.0},
				{Ticker: "BBB", TodayPct: -4.0, W1Pct: &w8, M1Pct: 15.0},
				{Ticker: "TGT", TodayPct: 1.0, W1Pct: &w1, M1Pct: 2.0},
			},
		},
	}
	resW1, ok := ComputeMMDAnalysis(AnalysisInput{
		Ticker: "TGT", Candles: []DailyCandle{c}, AllHistoryCandles: []DailyCandle{c}, GroupPerf: gpW1, AutoRetail: true, Today: d,
	})
	if !ok || resW1.AIScoreComponents.C4 != 10.0 {
		t.Fatalf("C4=%v want 10.0", resW1.AIScoreComponents.C4)
	}

	gpFallback := &GroupPerformance{
		Conglomerate: &GroupPerfBlock{
			GroupName: "TEST", Type: "Konglomerasi",
			Members: []GroupPerfMember{
				{Ticker: "AAA", TodayPct: -5.0, W1Pct: nil},
				{Ticker: "BBB", TodayPct: -3.0, W1Pct: nil},
				{Ticker: "TGT", TodayPct: 1.0, W1Pct: nil},
			},
		},
	}
	resFB, ok := ComputeMMDAnalysis(AnalysisInput{
		Ticker: "TGT", Candles: []DailyCandle{c}, AllHistoryCandles: []DailyCandle{c}, GroupPerf: gpFallback, AutoRetail: true, Today: d,
	})
	if !ok || resFB.AIScoreComponents.C4 != -4.0 {
		t.Fatalf("C4=%v want -4.0", resFB.AIScoreComponents.C4)
	}
}

func TestRetailRowsIncludesGRAndTwoSidedTurnover(t *testing.T) {
	d := time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC)
	retailBuys := make(map[string]float64, 9)
	for _, code := range CoreRetailList {
		retailBuys[code] = 1000
	}
	c := makeTestCandle(d, 1000, 1000, retailBuys, map[string]float64{"AK": 9000}, 9000)
	res, ok := ComputeMMDAnalysis(AnalysisInput{
		Ticker: "TEST", Candles: []DailyCandle{c}, AllHistoryCandles: []DailyCandle{c}, AutoRetail: false, Today: d,
	})
	if !ok || len(res.RetailRows) != 9 {
		t.Fatalf("retail_rows len=%d want 9", len(res.RetailRows))
	}
	hasGR := false
	sumPct := 0.0
	for _, rr := range res.RetailRows {
		if rr.Broker == "GR" {
			hasGR = true
		}
		sumPct += rr.PctTrans
	}
	if !hasGR {
		t.Fatal("GR missing from retail_rows")
	}
	if sumPct > 100.1 {
		t.Fatalf("sum pct_trans=%.2f exceeds 100.1", sumPct)
	}
}

func TestPengendaliDirectionalNetAndRetailCommentary(t *testing.T) {
	d := time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)
	c := makeTestCandle(d, 1000, 960, map[string]float64{"YP": 4000, "PD": 3000, "AK": 1000}, map[string]float64{"AK": 6000, "BK": 2000}, 8000)
	sr := computeSliceResume([]DailyCandle{c})
	if !strings.Contains(sr.PengendaliAction.Cls, "dist") {
		t.Fatalf("expected dist pengendali, got %+v", sr.PengendaliAction)
	}
	res, ok := ComputeMMDAnalysis(AnalysisInput{
		Ticker: "TEST", Candles: []DailyCandle{c}, AllHistoryCandles: []DailyCandle{c}, AutoRetail: false, Today: d,
	})
	if !ok || !strings.Contains(res.PengendaliStats.Action.Cls, "dist") {
		t.Fatalf("expected dist action, got %+v", res.PengendaliStats.Action)
	}

	cSell := makeTestCandle(d, 1000, 990, map[string]float64{"AK": 5000, "BK": 5000}, map[string]float64{"YP": 6000, "PD": 4000}, 10000)
	resSell, ok := ComputeMMDAnalysis(AnalysisInput{
		Ticker: "TEST", Candles: []DailyCandle{cSell}, AllHistoryCandles: []DailyCandle{cSell}, AutoRetail: false, Today: d,
	})
	if !ok || resSell.PengendaliStats.Ritel.Main != 50.0 || resSell.PengendaliStats.Ritel.Sub >= 0 {
		t.Fatalf("unexpected pengendali ritel stats: %+v", resSell.PengendaliStats.Ritel)
	}
	for _, msg := range resSell.AICommentary {
		if strings.Contains(msg, "tidak banyak bertransaksi") {
			t.Fatalf("commentary_section3_point_5 fired falsely: %q", msg)
		}
	}
}

func TestCommentaryIPOAndRetailOverrideAndScoreClamp(t *testing.T) {
	today := time.Date(2026, 11, 26, 0, 0, 0, 0, time.UTC)
	d3m := "2026-08-26"
	d9m := "2026-02-26"
	d24m := "2024-11-01"
	if got := commentaryIPO("ABCD", &d3m, today); !strings.Contains(got, "6 bulan") {
		t.Fatalf("got %q", got)
	}
	if got := commentaryIPO("ABCD", &d9m, today); !strings.Contains(got, "12 bulan") {
		t.Fatalf("got %q", got)
	}
	if got := commentaryIPO("ABCD", &d24m, today); got != "-" {
		t.Fatalf("got %q", got)
	}

	c := makeTestCandle(today, 1000, 1000, map[string]float64{"YP": 5000}, map[string]float64{"AK": 5000}, 5000)
	res, ok := ComputeMMDAnalysis(AnalysisInput{
		Ticker: "ABCD", Candles: []DailyCandle{c}, AllHistoryCandles: []DailyCandle{c},
		AutoRetail: true, RetailOverride: "yp", Today: today,
	})
	if !ok {
		t.Fatal("expected ok")
	}
	if res.ActiveRetailBandar == nil || *res.ActiveRetailBandar != "YP" {
		t.Fatalf("ActiveRetailBandar = %v, want YP", res.ActiveRetailBandar)
	}
	if res.AIScore < AIScoreMin || res.AIScore > AIScoreMax {
		t.Fatalf("AIScore = %d out of bounds [%d,%d]", res.AIScore, AIScoreMin, AIScoreMax)
	}
	if got := fmt.Sprintf("%.4f", VolumeRatioWeight(7.0)); got != "0.8385" {
		t.Fatalf("VolumeRatioWeight(7.0) = %s", got)
	}
}
