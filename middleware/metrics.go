package middleware

import (
	"net/http"
	"sync/atomic"
)

func MetricsMiddleware(hits *atomic.Int32) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hits.Add(1)
			next.ServeHTTP(w, r)
		})
	}
}
