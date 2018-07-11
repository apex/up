package config

import (
	"github.com/pkg/errors"
)

// Relay config.
type Relay struct {
	// Command run to start your server.
	Command string `json:"command"`

	// Timeout in seconds to wait for a response.
	Timeout int `json:"timeout"`

	// ListenTimeout in seconds when waiting for the app to bind to PORT.
	ListenTimeout int `json:"listen_timeout"`
}

// Default implementation.
func (r *Relay) Default() error {
	if r.Command == "" {
		r.Command = "./server"
	}

	if r.Timeout == 0 {
		r.Timeout = 15
	}

	if r.ListenTimeout == 0 {
		r.ListenTimeout = 15
	}

	return nil
}

// Validate will try to perform sanity checks for this Relay configuration.
func (r *Relay) Validate() error {
	if r.Command == "" {
		err := errors.New("should not be empty")
		return errors.Wrap(err, ".command")
	}

	if r.ListenTimeout <= 0 {
		err := errors.New("should be greater than 0")
		return errors.Wrap(err, ".listen_timeout")
	}

	if r.ListenTimeout > 25 {
		err := errors.New("should be <= 25")
		return errors.Wrap(err, ".listen_timeout")
	}

	if r.Timeout > 25 {
		err := errors.New("should be <= 25")
		return errors.Wrap(err, ".timeout")
	}

	return nil
}

// Override config.
func (r *Relay) Override(c *Config) {
	if r.Command != "" {
		c.Proxy.Command = r.Command
	}
}
