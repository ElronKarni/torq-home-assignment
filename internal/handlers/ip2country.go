package handlers

import (
	"errors"
	"net/http"

	"ip2country-api/internal/ip2country"
	"ip2country-api/internal/utils"
)

// FindCountryHandler creates an HTTP handler function for the find-country endpoint
func FindCountryHandler(ip2countryService ip2country.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract IP from query parameter
		ip := r.URL.Query().Get("ip")
		if ip == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing 'ip' parameter"})
			return
		}

		// Look up IP information
		result, err := ip2countryService.LookupIP(ip)
		if err != nil {
			// Handle specific error cases
			switch {
			case errors.Is(err, ip2country.ErrInvalidIP):
				utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid IP address"})
			case errors.Is(err, ip2country.ErrIPNotFound):
				utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "IP address not found"})
			default:
				utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to look up IP information"})
			}
			return
		}

		// Return JSON response
		utils.WriteJSON(w, http.StatusOK, result)

	}
}
