package portal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPortalExcludesBei1Persen(t *testing.T) {
	r := Routes()
	for _, path := range []string{"/status", "/entitlements"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s returned %d", path, rec.Code)
		}
		if strings.Contains(strings.ToLower(rec.Body.String()), "bei1persen") {
			t.Fatalf("%s unexpectedly contains bei1persen: %s", path, rec.Body.String())
		}
	}
}
