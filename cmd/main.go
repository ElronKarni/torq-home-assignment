package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ip2country-api/internal/config"
	"ip2country-api/internal/ip2country"
	"ip2country-api/internal/routes"
	"ip2country-api/pkg/ratelimit"
)

// setupServer initializes all components and returns the HTTP server
// This function is extracted to make it testable
func setupServer() (*http.Server, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	// init IP2country service with just the BackendConfig
	ip2countryService, err := ip2country.NewService(cfg.IP2Country)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize IP2Country service: %v", err)
	}

	// Initialize rate limiter
	limiter := ratelimit.NewLimiter(cfg.RateLimit)

	// Set up HTTP routes with middleware
	handler := routes.RegisterRoutes(ip2countryService, limiter, cfg.AllowedOrigins)

	// Create HTTP server
	addr := fmt.Sprintf(":%d", cfg.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      handler, // Use our custom handler with middleware
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Rate limit: %d requests per second", cfg.RateLimit)
	log.Printf("IP2Country backend: %#v", cfg.IP2Country)
	log.Printf("CORS allowed origins: %v", cfg.AllowedOrigins)

	return server, nil
}

func main() {
	server, err := setupServer()
	if err != nil {
		log.Fatalf("Server setup failed: %v", err)
	}

	// Create channel to listen for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped")
}
