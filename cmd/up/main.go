package main

import (
	"errors"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/stripe/stripe-go"
	"github.com/tj/go/env"
	"github.com/tj/go/term"

	// commands
	_ "github.com/apex/up/internal/cli/account"
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
	_ "github.com/apex/up/internal/cli/upgrade"
	_ "github.com/apex/up/internal/cli/url"
	_ "github.com/apex/up/internal/cli/version"

	"github.com/apex/up/internal/cli/app"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
)

var version = "master"

// TODO: yuck refactor this

func main() {
	trap()

	stripe.Key = env.GetDefault("STRIPE_KEY", "pk_live_23pGrHcZ2QpfX525XYmiyzmx")
	stripe.LogLevel = 0

	err := run()

	if err == nil {
		return
	}

	term.ShowCursor()

	if strings.Contains(err.Error(), "NoCredentialProviders") {
		util.Fatal(errors.New("Cannot find credentials, visit https://up.docs.apex.sh/#aws_credentials for help."))
	}

	util.Fatal(err)
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

// trap signals.
func trap() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	// TODO: abort with context
	go func() {
		<-sigs
		term.ShowCursor()
		println("\n")
		os.Exit(1)
	}()
}
