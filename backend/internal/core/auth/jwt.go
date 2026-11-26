package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	SessionCookieName = "lazytrader_session"
	SessionMaxAge     = 30 * 24 * time.Hour
)

// SessionUser holds the authenticated user claims stored in the HttpOnly JWT cookie.
type SessionUser struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name,omitempty"`
	Image      string `json:"image,omitempty"`
	Role       string `json:"role"`
	IsApproved bool   `json:"is_approved"`
}

type sessionClaims struct {
	User SessionUser `json:"user"`
	jwt.RegisteredClaims
}

// SignSession creates a signed HS256 JWT token valid for 30 days.
func SignSession(secret string, user SessionUser) (string, time.Time, error) {
	exp := time.Now().Add(SessionMaxAge)
	claims := sessionClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString([]byte(secret))
	return signed, exp, err
}

// VerifySession validates an HS256 JWT token and returns the SessionUser + expiry.
func VerifySession(secret, tokenStr string) (*SessionUser, time.Time, error) {
	if strings.TrimSpace(tokenStr) == "" {
		return nil, time.Time{}, errors.New("empty token")
	}
	tok, err := jwt.ParseWithClaims(tokenStr, &sessionClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !tok.Valid {
		return nil, time.Time{}, errors.New("invalid or expired session token")
	}
	claims, ok := tok.Claims.(*sessionClaims)
	if !ok || claims.User.Email == "" || claims.ExpiresAt == nil {
		return nil, time.Time{}, errors.New("malformed session claims")
	}
	return &claims.User, claims.ExpiresAt.Time, nil
}

// SetSessionCookie writes the HttpOnly session cookie on the response.
func SetSessionCookie(w http.ResponseWriter, r *http.Request, token string, maxAgeSeconds int) {
	secure := r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   maxAgeSeconds,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
	})
}
