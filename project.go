package up

import (
	"os"
	"os/exec"

	"github.com/apex/log"
	"github.com/pkg/errors"

	"github.com/apex/up/platform"
	"github.com/apex/up/platform/event"
)

// Project manager.
type Project struct {
	config   *Config
	platform platform.Interface
	events   event.Events
}

// New project.
func New(c *Config, events event.Events) *Project {
	return &Project{
		config: c,
		events: events,
	}
}

// WithPlatform to `platform`.
func (p *Project) WithPlatform(platform platform.Interface) *Project {
	p.platform = platform
	return p
}

// HookCommand returns a hook command by name or empty string.
func (p *Project) HookCommand(name string) string {
	switch name {
	case "build":
		return p.config.Hooks.Build
	case "clean":
		return p.config.Hooks.Clean
	default:
		return ""
	}
}

// RunHook runs a hook by name.
func (p *Project) RunHook(name string) error {
	command := p.HookCommand(name)

	if command == "" {
		log.Debugf("hook %s is not defined", name)
		return nil
	}

	defer p.events.Time("hook", event.Fields{
		"name": name,
	})()

	cmd := exec.Command("sh", "-c", command)
	cmd.Env = append(os.Environ(), "PATH=node_modules/.bin:"+os.Getenv("PATH"))
	b, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(b))
	}

	return nil
}

// Build the project.
func (p *Project) Build() error {
	defer p.events.Time("platform.build", nil)()

	if err := p.RunHook("build"); err != nil {
		return errors.Wrap(err, "build hook")
	}

	return p.platform.Build()
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
	return p.platform.Deploy(stage)
}

// Logs for the project.
func (p *Project) Logs(region, query string) platform.Logs {
	return p.platform.Logs(region, query)
}

// URL returns the endpoint.
func (p *Project) URL(region, stage string) (string, error) {
	return p.platform.URL(region, stage)
}

// CreateStack implementation.
func (p *Project) CreateStack(region, version string) error {
	defer p.events.Time("stack.create", event.Fields{
		"region":  region,
		"version": version,
	})()

	return p.platform.CreateStack(region, version)
}

// DeleteStack implementation.
func (p *Project) DeleteStack(region string, wait bool) error {
	defer p.events.Time("stack.delete", event.Fields{
		"region": region,
	})()

	return p.platform.DeleteStack(region, wait)
}

// ShowStack implementation.
func (p *Project) ShowStack(region string) error {
	defer p.events.Time("stack.show", event.Fields{
		"region": region,
	})()

	return p.platform.ShowStack(region)
}

// ShowMetrics implementation.
func (p *Project) ShowMetrics(region, stage string) error {
	defer p.events.Time("metrics", event.Fields{
		"region": region,
		"stage":  stage,
	})()

	return p.platform.ShowMetrics(region, stage)
}
