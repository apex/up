package config

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/tj/assert"
)

func TestDuration_UnmarshalJSON(t *testing.T) {
	t.Run("numeric seconds", func(t *testing.T) {
		s := `{
      "timeout": 5
    }`

		var c struct {
			Timeout Duration
		}

		err := json.Unmarshal([]byte(s), &c)
		assert.NoError(t, err, "unmarshal")

		assert.Equal(t, Duration(5*time.Second), c.Timeout)
	})

	t.Run("string duration", func(t *testing.T) {
		s := `{
      "timeout": "1.5m"
    }`

		var c struct {
			Timeout Duration
		}

		err := json.Unmarshal([]byte(s), &c)
		assert.NoError(t, err, "unmarshal")

		assert.Equal(t, Duration(90*time.Second), c.Timeout)
	})
}
