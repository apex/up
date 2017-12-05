package up

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/pkg/errors"

	"github.com/apex/up/config"
	"github.com/apex/up/internal/header"
	"github.com/apex/up/internal/inject"
	"github.com/apex/up/internal/redirect"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/internal/validate"
	"github.com/apex/up/platform/lambda/regions"
	"github.com/aws/aws-sdk-go/aws/session"
)

// TODO: refactor defaulting / validation with slices

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
	// Name of the project.
	Name string `json:"name"`

	// Description of the project.
	Description string `json:"description"`

	// Type of project.
	Type string `json:"type"`

	// Headers injection rules.
	Headers header.Rules `json:"headers"`

	// Redirects redirection rules.
	Redirects redirect.Rules `json:"redirects"`

	// Hooks defined for the project.
	Hooks config.Hooks `json:"hooks"`

	// Environment variables.
	Environment config.Environment `json:"environment"`

	// Regions is a list of regions to deploy to.
	Regions []string `json:"regions"`

	// Profile is the AWS profile name to reference for credentials.
	Profile string `json:"profile"`

	// Inject rules.
	Inject inject.Rules `json:"inject"`

	// Lambda provider configuration.
	Lambda config.Lambda `json:"lambda"`

	// CORS config.
	CORS *config.CORS `json:"cors"`

	// ErrorPages config.
	ErrorPages config.ErrorPages `json:"error_pages"`

	// Proxy config.
	Proxy config.Relay `json:"proxy"`

	// Static config.
	Static config.Static `json:"static"`

	// Logs config.
	Logs config.Logs `json:"logs"`

	// Stages config.
	Stages config.Stages `json:"stages"`

	// DNS config.
	DNS config.DNS `json:"dns"`
}

