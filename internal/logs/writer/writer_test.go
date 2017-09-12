package writer

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/tj/assert"
)

func init() {
	log.Now = func() time.Time {
		return time.Unix(0, 0)
	}
}

func TestWriter_plainTextFlat(t *testing.T) {
	var buf bytes.Buffer

	log.SetHandler(json.New(&buf))

	w := New(log.InfoLevel)

	input := `GET /
GET /account
GET /login
POST /
POST /logout
`

	_, err := io.Copy(w, strings.NewReader(input))
	assert.NoError(t, err, "copy")

	assert.NoError(t, w.Close(), `close`)

	expected := `{"fields":{"app":true},"level":"info","timestamp":"1969-12-31T16:00:00-08:00","message":"GET /"}
{"fields":{"app":true},"level":"info","timestamp":"1969-12-31T16:00:00-08:00","message":"GET /account"}
{"fields":{"app":true},"level":"info","timestamp":"1969-12-31T16:00:00-08:00","message":"GET /login"}
{"fields":{"app":true},"level":"info","timestamp":"1969-12-31T16:00:00-08:00","message":"POST /"}
{"fields":{"app":true},"level":"info","timestamp":"1969-12-31T16:00:00-08:00","message":"POST /logout"}
`

	assert.Equal(t, expected, buf.String())
}

func TestWriter_plainTextIndented(t *testing.T) {
	var buf bytes.Buffer

	log.SetHandler(json.New(&buf))

	w := New(log.InfoLevel)

	input := `GET /
GET /account
SomethingError: one
  at foo
  at bar
  at baz
  at raz
GET /login
SomethingError: two
  at foo
  at bar
  at baz
POST /
SomethingError: three
  at foo
  at bar
  at baz
`

	_, err := io.Copy(w, strings.NewReader(input))
	assert.NoError(t, err, "copy")

	assert.NoError(t, w.Close(), `close`)

	expected := `{"fields":{"app":true},"level":"info","timestamp":"1969-12-31T16:00:00-08:00","message":"GET /"}
{"fields":{"app":true},"level":"info","timestamp":"1969-12-31T16:00:00-08:00","message":"GET /account"}
{"fields":{"app":true},"level":"info","timestamp":"1969-12-31T16:00:00-08:00","message":"SomethingError: one\n  at foo\n  at bar\n  at baz\n  at raz"}
{"fields":{"app":true},"level":"info","timestamp":"1969-12-31T16:00:00-08:00","message":"GET /login"}
{"fields":{"app":true},"level":"info","timestamp":"1969-12-31T16:00:00-08:00","message":"SomethingError: two\n  at foo\n  at bar\n  at baz"}
{"fields":{"app":true},"level":"info","timestamp":"1969-12-31T16:00:00-08:00","message":"POST /"}
{"fields":{"app":true},"level":"info","timestamp":"1969-12-31T16:00:00-08:00","message":"SomethingError: three\n  at foo\n  at bar\n  at baz"}
`

	assert.Equal(t, expected, buf.String())
}

func TestWriter_json(t *testing.T) {
	var buf bytes.Buffer

	log.SetHandler(json.New(&buf))

	w := New(log.InfoLevel)

	input := `{ "level": "info", "message": "request", "method": "GET", "path": "/" }
{ "level": "info", "message": "request", "method": "GET", "path": "/login" }
{ "level": "info", "message": "request", "method": "POST", "path": "/login" }
`

	_, err := io.Copy(w, strings.NewReader(input))
	assert.NoError(t, err, "copy")

	assert.NoError(t, w.Close(), `close`)

	expected := `{"fields":{"app":true,"method":"GET","path":"/"},"level":"info","timestamp":"1969-12-31T16:00:00-08:00","message":"request"}
{"fields":{"app":true,"method":"GET","path":"/login"},"level":"info","timestamp":"1969-12-31T16:00:00-08:00","message":"request"}
{"fields":{"app":true,"method":"POST","path":"/login"},"level":"info","timestamp":"1969-12-31T16:00:00-08:00","message":"request"}
`

	assert.Equal(t, expected, buf.String())
}
