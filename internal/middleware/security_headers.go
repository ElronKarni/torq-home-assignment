package middleware

import (
	"net/http"

	"github.com/unrolled/secure"
)

// SecurityHeaders adds security headers to all responses
func SecurityHeaders(next http.Handler) http.Handler {
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		STSSeconds:            31536000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		ReferrerPolicy:        "same-origin",
		IsDevelopment:         false, // Ensure headers are set even in dev environment
	})

	return secureMiddleware.Handler(next)
}
