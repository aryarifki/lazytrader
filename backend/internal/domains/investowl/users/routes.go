package users

import (
	"github.com/aryarifki/lazytrader/backend/internal/core/auth"
	"github.com/go-chi/chi/v5"
)

// Routes registers /api/users endpoints with exact InvestOwl access rules.
func Routes(h *Handler, guard *auth.Guard) chi.Router {
	r := chi.NewRouter()
	r.Post("/sync", h.SyncUser)
	r.Get("/settings/auth-gate", h.GetAuthGate)

	r.With(guard.RequireAdmin).Get("/", h.ListUsers)
	r.With(guard.RequireAdmin).Patch("/{user_id}/approve", h.ApproveUser)
	r.With(guard.RequireAdmin).Delete("/{user_id}", h.DeleteUser)
	r.With(guard.RequireAdmin).Patch("/settings/auth-gate", h.UpdateAuthGate)

	r.With(guard.RequireApproved).Get("/{email}/watchlist", h.GetWatchlist)
	r.With(guard.RequireApproved).Put("/{email}/watchlist", h.UpdateWatchlist)
	return r
}
