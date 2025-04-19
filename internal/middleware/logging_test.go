package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	// Create a simple test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil) // Reset logger

	// Apply logger middleware
	handler := Logger(testHandler)

	// Create and send a test request
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Verify response is passed through
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Body.String() != "OK" {
		t.Errorf("Expected body %q, got %q", "OK", rec.Body.String())
	}

	// Verify log contains essential information
	logOutput := buf.String()
	requiredInfo := []string{"GET", "/test", "200"} // Basic info that should be in logs
	for _, info := range requiredInfo {
		if !strings.Contains(logOutput, info) {
			t.Errorf("Log missing required info %q, log: %s", info, logOutput)
		}
	}
}

// TestStatusRecorderDefaultStatus tests that the statusRecorder sets default status to 200 OK
// when WriteHeader is not called explicitly
func TestStatusRecorderDefaultStatus(t *testing.T) {
	// Create a test handler that doesn't call WriteHeader explicitly
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write directly without setting status code
		w.Write([]byte("OK"))
	})

	// Create a response recorder and wrap it in statusRecorder
	rec := httptest.NewRecorder()
	sr := &statusRecorder{ResponseWriter: rec}

	// Create and send a test request
	req := httptest.NewRequest("GET", "/test", nil)
	testHandler.ServeHTTP(sr, req)

	// Verify default status is set to 200 OK
	if sr.Status != http.StatusOK {
		t.Errorf("Expected default status code %d, got %d", http.StatusOK, sr.Status)
	}

	// Verify bytes are counted correctly
	if sr.Bytes != 2 { // "OK" is 2 bytes
		t.Errorf("Expected bytes count %d, got %d", 2, sr.Bytes)
	}
}
