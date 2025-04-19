package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"ip2country-api/internal/handlers"
	"ip2country-api/internal/ip2country"
	"ip2country-api/internal/middleware"
	"ip2country-api/pkg/ratelimit"
)

// MockIPService is a mock implementation of the ip2country.Service interface for testing
type MockIPService struct {
	lookupFunc func(ip string) (*ip2country.Result, error)
}

func (m *MockIPService) LookupIP(ip string) (*ip2country.Result, error) {
	return m.lookupFunc(ip)
}

// MockLimiter is a mock implementation of a rate limiter for testing
type MockLimiter struct {
	shouldAllow bool
}

// Allow implements the middleware.RateLimiter interface
func (m *MockLimiter) Allow() error {
	if !m.shouldAllow {
		return ratelimit.ErrRateLimitExceeded
	}
	return nil
}

// TestSetupServer tests the server setup functionality
func TestSetupServer(t *testing.T) {
	// Save original environment and restore after test
	origDataPath := os.Getenv("CSV_DATA_PATH")
	origPort := os.Getenv("PORT")
	origDBType := os.Getenv("IP2COUNTRY_DB_TYPE")

	// Set test environment variables
	testDataDir, err := os.MkdirTemp("", "ip2country-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(testDataDir)

	// Create a test data file
	testDataFile := testDataDir + "/ip2country.csv"
	testData := "1.1.1.1,Sydney,Australia\n8.8.8.8,Mountain View,United States"
	if err := os.WriteFile(testDataFile, []byte(testData), 0644); err != nil {
		t.Fatalf("Failed to write test data file: %v", err)
	}

	// Set environment variables for the test
	os.Setenv("CSV_DATA_PATH", testDataFile)
	os.Setenv("IP2COUNTRY_DB_TYPE", "csv")
	os.Setenv("PORT", "8081") // Use a different port than default
	defer func() {
		os.Setenv("CSV_DATA_PATH", origDataPath)
		os.Setenv("PORT", origPort)
		os.Setenv("IP2COUNTRY_DB_TYPE", origDBType)
	}()

	// Clear DefaultServeMux to avoid conflicts from previous tests
	http.DefaultServeMux = http.NewServeMux()

	// Test the setupServer function
	server, err := setupServer()
	if err != nil {
		t.Fatalf("setupServer() failed: %v", err)
	}

	// Verify the server was set up correctly
	if server == nil {
		t.Fatal("setupServer() returned nil server")
	}

	if server.Addr != ":8081" {
		t.Errorf("setupServer() configured wrong address: got %s, want :8081", server.Addr)
	}
}

// TestSetupServerErrors tests the error cases in setupServer
func TestSetupServerErrors(t *testing.T) {
	// Test case 1: Invalid configuration
	t.Run("ConfigurationError", func(t *testing.T) {
		// Save original environment and restore after test
		origPort := os.Getenv("PORT")

		// Set invalid PORT to trigger config error
		os.Setenv("PORT", "not-a-number")
		defer os.Setenv("PORT", origPort)

		// Test the setupServer function
		server, err := setupServer()

		// Verify error is returned
		if err == nil {
			t.Fatal("setupServer() should have failed with invalid PORT")
		}
		if server != nil {
			t.Fatal("setupServer() should return nil server on error")
		}
	})

	// Test case 2: Invalid data path
	t.Run("InvalidDataPath", func(t *testing.T) {
		// Save original environment and restore after test
		origDataPath := os.Getenv("CSV_DATA_PATH")
		origDBType := os.Getenv("IP2COUNTRY_DB_TYPE")

		// Set CSV as the database type
		os.Setenv("IP2COUNTRY_DB_TYPE", "csv")

		// Set non-existent data path
		os.Setenv("CSV_DATA_PATH", "/path/that/does/not/exist.csv")

		defer func() {
			os.Setenv("CSV_DATA_PATH", origDataPath)
			os.Setenv("IP2COUNTRY_DB_TYPE", origDBType)
		}()

		// Test the setupServer function
		server, err := setupServer()

		// Verify error is returned
		if err == nil {
			t.Fatal("setupServer() should have failed with invalid data path")
		}
		if server != nil {
			t.Fatal("setupServer() should return nil server on error")
		}
	})
}

// TestServerIntegration tests that the server starts and listens on the correct port
// This is a basic smoke test and won't actually start the server to avoid port conflicts
func TestServerIntegration(t *testing.T) {
	// Test that the handler function is correctly registered
	mux := http.NewServeMux()
	mockService := &MockIPService{
		lookupFunc: func(ip string) (*ip2country.Result, error) {
			return &ip2country.Result{Country: "Test", City: "Test"}, nil
		},
	}
	mockLimiter := &MockLimiter{shouldAllow: true}

	// Set up the handler without rate limiter
	mux.HandleFunc("/v1/find-country", handlers.FindCountryHandler(mockService))

	// Apply rate limit middleware
	handler := middleware.RateLimit(mockLimiter)(mux)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Send a test request to verify the endpoint works
	resp, err := http.Get(ts.URL + "/v1/find-country?ip=1.1.1.1")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Wrong status code: got %v, want %v", resp.StatusCode, http.StatusOK)
	}
}
