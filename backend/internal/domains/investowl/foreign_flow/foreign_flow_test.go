package foreign_flow

import (
	"math"
	"testing"
)

func strPtr(s string) *string { return &s }

func TestForeignForensicInvariants(t *testing.T) {
	if got := classifyCounterparty("Asing", strPtr("retail")); got != "retail" {
		t.Fatalf("expected PMA retail -> retail, got %s", got)
	}
	if got := classifyCounterparty("Asing", strPtr("semi_ritel")); got != "retail" {
		t.Fatalf("expected semi_ritel -> retail, got %s", got)
	}
	if got := classifyCounterparty("Asing", strPtr("institusi")); got != "foreign" {
		t.Fatalf("expected foreign, got %s", got)
	}
	if got := classifyCounterparty("Lokal", strPtr("institusi")); got != "local_inst" {
		t.Fatalf("expected local_inst, got %s", got)
	}

	verdict, tone := evaluateCounterpartyVerdict(10e9, 5e9, -15e9)
	if verdict != "Big Money Kompak Akumulasi ← Ritel" || tone != "positive" {
		t.Fatalf("unexpected verdict: %s / %s", verdict, tone)
	}

	state, _, _, _ := evaluateDivergence(5.0, 10_000_000.0, 20_000_000_000.0)
	if state != "NEUTRAL" {
		t.Fatalf("expected deadband NEUTRAL, got %s", state)
	}

	hhi := ComputeForeignHHI([]float64{50, 50})
	if math.Abs(hhi-0.5) > 1e-9 {
		t.Fatalf("HHI = %v, want 0.5", hhi)
	}

	hm := ComputeBrokerHeatmap([]MicroRow{
		{Date: "2026-11-01", BrokerCode: "AK", NetValue: 10},
		{Date: "2026-11-02", BrokerCode: "AK", NetValue: 20},
		{Date: "2026-11-03", BrokerCode: "AK", NetValue: 30},
	}, 15)
	// pandas ddof=1 std([10,20,30]) = 10 -> z = [-1, 0, 1]
	if len(hm.MatrixData) != 3 || hm.MatrixData[0][2] != -1.0 || hm.MatrixData[2][2] != 1.0 {
		t.Fatalf("unexpected heatmap matrix: %+v", hm.MatrixData)
	}
}
