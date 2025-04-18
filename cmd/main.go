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

// For testing purposes - allows tests to override server start behavior
var serverListenAndServe = func(server *http.Server) error {
	return server.ListenAndServe()
}

// setupServer initializes all components and returns the HTTP server
// This function is extracted to make it testable
func setupServer() (*http.Server, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	// Initialize IP2Country service using the factory function
	ipService, err := ip2country.NewService(cfg.DataPath, cfg.IP2CountryDBType)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize IP2Country service: %v", err)
	}

	// Initialize rate limiter
	limiter := ratelimit.NewLimiter(cfg.RateLimit)

	// Set up HTTP routes
	handlers.RegisterRoutes(ipService, limiter)

	// Create HTTP server
	addr := fmt.Sprintf(":%d", cfg.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: nil, // Use DefaultServeMux
	}

	log.Printf("Server configured on %s", addr)
	log.Printf("Rate limit: %d requests per second", cfg.RateLimit)
	log.Printf("IP2Country data path: %s", cfg.DataPath)
	log.Printf("IP2Country database type: %s", cfg.IP2CountryDBType)

	return server, nil
}

func main() {
	server, err := setupServer()
	if err != nil {
		log.Fatalf("Server setup failed: %v", err)
	}

	log.Printf("Starting server on %s", server.Addr)
	if err := serverListenAndServe(server); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
