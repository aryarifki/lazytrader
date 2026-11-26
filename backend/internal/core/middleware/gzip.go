package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	writer io.Writer
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.writer.Write(b)
}

// Gzip compresses responses when the client accepts gzip encoding.
func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, ".png") || strings.HasSuffix(r.URL.Path, ".jpg") || strings.HasSuffix(r.URL.Path, ".ico") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Add("Vary", "Accept-Encoding")
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, writer: gz}, r)
	})
}