// Validate implementation.
func (c *Config) Validate() error {
	if err := validate.Name(c.Name); err != nil {
		return errors.Wrapf(err, ".name %q", c.Name)
	}

	if err := validate.List(c.Type, []string{"static", "server"}); err != nil {
		return errors.Wrap(err, ".type")
	}

	if err := validate.Lists(c.Regions, regions.All); err != nil {
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

	return nil
}

// Default implementation.
func (c *Config) Default() error {
	// TODO: hack, move to the instantiation of aws clients
	if c.Profile != "" {
		os.Setenv("AWS_PROFILE", c.Profile)
	}

	// default type to server
	if c.Type == "" {
		c.Type = "server"
	}

	// runtime defaults
	if c.Type != "static" {
		if err := c.inferRuntime(); err != nil {
			return errors.Wrap(err, "runtime")
		}
	}

	// default .name
	if err := c.defaultName(); err != nil {
		return errors.Wrap(err, ".name")
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

	// default .inject
	if err := c.Inject.Default(); err != nil {
		return errors.Wrap(err, ".inject")
	}

	// default .static
	if err := c.Static.Default(); err != nil {
		return errors.Wrap(err, ".static")
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

// inferRuntime performs inferences based on what Up thinks the runtime is.
func (c *Config) inferRuntime() error {
	switch {
	case util.Exists("main.go"):
		golang(c)
	case util.Exists("pom.xml"):
		javaMaven(c)
	case util.Exists("build.gradle"):
		javaGradle(c)
	case util.Exists("main.cr"):
		crystal(c)
	case util.Exists("package.json"):
		if err := nodejs(c); err != nil {
			return err
		}
	case util.Exists("app.js"):
		c.Proxy.Command = "node app.js"
	case util.Exists("app.py"):
		python(c)
	case util.Exists("index.html"):
		c.Type = "static"
	}
	return nil
}

// defaultName infers the name from the CWD if it's not set.
func (c *Config) defaultName() error {
	if c.Name != "" {
		return nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	c.Name = filepath.Base(dir)
	log.Debugf("infer name from current working directory %q", c.Name)
	return nil
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

	if os.IsNotExist(err) {
		c := &Config{}

		if err := c.Default(); err != nil {
			return nil, errors.Wrap(err, "defaulting")
		}

		if err := c.Validate(); err != nil {
			return nil, errors.Wrap(err, "validating")
		}

		return c, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "reading file")
	}

	return ParseConfig(b)
}

// golang config.
func golang(c *Config) {
	if c.Hooks.Build.IsEmpty() {
		c.Hooks.Build = config.Hook{`GOOS=linux GOARCH=amd64 go build -o server *.go`}
	}

	if c.Hooks.Clean.IsEmpty() {
		c.Hooks.Clean = config.Hook{`rm server`}
	}
}

// java gradle config.
func javaGradle(c *Config) {
	if c.Proxy.Command == "" {
		c.Proxy.Command = "java -jar server.jar"
	}

	if c.Hooks.Build.IsEmpty() {
		// assumes build results in a shaded jar named server.jar
		if util.Exists("gradlew") {
			c.Hooks.Build = config.Hook{`./gradlew clean build && cp build/libs/server.jar .`}
		} else {
			c.Hooks.Build = config.Hook{`gradle clean build && cp build/libs/server.jar .`}
		}
	}

	if c.Hooks.Clean.IsEmpty() {
		c.Hooks.Clean = config.Hook{`rm server.jar && gradle clean`}
	}
}

// java maven config.
func javaMaven(c *Config) {
	if c.Proxy.Command == "" {
		c.Proxy.Command = "java -jar server.jar"
	}

	if c.Hooks.Build.IsEmpty() {
		// assumes package results in a shaded jar named server.jar
		if util.Exists("mvnw") {
			c.Hooks.Build = config.Hook{`./mvnw clean package && cp target/server.jar .`}
		} else {
			c.Hooks.Build = config.Hook{`mvn clean package && cp target/server.jar .`}
		}
	}

	if c.Hooks.Clean.IsEmpty() {
		c.Hooks.Clean = config.Hook{`rm server.jar && mvn clean`}
	}
}

// crystal config.
func crystal(c *Config) {
	if c.Hooks.Build.IsEmpty() {
		c.Hooks.Build = config.Hook{`docker run --rm -v $(PWD):/src -w /src tjholowaychuk/up-crystal crystal build --link-flags -static -o server main.cr`}
	}

	if c.Hooks.Clean.IsEmpty() {
		c.Hooks.Clean = config.Hook{`rm server`}
	}
}

// nodejs config.
func nodejs(c *Config) error {
	var pkg struct {
		Scripts struct {
			Start string `json:"start"`
			Build string `json:"build"`
		} `json:"scripts"`
	}

	// read package.json
	if err := util.ReadFileJSON("package.json", &pkg); err != nil {
		return err
	}

	// use "start" script unless explicitly defined in up.json
	if c.Proxy.Command == "" {
		if s := pkg.Scripts.Start; s == "" {
			c.Proxy.Command = `node app.js`
		} else {
			c.Proxy.Command = s
		}
	}

	// use "build" script unless explicitly defined in up.json
	if c.Hooks.Build.IsEmpty() {
		c.Hooks.Build = config.Hook{pkg.Scripts.Build}
	}

	return nil
}

// python config.
func python(c *Config) {
	if c.Proxy.Command == "" {
		c.Proxy.Command = "python app.py"
	}

	// Only add build & clean hooks if a requirements.txt exists
	if !util.Exists("requirements.txt") {
		return
	}

	// Set PYTHONPATH env
	if c.Environment == nil {
		c.Environment = config.Environment{}
	}
	c.Environment["PYTHONPATH"] = ".pypath/"

	// Copy libraries into .pypath/
	if c.Hooks.Build.IsEmpty() {
		c.Hooks.Build = config.Hook{`mkdir -p .pypath/ && pip install -r requirements.txt -t .pypath/`}
	}

	// Clean .pypath/
	if c.Hooks.Clean.IsEmpty() {
		c.Hooks.Clean = config.Hook{`rm -r .pypath/`}
	}
}
