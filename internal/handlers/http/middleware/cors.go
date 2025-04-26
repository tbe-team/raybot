package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

// Cors is a middleware handler that sets the CORS configuration.
func Cors() func(http.Handler) http.Handler {
	opts := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           86400,
	}

	return cors.Handler(opts)
}
