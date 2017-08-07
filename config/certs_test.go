package config

import (
	"testing"

	"github.com/tj/assert"
)

func TestCert_Validate(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		c := &Cert{}
		assert.EqualError(t, c.Validate(), `.domains: must have at least 1 value`)
	})
}

func TestCerts_Validate(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		c := Certs{
			Cert{Domains: []string{"apex.sh"}},
			Cert{},
		}

		assert.EqualError(t, c.Validate(), `cert 1: .domains: must have at least 1 value`)
	})
}
