package middleware

import (
	"context"
	"net/http"
	"time"
)

// RequestTimeoutMiddleware adds a 30-second timeout to each request
func RequestTimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a context with 30-second timeout
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		// Create a new request with the timeout context
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
