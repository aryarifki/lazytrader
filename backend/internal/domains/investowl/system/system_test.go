package system

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/auth"
	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
)

type stubUserLookup struct{}

func (stubUserLookup) SyncOrGetByEmail(_ context.Context, email, name, image, provider string) (*auth.SessionUser, error) {
	return &auth.SessionUser{Email: email, Name: name, Role: "user", IsApproved: true}, nil
}

func TestBuildJobCommand_AllJobs(t *testing.T) {
	mgr := &JobManager{SMtrackerDir: "/root/SMtracker", PythonBin: "/root/SMtracker/venv/bin/python3"}

	req := DefaultJobRunRequest()
	days3 := 3
	req.Days = &days3
	req.Force = true
	gotBackfill := mgr.BuildJobCommand("backfill", JobConfigs["backfill"], req)
	wantBackfill := []string{"/root/SMtracker/venv/bin/python3", "backfill_daily.py", "--universe", "all", "--days", "3", "--force"}
	if !reflect.DeepEqual(gotBackfill, wantBackfill) {
		t.Fatalf("backfill cmd = %v, want %v", gotBackfill, wantBackfill)
	}

	gotPipeline := mgr.BuildJobCommand("pipeline", JobConfigs["pipeline"], DefaultJobRunRequest())
	wantPipeline := []string{"/bin/bash", "/root/SMtracker/run_daily.sh"}
	if !reflect.DeepEqual(gotPipeline, wantPipeline) {
		t.Fatalf("pipeline cmd = %v, want %v", gotPipeline, wantPipeline)
	}

	gotETL := mgr.BuildJobCommand("etl_flow", JobConfigs["etl_flow"], DefaultJobRunRequest())
	wantETL := []string{"/root/SMtracker/venv/bin/python3", "etl_foreign_flow.py", "--universe", "watchlist", "--window", "20", "60"}
	if !reflect.DeepEqual(gotETL, wantETL) {
		t.Fatalf("etl_flow cmd = %v, want %v", gotETL, wantETL)
	}
}

func TestSanitizeLogFilename(t *testing.T) {
	if _, ok := SanitizeLogFilename("../.env"); ok {
		t.Fatalf("expected path traversal to be rejected")
	}
	if _, ok := SanitizeLogFilename("a/b.log"); ok {
		t.Fatalf("expected slash in filename to be rejected")
	}
	clean, ok := SanitizeLogFilename("scanner_20261125.log")
	if !ok || clean != "scanner_20261125.log" {
		t.Fatalf("expected clean filename, got %q ok=%v", clean, ok)
	}
}

func TestSystemAuthAndDryRun(t *testing.T) {
	secret := "test-secret-key"
	guard := &auth.Guard{
		Secret:      secret,
		AdminEmails: []string{"admin@example.com"},
		Users:       stubUserLookup{},
	}
	r := Routes(Deps{
		Cache:        &cache.Store{},
		Guard:        guard,
		SMtrackerDir: t.TempDir(),
	})

	// 1. Anonymous caller must get 401
	anonReq := httptest.NewRequest(http.MethodGet, "/jobs/status", nil)
	anonRec := httptest.NewRecorder()
	r.ServeHTTP(anonRec, anonReq)
	if anonRec.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous /jobs/status status = %d, want 401", anonRec.Code)
	}

	// 2. Signed admin cookie succeeds and dry_run works
	token, _, err := auth.SignSession(secret, auth.SessionUser{
		ID:         1,
		Email:      "admin@example.com",
		Name:       "Admin",
		Role:       "admin",
		IsApproved: true,
	})
	if err != nil {
		t.Fatalf("SignSession error: %v", err)
	}

	body, _ := json.Marshal(map[string]any{"dry_run": true, "universe": "watchlist"})
	dryReq := httptest.NewRequest(http.MethodPost, "/jobs/backfill/run", bytes.NewReader(body))
	dryReq.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: token})
	dryRec := httptest.NewRecorder()
	r.ServeHTTP(dryRec, dryReq)
	if dryRec.Code != http.StatusOK {
		t.Fatalf("admin dry_run status = %d, want 200 (%s)", dryRec.Code, dryRec.Body.String())
	}
	if !strings.Contains(dryRec.Body.String(), `"dry_run"`) {
		t.Fatalf("expected dry_run payload, got %s", dryRec.Body.String())
	}
}

