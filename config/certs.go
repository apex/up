package config

import (
	"github.com/pkg/errors"
	"github.com/apex/up/internal/validate"
)

// Cert config.
type Cert struct {
	Domains []string `json:"domains"`
}

// Validate implementation.
func (c *Cert) Validate() error {
	if err := validate.MinStrings(c.Domains, 1); err != nil {
		return errors.Wrap(err, ".domains")
	}

	return nil
}

// Certs config.
type Certs []Cert

// Validate implementation.
func (c Certs) Validate() error {
	for i, v := range c {
		if err := v.Validate(); err != nil {
			return errors.Wrapf(err, "cert %d", i)
		}
	}

	return nil
}
