package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/apex/up"
	"github.com/tj/assert"
)

func newHandler(t testing.TB, c *up.Config) http.Handler {
	h, err := FromConfig(c)
	assert.NoError(t, err, "FromConfig")

	h, err = New(c, h)
	assert.NoError(t, err, "New")

	return h
}

func TestNode(t *testing.T) {
	os.Chdir("testdata/node")
	defer os.Chdir("../..")

	c, err := up.ReadConfig("up.json")
	assert.NoError(t, err, "read config")

	h := newHandler(t, c)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	h.ServeHTTP(res, req)

	actual := res.Header()
	assert.NotEmpty(t, actual.Get("Date"), "date")
	actual.Del("Date")

	header := make(http.Header)
	header.Add("X-Powered-By", "up")
	header.Add("X-Robots-Tag", "none")
	header.Add("X-Foo", "bar")
	header.Add("Content-Length", "11")
	header.Add("Content-Type", "text/plain; charset=utf-8")
	header.Add("Vary", "Accept-Encoding")
	assert.Equal(t, header, actual)
}

func TestStatic(t *testing.T) {
	os.Chdir("testdata/static")
	defer os.Chdir("../..")

	c, err := up.ReadConfig("up.json")
	assert.NoError(t, err, "read config")

	h := newHandler(t, c)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	h.ServeHTTP(res, req)

	actual := res.Header()
	assert.NotEmpty(t, actual.Get("Last-Modified"), "last-modified")
	actual.Del("Last-Modified")

	header := make(http.Header)
	header.Add("X-Powered-By", "up")
	header.Add("X-Robots-Tag", "none")
	header.Add("Content-Length", "12")
	header.Add("Content-Type", "text/html; charset=utf-8")
	header.Add("Accept-Ranges", "bytes")
	header.Add("Vary", "Accept-Encoding")

	assert.Equal(t, header, actual)
}

func TestNodeWithPackage(t *testing.T) {
	os.Chdir("testdata/node-pkg")
	defer os.Chdir("../..")

	c, err := up.ReadConfig("up.json")
	assert.NoError(t, err, "read config")

	h := newHandler(t, c)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	h.ServeHTTP(res, req)

	assert.Equal(t, "Hello World", res.Body.String())
}

func TestNodeWithPackageStart(t *testing.T) {
	os.Chdir("testdata/node-pkg-start")
	defer os.Chdir("../..")

	c, err := up.ReadConfig("up.json")
	assert.NoError(t, err, "read config")

	h := newHandler(t, c)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	h.ServeHTTP(res, req)

	assert.Equal(t, "Hello World", res.Body.String())
}

func TestHandler_conditionalGet(t *testing.T) {
	os.Chdir("testdata/static")
	defer os.Chdir("../..")

	c, err := up.ReadConfig("up.json")
	assert.NoError(t, err, "read config")

	h := newHandler(t, c)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/style.css", nil)
	req.Header.Set("If-Modified-Since", "Thu, 27 Jul 2030 03:30:31 GMT")
	h.ServeHTTP(res, req)
	assert.Equal(t, 304, res.Code)
	assert.Equal(t, "", res.Header().Get("Content-Length"))
	assert.Equal(t, "", res.Body.String())
}

func TestHandler_rewrite(t *testing.T) {
	os.Chdir("testdata/static-rewrites")
	defer os.Chdir("../..")

	c, err := up.ReadConfig("up.json")
	assert.NoError(t, err, "read config")

	h := newHandler(t, c)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/docs/ping/guides/alerts", nil)
	h.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "14", res.Header().Get("Content-Length"))
	assert.Equal(t, "Alerting docs\n", res.Body.String())
}

func TestHandler_redirect(t *testing.T) {
	os.Chdir("testdata/static-redirects")
	defer os.Chdir("../..")

	c, err := up.ReadConfig("up.json")
	assert.NoError(t, err, "read config")

	h := newHandler(t, c)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/docs/ping/guides/alerts/", nil)
	h.ServeHTTP(res, req)
	assert.Equal(t, "/help/ping/alerts", res.Header().Get("Location"))
	assert.Equal(t, 302, res.Code)
	assert.Equal(t, "Found\n", res.Body.String())
}

func TestHandler_spa(t *testing.T) {
	os.Chdir("testdata/spa")
	defer os.Chdir("../..")

	c, err := up.ReadConfig("up.json")
	assert.NoError(t, err, "read config")

	h := newHandler(t, c)

	t.Run("index", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		h.ServeHTTP(res, req)
		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "Index\n", res.Body.String())
	})

	t.Run("redirect", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/index.html", nil)
		h.ServeHTTP(res, req)
		assert.Equal(t, 301, res.Code)
	})

	t.Run("file does not exist", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/something/here", nil)
		h.ServeHTTP(res, req)
		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "Index\n", res.Body.String())
	})

	t.Run("file exists", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/app.js", nil)
		h.ServeHTTP(res, req)
		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "app js\n", res.Body.String())
	})

	t.Run("file exists nested", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/css/bar.css", nil)
		h.ServeHTTP(res, req)
		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "bar css\n", res.Body.String())
	})
}
