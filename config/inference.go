package config

import "github.com/apex/up/internal/util"

// inferRuntime performs config inferences based on what Up thinks the runtime is.
func inferRuntime(c *Config) error {
	switch {
	case util.Exists("main.go"):
		golang(c)
	case util.Exists("project.clj"):
		clojureLein(c)
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

// golang config.
func golang(c *Config) {
	if c.Hooks.Build.IsEmpty() {
		c.Hooks.Build = Hook{`GOOS=linux GOARCH=amd64 go build -o server *.go`}
	}

	if c.Hooks.Clean.IsEmpty() {
		c.Hooks.Clean = Hook{`rm server`}
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
		c.Hooks.Clean = Hook{`rm server.jar && gradle clean`}
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
		c.Hooks.Build = Hook{`docker run --rm -v $(PWD):/src -w /src tjholowaychuk/up-crystal crystal build --link-flags -static -o server main.cr`}
	}

	if c.Hooks.Clean.IsEmpty() {
		c.Hooks.Clean = Hook{`rm server`}
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
