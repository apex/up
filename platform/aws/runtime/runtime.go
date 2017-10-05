package runtime

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/up"
	"github.com/apex/up/internal/secret"
	"github.com/pkg/errors"
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

	log.Info("loading secrets")

	// TODO: all regions
	secrets, err := NewSecrets(r.config.Name, stage, r.config.Regions[0]).Load()
	if err != nil {
		return errors.Wrap(err, "loading secrets")
	}

	secrets = secret.FilterByApp(secrets, r.config.Name)
	stages := secret.GroupByStage(secrets)

	// TODO: util to de-dupe first
	precedence := []string{
		"all",
		stage,
	}

	for _, name := range precedence {
		if secrets := stages[name]; len(secrets) > 0 {
			log.Infof("  %s %d variables", name, len(secrets))
			for _, s := range secrets {
				log.Infof("    - %s", s.Name)
				os.Setenv(s.Name, s.Value)
			}
		}
	}

	return nil
}
