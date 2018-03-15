// Package redirects provides redirection and URL rewriting.
package redirects

import (
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/apex/up"
	"github.com/apex/up/internal/logs"
	"github.com/apex/up/internal/redirect"
)

// TODO: tests for popagating 4xx / 5xx, dont mask all these
// TODO: load _redirects relative to .Static.Dir?
// TODO: add list of methods to match on

// log context.
var ctx = logs.Plugin("redirects")

type rewrite struct {
	http.ResponseWriter
	header     bool
	isNotFound bool
}

// WriteHeader implementation.
func (r *rewrite) WriteHeader(code int) {
	r.header = true
	r.isNotFound = code == 404

	if r.isNotFound {
		return
	}

	r.ResponseWriter.WriteHeader(code)
}

// Write implementation.
func (r *rewrite) Write(b []byte) (int, error) {
	if r.isNotFound {
		return len(b), nil
	}

	if !r.header {
		r.WriteHeader(200)
		return r.Write(b)
	}

	return r.ResponseWriter.Write(b)
}

// New redirects handler.
func New(c *up.Config, next http.Handler) (http.Handler, error) {
	if len(c.Redirects) == 0 {
		return next, nil
	}

	rules, err := redirect.Compile(c.Redirects)
	if err != nil {
		return nil, err
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rule := rules.Lookup(r.URL.Path)

		ctx := ctx.WithFields(log.Fields{
			"path": r.URL.Path,
		})

		// pass-through
		if rule == nil {
			ctx.Debug("no match")
			next.ServeHTTP(w, r)
			return
		}

		// destination path
		path := rule.URL(r.URL.Path)

		// forced rewrite
		if rule.IsRewrite() && rule.Force {
			ctx.WithField("dest", path).Info("forced rewrite")
			r.Header.Set("X-Original-Path", r.URL.Path)
			r.URL.Path = path
			next.ServeHTTP(w, r)
			return
		}

		// rewrite
		if rule.IsRewrite() {
			res := &rewrite{ResponseWriter: w}
			next.ServeHTTP(res, r)

			if res.isNotFound {
				ctx.WithField("dest", path).Info("rewrite")
				r.Header.Set("X-Original-Path", r.URL.Path)
				r.URL.Path = path
				// This hack is necessary for SPAs because the Go
				// static file server uses .html to set the correct mime,
				// ideally it uses the file's extension or magic number etc.
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				next.ServeHTTP(w, r)
			}
			return
		}

		// redirect
		ctx.WithField("dest", path).Info("redirect")
		w.Header().Set("Location", path)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(rule.Status)
		fmt.Fprintln(w, http.StatusText(rule.Status))
	})

	return h, nil
}
