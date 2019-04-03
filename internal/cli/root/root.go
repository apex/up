package root

import (
	"os"
	"runtime"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/delta"
	"github.com/pkg/errors"
	"github.com/tj/kingpin"

	"github.com/apex/up"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/platform/event"
	"github.com/apex/up/platform/lambda"
	"github.com/apex/up/reporter"
)

// Cmd is the root command.
var Cmd = kingpin.New("up", "")

// Command registers a command.
var Command = Cmd.Command

// Init function.
var Init func() (*up.Config, *up.Project, error)

func init() {
	log.SetHandler(cli.Default)

	Cmd.Example(`up`, "Deploy the project to the staging environment.")
	Cmd.Example(`up deploy production`, "Deploy the project to the production stage.")
	Cmd.Example(`up url`, "Show the staging endpoint url.")
	Cmd.Example(`up logs -f`, "Tail project logs.")
	Cmd.Example(`up logs 'error or fatal'`, "Show error or fatal level logs.")
	Cmd.Example(`up run build`, "Run build command manually.")
	Cmd.Example(`up help team`, "Show help and examples for a command.")
	Cmd.Example(`up help team members`, "Show help and examples for a sub-command.")

	workdir := Cmd.Flag("chdir", "Change working directory.").Default(".").Short('C').String()
	verbose := Cmd.Flag("verbose", "Enable verbose log output.").Short('v').Bool()
	format := Cmd.Flag("format", "Output formatter.").Default("text").String()
	region := Cmd.Flag("region", "Target region id.").String()

	Cmd.PreAction(func(ctx *kingpin.ParseContext) error {
		os.Chdir(*workdir)

		if *verbose {
			log.SetHandler(delta.Default)
			log.SetLevel(log.DebugLevel)
			log.Debugf("up version %s (os: %s, arch: %s)", Cmd.GetVersion(), runtime.GOOS, runtime.GOARCH)
		}

		Init = func() (*up.Config, *up.Project, error) {
			c, err := up.ReadConfig("up.json")
			if err != nil {
				return nil, nil, errors.Wrap(err, "reading config")
			}

			if *region != "" {
				c.Regions = []string{*region}
			}

			events := make(event.Events)
			p := up.New(c, events).WithPlatform(lambda.New(c, events))

			switch {
			case *verbose:
				go reporter.Discard(events)
			case *format == "plain" || util.IsCI():
				go reporter.Plain(events)
			default:
				go reporter.Text(events)
			}

			return c, p, nil
		}

		return nil
	})
}
