package config

import (
	"os"

	"github.com/pkg/errors"
)

// Static configuration.
type Static struct {
	// Dir containing static files.
	Dir string `json:"dir"`
}

// Default implementation.
func (s *Static) Default() error {
	if s.Dir == "" {
		s.Dir = "."
	}

	return nil
}

// Validate implementation.
func (s *Static) Validate() error {
	info, err := os.Stat(s.Dir)

	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return errors.Wrap(err, ".dir")
	}

	if !info.IsDir() {
		return errors.Errorf(".dir %s is not a directory", s.Dir)
	}

	return nil
}
