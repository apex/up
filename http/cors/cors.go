// Package cors provides CORS support.
package cors

import (
	"net/http"

	"github.com/rs/cors"

	"github.com/apex/up"
	"github.com/apex/up/config"
)

// New CORS handler.
func New(c *up.Config, next http.Handler) http.Handler {
	if c.CORS == nil {
		return next
	}

	return cors.New(options(c.CORS)).Handler(next)
}

// options returns the canonical options.
func options(c *config.CORS) cors.Options {
	return cors.Options{
		AllowedOrigins:   c.AllowedOrigins,
		AllowedMethods:   c.AllowedMethods,
		AllowedHeaders:   c.AllowedHeaders,
		ExposedHeaders:   c.ExposedHeaders,
		AllowCredentials: c.AllowCredentials,
		MaxAge:           c.MaxAge,
		Debug:            c.Debug,
	}
}
