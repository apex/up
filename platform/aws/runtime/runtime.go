package runtime

import (
	"os"
	"time"

	"github.com/apex/log"
	"github.com/apex/up"
	"github.com/apex/up/internal/secret"
	"github.com/apex/up/internal/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
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

	r.log.Info("loading secrets")
	if err := r.loadSecrets(stage); err != nil {
		return errors.Wrap(err, "loading secrets")
	}

	return nil
}

// Metric records a metric value.
func (r *Runtime) Metric(name string, value float64) error {
	// TODO: move
	c := cloudwatch.New(session.New(aws.NewConfig()))

	// TODO: conventions for Name? ByApp ?
	// TODO: timeouts or delegate
	// TODO: stage dim?
	_, err := c.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace: aws.String("up"),
		MetricData: []*cloudwatch.MetricDatum{
			{
				MetricName: &name,
				Value:      &value,
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("app"),
						Value: &r.config.Name,
					},
				},
			},
		},
	})

	return err
}

// loadSecrets loads secrets.
func (r *Runtime) loadSecrets(stage string) error {
	start := time.Now()
	initialEnv := util.EnvironMap()

	r.log.Info("initializing secrets")
	defer func() {
		r.log.WithField("duration", util.MillisecondsSince(start)).Info("initialized secrets")
	}()

	// TODO: all regions
	secrets, err := NewSecrets(r.config.Name, stage, r.config.Regions[0]).Load()
	if err != nil {
		return err
	}

	secrets = secret.FilterByApp(secrets, r.config.Name)
	stages := secret.GroupByStage(secrets)

	precedence := []string{
		"all",
		stage,
	}

	for _, name := range precedence {
		if secrets := stages[name]; len(secrets) > 0 {
			r.log.WithFields(log.Fields{
				"name":  name,
				"count": len(secrets),
			}).Info("initializing variables")

			for _, s := range secrets {
				ctx := r.log.WithFields(log.Fields{
					"name":  s.Name,
					"value": secret.String(s),
				})

				// in development we allow existing vars to override `up env`
				if _, ok := initialEnv[s.Name]; ok && stage == "development" {
					ctx.Debug("variable already defined")
					continue
				}

				ctx.Info("set variable")
				os.Setenv(s.Name, s.Value)
			}
		}
	}

	return nil
}
