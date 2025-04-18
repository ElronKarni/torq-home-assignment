package handlers

import (
	"net/http"

	"ip2country-api/pkg/ip2country"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(
	ipService ip2country.Service,
	limiter RateLimiter,
) {
	// IP-to-country API endpoints
	http.HandleFunc("/v1/find-country", FindCountryHandler(ipService, limiter))

	// Additional routes can be added here as the API grows
}
