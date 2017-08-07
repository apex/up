// Package logs provides HTTP request and response logging.
package logs

import (
	"net/http"
	"time"

	"github.com/apex/log"

	"github.com/apex/up"
	"github.com/apex/up/internal/logs"
)

// TODO: optional verbose mode with req/res header etc?

// log context.
var ctx = logs.Plugin("logs")

// response wrapper.
type response struct {
	http.ResponseWriter
	written int
	code    int
}

// Write implementation.
func (r *response) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.written += n
	return n, err
}

// WriteHeader implementation.
func (r *response) WriteHeader(code int) {
	r.code = code
	r.ResponseWriter.WriteHeader(code)
}

// New logs handler.
func New(c *up.Config, next http.Handler) (http.Handler, error) {
	if c.Logs.Disable {
		return next, nil
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		res := &response{ResponseWriter: w, code: 200}

		next.ServeHTTP(res, r)

		c := ctx.WithFields(log.Fields{
			"stage":    r.Header.Get("X-Stage"),
			"id":       r.Header.Get("X-Request-Id"),
			"method":   r.Method,
			"path":     r.URL.Path,
			"query":    r.URL.Query().Encode(),
			"duration": int(time.Since(start) / time.Millisecond),
			"size":     res.written,
			"ip":       r.RemoteAddr,
			"status":   res.code,
		})

		switch {
		case res.code >= 500:
			c.Error("response")
		case res.code >= 400:
			c.Warn("response")
		default:
			c.Info("response")
		}
	})

	return h, nil
}
