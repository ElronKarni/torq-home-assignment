package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {

	t.Run("valid data", func(t *testing.T) {
		rr := httptest.NewRecorder()

		WriteJSON(rr, http.StatusOK, "hello")

		if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
		}

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d OK, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("invalid data", func(t *testing.T) {
		rr := httptest.NewRecorder()

		WriteJSON(rr, http.StatusOK, func() {}) // func() is not encodable

		if rr.Body.Len() != 0 {
			t.Errorf("expected empty body on encoding error, got: %q", rr.Body.String())
		}

	})

}
