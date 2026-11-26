package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aryarifki/lazytrader/backend/internal/core/auth"
	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

// Handler implements all 9 InvestOwl /api/users/* endpoints.
type Handler struct {
	Repo        *Repository
	Redis       *redis.Client
	AdminEmails []string
}

func (h *Handler) SyncUser(w http.ResponseWriter, r *http.Request) {
	var req UserSyncRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Email) == "" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "Invalid user sync payload")
		return
	}
	u, err := h.Repo.SyncUser(r.Context(), req)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "Gagal menyinkronkan user")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, u)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	list, err := h.Repo.ListAll(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "Gagal memuat daftar user")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, list)
}

func (h *Handler) ApproveUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "user_id tidak valid")
		return
	}
	var req UserApprovalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "Payload approval tidak valid")
		return
	}
	u, err := h.Repo.SetApproval(r.Context(), id, req.IsApproved)
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.WriteError(w, http.StatusNotFound, "User tidak ditemukan")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "Gagal memperbarui status approval")
		return
	}
	httpx.SetNoStore(w)
	httpx.WriteJSON(w, http.StatusOK, u)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "user_id tidak valid")
		return
	}
	u, err := h.Repo.GetByID(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.WriteError(w, http.StatusNotFound, "User tidak ditemukan")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "Gagal membaca data user")
		return
	}
	if u.Role == "admin" || auth.IsAdminEmail(u.Email, h.AdminEmails) {
		httpx.WriteError(w, http.StatusBadRequest, "Akun Administrator tidak dapat dihapus")
		return
	}
	if err := h.Repo.Delete(r.Context(), id); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "Gagal menghapus user")
		return
	}
	httpx.SetNoStore(w)
	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": fmt.Sprintf("User %s berhasil dihapus", u.Email),
	})
}

func (h *Handler) GetWatchlist(w http.ResponseWriter, r *http.Request) {
	caller, _ := auth.UserFromContext(r.Context())
	targetEmail := strings.ToLower(strings.TrimSpace(chi.URLParam(r, "email")))
	callerEmail := strings.ToLower(strings.TrimSpace(caller.Email))
	if callerEmail != targetEmail && caller.Role != "admin" && !auth.IsAdminEmail(callerEmail, h.AdminEmails) {
		httpx.WriteError(w, http.StatusForbidden, "Tidak dapat melihat watchlist pengguna lain.")
		return
	}
	u, err := h.Repo.GetByEmail(r.Context(), targetEmail)
	if err != nil || u.Watchlist == nil {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"watchlist": []string{}})
		return
	}
	var items []string
	if err := json.Unmarshal([]byte(*u.Watchlist), &items); err != nil {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"watchlist": []string{}})
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"watchlist": items})
}

func (h *Handler) UpdateWatchlist(w http.ResponseWriter, r *http.Request) {
	caller, _ := auth.UserFromContext(r.Context())
	targetEmail := strings.ToLower(strings.TrimSpace(chi.URLParam(r, "email")))
	callerEmail := strings.ToLower(strings.TrimSpace(caller.Email))
	if callerEmail != targetEmail && caller.Role != "admin" && !auth.IsAdminEmail(callerEmail, h.AdminEmails) {
		httpx.WriteError(w, http.StatusForbidden, "Tidak dapat mengubah watchlist pengguna lain.")
		return
	}
	var req UserWatchlistUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "Payload watchlist tidak valid")
		return
	}
	if _, err := h.Repo.GetByEmail(r.Context(), targetEmail); errors.Is(err, pgx.ErrNoRows) {
		httpx.WriteError(w, http.StatusNotFound, "User tidak ditemukan")
		return
	}
	raw, _ := json.Marshal(req.Watchlist)
	if err := h.Repo.UpdateWatchlist(r.Context(), targetEmail, string(raw)); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "Gagal menyimpan watchlist")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"status": "ok", "watchlist": req.Watchlist})
}

func (h *Handler) GetAuthGate(w http.ResponseWriter, r *http.Request) {
	httpx.SetNoStore(w)
	redisKey := "setting:" + authGateKey
	if h.Redis != nil {
		if cached, err := h.Redis.Get(r.Context(), redisKey).Result(); err == nil {
			httpx.WriteJSON(w, http.StatusOK, AuthGateStatus{RequireLoginFullFeature: isTruthy(cached)})
			return
		}
	}
	val, err := h.Repo.GetAuthGate(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "Gagal membaca pengaturan auth-gate")
		return
	}
	if h.Redis != nil {
		strVal := "false"
		if val {
			strVal = "true"
		}
		_ = h.Redis.Set(r.Context(), redisKey, strVal, 0).Err()
	}
	httpx.WriteJSON(w, http.StatusOK, AuthGateStatus{RequireLoginFullFeature: val})
}

func (h *Handler) UpdateAuthGate(w http.ResponseWriter, r *http.Request) {
	var req AuthGateUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "Payload auth-gate tidak valid")
		return
	}
	if strings.TrimSpace(req.AdminEmail) != "" && !auth.IsAdminEmail(req.AdminEmail, h.AdminEmails) {
		u, err := h.Repo.GetByEmail(r.Context(), req.AdminEmail)
		if err != nil || u.Role != "admin" {
			httpx.WriteError(w, http.StatusForbidden, "Hanya Administrator yang berwenang mengubah pengaturan ini.")
			return
		}
	}
	if err := h.Repo.SetAuthGate(r.Context(), req.RequireLoginFullFeature); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "Gagal menyimpan pengaturan auth-gate")
		return
	}
	if h.Redis != nil {
		strVal := "false"
		if req.RequireLoginFullFeature {
			strVal = "true"
		}
		_ = h.Redis.Set(r.Context(), "setting:"+authGateKey, strVal, 0).Err()
	}
	httpx.SetNoStore(w)
	httpx.WriteJSON(w, http.StatusOK, AuthGateStatus{RequireLoginFullFeature: req.RequireLoginFullFeature})
}
