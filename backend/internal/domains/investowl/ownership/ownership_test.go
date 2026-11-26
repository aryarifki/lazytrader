package ownership

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRoundEvenExactCPythonParity(t *testing.T) {
	if got := roundEven(20.915, 2); got != 20.91 {
		t.Fatalf("roundEven(20.915, 2) = %v, want 20.91", got)
	}
	if got := roundEven(-0.0, 2); got != 0 || math.Signbit(got) {
		t.Fatalf("roundEven(-0.0, 2) = %v (sign=%v), want +0.0", got, math.Signbit(got))
	}
	if got := roundEven(math.NaN(), 2); got != 0 {
		t.Fatalf("roundEven(NaN, 2) = %v, want 0", got)
	}
	if got := roundEven(1234.5, -1); got != 1230 {
		t.Fatalf("roundEven(1234.5, -1) = %v, want 1230", got)
	}
}

func TestMarketCapProvenance(t *testing.T) {
	f64 := func(v float64) *float64 { return &v }
	sp := func(s string) *string { return &s }

	gotInd := marketCapProvenance(f64(1000), sp("2026-11-25"), sp("prices"), f64(5e6), sp("2026-10-23"), sp("company_shares"), nil)
	if gotInd.Basis != "computed" || !gotInd.Independent || *gotInd.Value != 5e9 || *gotInd.Source != "prices x company_shares" {
		t.Fatalf("unexpected computed independent provenance: %+v", gotInd)
	}

	gotMixed := marketCapProvenance(f64(1000), nil, sp("ct_baseline"), f64(5e6), nil, sp("ct_baseline.ticker_meta"), f64(5))
	if gotMixed.Basis != "computed" || gotMixed.Independent {
		t.Fatalf("expected non-independent computed provenance, got %+v", gotMixed)
	}

	gotRep := marketCapProvenance(nil, nil, nil, nil, nil, nil, f64(12.5))
	if gotRep.Basis != "reported" || gotRep.Independent || *gotRep.Value != 12.5e9 || *gotRep.Source != "ct_baseline.ticker_meta" {
		t.Fatalf("unexpected reported provenance: %+v", gotRep)
	}

	gotUnavail := marketCapProvenance(nil, nil, nil, nil, nil, nil, nil)
	if gotUnavail.Basis != "unavailable" || gotUnavail.Value != nil {
		t.Fatalf("unexpected unavailable provenance: %+v", gotUnavail)
	}
}

