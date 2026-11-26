package shareholders

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
)

func TestClassifyShareholder(t *testing.T) {
	cases := []struct {
		name   string
		pct    float64
		ctrl   bool
		expect string
	}{
		{"PT Dwimuria Investama Andalan", 54.94, true, "controlling"},
		{"Saham Treasury", 0.42, false, "treasury"},
		{"Masyarakat Non Warkat", 42.07, false, "public"},
		{"Publik Warkat", 6.50, false, "public"},
		{"Institusi Besar", 7.25, false, "substantial"},
		{"Jahja Setiaatmadja", 0.03, false, "management"},
	}
	for _, tc := range cases {
		if got := ClassifyShareholder(tc.name, tc.pct, tc.ctrl); got != tc.expect {
			t.Fatalf("ClassifyShareholder(%q, %.2f, %v) = %q, want %q", tc.name, tc.pct, tc.ctrl, got, tc.expect)
		}
	}
}

func TestSharedSlugAndGroupFallback(t *testing.T) {
	if got := shared.Slug("Barito Pacific / Prajogo Pangestu"); got != "barito-pacific-prajogo-pangestu" {
		t.Fatalf("Slug unexpected: %q", got)
	}
	if got := shared.Slug("   "); got != "group" {
		t.Fatalf("Slug blank = %q, want group", got)
	}
	ordered, byID := shared.LoadConglomerateGroups(context.Background(), nil)
	if len(ordered) != 0 || len(byID) != 0 {
		t.Fatalf("expected empty groups on nil DB, got %d / %d", len(ordered), len(byID))
	}
}

func TestSlugAndShortLabelTruncation(t *testing.T) {
	longOwner := "PT Perusahaan Kepemilikan Panjang Sekali Indonesia Tbk"
	runes := []rune(longOwner)
	short := string(runes[:20]) + "..."
	if short != "PT Perusahaan Kepemi..." {
		t.Fatalf("shortLabel = %q", short)
	}
	if got := roundEven(20.915, 2); got != 20.91 {
		t.Fatalf("roundEven(20.915, 2) = %v, want 20.91 (CPython parity)", got)
	}
}

func TestRefreshCacheEndpoint(t *testing.T) {
	r := Routes(nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status = %v, want ok", body["status"])
	}
	if _, ok := body["cleared_keys"]; !ok {
		t.Fatalf("missing cleared_keys: %v", body)
	}
}
