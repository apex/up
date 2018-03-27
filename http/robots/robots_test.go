package robots

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/apex/up"
	"github.com/tj/assert"

	"github.com/apex/up/config"
	"github.com/apex/up/http/static"
)

func TestRobots(t *testing.T) {
	c := &up.Config{
		Static: config.Static{
			Dir: "testdata",
		},
	}

	t.Run("should set X-Robots-Tag", func(t *testing.T) {
		h := New(c, static.New(c))

		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "none", res.Header().Get("X-Robots-Tag"))
		assert.Equal(t, "text/html; charset=utf-8", res.Header().Get("Content-Type"))
		assert.Equal(t, "Index HTML\n", res.Body.String())
	})

	t.Run("should not set X-Robots-Tag for production stage", func(t *testing.T) {
		os.Setenv("UP_STAGE", "production")
		defer os.Setenv("UP_STAGE", "")

		h := New(c, static.New(c))

		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "", res.Header().Get("X-Robots-Tag"))
		assert.Equal(t, "text/html; charset=utf-8", res.Header().Get("Content-Type"))
		assert.Equal(t, "Index HTML\n", res.Body.String())
	})
}
