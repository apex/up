package runtime

import (
	"os"
	"time"

	"github.com/apex/log"
	"github.com/apex/up"
	"github.com/apex/up/internal/secret"
	"github.com/apex/up/internal/util"
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
	if err := r.loadSecrets(stage); err != nil {
		return errors.Wrap(err, "loading secrets")
	}

	return nil
}

// loadSecrets loads secrets.
func (r *Runtime) loadSecrets(stage string) error {
	start := time.Now()

	log.Info("initializing secrets")
	defer func() {
		log.WithField("duration", util.MillisecondsSince(start)).Info("initialized secrets")
	}()

	// TODO: all regions
	secrets, err := NewSecrets(r.config.Name, stage, r.config.Regions[0]).Load()
	if err != nil {
		return err
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
			log.Infof("loaded %d variables for %s stage(s)", len(secrets), name)
			for _, s := range secrets {
				os.Setenv(s.Name, s.Value)
			}
		}
	}

	return nil
}
