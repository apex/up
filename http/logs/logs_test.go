package logs

import (
	"bytes"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/tj/assert"

	"github.com/apex/up"
	"github.com/apex/up/config"
	"github.com/apex/up/http/static"
)

func TestLogs(t *testing.T) {
	// TODO: refactor and pass in app name/version/region

	var buf bytes.Buffer
	log.SetOutput(&buf)

	c := &up.Config{
		Static: config.Static{
			Dir: "testdata",
		},
	}

	h, err := New(c, static.New(c))
	assert.NoError(t, err)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/?foo=bar", nil)

	h.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "text/html; charset=utf-8", res.Header().Get("Content-Type"))
	assert.Equal(t, "Index HTML\n", res.Body.String())

	s := buf.String()
	assert.Contains(t, s, `info response`)
	// assert.Contains(t, s, `app_name=api`)
	// assert.Contains(t, s, `app_version=5`)
	// assert.Contains(t, s, `app_region=us-west-2`)
	assert.Contains(t, s, `ip=192.0.2.1:1234`)
	assert.Contains(t, s, `method=GET`)
	assert.Contains(t, s, `path=/`)
	assert.Contains(t, s, `plugin=logs`)
	assert.Contains(t, s, `size=11`)
	assert.Contains(t, s, `status=200`)
}
