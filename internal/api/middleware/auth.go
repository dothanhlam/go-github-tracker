package middleware

import (
	"net/http"
	"strings"

	"github.com/dothanhlam/go-github-tracker/internal/api/response"
)

// APIKeyAuth creates middleware that validates API keys
func APIKeyAuth(validKeys []string) func(http.Handler) http.Handler {
	// Create a map for faster lookup
	keyMap := make(map[string]bool)
	for _, key := range validKeys {
		if key != "" {
			keyMap[key] = true
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for health check
			if r.URL.Path == "/api/v1/health" {
				next.ServeHTTP(w, r)
				return
			}

			// Get API key from header
			apiKey := r.Header.Get("X-API-Key")
			apiKey = strings.TrimSpace(apiKey)

			// Validate API key
			if apiKey == "" {
				response.Unauthorized(w, "API key is required")
				return
			}

			if !keyMap[apiKey] {
				response.Unauthorized(w, "Invalid API key")
				return
			}

			// API key is valid, proceed
			next.ServeHTTP(w, r)
		})
	}
}
