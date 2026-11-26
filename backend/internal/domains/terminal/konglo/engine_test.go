package konglo

import (
	"testing"
	"time"
)

func fp(v float64) *float64 { return &v }
func ip(v int64) *int64     { return &v }
func sp(v string) *string   { return &v }

func TestCalculateKongloReportAndMetadata(t *testing.T) {
	stocks := []StockRow{
		{
			Ticker: "BBCA", Name: sp("Bank Central Asia"), GroupName: sp("Djarum Group"),
			Prev: fp(5800), Last: fp(5825), Val: fp(955e9), Vol: ip(100000), Freq: ip(1000),
			SharesListed: fp(120e9),
		},
		{
			Ticker: "BREN", Name: sp("Barito Renewables"), GroupName: sp("Barito Group"),
			Prev: fp(2000), Last: fp(1900), Val: fp(50e9), Vol: ip(50000), Freq: ip(500),
			SharesListed: fp(100e9),
		},
		{
			Ticker: "DEAD", Name: sp("Inactive"), GroupName: sp("Djarum Group"),
			Prev: fp(100), Last: fp(0), SharesListed: fp(1e9),
		},
	}
	groupOrder := []string{"Barito Group", "Djarum Group"}
	groupsMap := map[string][]string{
		"Barito Group": {"BREN", "MISS"},
		"Djarum Group": {"BBCA", "BREN"}, // BREN shared -> deduplicated in overview/top5/top15
	}

	rep := CalculateKongloReport(stocks, groupOrder, groupsMap)
	if rep.Summary.TotalStocks != 3 || rep.Summary.TotalActiveStocks != 2 {
		t.Fatalf("unexpected summary: %+v", rep.Summary)
	}
	if rep.KongloCards.Overview.TotalKongloStocks != 2 || rep.KongloCards.Overview.MissingCount != 1 {
		t.Fatalf("unexpected overview deduplication: %+v", rep.KongloCards.Overview)
	}
	if rep.KongloCards.BestMC.Group != "Djarum Group" {
		t.Fatalf("expected BestMC Djarum Group, got %+v", rep.KongloCards.BestMC)
	}
	if len(rep.FundamentalReport) != 1 || rep.FundamentalReport[0].Ticker != "BBCA" {
		t.Fatalf("unexpected fundamental report: %+v", rep.FundamentalReport)
	}

	td := time.Date(2026, 11, 25, 0, 0, 0, 0, time.UTC)
	mu := time.Date(2026, 11, 25, 9, 0, 0, 0, time.UTC)
	meta := buildReportMetadata(&td, &mu, rep.Summary.TotalActiveStocks, rep.Summary.TotalStocks)
	if meta.DateDisplay != "Rabu, 25 Nov 2026" || meta.LastUpdatedStr != "Rab, 25/11/2026 16:00 WIB" {
		t.Fatalf("unexpected Indonesian metadata: %+v", meta)
	}
	if meta.Timestamp != "2026-11-25T09:00:00+00:00" {
		t.Fatalf("unexpected metadata timestamp: %q", meta.Timestamp)
	}
}
