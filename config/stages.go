package config

import (
	"github.com/apex/up/internal/validate"
	"github.com/pkg/errors"
)

// Stage config.
type Stage struct {
	Domain string `json:"domain"`
	Path   string `json:"path"`
	Cert   string `json:"cert"`
	Name   string `json:"-"`
}

// Stages config.
type Stages struct {
	Development *Stage `json:"development"`
	Staging     *Stage `json:"staging"`
	Production  *Stage `json:"production"`
}

// Default implementation.
func (s *Stages) Default() error {
	if s := s.Development; s != nil {
		s.Name = "development"
	}

	if s := s.Staging; s != nil {
		s.Name = "staging"
	}

	if s := s.Production; s != nil {
		s.Name = "production"
	}

	return nil
}

// Validate implementation.
func (s *Stages) Validate() error {
	if s := s.Development; s != nil {
		if err := validate.RequiredString(s.Domain); err != nil {
			return errors.Wrap(err, ".development: .domain")
		}
	}

	if s := s.Staging; s != nil {
		if err := validate.RequiredString(s.Domain); err != nil {
			return errors.Wrap(err, ".staging: .domain")
		}
	}

	if s := s.Production; s != nil {
		if err := validate.RequiredString(s.Domain); err != nil {
			return errors.Wrap(err, ".production: .domain")
		}
	}

	return nil
}

// List returns configured stages.
func (s *Stages) List() (v []*Stage) {
	if s := s.Development; s != nil {
		v = append(v, s)
	}

	if s := s.Staging; s != nil {
		v = append(v, s)
	}

	if s := s.Production; s != nil {
		v = append(v, s)
	}

	return
}

// Domains returns configured domains.
func (s *Stages) Domains() (v []string) {
	for _, s := range s.List() {
		if s.Domain != "" {
			v = append(v, s.Domain)
		}
	}

	return
}

// GetByDomain returns the stage by domain or nil.
func (s *Stages) GetByDomain(domain string) *Stage {
	for _, s := range s.List() {
		if s.Domain == domain {
			return s
		}
	}
	return nil
}