func TestPurgeCache_IncludesTerminalCache(t *testing.T) {
	secret := "test-secret-key"
	guard := &auth.Guard{
		Secret:      secret,
		AdminEmails: []string{"admin@example.com"},
		Users:       stubUserLookup{},
	}
	invCache := cache.NewMemoryOnlyStore(16)
	termCache := cache.NewMemoryOnlyStore(16)
	invCache.Set(t.Context(), "signal:daily", []byte(`{"ok":true}`), time.Minute)
	termCache.Set(t.Context(), "mmd:tickers", []byte(`{"ok":true}`), time.Minute)

	r := Routes(Deps{
		Cache:         invCache,
		TerminalCache: termCache,
		Guard:         guard,
		SMtrackerDir:  t.TempDir(),
	})

	token, _, err := auth.SignSession(secret, auth.SessionUser{
		ID:         1,
		Email:      "admin@example.com",
		Name:       "Admin",
		Role:       "admin",
		IsApproved: true,
	})
	if err != nil {
		t.Fatalf("SignSession error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/cache/purge?pattern=*&purge_cloudflare=false", nil)
	req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: token})
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("PurgeCache status = %d, want 200 (%s)", rec.Code, rec.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if int(resp["deleted_keys"].(float64)) != 2 || int(resp["terminal_deleted_keys"].(float64)) != 1 {
		t.Fatalf("unexpected purge counts: %v", resp)
	}
	if _, ok := invCache.Get(t.Context(), "signal:daily"); ok {
		t.Fatalf("expected invCache L1 key to be purged")
	}
	if _, ok := termCache.Get(t.Context(), "mmd:tickers"); ok {
		t.Fatalf("expected termCache L1 key to be purged")
	}
}

func TestRefreshTerminalCache_LoopbackOnly(t *testing.T) {
	termCache := cache.NewMemoryOnlyStore(16)
	termCache.Set(t.Context(), "konglo:summary", []byte(`{"ok":true}`), time.Minute)
	h := NewTerminalCacheRefreshHandler(termCache)

	// 1. Proxied CF-Connecting-IP must be rejected (403) even from 127.0.0.1
	cfReq := httptest.NewRequest(http.MethodPost, "/api/v1/terminal/cache/refresh", nil)
	cfReq.RemoteAddr = "127.0.0.1:12345"
	cfReq.Header.Set("CF-Connecting-IP", "1.2.3.4")
	cfRec := httptest.NewRecorder()
	h.ServeHTTP(cfRec, cfReq)
	if cfRec.Code != http.StatusForbidden {
		t.Fatalf("CF-Connecting-IP status = %d, want 403", cfRec.Code)
	}

	// 2. Non-loopback IP must be rejected (403)
	extReq := httptest.NewRequest(http.MethodPost, "/api/v1/terminal/cache/refresh", nil)
	extReq.RemoteAddr = "203.0.113.5:12345"
	extRec := httptest.NewRecorder()
	h.ServeHTTP(extRec, extReq)
	if extRec.Code != http.StatusForbidden {
		t.Fatalf("non-loopback status = %d, want 403", extRec.Code)
	}

	// 3. Direct loopback succeeds (200) and purges termCache L1
	okReq := httptest.NewRequest(http.MethodPost, "/api/v1/terminal/cache/refresh", nil)
	okReq.RemoteAddr = "127.0.0.1:12345"
	okRec := httptest.NewRecorder()
	h.ServeHTTP(okRec, okReq)
	if okRec.Code != http.StatusOK {
		t.Fatalf("loopback status = %d, want 200 (%s)", okRec.Code, okRec.Body.String())
	}
	if _, ok := termCache.Get(t.Context(), "konglo:summary"); ok {
		t.Fatalf("expected termCache L1 key to be cleared")
	}
}
