package shared

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
)

var wibLocation = time.FixedZone("WIB", 7*3600)

// WIBLocation returns the fixed Asia/Jakarta (UTC+7) location.
func WIBLocation() *time.Location {
	return wibLocation
}

// PyRound matches CPython decimal-exact round(v, digits) using round-half-to-even formatting.
func PyRound(v float64, digits int) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return v
	}
	if digits < 0 {
		scale := math.Pow10(-digits)
		return math.RoundToEven(v/scale) * scale
	}
	rounded, err := strconv.ParseFloat(strconv.FormatFloat(v, 'f', digits, 64), 64)
	if err != nil {
		return v
	}
	if rounded == 0 {
		return 0
	}
	return rounded
}

// PyRoundInt matches CPython single-argument round(x) using round-half-to-even.
func PyRoundInt(v float64) int64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return int64(math.RoundToEven(v))
}

// FormatISO8601UTC formats time in FastAPI/psycopg2 UTC ISO-8601 style (+00:00 and 6-digit microseconds when non-zero).
func FormatISO8601UTC(t time.Time) string {
	u := t.UTC()
	micro := u.Nanosecond() / 1000
	if micro == 0 {
		return u.Format("2006-01-02T15:04:05+00:00")
	}
	return fmt.Sprintf("%s.%06d+00:00", u.Format("2006-01-02T15:04:05"), micro)
}

// WriteCachedTerminalJSON writes terminal API JSON with strict no-store headers matching Avenium FastAPI.
func WriteCachedTerminalJSON(w http.ResponseWriter, body []byte) {
	httpx.SetNoStore(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

// GetOrComputeCached fetches JSON bytes from cache or computes, marshals, and stores them with ttl.
func GetOrComputeCached(ctx context.Context, store *cache.Store, key string, ttl time.Duration, computeFn func(context.Context) (any, int, string, error)) ([]byte, int, string, error) {
	if store != nil && key != "" {
		if raw, ok := store.Get(ctx, key); ok {
			return raw, http.StatusOK, "", nil
		}
	}
	val, status, detail, err := computeFn(ctx)
	if err != nil {
		return nil, status, detail, err
	}
	if status != 0 && status != http.StatusOK {
		return nil, status, detail, nil
	}
	raw, err := json.Marshal(val)
	if err != nil {
		return nil, http.StatusInternalServerError, err.Error(), err
	}
	if store != nil && key != "" {
		store.Set(ctx, key, raw, ttl)
	}
	return raw, http.StatusOK, "", nil
}

// ParseIntQueryRange parses an integer query parameter with FastAPI-compatible 422 range validation.
func ParseIntQueryRange(r *http.Request, name string, defVal, minVal, maxVal int) (int, error) {
	raw := strings.TrimSpace(r.URL.Query().Get(name))
	if raw == "" {
		return defVal, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("Input should be a valid integer for '%s'", name)
	}
	if v < minVal || v > maxVal {
		return 0, fmt.Errorf("Input for '%s' must be between %d and %d", name, minVal, maxVal)
	}
	return v, nil
}

// FormatComma2 formats a float with comma thousands separators and 2 decimals matching Python f"{v:,.2f}".
func FormatComma2(v float64) string {
	neg := v < 0
	absVal := math.Abs(v)
	raw := strconv.FormatFloat(absVal, 'f', 2, 64)
	dot := strings.IndexByte(raw, '.')
	intPart := raw
	fracPart := "00"
	if dot >= 0 {
		intPart = raw[:dot]
		fracPart = raw[dot+1:]
	}
	n := len(intPart)
	var b strings.Builder
	if neg && raw != "0.00" {
		b.WriteByte('-')
	}
	for i := 0; i < n; i++ {
		if i > 0 && (n-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteByte(intPart[i])
	}
	b.WriteByte('.')
	b.WriteString(fracPart)
	return b.String()
}
