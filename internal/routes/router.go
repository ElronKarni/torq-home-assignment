package routes

import (
	"net/http"

	"ip2country-api/internal/handlers"
	"ip2country-api/internal/ip2country"
	"ip2country-api/internal/middleware"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(
	ip2countryService ip2country.Service,
	limiter middleware.RateLimiter,
) http.Handler {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// IP-to-country API endpoints
	mux.HandleFunc("/v1/find-country", handlers.FindCountryHandler(ip2countryService))

	// Additional routes can be added here as the API grows

	// Apply middlewares to all routes
	// Start with the innermost middleware
	var handler http.Handler = mux

	// Apply rate limiting middleware
	handler = middleware.RateLimit(limiter)(handler)

	return handler
}
