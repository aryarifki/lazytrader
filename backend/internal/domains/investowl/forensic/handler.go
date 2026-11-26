package forensic

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
)

// Handler serves /api/forensic endpoints.
type Handler struct {
	Svc *Service
}

// Refresh handles POST /api/forensic/refresh.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	h.Svc.Refresh(r.Context())
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"message": "Forensic cache cleared successfully.",
	})
}

// GetForeignExodus handles GET /api/forensic/foreign-exodus.
func (h *Handler) GetForeignExodus(w http.ResponseWriter, r *http.Request) {
	days := 45
	if raw := strings.TrimSpace(r.URL.Query().Get("days")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 10 || n > 90 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "days must be between 10 and 90")
			return
		}
		days = n
	}

	res, err := h.Svc.GetForeignExodusRadar(r.Context(), days)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate foreign exodus radar: %v", err))
		return
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}
