package httpx

import (
	"encoding/json"
	"net/http"
)

// SetNoStore sets strict no-cache headers for auth, admin, and dynamic API responses.
func SetNoStore(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

// WriteJSON serializes v as JSON with the given HTTP status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteError returns a FastAPI-compatible error envelope: {"detail": "..."}.
func WriteError(w http.ResponseWriter, status int, detail string) {
	WriteJSON(w, status, map[string]string{"detail": detail})
}