func TestHelpersAndCTFunctions(t *testing.T) {
	if c := classifyHolder("PT X", 12, true); c != "controlling" {
		t.Fatalf("classify controlling = %s", c)
	}
	if c := classifyHolder("Saham Treasury PT X", 1, false); c != "treasury" {
		t.Fatalf("classify treasury = %s", c)
	}
	if c := classifyHolder("Masyarakat Warkat", 30, false); c != "public" {
		t.Fatalf("classify public = %s", c)
	}
	if c := classifyHolder("Vanguard", 5.0, false); c != "substantial" {
		t.Fatalf("classify substantial = %s", c)
	}
	if c := classifyHolder("Director", 1.2, false); c != "management" {
		t.Fatalf("classify management = %s", c)
	}

	d := time.Date(2026, time.August, 31, 0, 0, 0, 0, time.UTC)
	if lbl := periodLabel(&d); lbl != "Agu 2026" {
		t.Fatalf("periodLabel = %s, want Agu 2026", lbl)
	}
	if lbl := periodLabel(nil); lbl != "-" {
		t.Fatalf("periodLabel(nil) = %s, want -", lbl)
	}

	f64 := func(v float64) *float64 { return &v }
	if pc := pctChange(f64(110), f64(100)); pc == nil || *pc != 10.0 {
		t.Fatalf("pctChange(110, 100) = %v, want 10.0", pc)
	}
	if pc := pctChange(f64(110), f64(0)); pc != nil {
		t.Fatalf("pctChange(110, 0) = %v, want nil", *pc)
	}

	if nc := normalizeCompany("PT Bank Rakyat Indonesia (Persero) Tbk."); nc != "BANK RAKYAT INDONESIA" {
		t.Fatalf("normalizeCompany = %q", nc)
	}

	parsed, ok := parseFullDate("30 September 2026")
	if !ok || parsed.Format("2006-01-02") != "2026-09-30" {
		t.Fatalf("parseFullDate failed: ok=%v d=%v", ok, parsed)
	}
	if gl := glMonthLabel(map[string]any{"date_jun": "30 September 2026"}); gl != "Oktober 2026" {
		t.Fatalf("glMonthLabel = %q, want Oktober 2026", gl)
	}

	if bucketPct(10.1) != ">10" || bucketPct(5.0) != "5-10" || bucketPct(1.0) != "1-5" || bucketPct(0.5) != "<1" {
		t.Fatalf("bucketPct mismatch")
	}
	if bucketSIDPct(10.1) != ">10" || bucketSIDPct(3.0) != "3-10" || bucketSIDPct(1.0) != "1-3" || bucketSIDPct(0.9) != "<1" {
		t.Fatalf("bucketSIDPct mismatch")
	}
	if bucketValue(101) != ">100M" || bucketValue(5) != "5-100M" || bucketValue(4.9) != "<5M" {
		t.Fatalf("bucketValue mismatch")
	}
	if bucketMktcap(f64(250)) != "micro" || bucketMktcap(f64(500)) != "small" || bucketMktcap(f64(2500)) != "midsmall" || bucketMktcap(f64(20000)) != "mid" || bucketMktcap(f64(60000)) != "big" {
		t.Fatalf("bucketMktcap mismatch")
	}
	if !passInvestorType("Local ID", "Local", map[string]bool{"ID": true}) {
		t.Fatalf("expected Local ID to pass")
	}
	if passInvestorType("Foreign CP", "Local", map[string]bool{"ID": true, "CP": true}) {
		t.Fatalf("expected Foreign CP to fail Local filter")
	}
	if !passInvestorType("Foreign IS", "Foreign", map[string]bool{"OTHER": true}) {
		t.Fatalf("expected Foreign IS to map to OTHER")
	}
}

func TestInvestorGLPeakCapital(t *testing.T) {
	f64 := func(v float64) *float64 { return &v }
	iPtr := func(v int) *int { return &v }

	// Base: 1000 shares @ 1000 = 1,000,000 modal_awal.
	// T5 same month: buy 10 lot (1000 shares) @ 1000 -> value_m = 0.001 (1M Rp), so cash = -1,000,000 -> peak_basis = 2,000,000.
	// Closing price = 1500 -> market_now = 2000 * 1500 = 3,000,000; modal = 2,000,000; gl = 1,000,000; gl_pct = 50.0.
	ctx := &ctContext{
		Year:       2026,
		GLYear:     2026,
		GLMonthNum: 10,
		Tickers: map[string]*ctContextTicker{
			"AAAA": {Ticker: "AAAA", Price: f64(1000)},
			"BBBB": {Ticker: "BBBB", Price: f64(500)},
		},
		OHLC: map[string]ctContextOHLC{
			"AAAA": {Close: f64(1500)},
			"BBBB": {Close: f64(600)},
		},
		T5ByTicker: map[string][]ctContextT5{
			"AAAA": {
				{
					Ticker:    "AAAA",
					Name:      "Investor Alpha",
					Action:    "BUY",
					MonthNum:  iPtr(10),
					DayNum:    iPtr(5),
					ChangeLot: f64(10),
					ValueM:    f64(0.001),
				},
			},
		},
		T5ByName: map[string][]ctContextT5{
			"Investor Alpha": {
				{
					Ticker:    "AAAA",
					Name:      "Investor Alpha",
					Action:    "BUY",
					MonthNum:  iPtr(10),
					DayNum:    iPtr(5),
					ChangeLot: f64(10),
					ValueM:    f64(0.001),
				},
				{
					Ticker:    "BBBB",
					Name:      "Investor Alpha",
					Action:    "BUY",
					MonthNum:  iPtr(10),
					DayNum:    iPtr(8),
					ChangeLot: f64(10),
					ValueM:    f64(0.0005),
				},
			},
		},
	}

	holdings := []ctContextInvHolding{
		{Name: "Investor Alpha", Ticker: "AAAA", SLJun: f64(1000), Status: "both"},
	}
	res := investorGL(ctx, "Investor Alpha", holdings)
	if res == nil {
		t.Fatalf("expected non-nil GL result")
	}
	// Extra BBBB adds 1000 shares * 600 = 600,000 market_now, +500,000 modal, cash -> -1,500,000 -> peak_basis = 2,500,000.
	// Total market_now = 3,000,000 + 600,000 = 3,600,000; modal = 2,500,000; gl = 1,100,000; gl_pct = 44.0.
	if math.Abs(res.PeakBasis-2_500_000) > 1e-6 {
		t.Fatalf("PeakBasis = %v, want 2500000", res.PeakBasis)
	}
	if math.Abs(res.GL-1_100_000) > 1e-6 {
		t.Fatalf("GL = %v, want 1100000", res.GL)
	}
	if res.GLPct == nil || math.Abs(*res.GLPct-44.0) > 1e-6 {
		t.Fatalf("GLPct = %v, want 44.0", res.GLPct)
	}
}

