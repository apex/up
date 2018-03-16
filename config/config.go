package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/apex/log"
	"github.com/pkg/errors"

	"github.com/apex/up/internal/header"
	"github.com/apex/up/internal/inject"
	"github.com/apex/up/internal/redirect"
	"github.com/apex/up/internal/validate"
	"github.com/apex/up/platform/aws/regions"
	"github.com/aws/aws-sdk-go/aws/session"
)

// defaulter is the interface that provides config defaulting.
type defaulter interface {
	Default() error
}

// validator is the interface that provides config validation.
type validator interface {
	Validate() error
}

// Config for the project.
type Config struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        string         `json:"type"`
	Headers     header.Rules   `json:"headers"`
	Redirects   redirect.Rules `json:"redirects"`
	Hooks       Hooks          `json:"hooks"`
	Environment Environment    `json:"environment"`
	Regions     []string       `json:"regions"`
	Profile     string         `json:"profile"`
	Inject      inject.Rules   `json:"inject"`
	Lambda      Lambda         `json:"lambda"`
	CORS        *CORS          `json:"cors"`
	ErrorPages  ErrorPages     `json:"error_pages"`
	Proxy       Relay          `json:"proxy"`
	Static      Static         `json:"static"`
	Logs        Logs           `json:"logs"`
	Stages      Stages         `json:"stages"`
	DNS         DNS            `json:"dns"`
}

// Validate implementation.
func (c *Config) Validate() error {
	if err := validate.RequiredString(c.Name); err != nil {
		return errors.Wrap(err, ".name")
	}

	if err := validate.Name(c.Name); err != nil {
		return errors.Wrapf(err, ".name %q", c.Name)
	}

	if err := validate.List(c.Type, []string{"static", "server"}); err != nil {
		return errors.Wrap(err, ".type")
	}

	if err := validate.Lists(c.Regions, regions.IDs); err != nil {
		return errors.Wrap(err, ".regions")
	}

	if err := c.DNS.Validate(); err != nil {
		return errors.Wrap(err, ".dns")
	}

	if err := c.Static.Validate(); err != nil {
		return errors.Wrap(err, ".static")
	}

	if err := c.Inject.Validate(); err != nil {
		return errors.Wrap(err, ".inject")
	}

	if err := c.Lambda.Validate(); err != nil {
		return errors.Wrap(err, ".lambda")
	}

	if err := c.Proxy.Validate(); err != nil {
		return errors.Wrap(err, ".proxy")
	}

	if err := c.Stages.Validate(); err != nil {
		return errors.Wrap(err, ".stages")
	}

	if len(c.Regions) > 1 {
		return errors.New("multiple regions is not yet supported, see https://github.com/apex/up/issues/134")
	}

	return nil
}

// Default implementation.
func (c *Config) Default() error {
	if c.Stages == nil {
		c.Stages = make(Stages)
	}

	// we default stages here before others simply to
	// initialize the default stages such as "development"
	// allowing runtime inference to default values.
	if err := c.Stages.Default(); err != nil {
		return errors.Wrap(err, ".stages")
	}

	// TODO: hack, move to the instantiation of aws clients
	if c.Profile != "" {
		setProfile(c.Profile)
	}

	// default type to server
	if c.Type == "" {
		c.Type = "server"
	}

	// runtime defaults
	if c.Type != "static" {
		runtime := inferRuntime()
		log.WithField("type", runtime).Debug("inferred runtime")

		if err := runtimeConfig(runtime, c); err != nil {
			return errors.Wrap(err, "runtime")
		}
	}

	// default .regions
	if err := c.defaultRegions(); err != nil {
		return errors.Wrap(err, ".region")
	}

	// region globbing
	c.Regions = regions.Match(c.Regions)

	// default .proxy
	if err := c.Proxy.Default(); err != nil {
		return errors.Wrap(err, ".proxy")
	}

	// default .lambda
	if err := c.Lambda.Default(); err != nil {
		return errors.Wrap(err, ".lambda")
	}

	// default .dns
	if err := c.DNS.Default(); err != nil {
		return errors.Wrap(err, ".dns")
	}

	// default .logs
	if err := c.Logs.Default(); err != nil {
		return errors.Wrap(err, ".logs")
	}

	// default .inject
	if err := c.Inject.Default(); err != nil {
		return errors.Wrap(err, ".inject")
	}

	// default .error_pages
	if err := c.ErrorPages.Default(); err != nil {
		return errors.Wrap(err, ".error_pages")
	}

	// default .stages
	if err := c.Stages.Default(); err != nil {
		return errors.Wrap(err, ".stages")
	}

	return nil
}

// Override with stage config if present, and re-validate.
func (c *Config) Override(stage string) error {
	s := c.Stages.GetByName(stage)
	if s == nil {
		return nil
	}

	s.Override(c)

	return c.Validate()
}

// defaultRegions checks AWS_REGION and falls back on us-west-2.
func (c *Config) defaultRegions() error {
	if len(c.Regions) != 0 {
		log.Debugf("%d regions from config", len(c.Regions))
		return nil
	}

	s, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		return errors.Wrap(err, "creating session")
	}

	if r := *s.Config.Region; r != "" {
		log.Debugf("region from aws shared config %q", r)
		c.Regions = append(c.Regions, r)
		return nil
	}

	r := "us-west-2"
	log.Debugf("region defaulted to %q", r)
	c.Regions = append(c.Regions, r)
	return nil
}

// ParseConfig returns config from JSON bytes.
func ParseConfig(b []byte) (*Config, error) {
	c := &Config{}

	if err := json.Unmarshal(b, c); err != nil {
		return nil, errors.Wrap(err, "parsing json")
	}

	if err := c.Default(); err != nil {
		return nil, errors.Wrap(err, "defaulting")
	}

	if err := c.Validate(); err != nil {
		return nil, errors.Wrap(err, "validating")
	}

	return c, nil
}

// ParseConfigString returns config from JSON string.
func ParseConfigString(s string) (*Config, error) {
	return ParseConfig([]byte(s))
}

// MustParseConfigString returns config from JSON string.
func MustParseConfigString(s string) *Config {
	c, err := ParseConfigString(s)
	if err != nil {
		panic(err)
	}

	return c
}

// ReadConfig reads the configuration from `path`.
func ReadConfig(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseConfig(b)
}

// setProfile sets the AWS_PROFILE.
func setProfile(name string) {
	os.Setenv("AWS_PROFILE", name)
}
