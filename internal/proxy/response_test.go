package proxy

import (
	"bytes"
	"testing"

	"github.com/tj/assert"
)

func Test_JSON_isTextMime(t *testing.T) {
	assert.Equal(t, isTextMime("application/json"), true)
	assert.Equal(t, isTextMime("application/json; charset=utf-8"), true)
	assert.Equal(t, isTextMime("Application/JSON"), true)
}

func Test_XML_isTextMime(t *testing.T) {
	assert.Equal(t, isTextMime("application/xml"), true)
	assert.Equal(t, isTextMime("application/xml; charset=utf-8"), true)
	assert.Equal(t, isTextMime("ApPlicaTion/xMl"), true)
}

func TestResponseWriter_Header(t *testing.T) {
	w := NewResponse()
	w.Header().Set("Foo", "bar")
	w.Header().Set("Bar", "baz")

	var buf bytes.Buffer
	w.header.Write(&buf)

	assert.Equal(t, "Bar: baz\r\nFoo: bar\r\n", buf.String())
}

func TestResponseWriter_Write_text(t *testing.T) {
	types := []string{
		"text/x-custom",
		"text/plain",
		"text/plain; charset=utf-8",
		"application/json",
		"application/json; charset=utf-8",
		"application/xml",
		"image/svg+xml",
	}

	for _, kind := range types {
		t.Run(kind, func(t *testing.T) {
			w := NewResponse()
			w.Header().Set("Content-Type", kind)
			w.Write([]byte("hello world\n"))

			e := w.End()
			assert.Equal(t, 200, e.StatusCode)
			assert.Equal(t, "hello world\n", e.Body)
			assert.Equal(t, kind, e.Headers["Content-Type"])
			assert.False(t, e.IsBase64Encoded)
		})
	}
}

func TestResponseWriter_Write_binary(t *testing.T) {
	w := NewResponse()
	w.Header().Set("Content-Type", "image/png")
	w.Write([]byte("data"))

	e := w.End()
	assert.Equal(t, 200, e.StatusCode)
	assert.Equal(t, "ZGF0YQ==", e.Body)
	assert.Equal(t, "image/png", e.Headers["Content-Type"])
	assert.True(t, e.IsBase64Encoded)
}

func TestResponseWriter_Write_gzip(t *testing.T) {
	w := NewResponse()
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Encoding", "gzip")
	w.Write([]byte("data"))

	e := w.End()
	assert.Equal(t, 200, e.StatusCode)
	assert.Equal(t, "ZGF0YQ==", e.Body)
	assert.Equal(t, "text/plain", e.Headers["Content-Type"])
	assert.True(t, e.IsBase64Encoded)
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	w := NewResponse()
	w.WriteHeader(404)
	w.Write([]byte("Not Found\n"))

	e := w.End()
	assert.Equal(t, 404, e.StatusCode)
	assert.Equal(t, "Not Found\n", e.Body)
	assert.Equal(t, "text/plain; charset=utf8", e.Headers["Content-Type"])
}
