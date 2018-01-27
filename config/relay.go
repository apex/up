package config

import (
	"github.com/pkg/errors"
)

// Relay config.
type Relay struct {
	// Command run to start your server.
	Command string `json:"command"`

	// Backoff configuration.
	Backoff Backoff `json:"backoff"`

	// Retry enables idempotent request retries. Enabled by default.
	Retry *bool `json:"retry"`

	// Timeout in seconds to wait for a response.
	// This is also taken into account when performing
	// retries, as to not exceed the limit.
	Timeout int `json:"timeout"`

	// ListenTimeout in seconds when waiting for
	// the application to bind to PORT.
	ListenTimeout int `json:"listen_timeout"`

	// ShutdownTimeout in seconds to wait after
	// sending a SIGINT before sending a SIGKILL.
	ShutdownTimeout int `json:"shutdown_timeout"`

	// platform is a currently unexported designation of the target deploy platform for this Relay
	platform string
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

	if r.ShutdownTimeout == 0 {
		r.ShutdownTimeout = 15
	}

	if err := r.Backoff.Default(); err != nil {
		return errors.Wrap(err, ".backoff")
	}

	if r.Retry != nil && !*r.Retry {
		r.Backoff.Attempts = 0
	}

	if r.platform == "" {
		r.platform = "lambda"
	}

	return nil
}

// Validate will try to perform sanity checks for this Relay configuration.
func (r *Relay) Validate() error {
	if r.Command == "" {
		err := errors.New("should not be empty")
		return errors.Wrap(err, ".command")
	}

	if r.platform != "lambda" {
		err := errors.New("internal consistency error, platform should be lambda")
		return errors.Wrap(err, ".platform")
	}

	if r.ListenTimeout <= 0 {
		err := errors.New("should be greater than 0")
		return errors.Wrap(err, ".listen_timeout")
	}

	if r.platform == "lambda" && r.ListenTimeout > 25 {
		err := errors.New("should be <= 25")
		return errors.Wrap(err, ".listen_timeout")
	}

	if r.platform == "lambda" && r.Timeout > 25 {
		err := errors.New("should be <= 25")
		return errors.Wrap(err, ".timeout")
	}

	if r.ShutdownTimeout < 0 {
		err := errors.New("should be greater than 0")
		return errors.Wrap(err, ".shutdown_timeout")
	}

	return nil
}

// Override config.
func (r *Relay) Override(c *Config) {
	if r.Command != "" {
		c.Proxy.Command = r.Command
	}
}
