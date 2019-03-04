// Package static provides static file serving with HTTP cache support.
package static

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/up"
)

// New static handler.
func New(c *up.Config) http.Handler {
	next := http.FileServer(http.Dir(c.Static.Dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if blacklisted(r.URL.Path) {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// NewDynamic static handler for dynamic apps.
func NewDynamic(c *up.Config, next http.Handler) http.Handler {
	prefix := normalizePrefix(c.Static.Prefix)
	dir := c.Static.Dir

	if dir == "" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var skip bool
		path := r.URL.Path

		// prefix
		if prefix != "" {
			if strings.HasPrefix(path, prefix) {
				path = strings.Replace(path, prefix, "/", 1)
			} else {
				skip = true
			}
		}

		// convert
		path = filepath.FromSlash(path)

		// blacklist
		if blacklisted(r.URL.Path) {
			http.NotFound(w, r)
			return
		}

		// file exists, serve it
		if !skip {
			file := filepath.Join(dir, path)
			info, err := os.Stat(file)
			if !os.IsNotExist(err) && !info.IsDir() {
				http.ServeFile(w, r, file)
				return
			}
		}

		// delegate
		next.ServeHTTP(w, r)
	})
}

// normalizePrefix returns a prefix path normalized with leading and trailing "/".
func normalizePrefix(s string) string {
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}

	if !strings.HasSuffix(s, "/") {
		s = s + "/"
	}

	return s
}

// blacklisted returns true if the path is blacklisted.
func blacklisted(s string) bool {
	return strings.Contains(s, "up-env")
}
