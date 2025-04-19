package middleware

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeaders(t *testing.T) {
	// Create a test handler that just returns 200 OK
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Apply our middleware to the test handler
	secureHandler := SecurityHeaders(testHandler)

	// Create a test request with HTTPS to trigger HSTS
	req, err := http.NewRequest("GET", "https://example.com/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set TLS connection info to ensure HSTS is applied
	req.TLS = &tls.ConnectionState{}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	secureHandler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check security headers
	headers := map[string]string{
		"X-Frame-Options":           "DENY",
		"X-Content-Type-Options":    "nosniff",
		"X-Xss-Protection":          "1; mode=block",
		"Content-Security-Policy":   "default-src 'self'",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains; preload",
		"Referrer-Policy":           "same-origin",
	}

	for name, expected := range headers {
		if actual := rr.Header().Get(name); actual != expected {
			t.Errorf("%s header: got %s, want %s", name, actual, expected)
		}
	}
}
