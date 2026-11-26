package nyangkut

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestComputePBIRARScores(t *testing.T) {
	rar, pbi := ComputePBIRARScores(20.24, 1.90, 0.0, 1.73, 44.4)
	if rar != 53.2 || pbi != 0.0 {
		t.Fatalf("rallying case got rar=%.1f pbi=%.1f, want 53.2 / 0.0", rar, pbi)
	}

	rar2, pbi2 := ComputePBIRARScores(40.0, 2.5, 5.0, -35.0, 80.0)
	if rar2 <= 0 || pbi2 <= 0 || pbi2 > 100 {
		t.Fatalf("underwater case out of bounds: rar=%.1f pbi=%.1f", rar2, pbi2)
	}
}

func TestFormatIntDotAndRoundEven(t *testing.T) {
	if got := formatIntDot(759569); got != "759.569" {
		t.Fatalf("formatIntDot(759569) = %q", got)
	}
	if got := formatIntDot(1234567); got != "1.234.567" {
		t.Fatalf("formatIntDot(1234567) = %q", got)
	}
	if got := roundEven(-0.0001, 2); got != 0.0 {
		t.Fatalf("roundEven(-0.0001, 2) = %v", got)
	}
	if got := roundEven(80.9675, 1); got != 81.0 {
		t.Fatalf("roundEven(80.9675, 1) = %v", got)
	}
}

func TestRefreshCacheEndpoint(t *testing.T) {
	r := chi.NewRouter()
	Mount(r, nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/nyangkut/refresh", nil)
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
		t.Fatalf("status field = %v, want ok", body["status"])
	}
	if _, ok := body["cleared_keys"]; !ok {
		t.Fatalf("missing cleared_keys in response: %v", body)
	}
}
