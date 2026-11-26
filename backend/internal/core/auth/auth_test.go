package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsAdminEmail(t *testing.T) {
	whitelist := []string{"boss@example.com"}
	if !IsAdminEmail("boss@example.com", whitelist) {
		t.Fatal("expected whitelisted email to be admin")
	}
	if !IsAdminEmail("adryanrifky716@gmail.com", nil) {
		t.Fatal("expected adryan email to be admin")
	}
	if !IsAdminEmail("arya.dev@gmail.com", nil) {
		t.Fatal("expected arya email to be admin")
	}
	if IsAdminEmail("guest@example.com", whitelist) {
		t.Fatal("expected guest email NOT to be admin")
	}
}

func TestSignAndVerifySession(t *testing.T) {
	secret := "test-secret-32-bytes-minimum-key"
	u := SessionUser{ID: 42, Email: "user@example.com", Name: "User", Role: "user", IsApproved: true}
	tok, _, err := SignSession(secret, u)
	if err != nil {
		t.Fatalf("SignSession failed: %v", err)
	}
	got, _, err := VerifySession(secret, tok)
	if err != nil {
		t.Fatalf("VerifySession failed: %v", err)
	}
	if got.Email != u.Email || got.ID != u.ID || !got.IsApproved {
		t.Fatalf("unexpected claims: %+v", got)
	}
	if _, _, err := VerifySession("wrong-secret", tok); err == nil {
		t.Fatal("expected wrong secret to fail")
	}
}

func TestGuard_FailClosed(t *testing.T) {
	secret := "test-secret-32-bytes-minimum-key"
	g := &Guard{Secret: secret}
	protected := g.RequireApproved(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// 1. Anonymous -> 401
	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	rec := httptest.NewRecorder()
	protected.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for anonymous, got %d", rec.Code)
	}

	// 2. Unapproved user -> 403
	tok, _, _ := SignSession(secret, SessionUser{ID: 7, Email: "pending@example.com", Role: "user", IsApproved: false})
	req2 := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req2.AddCookie(&http.Cookie{Name: SessionCookieName, Value: tok})
	rec2 := httptest.NewRecorder()
	protected.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for unapproved user, got %d", rec2.Code)
	}

	// 3. Approved non-admin -> 403 on RequireAdmin
	tok3, _, _ := SignSession(secret, SessionUser{ID: 8, Email: "ok@example.com", Role: "user", IsApproved: true})
	adminOnly := g.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req3 := httptest.NewRequest(http.MethodGet, "/api/admin", nil)
	req3.AddCookie(&http.Cookie{Name: SessionCookieName, Value: tok3})
	rec3 := httptest.NewRecorder()
	adminOnly.ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for non-admin, got %d", rec3.Code)
	}
}
