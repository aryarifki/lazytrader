package spa

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var legacyRedirects = map[string]string{
	"/nyangkut":         "/ownership",
	"/whales":           "/group",
	"/backtest":         "/signal?tab=backtest",
	"/market-report":    "/admin/market-report",
	"/report-generator": "/admin/market-report",
}

// Handler serves built Vite static assets and falls back to index.html for SPA routes.
func Handler(distDir string) http.Handler {
	absDist, _ := filepath.Abs(distDir)
	fs := http.FileServer(http.Dir(absDist))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if target, ok := legacyRedirects[p]; ok {
			http.Redirect(w, r, target, http.StatusFound)
			return
		}
		if strings.HasPrefix(p, "/api/") {
			http.NotFound(w, r)
			return
		}

		clean := filepath.Clean(strings.TrimPrefix(p, "/"))
		candidate := filepath.Join(absDist, clean)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			if strings.HasPrefix(p, "/assets/") {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			fs.ServeHTTP(w, r)
			return
		}

		indexPath := filepath.Join(absDist, "index.html")
		if _, err := os.Stat(indexPath); err != nil {
			http.Error(w, "Frontend dist not built yet", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(w, r, indexPath)
	})
}
