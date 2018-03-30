// Package robots provides a way of dealing with robots exclusion protocol
package robots

import (
	"net/http"
	"os"

	"github.com/apex/up"
)

// New robots middleware.
func New(c *up.Config, next http.Handler) http.Handler {
	if os.Getenv("UP_STAGE") == "production" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Robots-Tag", "none")
		next.ServeHTTP(w, r)
	})
}
