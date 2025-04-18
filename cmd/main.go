package main

import (
	"fmt"
	"log"
	"net/http"

	"ip2country-api/pkg/config"
	"ip2country-api/pkg/handlers"
	"ip2country-api/pkg/ip2country"
	"ip2country-api/pkg/ratelimit"
)

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

	// Set up HTTP routes
	handlers.RegisterRoutes(ipService, limiter)

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
