package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// CORS returns a middleware that handles CORS
func CORS(origins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Create a CORS handler with our settings
		corsMiddleware := cors.New(cors.Options{
			AllowedOrigins:   origins,
			AllowedMethods:   []string{"GET", "OPTIONS"}, // Only GET and OPTIONS (required for preflight) since the service only handles GET requests
			AllowedHeaders:   []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           43200, // 12 hours in seconds
			Debug:            false,
		})

		// Use the cors Handler which automatically handles OPTIONS preflight
		return corsMiddleware.Handler(next)
	}
}
