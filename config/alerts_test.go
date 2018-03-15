package config

import (
	"testing"

	"github.com/tj/assert"
)

func TestAlert_Validate(t *testing.T) {
	t.Run("invalid operator", func(t *testing.T) {
		a := Alert{
			Operator: "===",
		}

		assert.EqualError(t, a.Validate(), `.operator: "===" is invalid, must be one of:

  • <
  • <=
  • >
  • >=`)
	})

	t.Run("invalid statistic", func(t *testing.T) {
		a := Alert{
			Operator:  ">=",
			Statistic: "minimumm",
		}

		assert.EqualError(t, a.Validate(), `.statistic: "minimumm" is invalid, must be one of:

  • average
  • avg
  • count
  • max
  • maximum
  • min
  • minimum
  • sum`)
	})

	t.Run("namespace explicit", func(t *testing.T) {
		a := Alert{
			Metric:    "5XXError",
			Statistic: "minimum",
			Namespace: "AWS/ApiGateway",
		}

		assert.NoError(t, a.Default(), "default")
		assert.NoError(t, a.Validate(), "default")

		assert.Equal(t, "AWS/ApiGateway", a.Namespace)
		assert.Equal(t, "GreaterThanThreshold", a.Operator)
		assert.Equal(t, "5XXError", a.Metric)
	})

	t.Run("namespace api", func(t *testing.T) {
		a := &Alert{
			Metric:    "http.5xx",
			Statistic: "min",
		}

		assert.NoError(t, a.Default(), "default")
		assert.NoError(t, a.Validate(), "default")

		assert.Equal(t, "AWS/ApiGateway", a.Namespace)
		assert.Equal(t, "GreaterThanThreshold", a.Operator)
		assert.Equal(t, "5XXError", a.Metric)
	})
}
