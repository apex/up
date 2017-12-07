package cors

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apex/up"
	"github.com/tj/assert"
)

var hello = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
})

func TestCORS_disabled(t *testing.T) {
	c, err := up.ParseConfigString(`{
		"name": "app"
	}`)

	assert.NoError(t, err, "config")

	h := New(c, hello)

	t.Run("GET", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		req.Header.Set("Origin", "https://example.com")

		h.ServeHTTP(res, req)

		header := make(http.Header)
		header.Add("Content-Type", "text/plain; charset=utf-8")

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, header, res.HeaderMap)
		assert.Equal(t, "Hello World", res.Body.String())
	})
}

func TestCORS_defaults(t *testing.T) {
	c, err := up.ParseConfigString(`{
		"name": "app",
		"cors": {}
	}`)

	assert.NoError(t, err, "config")

	h := New(c, hello)

	t.Run("GET", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		req.Header.Set("Origin", "https://example.com")

		h.ServeHTTP(res, req)

		header := make(http.Header)
		header.Add("Content-Type", "text/plain; charset=utf-8")
		header.Add("Vary", "Origin")
		header.Add("Access-Control-Allow-Origin", "*")

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, header, res.HeaderMap)
		assert.Equal(t, "Hello World", res.Body.String())
	})

	t.Run("OPTIONS", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("OPTIONS", "/", nil)

		req.Header.Set("Access-Control-Request-Method", "POST")
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Access-Control-Request-Headers", "Content-Type")

		h.ServeHTTP(res, req)

		header := make(http.Header)
		header.Add("Vary", "Origin")
		header.Add("Vary", "Access-Control-Request-Method")
		header.Add("Vary", "Access-Control-Request-Headers")
		header.Add("Access-Control-Allow-Methods", "POST")
		header.Add("Access-Control-Allow-Headers", "Content-Type")
		header.Add("Access-Control-Allow-Origin", "*")

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, header, res.HeaderMap)
		assert.Equal(t, "", res.Body.String())
	})
}

func TestCORS_options(t *testing.T) {
	c := up.MustParseConfigString(`{
		"name": "app",
		"cors": {
			"allowed_origins": ["https://apex.sh"],
			"allowed_methods": ["GET"],
			"allow_credentials": true,
			"max_age": 86400
		}
	}`)

	h := New(c, hello)

	t.Run("GET", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		req.Header.Set("Origin", "https://example.com")

		h.ServeHTTP(res, req)

		header := make(http.Header)
		header.Add("Content-Type", "text/plain; charset=utf-8")
		header.Add("Vary", "Origin")

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, header, res.HeaderMap)
		assert.Equal(t, "Hello World", res.Body.String())
	})

	t.Run("OPTIONS", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("OPTIONS", "/", nil)

		req.Header.Set("Access-Control-Request-Method", "POST")
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Access-Control-Request-Headers", "Content-Type")

		h.ServeHTTP(res, req)

		header := make(http.Header)
		header.Add("Vary", "Origin")
		header.Add("Vary", "Access-Control-Request-Method")
		header.Add("Vary", "Access-Control-Request-Headers")

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, header, res.HeaderMap)
		assert.Equal(t, "", res.Body.String())
	})
}
