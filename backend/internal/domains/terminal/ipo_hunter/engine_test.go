package ipohunter

import (
	"math"
	"testing"
	"time"
)

func dUTC(y, m, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}

func fPtr(v float64) *float64 { return &v }
func iPtr(v int64) *int64     { return &v }
func sPtr(v string) *string   { return &v }

func TestJSRoundSemantics(t *testing.T) {
	cases := []struct {
		in   float64
		want int64
	}{
		{0.5, 1},
		{2.5, 3},
		{3.5, 4},
		{-0.5, 0},
		{-2.5, -2},
		{-62.85, -63},
	}
	for _, tc := range cases {
		if got := JSRound(tc.in); got != tc.want {
			t.Fatalf("JSRound(%v) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestEffectiveRetailAndCategoryOverrides(t *testing.T) {
	effMNC := EffectiveRetailBrokers("BMTR", nil)
	for _, b := range effMNC {
		if b == "EP" {
			t.Fatal("MNC group BMTR should exclude EP from retail")
		}
	}
	if got := CategoryOf("EP", "BMTR", nil); got != "institusi" {
		t.Fatalf("CategoryOf(EP, BMTR) = %q, want institusi", got)
	}

	effSM := EffectiveRetailBrokers("DSSA", nil)
	for _, b := range effSM {
		if b == "DH" {
			t.Fatal("Sinarmas group DSSA should exclude DH from retail")
		}
	}
	if got := CategoryOf("DH", "DSSA", nil); got != "institusi" {
		t.Fatalf("CategoryOf(DH, DSSA) = %q, want institusi", got)
	}

	uwSet := map[string]struct{}{"PD": {}, "DH": {}, "LG": {}}
	effUW := EffectiveRetailBrokers("EMAS", uwSet)
	for _, b := range effUW {
		if b == "PD" || b == "DH" {
			t.Fatalf("UW broker %s leaked into EffectiveRetailBrokers", b)
		}
	}
	if got := CategoryOf("PD", "EMAS", uwSet); got != "uw" {
		t.Fatalf("CategoryOf(PD, EMAS, uw) = %q, want uw", got)
	}
}

func TestActionWeightAndRetailBandar(t *testing.T) {
	net := map[string]float64{"XL": 1000, "CC": -800, "YP": -200}
	vrShort, _ := ComputeActionWeight(net, 5000, 5000, 5)
	if math.Abs(vrShort-5.0) > 1e-9 {
		t.Fatalf("vrShort = %v, want 5.0", vrShort)
	}
	_, wLong := ComputeActionWeight(net, 10000, 10000, 20)
	if wLong < 0.25 || wLong > 1.0 {
		t.Fatalf("wLong = %v out of [0.25, 1.0]", wLong)
	}

	bt := map[string]BrokerTotal{
		"XL": {Buy: 3000, Sell: 1000},
		"CC": {Buy: 1000, Sell: 1000},
	}
	buyers := []BrokerNet{{Broker: "XL", Net: 2000}}
	sellers := []BrokerNet{{Broker: "YP", Net: -200}}
	best, reasons := DetermineRetailBandarOtomatis(
		map[string]float64{"XL": 2000, "YP": -200},
		bt,
		buyers,
		sellers,
		RetailBrokers,
		1100,
		6000,
	)
	if best == nil || *best != "XL" {
		t.Fatalf("best = %v, want XL", best)
	}
	if len(reasons["XL"]) == 0 {
		t.Fatalf("expected XL reasons to be populated")
	}
	if reasons["XL"][0] != "C1: Transaksi XL > 3.5x total retail lain" {
		t.Fatalf("unexpected C1 reason string: %q", reasons["XL"][0])
	}
}

func TestBuildAfterDaysAndGapWarning(t *testing.T) {
	listing := dUTC(2026, 7, 10)
	d0dt := dUTC(2026, 9, 30)
	master := &MasterRecord{Ticker: "RANS", ListingDate: &listing, IPOPrice: 170}
	d0Row := &DayRow{Date: d0dt, Ticker: "RANS"}
	later := []DayRow{
		{Date: dUTC(2026, 9, 29), Ticker: "RANS", Close: fPtr(150), Volume: iPtr(100)},
		{Date: dUTC(2026, 10, 8), Ticker: "RANS", Close: fPtr(160), Volume: iPtr(200)},
	}
	holidays := map[int]struct{}{
		DateOrdinal(dUTC(2026, 8, 17)): {},
	}
	hi := dUTC(2026, 10, 8)
	info := BuildAfterDays("RANS", master, d0Row, later, &hi, holidays)
	if !info.InD0 || info.NLater != 1 || len(info.Days) != 2 {
		t.Fatalf("unexpected AfterInfo: inD0=%v nLater=%d days=%d", info.InD0, info.NLater, len(info.Days))
	}
	if info.Gap == nil || info.Gap.From != "2026-10-01" || info.Gap.To != "2026-10-07" {
		t.Fatalf("unexpected gap range: %+v", info.Gap)
	}

	// hiDate <= d0Date leaves only D0 in Days
	hiEarly := dUTC(2026, 9, 25)
	infoEarly := BuildAfterDays("RANS", master, d0Row, later, &hiEarly, holidays)
	if infoEarly.NLater != 0 || len(infoEarly.Days) != 1 || !infoEarly.Days[0].IsD0 {
		t.Fatalf("expected only D0 row when hiDate <= d0Date")
	}
}

func TestAggregationZeroVolumeAndPieTopBrokers(t *testing.T) {
	days := []DayRow{
		{
			Date: dUTC(2026, 11, 10), Ticker: "TEST", IsD0: true, DayCount: 3,
			Open: fPtr(100), High: fPtr(110), Low: fPtr(90), Close: fPtr(100), Prev: fPtr(100), Tavg: fPtr(102),
			Volume: iPtr(1000), BuyBrokers: map[string]float64{"CC": 600.5, "XL": 399.5}, SellBrokers: map[string]float64{"YP": 1000},
		},
		{
			Date: dUTC(2026, 11, 11), Ticker: "TEST", DayCount: 1,
			Open: fPtr(100), High: fPtr(100), Low: fPtr(100), Close: fPtr(105), Prev: fPtr(100), Tavg: fPtr(0),
			Volume: iPtr(0), BuyBrokers: map[string]float64{}, SellBrokers: map[string]float64{},
		},
	}
	agg := BuildAggregation("TEST", days)
	if agg == nil {
		t.Fatal("expected non-nil agg")
	}
	// avg_daily_value divides by 2 (includes 0-volume day with tavg=0!)
	wantAvgVal := (1000.0 * 100.0 * 102.0) / 2.0
	if agg.AvgDailyValue == nil || math.Abs(*agg.AvgDailyValue-wantAvgVal) > 1e-9 {
		t.Fatalf("AvgDailyValue = %v, want %v", agg.AvgDailyValue, wantAvgVal)
	}
	// AfterVWAP skips 0-volume day and returns 102
	vw := AfterVWAP(days)
	if vw == nil || math.Abs(*vw-102.0) > 1e-9 {
		t.Fatalf("AfterVWAP = %v, want 102.0", vw)
	}
	// LastCloseIn skips IsD0 and picks 105
	lc := LastCloseIn(days)
	if lc == nil || *lc != 105 {
		t.Fatalf("LastCloseIn = %v, want 105", lc)
	}
}

func TestIPOListNullsLastAndSummary(t *testing.T) {
	d1 := dUTC(2026, 1, 10)
	d2 := dUTC(2025, 5, 1)
	masters := []MasterRecord{
		{
			Ticker: "AAAA", CompanyName: sPtr("Alpha"), ListingDate: &d1, IPOPrice: 100,
			NilaiIPO: fPtr(200e9), Underwriters: []Underwriter{{Broker: "CC", Role: "LEAD"}},
			DayReturns: []*float64{fPtr(10.0), nil}, Perf2W: fPtr(20.0),
		},
		{
			Ticker: "BBBB", CompanyName: sPtr(""), ListingDate: &d2, IPOPrice: 150,
			NilaiIPO: fPtr(600e9), Underwriters: []Underwriter{{Broker: "EP", Role: "CO"}},
			DayReturns: []*float64{fPtr(-5.0), fPtr(5.0)}, Perf2W: fPtr(1.0),
		},
		{
			Ticker: "CCCC", CompanyName: sPtr("Charlie"), ListingDate: &d2, IPOPrice: 200,
			NilaiIPO: nil, Underwriters: []Underwriter{{Broker: "CC", Role: "CO"}},
			DayReturns: []*float64{nil, fPtr(15.0)},
		},
	}
	res := BuildIPOListData(masters, ListParams{SortBy: "nilai", SortDir: 1})
	if res.Total != 3 || res.Rows[2].Ticker != "CCCC" {
		t.Fatalf("null nilai_ipo should sort last: %+v", res.Rows)
	}
	if res.KinerjaSummary.RataRata["d1"] == nil || *res.KinerjaSummary.RataRata["d1"] != 3 {
		t.Fatalf("rata_rata[d1] = %v, want 3 (js_round(2.5))", res.KinerjaSummary.RataRata["d1"])
	}
	if res.KinerjaSummary.KemungkinanNaik["d1"] == nil || *res.KinerjaSummary.KemungkinanNaik["d1"] != 50 {
		t.Fatalf("kemungkinan_naik[d1] = %v, want 50", res.KinerjaSummary.KemungkinanNaik["d1"])
	}
	// Empty string nama sorts last even when asc
	resNama := BuildIPOListData(masters, ListParams{SortBy: "nama", SortDir: 1})
	if resNama.Rows[2].Ticker != "BBBB" {
		t.Fatalf("empty company_name should sort last, got %s", resNama.Rows[2].Ticker)
	}
}

func TestScreenerFilterIndonesianCommaAndPresets(t *testing.T) {
	cands := []ScreenerRow{
		{Kode: "T1", Chg: fPtr(25), Ratrans: fPtr(2.5), VolRatio: fPtr(1.2)},
		{Kode: "T2", Chg: fPtr(-10), Ratrans: fPtr(1.5), VolRatio: fPtr(4.0)},
		{Kode: "T3", Chg: nil, Ratrans: fPtr(3.0), VolRatio: fPtr(2.0)},
	}
	conds := []ScreenerCondition{{K: "ratrans", Op: ">=", V: "1,5"}}
	ordered, activeCol, sk, sdir := FilterAndSortScreenerRows(cands, "all", conds, "volSurge", nil, nil)
	if activeCol != "volRatio" || sk != "volRatio" || sdir != -1 {
		t.Fatalf("unexpected sort meta: activeCol=%s sk=%s sdir=%d", activeCol, sk, sdir)
	}
	if len(ordered) != 3 || ordered[0].Kode != "T2" {
		t.Fatalf("unexpected volSurge ordering: %+v", ordered)
	}

	// Null chg excluded by chg=up
	upOnly, _, _, _ := FilterAndSortScreenerRows(cands, "up", nil, "topPerf", nil, nil)
	if len(upOnly) != 1 || upOnly[0].Kode != "T1" {
		t.Fatalf("expected only T1 for chg=up, got %+v", upOnly)
	}
}
