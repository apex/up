package poweredby

import (
	"net/http/httptest"
	"testing"

	"github.com/tj/assert"
	"github.com/apex/up"

	"github.com/apex/up/config"
	"github.com/apex/up/http/static"
)

func TestPoweredby(t *testing.T) {
	c := &up.Config{
		Static: config.Static{
			Dir: "testdata",
		},
	}

	h := New("up", static.New(c))

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	h.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "up", res.Header().Get("X-Powered-By"))
	assert.Equal(t, "text/html; charset=utf-8", res.Header().Get("Content-Type"))
	assert.Equal(t, "Index HTML\n", res.Body.String())
}
