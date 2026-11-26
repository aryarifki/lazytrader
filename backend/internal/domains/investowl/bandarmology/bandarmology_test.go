package bandarmology

import (
	"context"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/auth"
)

func f64Ptr(v float64) *float64 { return &v }

func TestGreedyBrokerDistributionAndSummary(t *testing.T) {
	d := time.Date(2026, 11, 20, 0, 0, 0, 0, time.UTC)
	acts := []ActivityRow{
		{Date: d, Ticker: "BBCA", BrokerCode: "AK", ParticipantType: "Asing", BuyValue: 100e9, NetValue: 100e9, BuyLot: f64Ptr(1e6), Frequency: f64Ptr(100)},
		{Date: d, Ticker: "BBCA", BrokerCode: "CC", ParticipantType: "Pemerintah", SellValue: 60e9, NetValue: -60e9, SellLot: f64Ptr(6e5), Frequency: f64Ptr(80)},
		{Date: d, Ticker: "BBCA", BrokerCode: "YP", ParticipantType: "Lokal", SellValue: 40e9, NetValue: -40e9, SellLot: f64Ptr(4e5), Frequency: f64Ptr(500)},
	}
	dist := BrokerDistributionDataRange(acts, d, d)
	edges := dist["edges"].([]map[string]any)
	if len(edges) != 2 {
		t.Fatalf("expected 2 greedy edges, got %d", len(edges))
	}
	if edges[0]["buyer_code"] != "AK" || edges[0]["seller_code"] != "CC" || edges[0]["matched_value"].(float64) != 60e9 {
		t.Fatalf("unexpected first edge: %+v", edges[0])
	}
	sum := BrokerSummaryTable(dist)
	if len(sum) != 2 || sum[0]["buy_broker"] != "AK" || sum[0]["sell_broker"] != "CC" {
		t.Fatalf("unexpected summary table: %+v", sum)
	}
}

func TestAffiliationAndPricePerformance(t *testing.T) {
	aff := GetBrokerAffiliation(context.Background(), nil, "BULL", "II")
	if aff == nil || aff["badge"] != "AFILIASI PENGENDALI & DIREKSI" {
		t.Fatalf("expected BULL-II documented affiliation, got %+v", aff)
	}
	var prices []PriceRow
	base := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 10; i++ {
		c := 1000.0 + float64(i)*10.0
		prices = append(prices, PriceRow{Date: base.AddDate(0, 0, i), Ticker: "BBCA", Close: &c})
	}
	perf := PricePerformanceTable(prices)
	if len(perf) != 10 || perf[0]["timeframe"] != "1D" || perf[5]["timeframe"] != "YTD" {
		t.Fatalf("unexpected performance table layout: %+v", perf)
	}
	ret1D := perf[0]["return"].(*float64)
	want := 1090.0/1080.0 - 1.0
	if ret1D == nil || math.Abs(*ret1D-want) > 1e-9 {
		t.Fatalf("1D return mismatch: got %v want %v", ret1D, want)
	}
}

func TestPipelineFailClosedAuth(t *testing.T) {
	guard := &auth.Guard{Secret: "test-secret"}
	r := Routes(Deps{Guard: guard})
	for _, path := range []string{"/pipeline/run", "/pipeline/backfill"} {
		req := httptest.NewRequest(http.MethodPost, path, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("anonymous %s status = %d, want 401", path, rec.Code)
		}
	}
}
