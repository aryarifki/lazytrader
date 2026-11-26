package auth

import (
	"context"
	"net/http"

	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
)

type ctxKey string

const userCtxKey ctxKey = "lazytrader_session_user"

// UserLookup refreshes approval/role status from the DB if needed.
type UserLookup interface {
	SyncOrGetByEmail(ctx context.Context, email, name, image, provider string) (*SessionUser, error)
}

// Guard provides fail-closed authentication and authorization middlewares.
type Guard struct {
	Secret      string
	AdminEmails []string
	Users       UserLookup
}

// UserFromContext extracts the authenticated SessionUser from request context.
func UserFromContext(ctx context.Context) (*SessionUser, bool) {
	u, ok := ctx.Value(userCtxKey).(*SessionUser)
	return u, ok && u != nil
}

// IsAdminSession checks whether the request has a valid InvestOwl admin session cookie.
func (g *Guard) IsAdminSession(r *http.Request) (isAdmin bool, hasCookie bool) {
	if g == nil {
		return false, false
	}
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil || cookie.Value == "" {
		return false, false
	}
	user, status, _ := g.extractUser(r)
	if user == nil || status != http.StatusOK {
		return false, true
	}
	return user.Role == "admin" || IsAdminEmail(user.Email, g.AdminEmails), true
}

func (g *Guard) extractUser(r *http.Request) (*SessionUser, int, string) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil || cookie.Value == "" {
		return nil, http.StatusUnauthorized, "Login diperlukan untuk mengakses fitur ini."
	}
	user, _, err := VerifySession(g.Secret, cookie.Value)
	if err != nil {
		return nil, http.StatusUnauthorized, "Sesi tidak valid atau kedaluwarsa."
	}
	if IsAdminEmail(user.Email, g.AdminEmails) {
		user.Role = "admin"
		user.IsApproved = true
	} else if !user.IsApproved && g.Users != nil {
		if fresh, err := g.Users.SyncOrGetByEmail(r.Context(), user.Email, user.Name, user.Image, "google"); err == nil && fresh != nil {
			user = fresh
		}
	}
	if user.Role != "admin" && !user.IsApproved {
		return nil, http.StatusForbidden, "Akun belum disetujui admin."
	}
	return user, http.StatusOK, ""
}

// RequireApproved rejects anonymous (401) or unapproved (403) callers.
func (g *Guard) RequireApproved(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, status, msg := g.extractUser(r)
		if user == nil {
			httpx.WriteError(w, status, msg)
			return
		}
		ctx := context.WithValue(r.Context(), userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin rejects non-admin callers with 403.
func (g *Guard) RequireAdmin(next http.Handler) http.Handler {
	return g.RequireApproved(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r.Context())
		if !ok || (user.Role != "admin" && !IsAdminEmail(user.Email, g.AdminEmails)) {
			httpx.WriteError(w, http.StatusForbidden, "Hanya Administrator yang berwenang mengakses fitur ini.")
			return
		}
		next.ServeHTTP(w, r)
	}))
}
