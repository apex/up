package redirects

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/apex/up"
	"github.com/tj/assert"

	"github.com/apex/up/internal/redirect"
)

func TestRedirects(t *testing.T) {
	t.Run("from config", func(t *testing.T) {
		c := &up.Config{
			Redirects: redirect.Rules{
				"/blog": {
					Location: "https://blog.apex.sh",
					Status:   301,
				},
				"/enterprise": {
					Location: "/docs/enterprise",
					Status:   302,
				},
				"/api": {
					Location: "/api/v1",
					Status:   200,
				},
				"/products": {
					Location: "/store",
					Status:   301,
				},
				"/app/*": {
					Location: "/",
				},
				"/app/login": {
					Location: "https://app.apex.sh",
					Status:   301,
				},
				"/documentation/:product/guides/:guide": {
					Location: "/docs/:product/:guide",
					Status:   200,
				},
				"/shop/:brand": {
					Location: "/products/:brand",
					Status:   301,
				},
				"/settings/*": {
					Location: "/admin/:splat",
					Status:   200,
					Force:    true,
				},
			},
		}

		handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/":
				fmt.Fprintln(w, "Index")
			case r.URL.Path == "/api/v1":
				fmt.Fprintln(w, "API V1")
			case r.URL.Path == "/products":
				fmt.Fprintln(w, "products")
			case strings.Contains(r.URL.Path, "/docs"):
				fmt.Fprintf(w, "docs %s", r.URL.Path)
			case strings.HasPrefix(r.URL.Path, "/brand"):
				fmt.Fprintf(w, "shop %s", r.URL.Path)
			case strings.HasPrefix(r.URL.Path, "/setting"):
				fmt.Fprintf(w, "settings %s", r.URL.Path)
			case strings.HasPrefix(r.URL.Path, "/admin"):
				fmt.Fprintf(w, "admin %s", r.URL.Path)
			default:
				http.NotFound(w, r)
			}
		})

		h, err := New(c, handle)
		assert.NoError(t, err, "init")

		test(t, h)
	})
}

func test(t *testing.T, h http.Handler) {
	os.Chdir("testdata")
	defer os.Chdir("..")

	t.Run("mismatch", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))
		assert.Equal(t, "Index\n", res.Body.String())
	})

	t.Run("exact match", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/blog", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 301, res.Code)
		assert.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))
		assert.Equal(t, "Moved Permanently\n", res.Body.String())
	})

	t.Run("exact match status", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/enterprise", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 302, res.Code)
		assert.Equal(t, "/docs/enterprise", res.Header().Get("Location"))
		assert.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))
		assert.Equal(t, "Found\n", res.Body.String())
	})

	t.Run("exact match rewrite", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "/api", req.Header.Get("X-Original-Path"))
		assert.Empty(t, res.Header().Get("Location"), "location")
		assert.Equal(t, "API V1\n", res.Body.String())
	})

	t.Run("shadowed exact", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/products", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 301, res.Code)
		assert.Equal(t, "/store", res.Header().Get("Location"))
		assert.Equal(t, "Moved Permanently\n", res.Body.String())
	})

	t.Run("shadowed dynamic", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/app/contact", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "Index\n", res.Body.String())
	})

	t.Run("match precedence", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/app/login", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 301, res.Code)
		assert.Equal(t, "https://app.apex.sh", res.Header().Get("Location"))
		assert.Equal(t, "Moved Permanently\n", res.Body.String())
	})

	t.Run("rewrite with placeholders", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/documentation/ping/guides/alerting", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "docs /docs/ping/alerting", res.Body.String())
	})

	t.Run("redirect with placeholders", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/shop/apple", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 301, res.Code)
		assert.Equal(t, "Moved Permanently\n", res.Body.String())
	})

	t.Run("forced rewrite", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/settings/login", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "admin /admin/login", res.Body.String())
	})
}
