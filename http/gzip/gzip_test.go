package gzip

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apex/up"
	"github.com/tj/assert"
)

var body = strings.Repeat("так", 5000)

var hello = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, body)
})

func TestGzip(t *testing.T) {
	c, err := up.ParseConfigString(`{ "name": "app" }`)
	assert.NoError(t, err, "config")

	h := New(c, hello)

	t.Run("accepts gzip", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept-Encoding", "gzip")

		h.ServeHTTP(res, req)

		header := make(http.Header)
		header.Add("Content-Type", "text/plain; charset=utf-8")
		header.Add("Content-Encoding", "gzip")
		header.Add("Vary", "Accept-Encoding")

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, header, res.HeaderMap)

		gz, err := gzip.NewReader(res.Body)
		assert.NoError(t, err, "reader")

		b, err := ioutil.ReadAll(gz)
		assert.NoError(t, err, "reading")
		assert.NoError(t, gz.Close(), "close")

		assert.Equal(t, body, string(b))
	})

	t.Run("accepts identity", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		h.ServeHTTP(res, req)

		header := make(http.Header)
		header.Add("Content-Type", "text/plain; charset=utf-8")
		header.Add("Vary", "Accept-Encoding")

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, header, res.HeaderMap)

		assert.Equal(t, body, res.Body.String())
	})
}
