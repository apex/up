package runtime

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/up"
)

// Runtime implementation.
type Runtime struct {
	config *up.Config
	log    log.Interface
}

// Option function.
type Option func(*Runtime)

// New with the given options.
func New(c *up.Config, options ...Option) *Runtime {
	var v Runtime
	v.config = c
	v.log = log.Log
	for _, o := range options {
		o(&v)
	}
	return &v
}

// WithLogger option.
func WithLogger(l log.Interface) Option {
	return func(v *Runtime) {
		v.log = l
	}
}

// Init implementation.
func (r *Runtime) Init(stage string) error {
	os.Setenv("UP_STAGE", stage)

	if os.Getenv("NODE_ENV") == "" {
		os.Setenv("NODE_ENV", stage)
	}

	return nil
}
