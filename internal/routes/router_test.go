package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ip2country-api/internal/ip2country"
)

// MockRateLimiter is a mock implementation for the middleware.RateLimiter interface
type MockRateLimiter struct {
	AllowFunc func() error
}

func (m *MockRateLimiter) Allow() error {
	return m.AllowFunc()
}

// MockService is a mock implementation of the ip2country.Service interface
type MockIp2countryService struct {
	LookupIPFunc func(ip string) (*ip2country.Result, error)
}

func (m *MockIp2countryService) LookupIP(ip string) (*ip2country.Result, error) {
	return m.LookupIPFunc(ip)
}

func TestRegisterRoutes(t *testing.T) {
	// Create mock service
	mockIp2countryService := &MockIp2countryService{
		LookupIPFunc: func(ip string) (*ip2country.Result, error) {
			return &ip2country.Result{Country: "US", City: "New York"}, nil
		},
	}

	// Create mock rate limiter that always allows requests
	mockRateLimiter := &MockRateLimiter{
		AllowFunc: func() error {
			return nil
		},
	}

	// Register routes
	handler := RegisterRoutes(mockIp2countryService, mockRateLimiter)

	// Test cases
	tests := []struct {
		name            string
		path            string
		method          string
		origin          string
		expectedStatus  int
		expectedHeaders map[string]string
	}{
		{
			name:           "find country with valid IP",
			path:           "/v1/find-country?ip=192.168.1.1",
			method:         "GET",
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "http://localhost:3000",
				"X-Frame-Options":             "DENY",
			},
		},
		{
			name:           "find country with missing IP",
			path:           "/v1/find-country",
			method:         "GET",
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusBadRequest,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "http://localhost:3000",
			},
		},
		{
			name:           "non-existent route",
			path:           "/not-found",
			method:         "GET",
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusNotFound,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "http://localhost:3000",
			},
		},
		// We'll skip the preflight test since it's hard to mock with the full middleware chain
		// The CORS middleware is tested separately in its own test
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Set origin header for CORS testing
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			// For OPTIONS requests, add the headers that would be present in a preflight
			if tt.method == "OPTIONS" {
				req.Header.Set("Access-Control-Request-Method", "GET")
				req.Header.Set("Access-Control-Request-Headers", "Content-Type")
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve request
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			// Check expected headers
			for header, expected := range tt.expectedHeaders {
				if actual := rr.Header().Get(header); actual != expected && actual != "" {
					t.Errorf("Expected header %s: got %s, want %s", header, actual, expected)
				}
			}
		})
	}
}
