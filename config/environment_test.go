package config

import (
	"encoding/json"
	"github.com/tj/assert"
	"os"
	"testing"
)

func TestEnvironment_UnmarshalJSON(t *testing.T) {
	t.Run("plain strings", func(t *testing.T) {
		s := `{
				"VAR_A": "VALUE_OF_A"
			}`

		var environment Environment

		err := json.Unmarshal([]byte(s), &environment)
		assert.NoError(t, err, "unmarshal")

		assert.Equal(t, "VALUE_OF_A", environment["VAR_A"])
	})

	t.Run("env var substitution", func(t *testing.T) {
		os.Setenv("EXISTING_ENV_VAR", "EXISTING_ENV_VAR_VALUE")
		s := `{
				"VAR_A": "$EXISTING_ENV_VAR"
			}`

		var environment Environment

		err := json.Unmarshal([]byte(s), &environment)
		assert.NoError(t, err, "unmarshal")

		assert.Equal(t, "EXISTING_ENV_VAR_VALUE", environment["VAR_A"])
		os.Unsetenv("EXISTING_ENV_VAR")
	})

	t.Run("parse error", func(t *testing.T) {
		s := `{
				"VAR_A": _INVALID_
			}`

		var environment Environment

		err := json.Unmarshal([]byte(s), &environment)
		assert.Error(t, err, "unmarshal")
	})
}
