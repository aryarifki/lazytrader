package market_report

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aryarifki/lazytrader/backend/internal/core/auth"
	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/go-chi/chi/v5"
)

// Handler serves the 8 Market Report Generator endpoints with dual admin authorization.
type Handler struct {
	Engine           *Engine
	Guard            *auth.Guard
	TerminalAdminKey string
}

// Routes mounts /api/v1/terminal/market-report/*.
func Routes(deps Deps) chi.Router {
	h := &Handler{
		Engine:           NewEngine(deps),
		Guard:            deps.Guard,
		TerminalAdminKey: strings.TrimSpace(deps.TerminalAdminKey),
	}

	r := chi.NewRouter()
	r.Get("/status", h.GetStatus)
	r.Get("/history", h.GetHistory)
	r.Get("/detail", h.GetDetail)

	r.Group(func(ar chi.Router) {
		ar.Use(h.requireAdmin)
		ar.Post("/generate", h.Generate)
		ar.Post("/update-token", h.UpdateToken)
		ar.Post("/auto-renew", h.AutoRenew)
		ar.Get("/logs", h.GetLogs)
		ar.Post("/reset-status", h.ResetStatus)
	})

	return r
}

func extractAdminKey(r *http.Request) string {
	if cand := strings.TrimSpace(r.Header.Get("X-Admin-Key")); cand != "" {
		return cand
	}
	if cand := strings.TrimSpace(r.URL.Query().Get("admin_key")); cand != "" {
		return cand
	}
	authHdr := strings.TrimSpace(r.Header.Get("Authorization"))
	if len(authHdr) > 7 && strings.EqualFold(authHdr[:7], "bearer ") {
		return strings.TrimSpace(authHdr[7:])
	}
	return ""
}

func (h *Handler) isAdminAuthorized(r *http.Request) bool {
	if h.Guard != nil {
		if isAdmin, _ := h.Guard.IsAdminSession(r); isAdmin {
			return true
		}
	}
	if h.TerminalAdminKey == "" {
		return false
	}
	cand := extractAdminKey(r)
	if cand == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(cand), []byte(h.TerminalAdminKey)) == 1
}

func (h *Handler) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !h.isAdminAuthorized(r) {
			httpx.SetNoStore(w)
			w.Header().Set("WWW-Authenticate", "Bearer")
			httpx.WriteError(w, http.StatusUnauthorized, "Akses ditolak: Kunci Admin Terminal (X-Admin-Key) diperlukan untuk aksi administratif ini.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	data := h.Engine.GetStatus(h.isAdminAuthorized(r))
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    data,
	})
}

func (h *Handler) Generate(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	var req GenerateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "Invalid JSON payload.")
		return
	}
	session := 1
	if req.Session != nil {
		session = *req.Session
	}
	rateLimit := 180.0
	if req.RateLimit != nil {
		rateLimit = *req.RateLimit
	}
	workers := 8
	if req.Workers != nil {
		workers = *req.Workers
	}
	if session < 1 || session > 2 {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "Input for 'session' must be between 1 and 2.")
		return
	}
	if rateLimit < 30.0 || rateLimit > 300.0 {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "Input for 'rate_limit' must be between 30.0 and 300.0.")
		return
	}
	if workers < 1 || workers > 16 {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "Input for 'workers' must be between 1 and 16.")
		return
	}

	res := h.Engine.StartGeneration(session, rateLimit, workers)
	if ok, _ := res["success"].(bool); !ok {
		msg, _ := res["message"].(string)
		httpx.WriteError(w, http.StatusBadRequest, msg)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}

func (h *Handler) UpdateToken(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	var req UpdateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "Field 'token' is required.")
		return
	}
	res := h.Engine.UpdateToken(*req.Token)
	if ok, _ := res["success"].(bool); !ok {
		msg, _ := res["message"].(string)
		httpx.WriteError(w, http.StatusBadRequest, msg)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}

func (h *Handler) AutoRenew(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	res := h.Engine.TriggerAutoRenewToken()
	if ok, _ := res["success"].(bool); !ok {
		msg, _ := res["message"].(string)
		httpx.WriteError(w, http.StatusBadRequest, msg)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}

func (h *Handler) GetLogs(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"logs":    h.Engine.GetLogs(),
	})
}

func (h *Handler) ResetStatus(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	h.Engine.ResetError()
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Status engine berhasil di-reset ke idle.",
	})
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"history": h.Engine.GetHistory(),
	})
}

func (h *Handler) GetDetail(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	filename := r.URL.Query().Get("filename")
	if strings.TrimSpace(filename) == "" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "Query parameter 'filename' is required.")
		return
	}
	report := h.Engine.ReadReport(filename, h.isAdminAuthorized(r))
	if report == nil {
		httpx.WriteError(w, http.StatusNotFound, "File laporan tidak ditemukan.")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"report":  report,
	})
}
