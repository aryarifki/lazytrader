package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/redis/go-redis/v9"
)

var dataRoutePrefixes = []string{
	"/api/bandar",
	"/api/foreign-flow",
	"/api/broker-flow",
	"/api/signal",
	"/api/forensic",
	"/api/nyangkut",
	"/api/stocks",
	"/api/v1/ownership",
	"/api/v1/terminal",
}

// ExtractClientIP extracts the real visitor IP from Cloudflare / reverse-proxy headers.
func ExtractClientIP(r *http.Request) string {
	if cf := strings.TrimSpace(r.Header.Get("CF-Connecting-IP")); cf != "" {
		return cf
	}
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	if xri := strings.TrimSpace(r.Header.Get("X-Real-IP")); xri != "" {
		return xri
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}

// RateLimit mirrors InvestOwl's Redis fixed-window rate limiter (20/min auth, 80/min data).
func RateLimit(rdb *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rdb == nil {
				next.ServeHTTP(w, r)
				return
			}
			path := r.URL.Path
			limit := 0
			group := ""
			if strings.HasPrefix(path, "/api/users/sync") || strings.HasPrefix(path, "/api/v1/auth") {
				limit = 20
				group = "auth"
			} else {
				for _, p := range dataRoutePrefixes {
					if strings.HasPrefix(path, p) {
						limit = 80
						group = "data"
						break
					}
				}
			}
			if limit > 0 {
				ip := ExtractClientIP(r)
				key := fmt.Sprintf("rate_limit:%s:%s", group, ip)
				if count, err := rdb.Incr(r.Context(), key).Result(); err == nil {
					if count == 1 {
						_ = rdb.Expire(r.Context(), key, time.Minute).Err()
					}
					if int(count) > limit {
						w.Header().Set("Retry-After", "60")
						httpx.WriteError(w, http.StatusTooManyRequests, fmt.Sprintf("Rate limit exceeded. Terlalu banyak permintaan pada grup '%s'. Silakan tunggu 1 menit.", group))
						return
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
