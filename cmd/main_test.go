package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"ip2country-api/pkg/ip2country"
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

func (m *MockLimiter) Allow() error {
	if !m.shouldAllow {
		return ratelimit.ErrRateLimitExceeded
	}
	return nil
}

func TestFindCountryHandler(t *testing.T) {
	testCases := []struct {
		name           string
		ipParam        string
		mockLookupFunc func(ip string) (*ip2country.Result, error)
		expectStatus   int
		expectBody     map[string]string
		limitExceeded  bool
	}{
		{
			name:    "Valid IP - Found",
			ipParam: "1.1.1.1",
			mockLookupFunc: func(ip string) (*ip2country.Result, error) {
				return &ip2country.Result{Country: "Australia", City: "Sydney"}, nil
			},
			expectStatus:  http.StatusOK,
			expectBody:    map[string]string{"country": "Australia", "city": "Sydney"},
			limitExceeded: false,
		},
		{
			name:    "Valid IP - Not Found",
			ipParam: "9.9.9.9",
			mockLookupFunc: func(ip string) (*ip2country.Result, error) {
				return nil, ip2country.ErrIPNotFound
			},
			expectStatus:  http.StatusNotFound,
			expectBody:    map[string]string{"error": "IP address not found"},
			limitExceeded: false,
		},
		{
			name:    "Invalid IP",
			ipParam: "not-an-ip",
			mockLookupFunc: func(ip string) (*ip2country.Result, error) {
				return nil, ip2country.ErrInvalidIP
			},
			expectStatus:  http.StatusBadRequest,
			expectBody:    map[string]string{"error": "Invalid IP address"},
			limitExceeded: false,
		},
		{
			name:    "Missing IP Parameter",
			ipParam: "",
			mockLookupFunc: func(ip string) (*ip2country.Result, error) {
				t.Fatal("LookupIP should not be called when IP parameter is missing")
				return nil, nil
			},
			expectStatus:  http.StatusBadRequest,
			expectBody:    map[string]string{"error": "Missing 'ip' parameter"},
			limitExceeded: false,
		},
		{
			name:    "Rate Limit Exceeded",
			ipParam: "1.1.1.1",
			mockLookupFunc: func(ip string) (*ip2country.Result, error) {
				t.Fatal("LookupIP should not be called when rate limit is exceeded")
				return nil, nil
			},
			expectStatus:  http.StatusTooManyRequests,
			expectBody:    map[string]string{"error": "Too many requests"},
			limitExceeded: true,
		},
		{
			name:    "Server Error",
			ipParam: "1.1.1.1",
			mockLookupFunc: func(ip string) (*ip2country.Result, error) {
				return nil, errors.New("unexpected server error")
			},
			expectStatus:  http.StatusInternalServerError,
			expectBody:    map[string]string{"error": "Failed to look up IP information"},
			limitExceeded: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock IP service
			mockService := &MockIPService{
				lookupFunc: tc.mockLookupFunc,
			}

			// Create a mock rate limiter
			mockLimiter := &MockLimiter{
				shouldAllow: !tc.limitExceeded,
			}

			// Create handler
			handler := createFindCountryHandler(mockService, mockLimiter)

			// Create test request
			req, err := http.NewRequest(http.MethodGet, "/v1/find-country", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Add IP parameter if it's not empty
			if tc.ipParam != "" {
				q := req.URL.Query()
				q.Add("ip", tc.ipParam)
				req.URL.RawQuery = q.Encode()
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tc.expectStatus {
				t.Errorf("Handler returned wrong status code: got %v, want %v",
					rr.Code, tc.expectStatus)
			}

			// Check response body
			var body map[string]string
			if err := json.NewDecoder(rr.Body).Decode(&body); err != nil && err != io.EOF {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			// Compare body contents
			for k, v := range tc.expectBody {
				if body[k] != v {
					t.Errorf("Handler returned wrong body: got %v, want %v for key %q",
						body[k], v, k)
				}
			}
		})
	}
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

	mux.HandleFunc("/v1/find-country", createFindCountryHandler(mockService, mockLimiter))

	ts := httptest.NewServer(mux)
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
