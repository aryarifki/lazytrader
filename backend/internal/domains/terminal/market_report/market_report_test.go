package market_report

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/auth"
)

func makeTestJWT(iatOffsetSec, expOffsetSec int64) string {
	now := time.Now().Unix()
	hdr := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	pl, _ := json.Marshal(map[string]any{
		"iat": now + iatOffsetSec,
		"exp": now + expOffsetSec,
	})
	body := base64.RawURLEncoding.EncodeToString(pl)
	return hdr + "." + body + ".dummysignaturepadded00000000000000000000"
}

func setupTempMarketReport(t *testing.T, verifyHandler http.HandlerFunc, runner CommandRunner) (http.Handler, string, string) {
	t.Helper()
	tmp := t.TempDir()
	reportDir := filepath.Join(tmp, "reports")
	smDir := filepath.Join(tmp, "smtracker")
	if err := os.MkdirAll(reportDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(smDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(reportDir, "market_report_sesi1_2026-11-12.txt"), []byte("masked report body"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(reportDir, "market_report_sesi1_2026-11-12_unmasked.txt"), []byte("secret unmasked body"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(smDir, ".env"), []byte("BROKER_API_TOKEN=\"initial-valid-token-over-20-chars\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var verifyURL string
	if verifyHandler != nil {
		srv := httptest.NewServer(verifyHandler)
		t.Cleanup(srv.Close)
		verifyURL = srv.URL
	}

	guard := &auth.Guard{Secret: "unit-test-secret", AdminEmails: []string{"owner@example.com"}}
	r := Routes(Deps{
		Guard:             guard,
		TerminalAdminKey:  "test-admin-key-123",
		ReportDir:         reportDir,
		SMtrackerDir:      smDir,
		ScriptPath:        filepath.Join(tmp, "gen.py"),
		StockbitVerifyURL: verifyURL,
		Runner:            runner,
	})
	return r, reportDir, filepath.Join(smDir, ".env")
}

func TestMarketReportAuthAndMasking(t *testing.T) {
	r, _, _ := setupTempMarketReport(t, nil, nil)

	// 1. Anonymous rejected on all 5 admin routes with 401
	for _, tc := range []struct{ method, path, body string }{
		{"POST", "/generate", `{}`},
		{"POST", "/update-token", `{"token":"abc"}`},
		{"POST", "/auto-renew", `{}`},
		{"GET", "/logs", ``},
		{"POST", "/reset-status", `{}`},
	} {
		req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("%s %s anon status = %d, want 401", tc.method, tc.path, rec.Code)
		}
		if rec.Header().Get("WWW-Authenticate") != "Bearer" {
			t.Fatalf("%s %s missing WWW-Authenticate", tc.method, tc.path)
		}
	}

	// 2. Public /status & /detail strip unmasked_text
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest("GET", "/status", nil))
	if rec.Code != 200 {
		t.Fatalf("/status code = %d", rec.Code)
	}
	var stResp map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &stResp)
	lr := stResp["data"].(map[string]any)["latest_report"].(map[string]any)
	if lr["unmasked_text"] != nil {
		t.Fatalf("expected nil unmasked_text for anon /status, got %v", lr["unmasked_text"])
	}

	// 3. InvestOwl Admin Session cookie grants access to unmasked_text and /logs
	tok, _, err := auth.SignSession("unit-test-secret", auth.SessionUser{
		Email: "owner@example.com", Name: "Owner", Role: "admin", IsApproved: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	reqDet := httptest.NewRequest("GET", "/detail?filename=market_report_sesi1_2026-11-12.txt", nil)
	reqDet.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: tok})
	recDet := httptest.NewRecorder()
	r.ServeHTTP(recDet, reqDet)
	var detResp map[string]any
	_ = json.Unmarshal(recDet.Body.Bytes(), &detResp)
	repObj := detResp["report"].(map[string]any)
	if repObj["unmasked_text"] != "secret unmasked body" {
		t.Fatalf("admin cookie did not unlock unmasked_text: %v", repObj["unmasked_text"])
	}

	// 4. X-Admin-Key, Bearer, and admin_key query also unlock /logs
	for _, modReq := range []*http.Request{
		func() *http.Request {
			req := httptest.NewRequest("GET", "/logs", nil)
			req.Header.Set("X-Admin-Key", "test-admin-key-123")
			return req
		}(),
		func() *http.Request {
			req := httptest.NewRequest("GET", "/logs", nil)
			req.Header.Set("Authorization", "Bearer test-admin-key-123")
			return req
		}(),
		httptest.NewRequest("GET", "/logs?admin_key=test-admin-key-123", nil),
	} {
		recLogs := httptest.NewRecorder()
		r.ServeHTTP(recLogs, modReq)
		if recLogs.Code != 200 {
			t.Fatalf("authorized /logs got %d", recLogs.Code)
		}
	}

	// 5. Path traversal returns 404
	recTrav := httptest.NewRecorder()
	r.ServeHTTP(recTrav, httptest.NewRequest("GET", "/detail?filename=../../etc/passwd", nil))
	if recTrav.Code != http.StatusNotFound {
		t.Fatalf("path traversal status = %d, want 404", recTrav.Code)
	}
}

func TestMarketReportUpdateTokenAndRunner(t *testing.T) {
	stockbitStatus := http.StatusOK
	verifyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(stockbitStatus)
	})
	mockRunner := func(ctx context.Context, bin string, args []string, extraEnv map[string]string, onLine func(string)) (int, error) {
		onLine("[PROGRESS] 481/962 (50.0%) - Halfway")
		return 0, nil
	}
	r, _, envFile := setupTempMarketReport(t, verifyHandler, mockRunner)

	// Expired JWT -> 400
	expiredJWT := makeTestJWT(-7200, -10)
	req := httptest.NewRequest("POST", "/update-token", bytes.NewBufferString(fmt.Sprintf(`{"token":%q}`, expiredJWT)))
	req.Header.Set("X-Admin-Key", "test-admin-key-123")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != 400 {
		t.Fatalf("expired token status = %d, want 400", rec.Code)
	}

	// Stockbit 401 -> 400
	stockbitStatus = http.StatusUnauthorized
	freshJWT := makeTestJWT(-60, 3600)
	req2 := httptest.NewRequest("POST", "/update-token", bytes.NewBufferString(fmt.Sprintf(`{"token":"Bearer %s"}`, freshJWT)))
	req2.Header.Set("X-Admin-Key", "test-admin-key-123")
	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, req2)
	if rec2.Code != 400 {
		t.Fatalf("stockbit 401 token status = %d, want 400", rec2.Code)
	}

	// Valid -> 200 and updates .env
	stockbitStatus = http.StatusOK
	req3 := httptest.NewRequest("POST", "/update-token", bytes.NewBufferString(fmt.Sprintf(`{"token":%q}`, freshJWT)))
	req3.Header.Set("X-Admin-Key", "test-admin-key-123")
	rec3 := httptest.NewRecorder()
	r.ServeHTTP(rec3, req3)
	if rec3.Code != 200 {
		t.Fatalf("valid token status = %d, body=%s", rec3.Code, rec3.Body.String())
	}
	rawEnv, _ := os.ReadFile(envFile)
	if !strings.Contains(string(rawEnv), freshJWT) {
		t.Fatalf("expected updated .env to contain token")
	}

	// Generate with mock runner -> completes asynchronously
	reqGen := httptest.NewRequest("POST", "/generate", strings.NewReader(`{"session":1}`))
	reqGen.Header.Set("X-Admin-Key", "test-admin-key-123")
	recGen := httptest.NewRecorder()
	r.ServeHTTP(recGen, reqGen)
	if recGen.Code != 200 {
		t.Fatalf("/generate status = %d, body=%s", recGen.Code, recGen.Body.String())
	}
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		recSt := httptest.NewRecorder()
		r.ServeHTTP(recSt, httptest.NewRequest("GET", "/status", nil))
		var parsed map[string]any
		_ = json.Unmarshal(recSt.Body.Bytes(), &parsed)
		if parsed["data"].(map[string]any)["status"] == "completed" {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("mock generation did not reach completed status in time")
}
