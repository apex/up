package config

import (
	"encoding/json"
	"testing"

	"github.com/tj/assert"
)

func TestHook(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		s := []byte(`{}`)

		var c struct {
			Build Hook
		}

		err := json.Unmarshal(s, &c)
		assert.NoError(t, err, "unmarshal")

		assert.Equal(t, Hook(nil), c.Build)
	})

	t.Run("invalid type", func(t *testing.T) {
		s := []byte(`
      {
        "build": 5
      }
    `)

		var c struct {
			Build Hook
		}

		err := json.Unmarshal(s, &c)
		assert.EqualError(t, err, `hook must be a string or array of strings`)
	})

	t.Run("string", func(t *testing.T) {
		s := []byte(`
      {
        "build": "go build main.go"
      }
    `)

		var c struct {
			Build Hook
		}

		err := json.Unmarshal(s, &c)
		assert.NoError(t, err, "unmarshal")

		assert.Equal(t, Hook{"go build main.go"}, c.Build)
	})

	t.Run("array", func(t *testing.T) {
		s := []byte(`
      {
        "build": [
          "go build main.go",
          "browserify src/index.js > app.js"
        ]
      }
    `)

		var c struct {
			Build Hook
		}

		err := json.Unmarshal(s, &c)
		assert.NoError(t, err, "unmarshal")

		assert.Equal(t, Hook{
			"go build main.go",
			"browserify src/index.js > app.js",
		}, c.Build)
	})
}
