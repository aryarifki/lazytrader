package forensic

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/aryarifki/lazytrader/backend/internal/core/flags"
)

const triggerOutflow = -5.0e12

func buildTL(flows, prices []float64) []TimelinePoint {
	tl := make([]TimelinePoint, len(flows))
	for i := range flows {
		tl[i] = TimelinePoint{
			Date:      fmt.Sprintf("2026-01-%02d", i+1),
			IHSGClose: prices[i],
			DailyFlow: flows[i],
		}
	}
	return tl
}

func prefixStates(flows, prices []float64, f *flags.Store) []string {
	tl := buildTL(flows, prices)
	out := make([]string, len(tl))
	for i := range tl {
		cp := make([]TimelinePoint, i+1)
		copy(cp, tl[:i+1])
		out[i] = EvaluateMarketState(cp, f)
	}
	return out
}

func TestExodusExitFix_UnitCases(t *testing.T) {
	t.Setenv("FLAG_EXODUS_EXIT_FIX", "1")
	t.Setenv("FLAG_EXODUS_RECOVERY_FIX", "1")
	t.Setenv("FLAG_EXODUS_THRESHOLD_MODE", "absolute")
	f := &flags.Store{}

	// 1. Latch CRASH exits after 8 quiet sessions
	flows := append([]float64{triggerOutflow}, make([]float64, 15)...)
	prices := make([]float64, 16)
	prices[0] = 1000
	for i := 1; i < 16; i++ {
		prices[i] = 1000 + float64(5*i)
	}
	st := prefixStates(flows, prices, f)
	if st[0] != "CRASH" || st[9] != "CRASH" || st[10] != "WARNING" {
		t.Fatalf("unexpected exit states: [0]=%s [9]=%s [10]=%s", st[0], st[9], st[10])
	}

	// 2. Printing new lows keeps CRASH locked
	flows2 := append([]float64{triggerOutflow}, make([]float64, 20)...)
	prices2 := make([]float64, 21)
	for i := range prices2 {
		prices2[i] = 1000 - float64(5*i)
	}
	st2 := prefixStates(flows2, prices2, f)
	for i, s := range st2 {
		if s != "CRASH" {
			t.Fatalf("expected CRASH while printing new lows at idx %d, got %s", i, s)
		}
	}

	// 3. Recovery failure (< 0.98 * episodeLow) re-enters CRASH
	flows3 := append([]float64{triggerOutflow}, make([]float64, 12)...)
	prices3 := make([]float64, 0, 13)
	prices3 = append(prices3, 1000)
	for i := 0; i < 10; i++ {
		prices3 = append(prices3, 1005+float64(5*i))
	}
	prices3 = append(prices3, 970, 968)
	st3 := prefixStates(flows3, prices3, f)
	if st3[10] != "WARNING" || st3[11] != "CRASH" {
		t.Fatalf("expected [10]=WARNING [11]=CRASH, got %s / %s", st3[10], st3[11])
	}

	// 4. Flag OFF keeps old latch
	t.Setenv("FLAG_EXODUS_EXIT_FIX", "0")
	prices4 := make([]float64, len(flows2))
	prices4[0] = 1000
	for i := 1; i < len(flows2); i++ {
		prices4[i] = 1000 + float64(5*i)
	}
	st4 := prefixStates(flows2, prices4, f)
	for i, s := range st4 {
		if s != "CRASH" {
			t.Fatalf("flag=0 expected permanent CRASH at idx %d, got %s", i, s)
		}
	}
}

func TestCanarySnapshotsParity_BothModes(t *testing.T) {
	tlBytes, err := os.ReadFile("/root/SMtracker/tools/canary/snapshots/timeline.json")
	if err != nil {
		t.Skipf("canary snapshots unavailable: %v", err)
	}
	var timeline []TimelinePoint
	if err := json.Unmarshal(tlBytes, &timeline); err != nil {
		t.Fatalf("parse timeline.json: %v", err)
	}
	if len(timeline) < 1000 {
		t.Fatalf("unexpectedly short canary timeline: %d", len(timeline))
	}

	t.Setenv("FLAG_EXODUS_EXIT_FIX", "1")
	t.Setenv("FLAG_EXODUS_RECOVERY_FIX", "1")
	f := &flags.Store{}

	for _, mode := range []string{"absolute", "hybrid_p5_p15"} {
		t.Setenv("FLAG_EXODUS_THRESHOLD_MODE", mode)
		expBytes, err := os.ReadFile(fmt.Sprintf("/root/SMtracker/tools/canary/snapshots/states_%s.json", mode))
		if err != nil {
			t.Fatalf("read states_%s.json: %v", mode, err)
		}
		var expected map[string]string
		if err := json.Unmarshal(expBytes, &expected); err != nil {
			t.Fatalf("parse states_%s.json: %v", mode, err)
		}

		mismatches := 0
		for i := range timeline {
			cp := make([]TimelinePoint, i+1)
			copy(cp, timeline[:i+1])
			got := EvaluateMarketState(cp, f)
			want := expected[timeline[i].Date]
			if got != want {
				if mismatches < 5 {
					t.Errorf("mode=%s date=%s idx=%d: got %s, want %s", mode, timeline[i].Date, i, got, want)
				}
				mismatches++
			}
		}
		if mismatches > 0 {
			t.Fatalf("mode=%s had %d/%d state mismatches", mode, mismatches, len(timeline))
		}
		t.Logf("Canary mode=%s: 0/%d mismatches (100%% parity)", mode, len(timeline))
	}
}
