package config

import (
	"os"

	"github.com/apex/up/internal/util"
	"github.com/pkg/errors"
)

// Runtime is an app runtime.
type Runtime string

// Runtimes available.
const (
	RuntimeUnknown    Runtime = "unknown"
	RuntimeGo                 = "go"
	RuntimeNode               = "node"
	RuntimeClojure            = "clojure"
	RuntimeCrystal            = "crystal"
	RuntimePython             = "python"
	RuntimeStatic             = "static"
	RuntimeJavaMaven          = "java maven"
	RuntimeJavaGradle         = "java gradle"
)

// inferRuntime returns the runtime based on files present in the CWD.
func inferRuntime() Runtime {
	switch {
	case util.Exists("main.go"):
		return RuntimeGo
	case util.Exists("main.cr"):
		return RuntimeCrystal
	case util.Exists("package.json"):
		return RuntimeNode
	case util.Exists("app.js"):
		return RuntimeNode
	case util.Exists("project.clj"):
		return RuntimeClojure
	case util.Exists("pom.xml"):
		return RuntimeJavaMaven
	case util.Exists("build.gradle"):
		return RuntimeJavaGradle
	case util.Exists("app.py"):
		return RuntimePython
	case util.Exists("index.html"):
		return RuntimeStatic
	default:
		return RuntimeUnknown
	}
}

// runtimeConfig performs config inferences based on what Up thinks the runtime is.
func runtimeConfig(runtime Runtime, c *Config) error {
	switch runtime {
	case RuntimeGo:
		golang(c)
	case RuntimeClojure:
		clojureLein(c)
	case RuntimeJavaMaven:
		javaMaven(c)
	case RuntimeJavaGradle:
		javaGradle(c)
	case RuntimeCrystal:
		crystal(c)
	case RuntimePython:
		python(c)
	case RuntimeStatic:
		c.Type = "static"
	case RuntimeNode:
		if err := nodejs(c); err != nil {
			return err
		}
	}
	return nil
}

// golang config.
func golang(c *Config) {
	if c.Hooks.Build.IsEmpty() {
		c.Hooks.Build = Hook{`GOOS=linux GOARCH=amd64 go build -o server *.go`}
	}

	if c.Hooks.Clean.IsEmpty() {
		c.Hooks.Clean = Hook{`rm server`}
	}

	if s := c.Stages.GetByName("development"); s != nil {
		if s.Proxy.Command == "" {
			s.Proxy.Command = "go run *.go"
		}
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
			c.Hooks.Build = Hook{`./gradlew clean build && cp build/libs/server.jar .`}
		} else {
			c.Hooks.Build = Hook{`gradle clean build && cp build/libs/server.jar .`}
		}
	}

	if c.Hooks.Clean.IsEmpty() {
		if util.Exists("gradlew") {
			c.Hooks.Clean = Hook{`rm server.jar && ./gradlew clean`}
		} else {
			c.Hooks.Clean = Hook{`rm server.jar && gradle clean`}
		}
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
			c.Hooks.Build = Hook{`./mvnw clean package && cp target/server.jar .`}
		} else {
			c.Hooks.Build = Hook{`mvn clean package && cp target/server.jar .`}
		}
	}

	if c.Hooks.Clean.IsEmpty() {
		c.Hooks.Clean = Hook{`rm server.jar && mvn clean`}
	}
}

// clojure lein config.
func clojureLein(c *Config) {
	if c.Proxy.Command == "" {
		c.Proxy.Command = "java -jar server.jar"
	}

	if c.Hooks.Build.IsEmpty() {
		// assumes package results in a shaded jar named server.jar
		c.Hooks.Build = Hook{`lein uberjar && cp target/*-standalone.jar server.jar`}
	}

	if c.Hooks.Clean.IsEmpty() {
		c.Hooks.Clean = Hook{`lein clean && rm server.jar`}
	}
}

// crystal config.
func crystal(c *Config) {
	if c.Hooks.Build.IsEmpty() {
		c.Hooks.Build = Hook{`docker run --rm -v $(pwd):/src -w /src crystallang/crystal crystal build -o server main.cr --release --static`}
	}

	if c.Hooks.Clean.IsEmpty() {
		c.Hooks.Clean = Hook{`rm server`}
	}

	if s := c.Stages.GetByName("development"); s != nil {
		if s.Proxy.Command == "" {
			s.Proxy.Command = "crystal run main.cr"
		}
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
	if err := util.ReadFileJSON("package.json", &pkg); err != nil && !os.IsNotExist(errors.Cause(err)) {
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
		c.Hooks.Build = Hook{pkg.Scripts.Build}
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
		c.Environment = Environment{}
	}
	c.Environment["PYTHONPATH"] = ".pypath/"

	// Copy libraries into .pypath/
	if c.Hooks.Build.IsEmpty() {
		c.Hooks.Build = Hook{`mkdir -p .pypath/ && pip install -r requirements.txt -t .pypath/`}
	}

	// Clean .pypath/
	if c.Hooks.Clean.IsEmpty() {
		c.Hooks.Clean = Hook{`rm -r .pypath/`}
	}
}
