package relay

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/tj/assert"

	"github.com/apex/up"
	"github.com/apex/up/config"
)

func TestRelay(t *testing.T) {
	os.Chdir("testdata/basic")
	defer os.Chdir("../..")

	c := &up.Config{
		Proxy: config.Relay{
			ListenTimeout:   2,
			ShutdownTimeout: 2,
		},
	}

	assert.NoError(t, c.Default(), "default")

	var h http.Handler
	newHandler := func(t *testing.T) {
		v, err := New(c)
		assert.NoError(t, err, "init")
		h = v
	}

	// newRequestBody is a helper to fetch the body from a child process route
	newRequestBody := func(t *testing.T, method, path string, body io.Reader) string {
		t.Helper()
		res := httptest.NewRecorder()
		req := httptest.NewRequest(method, path, body)
		h.ServeHTTP(res, req)
		assert.Equal(t, 200, res.Code)
		return res.Body.String()
	}

	// numRestarts fetches the numer of restarts for this relay handler (UP_RESTARTS env var)
	numRestarts := func() int {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/env?key=UP_RESTARTS", nil)
		h.ServeHTTP(res, req)
		body := res.Body.String()
		r, err := strconv.ParseInt(body, 10, 32)
		if err != nil {
			panic(err)
		}

		return int(r)
	}

	// childPID returns the pid of the child process currently running in the proxy
	childPID := func(t *testing.T) int {
		t.Helper()
		body := newRequestBody(t, "GET", "/pid", nil)
		r, err := strconv.ParseInt(body, 10, 32)
		if err != nil {
			panic(err)
		}

		return int(r)
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
    "accept-encoding": "gzip"
  },
  "url": "/echo/01BM82CJ9K1WK6EFJX8C1R4YH7/foo%20%25%20bar%20&%20baz%20=%20raz",
  "body": ""
}`

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
		assertString(t, body, res.Body.String())
	})

	t.Run("POST basic", func(t *testing.T) {
		newHandler(t)

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

	closeApp := func(t *testing.T) {
		t.Helper()
		r1 := numRestarts()
		body := newRequestBody(t, "GET", "/close", nil)
		assertString(t, "closed", body)

		r2 := numRestarts()
		assert.Equal(t, r2, r1+1)
	}

	// A bad route (such as /throw) eventually stops retrying
	t.Run("bad route", func(t *testing.T) {
		newHandler(t)

		r1 := numRestarts()
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/throw", nil)
		h.ServeHTTP(res, req)
		r2 := numRestarts()

		// Be a good proxy now
		assert.Equal(t, 502, res.Code)

		// 3 retries == 4 tries
		assert.Equal(t, r2-r1, 4)
	})

	t.Run("restart server", func(t *testing.T) {
		newHandler(t)

		t.Run("200", func(t *testing.T) {
			body := newRequestBody(t, "GET", "/throw/env", nil)
			assertString(t, "Hello", body)
		})

		t.Run("use after close", func(t *testing.T) {
			closeApp(t)
		})

		t.Run("use after close restart loop", func(t *testing.T) {
			for i := 0; i < 7; i++ {
				closeApp(t)
			}
		})
	})

	t.Run("child process cleanup", func(t *testing.T) {

		// Test that a child process who stops accepting network connections
		// (but hasn't crashed) is gracefully asked to be shut down
		t.Run("net offline zombie", func(t *testing.T) {
			newHandler(t)

			pid := childPID(t)
			process, err := os.FindProcess(pid)
			assert.NoError(t, err, "find process")

			start := time.Now()
			closeApp(t)
			ps, err := process.Wait()
			if err != nil {
				// This process might have completed before this test progressed this far
				assert.Contains(t, err.Error(), "no child processes")
				return
			}

			assert.True(t, time.Since(start).Seconds() >= 2)
			assert.True(t, time.Since(start).Seconds() < 4)
			assert.NoError(t, err, "zombie wait")
			assert.True(t, ps.Exited())
		})

		// Test that a child process who swallows the "nice" shutdown signal
		// will eventually be sent a SIGKILL and shut down
		t.Run("signal swallower", func(t *testing.T) {
			newHandler(t)

			pid := childPID(t)
			process, err := os.FindProcess(pid)
			assert.NoError(t, err, "find process")

			// First, swallow any 'normal' signals
			newRequestBody(t, "GET", "/swallowSignals", nil)

			// Then close the app (this triggers a restart)
			start := time.Now()
			closeApp(t)
			_, err = process.Wait()
			if err != nil {
				assert.Contains(t, err.Error(), "no child processes")
			}

			assert.True(t, time.Since(start).Seconds() >= 2)
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
	t.Helper()
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
