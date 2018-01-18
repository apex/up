package up

import (
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/apex/log"
	"github.com/pkg/errors"

	"github.com/apex/up/config"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/platform/event"
)

// Config for a project.
type Config = config.Config

// ReadConfig reads the configuration from `path`.
var ReadConfig = config.ReadConfig

// ParseConfigString returns config from JSON string.
var ParseConfigString = config.ParseConfigString

// MustParseConfigString returns config from JSON string.
var MustParseConfigString = config.MustParseConfigString

// Project manager.
type Project struct {
	Platform
	config *Config
	events event.Events
}

// New project.
func New(c *Config, events event.Events) *Project {
	return &Project{
		config: c,
		events: events,
	}
}

// WithPlatform to `platform`.
func (p *Project) WithPlatform(platform Platform) *Project {
	p.Platform = platform
	return p
}

// RunHook runs a hook by name.
func (p *Project) RunHook(name string) error {
	hook := p.config.Hooks.Get(name)

	if hook.IsEmpty() {
		log.Debugf("hook %s is not defined", name)
		return nil
	}

	defer p.events.Time("hook", event.Fields{
		"name": name,
		"hook": hook,
	})()

	for _, command := range hook {
		log.Debugf("hook %q command %q", name, command)

		cmd := exec.Command("sh", "-c", command)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, util.Env(p.config.Environment)...)
		cmd.Env = append(cmd.Env, "PATH=node_modules/.bin:"+os.Getenv("PATH"))

		b, err := cmd.CombinedOutput()
		if err != nil {
			return errors.Errorf("%q: %s", command, b)
		}
	}

	return nil
}

// RunHooks runs hooks by name.
func (p *Project) RunHooks(names ...string) error {
	for _, n := range names {
		if err := p.RunHook(n); err != nil {
			return errors.Wrapf(err, "%q hook", n)
		}
	}
	return nil
}

// Build the project.
func (p *Project) Build() error {
	defer p.events.Time("platform.build", nil)()

	if err := p.RunHooks("prebuild", "build"); err != nil {
		return err
	}

	if err := p.Platform.Build(); err != nil {
		return errors.Wrap(err, "building")
	}

	if err := p.RunHooks("postbuild"); err != nil {
		return err
	}

	return nil
}

// Deploy the project.
func (p *Project) Deploy(stage string) error {
	defer p.events.Time("deploy", nil)()

	if err := p.Build(); err != nil {
		return errors.Wrap(err, "building")
	}

	if err := p.deploy(stage); err != nil {
		return errors.Wrap(err, "deploying")
	}

	if err := p.RunHook("clean"); err != nil {
		return errors.Wrap(err, "clean hook")
	}

	return nil
}

// deploy stage.
func (p *Project) deploy(stage string) error {
	if err := p.RunHooks("predeploy", "deploy"); err != nil {
		return err
	}

	if err := p.Platform.Deploy(stage); err != nil {
		return err
	}

	if err := p.RunHooks("postdeploy"); err != nil {
		return err
	}

	return nil
}

// Zip returns the zip if supported by the platform.
func (p *Project) Zip() (io.Reader, error) {
	z, ok := p.Platform.(Zipper)
	if !ok {
		return nil, errors.Errorf("platform does not support zips")
	}

	return z.Zip(), nil
}

// Init initializes the runtime such as remote environment variables.
func (p *Project) Init(stage string) error {
	r, ok := p.Platform.(Runtime)
	if !ok {
		return nil
	}

	return r.Init(stage)
}

// CreateStack implementation.
func (p *Project) CreateStack(region, version string) error {
	defer p.events.Time("stack.create", event.Fields{
		"region":  region,
		"version": version,
	})()

	return p.Platform.CreateStack(region, version)
}

// DeleteStack implementation.
func (p *Project) DeleteStack(region string, wait bool) error {
	defer p.events.Time("stack.delete", event.Fields{
		"region": region,
	})()

	return p.Platform.DeleteStack(region, wait)
}

// ShowStack implementation.
func (p *Project) ShowStack(region string) error {
	defer p.events.Time("stack.show", event.Fields{
		"region": region,
	})()

	return p.Platform.ShowStack(region)
}

// ShowMetrics implementation.
func (p *Project) ShowMetrics(region, stage string, start time.Time) error {
	defer p.events.Time("metrics", event.Fields{
		"region": region,
		"stage":  stage,
		"start":  start,
	})()

	return p.Platform.ShowMetrics(region, stage, start)
}

// PlanStack implementation.
func (p *Project) PlanStack(region string) error {
	defer p.events.Time("stack.plan", event.Fields{
		"region": region,
	})()

	return p.Platform.PlanStack(region)
}

// ApplyStack implementation.
func (p *Project) ApplyStack(region string) error {
	defer p.events.Time("stack.apply", event.Fields{
		"region": region,
	})()

	return p.Platform.ApplyStack(region)
}
