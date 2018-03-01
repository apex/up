// Package logs provides HTTP request and response logging.
package logs

import (
	"net/http"
	"strconv"
	"time"

	"github.com/apex/log"

	"github.com/apex/up"
	"github.com/apex/up/internal/logs"
	"github.com/apex/up/internal/util"
)

// TODO: optional verbose mode with req/res header etc?

// log context.
var ctx = logs.Plugin("logs")

// response wrapper.
type response struct {
	http.ResponseWriter
	written  int
	code     int
	duration time.Duration
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
		ctx := logContext(r)
		logRequest(ctx, r)

		start := time.Now()
		res := &response{ResponseWriter: w, code: 200}
		next.ServeHTTP(res, r)
		res.duration = time.Since(start)

		logResponse(ctx, res, r)
	})

	return h, nil
}

// logContext returns the common log context for a request.
func logContext(r *http.Request) log.Interface {
	return ctx.WithFields(log.Fields{
		"id":     r.Header.Get("X-Request-Id"),
		"method": r.Method,
		"path":   r.URL.Path,
		"query":  r.URL.Query().Encode(),
		"ip":     r.RemoteAddr,
	})
}

// logRequest logs the request.
func logRequest(ctx log.Interface, r *http.Request) {
	if s := r.Header.Get("Content-Length"); s != "" {
		n, err := strconv.Atoi(s)
		if err == nil {
			ctx = ctx.WithField("size", n)
		}
	}

	ctx.Info("request")
}

// logResponse logs the response.
func logResponse(ctx log.Interface, res *response, r *http.Request) {
	ctx = ctx.WithFields(log.Fields{
		"duration": util.Milliseconds(res.duration),
		"size":     res.written,
		"status":   res.code,
	})

	switch {
	case res.code >= 500:
		ctx.Error("response")
	case res.code >= 400:
		ctx.Warn("response")
	default:
		ctx.Info("response")
	}
}
