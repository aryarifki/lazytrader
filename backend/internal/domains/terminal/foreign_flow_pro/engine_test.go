package foreign_flow_pro

import (
	"fmt"
	"math"
	"testing"
)

func floatPtr(v float64) *float64 { return &v }

func TestGetFPColorAndBadge(t *testing.T) {
	hfp := GetFPColorAndBadge(42.56)
	if hfp.Badge != "HFP" || hfp.Color != "#FF6820" || hfp.Score != 42.6 {
		t.Fatalf("HFP mismatch: %+v", hfp)
	}
	mfp := GetFPColorAndBadge(20.0)
	if mfp.Badge != "MFP" || mfp.Color != "rgb(211,211,211)" {
		t.Fatalf("MFP mismatch: %+v", mfp)
	}
	lfp := GetFPColorAndBadge(0.0)
	if lfp.Badge != "LFP" || lfp.Color != "#FFB6C1" {
		t.Fatalf("LFP mismatch: %+v", lfp)
	}
}

func TestBBRIIndex10FPScoreNoFMADrift(t *testing.T) {
	// Reproduces BBRI 2026-06-23 window where ARM64 FMA could flip 50.65 -> 50.7 vs CPython 50.6.
	parts := []float64{0, 0, 0, 0, 0, 0, 55.16, 48.38, 72.9, 51.45, 61.74}
	pts := make([]DailyPoint, len(parts))
	for i, part := range parts {
		c := 2910.0
		v := 50000000.0
		pts[i] = DailyPoint{
			Date:             fmt.Sprintf("2026-06-%02d", 10+i),
			Open:             floatPtr(c),
			High:             floatPtr(c),
			Low:              floatPtr(c),
			Close:            floatPtr(c),
			Volume:           floatPtr(v),
			ForeignNetVal:    10.0,
			PartisipasiAsing: floatPtr(part),
			RawSeries: []any{
				nil, nil, nil, nil, nil,
				790.655, 0.0, nil, nil,
				part, 58.63, 20.45, 56.24, 0.0, 10.0, 0.0,
			},
		}
	}
	res := ProcessForeignFlowAnalytics(pts)
	if res.FPSeries[10].Score != 50.6 {
		t.Fatalf("FPSeries[10].Score = %v, want 50.6", res.FPSeries[10].Score)
	}
}

func TestRollingCorrelationZeroDenom(t *testing.T) {
	x := []float64{10, 10, 10, 10, 10}
	y := []float64{1, 2, 3, 4, 5}
	got := CalculateRollingCorrelation(x, y, 5)
	if got[4] == nil || *got[4] != 0.0 {
		t.Fatalf("expected 0.0 on zero variance, got %+v", got)
	}
}

func TestProcessForeignFlowAnalyticsSignalsAndAkumulasi(t *testing.T) {
	n := 45
	pts := make([]DailyPoint, 0, n)
	for i := 0; i < n; i++ {
		c := 1000.0 + float64(i)*10.0
		ff := 5.0
		kf := -1.0
		part := 60.0
		if i == 24 {
			ff = 80.0 // should trigger ORANGE
		}
		if i == 30 {
			c = 1400.0 // sharp 1d and 3d rise with ff < sff5 -> YELLOW candidate
			ff = 1.0
		}
		if i == 35 {
			ff = -30.0 // lowest in 40d and < -10 -> PURPLE
		}
		if i == 42 {
			kf = 50.0 // positive 20d kf sum and drop >5% -> RED
			c = 800.0
		}
		vol := 10000000.0
		pts = append(pts, DailyPoint{
			Date:             fmt.Sprintf("2026-10-%02d", (i%28)+1),
			Open:             floatPtr(c),
			High:             floatPtr(c + 10),
			Low:              floatPtr(c - 10),
			Close:            floatPtr(c),
			Volume:           floatPtr(vol),
			ForeignNetVal:    ff,
			PartisipasiAsing: floatPtr(part),
			RawSeries: []any{
				nil, nil, nil, nil, nil,
				15.0,   // [5] val_miliar > 3
				1000.0, // [6] net shares
				nil, nil,
				part, // [9] partisipasi
				65.0, // [10] kat
				20.0, // [11] kret
				30.0, // [12] kfh
				nil,
				ff,  // [14] fco
				kf,  // [15] kf
				nil, // 16
				nil, // 17
				nil, // 18
				nil, // 19
				-ff, // [20] xl
				0.0, // [21] sr
			},
		})
	}
	res := ProcessForeignFlowAnalytics(pts)
	if len(res.Dates) != n || len(res.FPSeries) != n {
		t.Fatalf("length mismatch: dates=%d fp=%d", len(res.Dates), len(res.FPSeries))
	}
	if len(res.SpecialCandles) == 0 {
		t.Fatalf("expected special candles detected")
	}
	seenTypes := map[string]bool{}
	for _, sc := range res.SpecialCandles {
		seenTypes[sc.Type] = true
	}
	if !seenTypes["ORANGE"] || !seenTypes["PURPLE"] || !seenTypes["RED"] {
		t.Fatalf("expected ORANGE, PURPLE, RED candles, got %+v", seenTypes)
	}
	if res.Statistics.AkumulasiDefault.Status != "ACCUMULATION" || res.Statistics.AkumulasiDefault.AvgPrice == nil {
		t.Fatalf("unexpected akumulasi: %+v", res.Statistics.AkumulasiDefault)
	}
	if math.IsNaN(res.PriceChange) {
		t.Fatalf("unexpected NaN price change")
	}
}
