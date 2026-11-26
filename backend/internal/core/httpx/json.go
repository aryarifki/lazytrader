package httpx

import (
	"encoding/json"
	"fmt"
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
	payload, err := json.Marshal(v)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"detail":%q}`, "json encode error: "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(payload)
}

// WriteError returns a FastAPI-compatible error envelope: {"detail": "..."}.
func WriteError(w http.ResponseWriter, status int, detail string) {
	WriteJSON(w, status, map[string]string{"detail": detail})
}

// WriteCachedJSON writes pre-serialized JSON with InvestOwl Cache-Control & X-Cache headers.
func WriteCachedJSON(w http.ResponseWriter, body []byte, clientMaxAge, sMaxAge int, hit bool) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set(
		"Cache-Control",
		fmt.Sprintf("public, max-age=%d, s-maxage=%d, stale-while-revalidate=30", clientMaxAge, sMaxAge),
	)
	if hit {
		w.Header().Set("X-Cache", "HIT")
	} else {
		w.Header().Set("X-Cache", "MISS")
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
