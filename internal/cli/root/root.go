package root

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/delta"
	"github.com/pkg/errors"
	"github.com/tj/kingpin"

	"github.com/apex/up"
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

	Cmd.Example(`up`, "Deploy the project to the development stage.")
	Cmd.Example(`up deploy production`, "Deploy the project to the production stage.")
	Cmd.Example(`up url`, "Show the development endpoint url.")
	Cmd.Example(`up logs -f`, "Tail project logs.")
	Cmd.Example(`up logs 'error or fatal'`, "Show error or fatal level logs.")
	Cmd.Example(`up run build`, "Run build command manually.")
	Cmd.Example(`up help logs`, "Show help and examples for a sub-command.")

	region := Cmd.Flag("region", "Override the region.").Short('r').String()
	workdir := Cmd.Flag("chdir", "Change working directory.").Default(".").Short('C').String()
	verbose := Cmd.Flag("verbose", "Enable verbose log output.").Short('v').Bool()
	endpoint := Cmd.Flag("endpoint", "Specify custom AWS API endpoint.").Short('e').String()

	Cmd.PreAction(func(ctx *kingpin.ParseContext) error {
		os.Chdir(*workdir)

		if *verbose {
			log.SetHandler(delta.Default)
			log.SetLevel(log.DebugLevel)
		}

		Init = func() (*up.Config, *up.Project, error) {
			c, err := up.ReadConfig("up.json")
			if err != nil {
				return nil, nil, errors.Wrap(err, "reading config")
			}

			if *region != "" {
				c.Regions = []string{*region}
			}

			if *endpoint != "" {
				c.Endpoint = *endpoint
			}

			events := make(event.Events)
			p := up.New(c, events).WithPlatform(lambda.New(c, events))

			if *verbose {
				go reporter.Discard(events)
			} else {
				go reporter.Text(events)
			}

			return c, p, nil
		}

		return nil
	})
}
