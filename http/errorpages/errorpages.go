// Package errorpages provides default and customizable
// error pages, via error.html, 5xx.html, or 500.html
// for example.
package errorpages

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
	accept "github.com/timewasted/go-accept-headers"

	"github.com/apex/up"
	"github.com/apex/up/internal/errorpage"
	"github.com/apex/up/internal/logs"
	"github.com/apex/up/internal/util"
)

// log context.
var ctx = logs.Plugin("errorpages")

// response wrapper.
type response struct {
	http.ResponseWriter
	config *up.Config
	pages  errorpage.Pages
	header bool
	ignore bool
}

// WriteHeader implementation.
func (r *response) WriteHeader(code int) {
	w := r.ResponseWriter

	r.header = true
	page := r.pages.Match(code)

	if page == nil {
		ctx.Debugf("did not match %d", code)
		w.WriteHeader(code)
		return
	}

	ctx.Debugf("matched %d with %q", code, page.Name)

	data := struct {
		StatusText string
		StatusCode int
		Variables  map[string]interface{}
	}{
		StatusText: http.StatusText(code),
		StatusCode: code,
		Variables:  r.config.ErrorPages.Variables,
	}

	html, err := page.Render(data)
	if err != nil {
		ctx.WithError(err).Error("rendering error page")
		http.Error(w, "Error rendering error page.", http.StatusInternalServerError)
		return
	}

	r.ignore = true
	util.ClearHeader(w.Header())
	w.Header().Set("Vary", "Accept")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	io.WriteString(w, html)
}

// Write implementation.
func (r *response) Write(b []byte) (int, error) {
	if r.ignore {
		return len(b), nil
	}

	if !r.header {
		r.WriteHeader(200)
		return r.Write(b)
	}

	return r.ResponseWriter.Write(b)
}

// Errors handles error page support.
type Errors struct {
	next  http.Handler
	pages errorpage.Pages
}

// New error pages handler.
func New(c *up.Config, next http.Handler) (http.Handler, error) {
	pages, err := errorpage.Load(c.ErrorPages.Dir)
	if err != nil {
		return nil, errors.Wrap(err, "loading error pages")
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mime, _ := accept.Negotiate(r.Header.Get("Accept"), "text/html")

		if mime == "" {
			next.ServeHTTP(w, r)
			return
		}

		res := &response{ResponseWriter: w, pages: pages, config: c}
		next.ServeHTTP(res, r)
	})

	return h, nil
}
