package handlers

import (
	"net/http"

	"ip2country-api/pkg/ip2country"
	"ip2country-api/pkg/middleware"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(
	ipService ip2country.Service,
	limiter middleware.RateLimiter,
) http.Handler {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// IP-to-country API endpoints
	mux.HandleFunc("/v1/find-country", FindCountryHandler(ipService))

	// Additional routes can be added here as the API grows

	// Apply middlewares to all routes
	// Start with the innermost middleware
	var handler http.Handler = mux

	// Apply rate limiting middleware
	handler = middleware.RateLimit(limiter)(handler)

	return handler
}
