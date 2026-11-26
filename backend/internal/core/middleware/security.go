package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// SecurityHeaders applies InvestOwl + Avenium security headers and X-Process-Time-Ms.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h := w.Header()
		h.Set("X-Frame-Options", "SAMEORIGIN")
		h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		h.Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline' https:; img-src 'self' data: https: blob:; font-src 'self' data: https:; connect-src 'self' https://api.iconify.design https://*.avenia.my.id; worker-src 'self' blob:;")

		next.ServeHTTP(w, r)

		elapsedMs := float64(time.Since(start).Microseconds()) / 1000.0
		if w.Header().Get("X-Process-Time-Ms") == "" && !strings.HasPrefix(r.URL.Path, "/assets/") {
			w.Header().Set("X-Process-Time-Ms", fmt.Sprintf("%.2f", elapsedMs))
		}
	})
}
