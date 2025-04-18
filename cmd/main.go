package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"ip2country-api/pkg/config"
	"ip2country-api/pkg/ip2country"
	"ip2country-api/pkg/ratelimit"
)

// RateLimiter defines the interface for rate limiting
type RateLimiter interface {
	Allow() error
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize IP2Country service using the factory function
	ipService, err := ip2country.NewService(cfg.DataPath, cfg.IP2CountryDBType)
	if err != nil {
		log.Fatalf("Failed to initialize IP2Country service: %v", err)
	}

	// Initialize rate limiter
	limiter := ratelimit.NewLimiter(cfg.RateLimit)

	// Set up HTTP handler
	http.HandleFunc("/v1/find-country", createFindCountryHandler(ipService, limiter))

	// Start HTTP server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting server on %s", addr)
	log.Printf("Rate limit: %d requests per second", cfg.RateLimit)
	log.Printf("IP2Country data path: %s", cfg.DataPath)
	log.Printf("IP2Country database type: %s", cfg.IP2CountryDBType)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createFindCountryHandler(ipService ip2country.Service, limiter RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Apply rate limiting
		if err := limiter.Allow(); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"error": "Too many requests"})
			return
		}

		// Extract IP from query parameter
		ip := r.URL.Query().Get("ip")
		if ip == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Missing 'ip' parameter"})
			return
		}

		// Look up IP information
		result, err := ipService.LookupIP(ip)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")

			// Handle specific error cases
			switch {
			case errors.Is(err, ip2country.ErrInvalidIP):
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid IP address"})
			case errors.Is(err, ip2country.ErrIPNotFound):
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "IP address not found"})
			default:
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Failed to look up IP information"})
			}
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}
