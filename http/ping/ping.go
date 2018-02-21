// Package ping provides the /_ping no-op route.
package ping

import (
	"fmt"
	"net/http"
	"time"

	"github.com/apex/up"
)

// New ping handler.
func New(c *up.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/_ping" {
			time.Sleep(time.Second)
			fmt.Fprintln(w, ":)")
			return
		}

		next.ServeHTTP(w, r)
	})
}
