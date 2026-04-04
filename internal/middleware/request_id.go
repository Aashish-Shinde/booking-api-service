package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"
)

// RequestIDKey is the context key for request ID
type RequestIDKey struct{}

// RequestTimeoutKey is the context key for request timeout
type RequestTimeoutKey struct{}

// RequestIDMiddleware adds a request ID to each request context
func RequestIDMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate request ID
			b := make([]byte, 8)
			rand.Read(b)
			requestID := hex.EncodeToString(b)

			// Add request ID to context
			ctx := context.WithValue(r.Context(), RequestIDKey{}, requestID)

			// Add request timeout (60 seconds default)
			timeout := 60 * time.Second
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			// Log request start
			log.Info("request started",
				slog.String("request_id", requestID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)

			// Create a response writer wrapper to capture status code
			wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

			// Call next handler
			next.ServeHTTP(wrapper, r.WithContext(ctx))

			// Log request end
			log.Info("request completed",
				slog.String("request_id", requestID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", wrapper.statusCode),
			)
		})
	}
}

// responseWriterWrapper wraps http.ResponseWriter to capture status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// WriteHeader captures the status code
func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	if !w.written {
		w.statusCode = statusCode
		w.written = true
		w.ResponseWriter.WriteHeader(statusCode)
	}
}

// Write wraps the Write method
func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	if !w.written {
		w.statusCode = http.StatusOK
		w.written = true
	}
	return w.ResponseWriter.Write(b)
}

// GetRequestID extracts the request ID from context
func GetRequestID(ctx context.Context) string {
	if id := ctx.Value(RequestIDKey{}); id != nil {
		if idStr, ok := id.(string); ok {
			return idStr
		}
	}
	return "unknown"
}
