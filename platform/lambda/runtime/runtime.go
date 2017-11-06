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
}

// New platform.
func New(c *up.Config) *Runtime {
	return &Runtime{
		config: c,
	}
}

// Init implementation.
func (r *Runtime) Init(stage string) error {
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
