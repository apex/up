// Package inject provides script and style injection.
package inject

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/apex/up"
	"github.com/apex/up/internal/inject"
)

// response wrapper.
type response struct {
	http.ResponseWriter
	rules  inject.Rules
	body   bytes.Buffer
	header bool
	ignore bool
	code   int
}

// Write implementation.
func (r *response) Write(b []byte) (int, error) {
	if !r.header {
		r.WriteHeader(200)
		return r.Write(b)
	}

	return r.body.Write(b)
}

// WriteHeader implementation.
func (r *response) WriteHeader(code int) {
	r.header = true
	w := r.ResponseWriter
	kind := w.Header().Get("Content-Type")
	r.ignore = !strings.HasPrefix(kind, "text/html") || code >= 300
	r.code = code
}

// end injects if necessary.
func (r *response) end() {
	w := r.ResponseWriter

	if r.ignore {
		w.WriteHeader(r.code)
		r.body.WriteTo(w)
		return
	}

	body := r.rules.Apply(r.body.String())
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	io.WriteString(w, body)
}

// New inject handler.
func New(c *up.Config, next http.Handler) (http.Handler, error) {
	if len(c.Inject) == 0 {
		return next, nil
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := &response{ResponseWriter: w, rules: c.Inject}
		next.ServeHTTP(res, r)
		res.end()
	})

	return h, nil
}
