package config

import (
	"time"

	"github.com/pkg/errors"
	"github.com/tj/backoff"
)

// Relay config.
type Relay struct {
	// Command run to start your server.
	Command string `json:"command"`

	// Backoff configuration.
	Backoff Backoff `json:"backoff"`

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

	if r.ListenTimeout == 0 {
		r.ListenTimeout = 15
	}

	if r.ShutdownTimeout == 0 {
		r.ShutdownTimeout = 15
	}

	if err := r.Backoff.Default(); err != nil {
		return errors.Wrap(err, ".backoff")
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

	if r.ShutdownTimeout < 0 {
		err := errors.New("should be greater than 0")
		return errors.Wrap(err, ".shutdown_timeout")
	}

	return nil
}

// Backoff config.
type Backoff struct {
	// Min time in milliseconds.
	Min int `json:"min"`

	// Max time in milliseconds.
	Max int `json:"max"`

	// Factor applied for every attempt.
	Factor float64 `json:"factor"`

	// Attempts performed before failing.
	Attempts int `json:"attempts"`

	// Jitter is applied when true.
	Jitter bool `json:"jitter"`
}

// Default implementation.
func (b *Backoff) Default() error {
	if b.Min == 0 {
		b.Min = 100
	}

	if b.Max == 0 {
		b.Max = 500
	}

	if b.Factor == 0 {
		b.Factor = 2
	}

	if b.Attempts == 0 {
		b.Attempts = 3
	}

	return nil
}

// Backoff returns the backoff from config.
func (b *Backoff) Backoff() *backoff.Backoff {
	return &backoff.Backoff{
		Min:    time.Duration(b.Min) * time.Millisecond,
		Max:    time.Duration(b.Max) * time.Millisecond,
		Factor: b.Factor,
		Jitter: b.Jitter,
	}
}
