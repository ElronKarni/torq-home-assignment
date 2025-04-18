package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"ip2country-api/internal/ip2country"
)

// FindCountryHandler creates an HTTP handler function for the find-country endpoint
func FindCountryHandler(ipService ip2country.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract IP from query parameter
		ip := r.URL.Query().Get("ip")
		if ip == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Missing 'ip' parameter"})
			return
		}

		// Look up IP information
		result, err := ipService.LookupIP(ip)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")

			// Handle specific error cases
			switch {
			case errors.Is(err, ip2country.ErrInvalidIP):
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid IP address"})
			case errors.Is(err, ip2country.ErrIPNotFound):
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "IP address not found"})
			default:
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Failed to look up IP information"})
			}
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}
