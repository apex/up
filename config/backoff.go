package config

import (
	"time"
	
	"github.com/tj/backoff"
)

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
