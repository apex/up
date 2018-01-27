package config

// Stage config.
type Stage struct {
	Domain string `json:"domain"`
	Path   string `json:"path"`
	Cert   string `json:"cert"`
	Name   string `json:"-"`
	StageOverrides
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

// GetByName returns the stage by name or nil.
func (s *Stages) GetByName(name string) *Stage {
	for _, s := range s.List() {
		if s.Name == name {
			return s
		}
	}
	return nil
}
