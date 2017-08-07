package static

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/tj/assert"
	"github.com/apex/up"
	"github.com/apex/up/config"
)

func TestStatic_defaults(t *testing.T) {
	os.Chdir("testdata")
	defer os.Chdir("..")

	c := &up.Config{}
	assert.NoError(t, c.Default(), "default")
	assert.NoError(t, c.Validate(), "validate")
	test(t, c)
}

func TestStatic_dir(t *testing.T) {
	c := &up.Config{
		Static: config.Static{
			Dir: "testdata",
		},
	}

	assert.NoError(t, c.Default(), "default")
	assert.NoError(t, c.Validate(), "validate")
	test(t, c)
}

func test(t *testing.T, c *up.Config) {
	h := New(c)

	t.Run("index.html", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "text/html; charset=utf-8", res.Header().Get("Content-Type"))
		assert.Equal(t, "Index HTML\n", res.Body.String())
	})

	t.Run("file", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/style.css", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "text/css; charset=utf-8", res.Header().Get("Content-Type"))
		assert.Equal(t, "body { background: whatever }\n", res.Body.String())
	})

	t.Run("missing", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/notfound", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 404, res.Code)
		assert.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))
		assert.Equal(t, "404 page not found\n", res.Body.String())
	})

	t.Run("conditional get", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/style.css", nil)
		req.Header.Set("If-Modified-Since", "Thu, 27 Jul 2030 03:30:31 GMT")
		h.ServeHTTP(res, req)
		assert.Equal(t, 304, res.Code)
		assert.Equal(t, "", res.Header().Get("Content-Length"))
		assert.Equal(t, "", res.Body.String())
	})
}
