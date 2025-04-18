package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockRateLimiter implements the RateLimiter interface for testing
type MockRateLimiter struct {
	AllowFunc func() error
}

func (m *MockRateLimiter) Allow() error {
	return m.AllowFunc()
}

func TestRateLimit(t *testing.T) {
	tests := []struct {
		name            string
		allowResult     error
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:           "allowed request",
			allowResult:    nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:            "rate limited request",
			allowResult:     errors.New("rate limit exceeded"),
			expectedStatus:  http.StatusTooManyRequests,
			expectedMessage: "Too many requests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock rate limiter
			mockLimiter := &MockRateLimiter{
				AllowFunc: func() error {
					return tt.allowResult
				},
			}

			// Create a test handler that always returns 200 OK
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Apply the middleware
			handler := RateLimit(mockLimiter)(nextHandler)

			// Create request
			req, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve request
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("middleware returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			// Check response message if rate limited
			if tt.allowResult != nil {
				var response map[string]string
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("could not parse response body: %v", err)
				}
				if msg, exists := response["error"]; !exists || msg != tt.expectedMessage {
					t.Errorf("expected error message %q, got %q", tt.expectedMessage, msg)
				}
			}
		})
	}
}
