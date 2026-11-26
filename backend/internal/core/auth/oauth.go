package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	oauthStateCookie    = "lazytrader_oauth_state"
	oauthCallbackCookie = "lazytrader_oauth_cb"
)

// OAuthHandler handles Google OAuth2 login/callback/logout/session directly in Go.
type OAuthHandler struct {
	Guard       *Guard
	OAuthConfig *oauth2.Config
}

// NewOAuthHandler constructs the OAuth2 handler if Google credentials are configured.
func NewOAuthHandler(guard *Guard, clientID, clientSecret, appBaseURL string) *OAuthHandler {
	var cfg *oauth2.Config
	if clientID != "" && clientSecret != "" {
		cfg = &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  strings.TrimRight(appBaseURL, "/") + "/api/auth/callback/google",
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		}
	}
	return &OAuthHandler{Guard: guard, OAuthConfig: cfg}
}

// Routes registers /api/auth/* endpoints.
func (h *OAuthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/login", h.Login)
	r.Get("/callback/google", h.CallbackGoogle)
	r.Get("/session", h.Session)
	r.Post("/logout", h.Logout)
	r.Get("/logout", h.Logout)
	return r
}

func (h *OAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if h.OAuthConfig == nil {
		httpx.WriteError(w, http.StatusServiceUnavailable, "Google OAuth2 belum dikonfigurasi pada server.")
		return
	}
	state := randomState()
	cb := sanitizeCallbackURL(r.URL.Query().Get("callbackUrl"))

	secure := r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
	http.SetCookie(w, &http.Cookie{Name: oauthStateCookie, Value: state, Path: "/", MaxAge: 300, HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: secure})
	http.SetCookie(w, &http.Cookie{Name: oauthCallbackCookie, Value: cb, Path: "/", MaxAge: 300, HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: secure})

	authURL := h.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *OAuthHandler) CallbackGoogle(w http.ResponseWriter, r *http.Request) {
	if h.OAuthConfig == nil {
		http.Redirect(w, r, "/login?error=OAuthNotConfigured", http.StatusFound)
		return
	}
	stateCookie, err := r.Cookie(oauthStateCookie)
	if err != nil || stateCookie.Value == "" || stateCookie.Value != r.URL.Query().Get("state") {
		http.Redirect(w, r, "/login?error=InvalidState", http.StatusFound)
		return
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, "/login?error=AccessDenied", http.StatusFound)
		return
	}

	tok, err := h.OAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Redirect(w, r, "/login?error=OAuthExchangeFailed", http.StatusFound)
		return
	}
	client := h.OAuthConfig.Client(r.Context(), tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			_ = resp.Body.Close()
		}
		http.Redirect(w, r, "/login?error=UserInfoFailed", http.StatusFound)
		return
	}
	defer resp.Body.Close()

	var info struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil || info.Email == "" {
		http.Redirect(w, r, "/login?error=InvalidUserInfo", http.StatusFound)
		return
	}

	user := &SessionUser{
		Email:      info.Email,
		Name:       info.Name,
		Image:      info.Picture,
		Role:       "user",
		IsApproved: false,
	}
	if h.Guard.Users != nil {
		synced, err := h.Guard.Users.SyncOrGetByEmail(r.Context(), info.Email, info.Name, info.Picture, "google")
		if err == nil && synced != nil {
			user = synced
		}
	}
	if IsAdminEmail(user.Email, h.Guard.AdminEmails) {
		user.Role = "admin"
		user.IsApproved = true
	}

	signed, _, err := SignSession(h.Guard.Secret, *user)
	if err != nil {
		http.Redirect(w, r, "/login?error=SessionSignFailed", http.StatusFound)
		return
	}

	SetSessionCookie(w, r, signed, int(SessionMaxAge.Seconds()))
	// Clear temporary OAuth cookies
	http.SetCookie(w, &http.Cookie{Name: oauthStateCookie, Value: "", Path: "/", MaxAge: -1, HttpOnly: true})
	http.SetCookie(w, &http.Cookie{Name: oauthCallbackCookie, Value: "", Path: "/", MaxAge: -1, HttpOnly: true})

	target := "/"
	if cbCookie, err := r.Cookie(oauthCallbackCookie); err == nil {
		target = sanitizeCallbackURL(cbCookie.Value)
	}
	if !user.IsApproved && user.Role != "admin" {
		target = "/pending"
	}
	http.Redirect(w, r, target, http.StatusFound)
}

func (h *OAuthHandler) Session(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil || cookie.Value == "" {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{})
		return
	}
	user, exp, err := VerifySession(h.Guard.Secret, cookie.Value)
	if err != nil {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{})
		return
	}
	if IsAdminEmail(user.Email, h.Guard.AdminEmails) {
		user.Role = "admin"
		user.IsApproved = true
	} else if !user.IsApproved && h.Guard.Users != nil {
		if fresh, err := h.Guard.Users.SyncOrGetByEmail(r.Context(), user.Email, user.Name, user.Image, "google"); err == nil && fresh != nil {
			user = fresh
			if signed, _, err := SignSession(h.Guard.Secret, *user); err == nil {
				SetSessionCookie(w, r, signed, int(SessionMaxAge.Seconds()))
			}
		}
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"user":    user,
		"expires": exp.UTC().Format(time.RFC3339),
	})
}

func (h *OAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	SetSessionCookie(w, r, "", -1)
	if r.Method == http.MethodGet {
		cb := sanitizeCallbackURL(r.URL.Query().Get("callbackUrl"))
		http.Redirect(w, r, cb, http.StatusFound)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func randomState() string {
	b := make([]byte, 24)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func sanitizeCallbackURL(raw string) string {
	if raw == "" {
		return "/"
	}
	u, err := url.Parse(raw)
	if err != nil || u.IsAbs() || !strings.HasPrefix(u.Path, "/") || strings.HasPrefix(u.Path, "//") {
		return "/"
	}
	if u.RawQuery != "" {
		return u.Path + "?" + u.RawQuery
	}
	return u.Path
}
