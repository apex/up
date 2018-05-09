package config

import (
	"sort"

	"github.com/apex/up/internal/validate"
	"github.com/pkg/errors"
)

// defaultStages is a list of default stage names.
var defaultStages = []string{
	"development",
	"staging",
	"production",
}

// Stage config.
type Stage struct {
	Domain string      `json:"domain"`
	Zone   interface{} `json:"zone"`
	Path   string      `json:"path"`
	Cert   string      `json:"cert"`
	Name   string      `json:"-"`
	StageOverrides
}

// IsLocal returns true if the stage represents a local environment.
func (s *Stage) IsLocal() bool {
	return s.Name == "development"
}

// IsRemote returns true if the stage represents a remote environment.
func (s *Stage) IsRemote() bool {
	return !s.IsLocal()
}

// Validate implementation.
func (s *Stage) Validate() error {
	if err := validate.Stage(s.Name); err != nil {
		return errors.Wrap(err, ".name")
	}

	switch s.Zone.(type) {
	case bool, string:
		return nil
	default:
		return errors.Errorf(".zone is an invalid type, must be string or boolean")
	}
}

// Default implementation.
func (s *Stage) Default() error {
	if s.Zone == nil {
		s.Zone = true
	}

	return nil
}

// StageOverrides config.
type StageOverrides struct {
	Hooks  Hooks  `json:"hooks"`
	Lambda Lambda `json:"lambda"`
	Proxy  Relay  `json:"proxy"`
}

// Override config.
func (s *StageOverrides) Override(c *Config) {
	s.Hooks.Override(c)
	s.Lambda.Override(c)
	s.Proxy.Override(c)
}

// Stages config.
type Stages map[string]*Stage

// Default implementation.
func (s Stages) Default() error {
	// defaults
	for _, name := range defaultStages {
		if _, ok := s[name]; !ok {
			s[name] = &Stage{}
		}
	}

	// assign names
	for name, s := range s {
		s.Name = name
	}

	// defaults
	for _, s := range s {
		if err := s.Default(); err != nil {
			return errors.Wrapf(err, "stage %q", s.Name)
		}
	}

	return nil
}

// Validate implementation.
func (s Stages) Validate() error {
	for _, s := range s {
		if err := s.Validate(); err != nil {
			return errors.Wrapf(err, "stage %q", s.Name)
		}
	}
	return nil
}

// List returns configured stages.
func (s Stages) List() (v []*Stage) {
	for _, s := range s {
		v = append(v, s)
	}

	return
}

// Domains returns configured domains.
func (s Stages) Domains() (v []string) {
	for _, s := range s.List() {
		if s.Domain != "" {
			v = append(v, s.Domain)
		}
	}

	return
}

// Names returns configured stage names.
func (s Stages) Names() (v []string) {
	for _, s := range s.List() {
		v = append(v, s.Name)
	}

	sort.Strings(v)
	return
}

// RemoteNames returns configured remote stage names.
func (s Stages) RemoteNames() (v []string) {
	for _, s := range s.List() {
		if s.IsRemote() {
			v = append(v, s.Name)
		}
	}

	sort.Strings(v)
	return
}

// GetByDomain returns the stage by domain or nil.
func (s Stages) GetByDomain(domain string) *Stage {
	for _, s := range s.List() {
		if s.Domain == domain {
			return s
		}
	}

	return nil
}

// GetByName returns the stage by name or nil.
func (s Stages) GetByName(name string) *Stage {
	for _, s := range s.List() {
		if s.Name == name {
			return s
		}
	}

	return nil
}
