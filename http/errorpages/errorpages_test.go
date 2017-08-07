package errorpages

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/tj/assert"
	"github.com/apex/up"
	"github.com/apex/up/config"
)

var server = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/404" {
		w.Header().Set("X-Foo", "bar")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if r.URL.Path == "/400" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if r.URL.Path == "/500" {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-Foo", "bar")
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello")
	fmt.Fprintf(w, " ")
	fmt.Fprintf(w, "World")
})

func TestErrors_defaults(t *testing.T) {
	os.Chdir("testdata")
	defer os.Chdir("..")

	c := &up.Config{}
	assert.NoError(t, c.Default(), "default")
	assert.NoError(t, c.Validate(), "validate")

	test(t, c)
}

func TestErrors_dir(t *testing.T) {
	c := &up.Config{
		ErrorPages: config.ErrorPages{
			Dir: "testdata",
		},
	}

	assert.NoError(t, c.Default(), "default")
	assert.NoError(t, c.Validate(), "validate")

	test(t, c)
}

func test(t *testing.T, c *up.Config) {
	h, err := New(c, server)
	assert.NoError(t, err, "init")

	t.Run("200", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "bar", res.Header().Get("X-Foo"))
		assert.Equal(t, "text/plain", res.Header().Get("Content-Type"))
		assert.Equal(t, "Hello World", res.Body.String())
	})

	t.Run("exact", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/404", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 404, res.Code)
		assert.Equal(t, "", res.Header().Get("X-Foo"))
		assert.Equal(t, "text/html; charset=utf-8", res.Header().Get("Content-Type"))
		assert.Equal(t, "Sorry! Can't find that.\n", res.Body.String())
	})

	t.Run("range", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/500", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 500, res.Code)
		assert.Equal(t, "Accept", res.Header().Get("Vary"))
		assert.Equal(t, "", res.Header().Get("X-Foo"))
		assert.Equal(t, "text/html; charset=utf-8", res.Header().Get("Content-Type"))
		assert.Equal(t, "500 – Internal Server Error\n", res.Body.String())
	})

	t.Run("default text/html", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/400", nil)
		req.Header.Set("Accept", "text/html")

		h.ServeHTTP(res, req)

		assert.Equal(t, 400, res.Code)
		assert.Equal(t, "Accept", res.Header().Get("Vary"))
		assert.Equal(t, "", res.Header().Get("X-Foo"))
		assert.Equal(t, "text/html; charset=utf-8", res.Header().Get("Content-Type"))
		assert.Contains(t, res.Body.String(), "<title>Bad Request – 400</title>", "title")
		assert.Contains(t, res.Body.String(), `<span class="status">Bad Request</span>`, "status text")
		assert.Contains(t, res.Body.String(), `<span class="code">400</span>`, "status code")
	})

	t.Run("default text/*", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/400", nil)
		req.Header.Set("Accept", "text/*")

		h.ServeHTTP(res, req)

		assert.Equal(t, 400, res.Code)
		assert.Equal(t, "Accept", res.Header().Get("Vary"))
		assert.Equal(t, "", res.Header().Get("X-Foo"))
		assert.Equal(t, "text/html; charset=utf-8", res.Header().Get("Content-Type"))
		assert.Contains(t, res.Body.String(), "<title>Bad Request – 400</title>", "title")
		assert.Contains(t, res.Body.String(), `<span class="status">Bad Request</span>`, "status text")
		assert.Contains(t, res.Body.String(), `<span class="code">400</span>`, "status code")
	})
}
