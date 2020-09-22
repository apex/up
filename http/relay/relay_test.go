package relay

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/tj/assert"

	"github.com/apex/up"
	"github.com/apex/up/config"
	"github.com/apex/up/internal/util"
)

func skipCI(t testing.TB) {
	if util.IsCI() {
		t.SkipNow()
	}
}

func TestRelay(t *testing.T) {
	os.Chdir("testdata/basic")
	defer os.Chdir("../..")

	c := &up.Config{
		Proxy: config.Relay{
			Timeout:       2,
			ListenTimeout: 2,
		},
	}

	assert.NoError(t, c.Default(), "default")

	var h http.Handler
	newHandler := func(t *testing.T) {
		v, err := New(c)
		assert.NoError(t, err, "init")
		h = v
	}

	t.Run("GET simple", func(t *testing.T) {
		newHandler(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/hello", nil)
		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "text/plain", res.Header().Get("Content-Type"))
		assert.Equal(t, "Hello World", res.Body.String())
	})

	t.Run("GET encoded path", func(t *testing.T) {
		newHandler(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/echo/01BM82CJ9K1WK6EFJX8C1R4YH7/foo%20%25%20bar%20&%20baz%20=%20raz", nil)
		req.Header.Set("Host", "example.com")
		req.Header.Set("User-Agent", "tobi")
		h.ServeHTTP(res, req)

		body := `{
  "header": {
    "host": "example.com",
    "user-agent": "tobi",
    "x-forwarded-for": "192.0.2.1",
    "accept-encoding": "gzip",
    "connection": "close"
  },
  "url": "/echo/01BM82CJ9K1WK6EFJX8C1R4YH7/foo%20%25%20bar%20&%20baz%20=%20raz",
  "body": ""
}`

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
		assertString(t, body, res.Body.String())
	})

	t.Run("POST simple", func(t *testing.T) {
		newHandler(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/echo/something", strings.NewReader("Some body here"))
		h.ServeHTTP(res, req)

		body := `{
  "header": {
    "host": "example.com",
    "content-length": "14",
    "x-forwarded-for": "192.0.2.1",
    "accept-encoding": "gzip",
    "connection": "close"
  },
  "url": "/echo/something",
  "body": "Some body here"
}`

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
		assertString(t, body, res.Body.String())
	})

	t.Run("crash", func(t *testing.T) {
		newHandler(t)

		// first
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/throw", nil)
		h.ServeHTTP(res, req)

		assert.Equal(t, 502, res.Code)
		assertString(t, "", res.Body.String())

		// wait for restart
		time.Sleep(time.Second)

		// second
		res = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "/hello", nil)
		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assertString(t, "Hello World", res.Body.String())
	})

	t.Run("timeout", func(t *testing.T) {
		newHandler(t)

		// first
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/timeout", nil)
		h.ServeHTTP(res, req)

		assert.Equal(t, 502, res.Code)
		assertString(t, "", res.Body.String())

		// second
		res = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "/hello", nil)
		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assertString(t, "Hello World", res.Body.String())

		// third
		res = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "/timeout", nil)
		req.Header.Add("UP-TIMEOUT", "0");
		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assertString(t, "Hello", res.Body.String())
	})
}

func assertString(t testing.TB, want, got string) {
	t.Helper()
	if want != got {
		t.Fatalf("\nwant:\n\n%s\n\ngot:\n\n%s\n", want, got)
	}
}
