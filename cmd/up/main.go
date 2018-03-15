package main

import (
	"errors"
	"os"
	"runtime"

	"github.com/stripe/stripe-go"
	"github.com/tj/go/env"
	"github.com/tj/go/term"

	// commands
	_ "github.com/apex/up/internal/cli/build"
	_ "github.com/apex/up/internal/cli/config"
	_ "github.com/apex/up/internal/cli/deploy"
	_ "github.com/apex/up/internal/cli/docs"
	_ "github.com/apex/up/internal/cli/domains"
	_ "github.com/apex/up/internal/cli/logs"
	_ "github.com/apex/up/internal/cli/metrics"
	_ "github.com/apex/up/internal/cli/run"
	_ "github.com/apex/up/internal/cli/stack"
	_ "github.com/apex/up/internal/cli/start"
	_ "github.com/apex/up/internal/cli/team"
	_ "github.com/apex/up/internal/cli/upgrade"
	_ "github.com/apex/up/internal/cli/url"
	_ "github.com/apex/up/internal/cli/version"

	"github.com/apex/up/internal/cli/app"
	"github.com/apex/up/internal/signal"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
)

var version = "master"

func main() {
	signal.Add(reset)
	stripe.Key = env.GetDefault("STRIPE_KEY", "pk_live_23pGrHcZ2QpfX525XYmiyzmx")
	stripe.LogLevel = 0

	err := run()

	if err == nil {
		return
	}

	term.ShowCursor()

	switch {
	case util.IsNoCredentials(err):
		util.Fatal(errors.New("Cannot find credentials, visit https://up.docs.apex.sh/#aws_credentials for help."))
	default:
		util.Fatal(err)
	}
}

// run the cli.
func run() error {
	stats.SetProperties(map[string]interface{}{
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
		"version": version,
		"ci":      os.Getenv("CI") == "true" || os.Getenv("CI") == "1",
	})

	return app.Run(version)
}

// reset cursor.
func reset() error {
	term.ShowCursor()
	println()
	return nil
}
