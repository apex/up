package main

import (
	"errors"
	"os"
	"runtime"
	"strings"

	"github.com/tj/go/term"

	// commands
	_ "github.com/apex/up/internal/cli/config"
	_ "github.com/apex/up/internal/cli/deploy"
	_ "github.com/apex/up/internal/cli/docs"
	_ "github.com/apex/up/internal/cli/logs"
	_ "github.com/apex/up/internal/cli/metrics"
	_ "github.com/apex/up/internal/cli/run"
	_ "github.com/apex/up/internal/cli/stack"
	_ "github.com/apex/up/internal/cli/start"
	_ "github.com/apex/up/internal/cli/upgrade"
	_ "github.com/apex/up/internal/cli/url"
	_ "github.com/apex/up/internal/cli/version"

	"github.com/apex/up/internal/cli/app"
	"github.com/apex/up/internal/signals"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
)

var version = "master"

func main() {
	signals.Capture()
	signals.AddCloser(func(_ os.Signal) {
		// Always exit w/clean terminal
		term.ShowCursor()
		println("\n")
		os.Exit(1)
	})

	err := run()

	if err == nil {
		return
	}

	term.ShowCursor()

	if strings.Contains(err.Error(), "NoCredentialProviders") {
		util.Fatal(errors.New("Cannot find credentials, visit https://github.com/apex/up/blob/master/docs/aws-credentials.md for help."))
	}

	util.Fatal(err)
}

// run the cli.
func run() error {
	stats.SetProperties(map[string]interface{}{
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
		"version": version,
	})

	return app.Run(version)
}
