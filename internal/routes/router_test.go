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
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "find country with valid IP",
			path:           "/v1/find-country?ip=192.168.1.1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "find country with missing IP",
			path:           "/v1/find-country",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "non-existent route",
			path:           "/not-found",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve request
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
