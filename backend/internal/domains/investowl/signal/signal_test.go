package signal

import (
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func fPtr(v float64) *float64 { return &v }

func TestNormalizeProbHonestNull(t *testing.T) {
	if got := normalizeProb(nil); got != nil {
		t.Fatalf("expected nil, got %v", *got)
	}
	if got := normalizeProb(fPtr(62.5)); got == nil || math.Abs(*got-0.625) > 1e-6 {
		t.Fatalf("normalizeProb(62.5) = %v, want 0.625", got)
	}
	if got := normalizeProb(fPtr(0.335)); got == nil || math.Abs(*got-0.335) > 1e-6 {
		t.Fatalf("normalizeProb(0.335) = %v, want 0.335", got)
	}
}

func TestRefreshSignalCacheEndpoint(t *testing.T) {
	r := Routes(nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /refresh status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"status":"ok"`) || !strings.Contains(body, `"cleared_keys":0`) {
		t.Fatalf("unexpected /refresh response: %s", body)
	}
}
