package middleware

import (
	"encoding/json"
	"net/http"
)

// RateLimiter defines the interface for rate limiting
type RateLimiter interface {
	Allow() error
}

// RateLimit creates a middleware that applies rate limiting to all requests
func RateLimit(limiter RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Apply rate limiting
			if err := limiter.Allow(); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]string{"error": "Too many requests"})
				return
			}

			// Pass to the next handler if rate limit not exceeded
			next.ServeHTTP(w, r)
		})
	}
}
