package relay

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/tj/assert"

	"github.com/apex/up"
)

// TODO: fix hanging

func TestRelay(t *testing.T) {
	os.Chdir("testdata/basic")
	defer os.Chdir("../..")

	c := &up.Config{}
	assert.NoError(t, c.Default(), "default")

	h, err := New(c)
	assert.NoError(t, err, "init")

	t.Run("GET simple", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/hello", nil)

		start := time.Now()
		h.ServeHTTP(res, req)
		t.Logf("latency = %s", time.Since(start))

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "text/plain", res.Header().Get("Content-Type"))
		assert.Equal(t, "Hello World", res.Body.String())
	})

	t.Run("GET encoded path", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
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
    "accept-encoding": "gzip"
  },
  "url": "/echo/01BM82CJ9K1WK6EFJX8C1R4YH7/foo%20%25%20bar%20&%20baz%20=%20raz",
  "body": ""
}`

			assert.Equal(t, 200, res.Code)
			assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
			assertString(t, body, res.Body.String())
		})
	})

	t.Run("POST basic", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			res := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/echo/something", strings.NewReader("Some body here"))

			h.ServeHTTP(res, req)

			body := `{
  "header": {
    "host": "example.com",
    "content-length": "14",
    "x-forwarded-for": "192.0.2.1",
    "accept-encoding": "gzip"
  },
  "url": "/echo/something",
  "body": "Some body here"
}`

			assert.Equal(t, 200, res.Code)
			assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
			assertString(t, body, res.Body.String())
		})
	})

	t.Run("restart server", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			res := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/throw/env", nil)

			h.ServeHTTP(res, req)

			assert.Equal(t, 200, res.Code)
			assertString(t, "Hello", res.Body.String())
		})
	})

}

// TODO: move test to handler pkg or just test the config
func TestRelay_node(t *testing.T) {
	os.Chdir("testdata/node")
	defer os.Chdir("../..")

	c, err := up.ReadConfig("up.json")
	assert.NoError(t, err, "config")

	h, err := New(c)
	assert.NoError(t, err)

	t.Run("npm start", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		h.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assertString(t, "Node", res.Body.String())
	})
}

func assertString(t testing.TB, want, got string) {
	if want != got {
		t.Fatalf("\nwant:\n\n%s\n\ngot:\n\n%s\n", want, got)
	}
}

func BenchmarkRelay(b *testing.B) {
	os.Chdir("testdata/basic")
	defer os.Chdir("../..")

	b.ReportAllocs()

	c := &up.Config{}
	assert.NoError(b, c.Default(), "default")

	h, err := New(c)
	assert.NoError(b, err)

	b.ResetTimer()
	b.SetParallelism(30)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			res := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)
			h.ServeHTTP(res, req)
		}
	})
}
