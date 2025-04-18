package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ip2country-api/pkg/ip2country"
)

// MockService is a mock implementation of the ip2country.Service interface
type MockService struct {
	LookupIPFunc func(ip string) (*ip2country.Result, error)
}

func (m *MockService) LookupIP(ip string) (*ip2country.Result, error) {
	return m.LookupIPFunc(ip)
}

func TestFindCountryHandler(t *testing.T) {
	tests := []struct {
		name            string
		ip              string
		mockLookupIP    func(ip string) (*ip2country.Result, error)
		expectedStatus  int
		expectedMessage string
	}{
		{
			name: "successful lookup",
			ip:   "192.168.1.1",
			mockLookupIP: func(ip string) (*ip2country.Result, error) {
				return &ip2country.Result{Country: "US", City: "New York"}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "missing ip parameter",
			ip:   "",
			mockLookupIP: func(ip string) (*ip2country.Result, error) {
				return nil, nil
			},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Missing 'ip' parameter",
		},
		{
			name: "invalid ip address",
			ip:   "invalid-ip",
			mockLookupIP: func(ip string) (*ip2country.Result, error) {
				return nil, ip2country.ErrInvalidIP
			},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid IP address",
		},
		{
			name: "ip not found",
			ip:   "10.0.0.1",
			mockLookupIP: func(ip string) (*ip2country.Result, error) {
				return nil, ip2country.ErrIPNotFound
			},
			expectedStatus:  http.StatusNotFound,
			expectedMessage: "IP address not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := &MockService{
				LookupIPFunc: tt.mockLookupIP,
			}

			// Create request
			req, err := http.NewRequest("GET", "/v1/find-country", nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.ip != "" {
				q := req.URL.Query()
				q.Add("ip", tt.ip)
				req.URL.RawQuery = q.Encode()
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Create handler
			handler := FindCountryHandler(mockService)

			// Serve request
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			// Check response message if expected
			if tt.expectedMessage != "" {
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