func TestSortCTRowsNullPlacement(t *testing.T) {
	f64 := func(v float64) *float64 { return &v }
	rows := []map[string]any{
		{"ticker": "AAA", "mtd_pct": (*float64)(nil)},
		{"ticker": "BBB", "mtd_pct": f64(2.5)},
		{"ticker": "CCC", "mtd_pct": f64(5.0)},
	}
	_, getters := baseExplorerGetters()
	sorted, err := sortCTRows(rows, "mtd", "desc", "ticker", getters)
	if err != nil {
		t.Fatalf("sortCTRows failed: %v", err)
	}
	if sorted[0]["ticker"] != "CCC" || sorted[1]["ticker"] != "BBB" || sorted[2]["ticker"] != "AAA" {
		t.Fatalf("unexpected sort order: %+v", sorted)
	}
}

func TestRoutesAndSyncCatalogHandler(t *testing.T) {
	r := Routes(nil, nil)

	// Query validation check (no DB needed when validation rejects early)
	for _, tc := range []struct {
		method string
		path   string
		want   int
	}{
		{"GET", "/investors?limit=0", http.StatusUnprocessableEntity},
		{"GET", "/investors/detail?name=A", http.StatusUnprocessableEntity},
		{"GET", "/investors/ct/detail?name=A", http.StatusUnprocessableEntity},
		{"GET", "/market?history_limit=5", http.StatusUnprocessableEntity},
		{"GET", "/explorer?mode=invalid", http.StatusUnprocessableEntity},
		{"GET", "/explorer/ct?category=bad", http.StatusUnprocessableEntity},
	} {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != tc.want {
			t.Fatalf("%s %s -> status %d, want %d", tc.method, tc.path, rec.Code, tc.want)
		}
	}

	// Mock sync catalog success and error fallback
	hOK := &Handler{
		SyncCatalogRunner: func(r *http.Request) (map[string]any, error) {
			return map[string]any{"status": "success", "catalogs_synced": 3}, nil
		},
	}
	recOK := httptest.NewRecorder()
	hOK.SyncGroupsCatalog(recOK, httptest.NewRequest("POST", "/groups/sync-catalog", nil))
	if recOK.Code != http.StatusOK {
		t.Fatalf("SyncGroupsCatalog ok status = %d", recOK.Code)
	}
	var bodyOK map[string]any
	_ = json.Unmarshal(recOK.Body.Bytes(), &bodyOK)
	if bodyOK["status"] != "success" || int(bodyOK["catalogs_synced"].(float64)) != 3 {
		t.Fatalf("unexpected sync ok payload: %+v", bodyOK)
	}

	hErr := &Handler{
		SyncCatalogRunner: func(r *http.Request) (map[string]any, error) {
			return nil, errors.New("token expired")
		},
	}
	recErr := httptest.NewRecorder()
	hErr.SyncGroupsCatalog(recErr, httptest.NewRequest("POST", "/groups/sync-catalog", nil))
	if recErr.Code != http.StatusOK {
		t.Fatalf("SyncGroupsCatalog error status = %d", recErr.Code)
	}
	var bodyErr map[string]any
	_ = json.Unmarshal(recErr.Body.Bytes(), &bodyErr)
	if bodyErr["status"] != "error" || bodyErr["message"] != "token expired" {
		t.Fatalf("unexpected sync error payload: %+v", bodyErr)
	}
}
