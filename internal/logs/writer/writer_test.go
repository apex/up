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
		return time.Unix(0, 0).UTC()
	}
}

func TestWriter_plainTextFlat(t *testing.T) {
	var buf bytes.Buffer

	log.SetHandler(json.New(&buf))

	w := New(log.InfoLevel, log.Log)

	input := `GET /
GET /account
GET /login
POST /
POST /logout
`

	_, err := io.Copy(w, strings.NewReader(input))
	assert.NoError(t, err, "copy")

	expected := `{"fields":{},"level":"info","timestamp":"1970-01-01T00:00:00Z","message":"GET /"}
{"fields":{},"level":"info","timestamp":"1970-01-01T00:00:00Z","message":"GET /account"}
{"fields":{},"level":"info","timestamp":"1970-01-01T00:00:00Z","message":"GET /login"}
{"fields":{},"level":"info","timestamp":"1970-01-01T00:00:00Z","message":"POST /"}
{"fields":{},"level":"info","timestamp":"1970-01-01T00:00:00Z","message":"POST /logout"}
`

	assert.Equal(t, expected, buf.String())
}

func TestWriter_json(t *testing.T) {
	var buf bytes.Buffer

	log.SetHandler(json.New(&buf))

	w := New(log.InfoLevel, log.Log)

	input := `{ "level": "info", "message": "request", "fields": { "method": "GET", "path": "/" } }
{ "level": "info", "message": "request", "fields": { "method": "GET", "path": "/login" } }
{ "level": "info", "message": "request", "fields": { "method": "POST", "path": "/login" } }
`

	_, err := io.Copy(w, strings.NewReader(input))
	assert.NoError(t, err, "copy")

	expected := `{"fields":{"method":"GET","path":"/"},"level":"info","timestamp":"1970-01-01T00:00:00Z","message":"request"}
{"fields":{"method":"GET","path":"/login"},"level":"info","timestamp":"1970-01-01T00:00:00Z","message":"request"}
{"fields":{"method":"POST","path":"/login"},"level":"info","timestamp":"1970-01-01T00:00:00Z","message":"request"}
`

	assert.Equal(t, expected, buf.String())
}
