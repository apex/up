package headers

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/tj/assert"
	"github.com/apex/up"

	"github.com/apex/up/http/static"
	"github.com/apex/up/internal/header"
)

func TestHeaders(t *testing.T) {
	os.Chdir("testdata")
	defer os.Chdir("..")

	c := &up.Config{
		Headers: header.Rules{
			"/*.css": {
				"Cache-Control": "public, max-age=999999",
			},
		},
	}

	h, err := New(c, static.New(c))
	assert.NoError(t, err, "init")

	t.Run("mismatch", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "", res.Header().Get("Cache-Control"))
		assert.Equal(t, "text/html; charset=utf-8", res.Header().Get("Content-Type"))
		assert.Equal(t, "Index HTML\n", res.Body.String())
	})

	t.Run("matched exact", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/style.css", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "public, max-age=999999", res.Header().Get("Cache-Control"))
		assert.Equal(t, "css", res.Header().Get("X-Type"))
		assert.Equal(t, "text/css; charset=utf-8", res.Header().Get("Content-Type"))
		assert.Equal(t, "body { color: red }\n", res.Body.String())
	})
}
