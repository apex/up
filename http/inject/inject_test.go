package inject

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/apex/up"
	"github.com/tj/assert"

	"github.com/apex/up/config"
	"github.com/apex/up/http/errorpages"
	"github.com/apex/up/http/static"
	"github.com/apex/up/internal/inject"
)

func TestInject(t *testing.T) {
	os.Chdir("testdata")
	defer os.Chdir("..")

	c := &up.Config{
		Name: "app",
		ErrorPages: config.ErrorPages{
			Enable: true,
		},
		Inject: inject.Rules{
			"head": []*inject.Rule{
				{
					Type:  "script",
					Value: "/whatever.js",
				},
			},
		},
	}

	assert.NoError(t, c.Default(), "default")
	assert.NoError(t, c.Validate(), "validate")

	h, err := New(c, static.New(c))
	assert.NoError(t, err, "init")

	h, err = errorpages.New(c, h)
	assert.NoError(t, err, "init")

	t.Run("2xx", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		h.ServeHTTP(res, req)

		html := `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <script src="/whatever.js"></script>
  </head>
  <body>

  </body>
</html>
`

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "text/html; charset=utf-8", res.Header().Get("Content-Type"))
		assert.Equal(t, html, res.Body.String())
	})

	t.Run("4xx", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/missing", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 404, res.Code)
		assert.Equal(t, "text/html; charset=utf-8", res.Header().Get("Content-Type"))
		assert.Equal(t, "<p>Not Found</p>\n", res.Body.String())
	})

	t.Run("non-html", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/style.css", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "text/css; charset=utf-8", res.Header().Get("Content-Type"))
		assert.Equal(t, "body{}\n", res.Body.String())
	})

	t.Run("write before header", func(t *testing.T) {
		s := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			io.WriteString(w, "<html>")
			io.WriteString(w, "<head>")
			io.WriteString(w, "</head>")
			io.WriteString(w, "<body>")
			io.WriteString(w, "</body>")
			io.WriteString(w, "</html>")
		})

		h, err := New(c, s)
		assert.NoError(t, err, "initialize")

		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "<html><head>  <script src=\"/whatever.js\"></script>\n  </head><body></body></html>", res.Body.String())
	})
}
