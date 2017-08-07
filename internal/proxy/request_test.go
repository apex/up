package proxy

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/tj/assert"
)

func TestNewRequest(t *testing.T) {
	t.Run("GET", func(t *testing.T) {
		var in Input
		err := json.Unmarshal([]byte(getEvent), &in)
		assert.NoError(t, err, "unmarshal")

		req, err := NewRequest(&in)
		assert.NoError(t, err, "new request")

		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, "apex-ping.com", req.Host)
		assert.Equal(t, "/pets/tobi", req.URL.Path)
		assert.Equal(t, "format=json", req.URL.Query().Encode())
		assert.Equal(t, "207.102.57.26", req.RemoteAddr)
	})

	t.Run("POST", func(t *testing.T) {
		var in Input
		err := json.Unmarshal([]byte(postEvent), &in)
		assert.NoError(t, err, "unmarshal")

		req, err := NewRequest(&in)
		assert.NoError(t, err, "new request")

		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "apex-ping.com", req.Host)
		assert.Equal(t, "/pets/tobi", req.URL.Path)
		assert.Equal(t, "", req.URL.Query().Encode())
		assert.Equal(t, "207.102.57.26", req.RemoteAddr)

		b, err := ioutil.ReadAll(req.Body)
		assert.NoError(t, err, "read body")

		assert.Equal(t, `{ "name": "Tobi" }`, string(b))
	})

	t.Run("POST binary", func(t *testing.T) {
		var in Input
		err := json.Unmarshal([]byte(postEventBinary), &in)
		assert.NoError(t, err, "unmarshal")

		req, err := NewRequest(&in)
		assert.NoError(t, err, "new request")

		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "/pets/tobi", req.URL.Path)
		assert.Equal(t, "", req.URL.Query().Encode())
		assert.Equal(t, "207.102.57.26", req.RemoteAddr)

		b, err := ioutil.ReadAll(req.Body)
		assert.NoError(t, err, "read body")

		assert.Equal(t, `Hello World`, string(b))
	})
}
