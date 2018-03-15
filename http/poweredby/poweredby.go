// Package poweredby provides nothing :).
package poweredby

import (
	"net/http"
)

// New powered-by middleware.
func New(name string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", name)
		next.ServeHTTP(w, r)
	})
}
