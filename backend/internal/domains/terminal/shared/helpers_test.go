package shared

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestPyRoundAndInt(t *testing.T) {
	if got := PyRound(20.915, 2); got != 20.91 {
		t.Fatalf("PyRound(20.915, 2) = %v, want 20.91", got)
	}
	if got := PyRound(2.55, 1); got != 2.5 {
		t.Fatalf("PyRound(2.55, 1) = %v, want 2.5", got)
	}
	if got := PyRound(2.65, 1); got != 2.6 {
		t.Fatalf("PyRound(2.65, 1) = %v, want 2.6", got)
	}
	if got := PyRound(2.75, 1); got != 2.8 {
		t.Fatalf("PyRound(2.75, 1) = %v, want 2.8", got)
	}
	if got := PyRoundInt(2.5); got != 2 {
		t.Fatalf("PyRoundInt(2.5) = %d, want 2", got)
	}
	if got := PyRoundInt(3.5); got != 4 {
		t.Fatalf("PyRoundInt(3.5) = %d, want 4", got)
	}
}

func TestFormatISO8601UTCAndComma2(t *testing.T) {
	tZero := time.Date(2026, 11, 25, 9, 0, 0, 0, time.UTC)
	if got := FormatISO8601UTC(tZero); got != "2026-11-25T09:00:00+00:00" {
		t.Fatalf("FormatISO8601UTC zero = %q", got)
	}
	tMicro := time.Date(2026, 11, 25, 12, 51, 2, 814276000, time.UTC)
	if got := FormatISO8601UTC(tMicro); got != "2026-11-25T12:51:02.814276+00:00" {
		t.Fatalf("FormatISO8601UTC micro = %q", got)
	}
	if got := FormatComma2(5197.614); got != "5,197.61" {
		t.Fatalf("FormatComma2(5197.614) = %q", got)
	}
	if got := FormatComma2(-1234.5); got != "-1,234.50" {
		t.Fatalf("FormatComma2(-1234.5) = %q", got)
	}
}

func TestParseIntQueryRange(t *testing.T) {
	req := httptest.NewRequest("GET", "/?limit=50", nil)
	v, err := ParseIntQueryRange(req, "limit", 200, 1, 1000)
	if err != nil || v != 50 {
		t.Fatalf("got (%d, %v), want (50, nil)", v, err)
	}
	reqBad := httptest.NewRequest("GET", "/?limit=5000", nil)
	if _, err := ParseIntQueryRange(reqBad, "limit", 200, 1, 1000); err == nil {
		t.Fatalf("expected range error")
	}
}
