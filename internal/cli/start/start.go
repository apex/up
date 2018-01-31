package start

import (
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/pkg/errors"
	"github.com/tj/kingpin"

	"github.com/apex/up/handler"
	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/logs/text"
	"github.com/apex/up/internal/stats"
)

func init() {
	cmd := root.Command("start", "Start development server.")
	cmd.Example(`up start`, "Start development server on port 3000.")
	cmd.Example(`up start --address :5000`, "Start development server on port 5000.")
	cmd.Example(`up start -c 'go run main.go'`, "Override proxy command.")
	command := cmd.Flag("command", "Proxy command override").Short('c').String()

	addr := cmd.Flag("address", "Address for server.").Default(":3000").String()

	cmd.Action(func(_ *kingpin.ParseContext) error {
		log.SetHandler(text.New(os.Stdout))

		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		for k, v := range c.Environment {
			os.Setenv(k, v)
		}

		stats.Track("Start", map[string]interface{}{
			"address":     *addr,
			"has_command": *command != "",
		})

		if err := p.Init("development"); err != nil {
			return errors.Wrap(err, "initializing")
		}

		if err := c.Override("development"); err != nil {
			return errors.Wrap(err, "overriding")
		}

		if s := *command; s != "" {
			c.Proxy.Command = s
		}

		h, err := handler.FromConfig(c)
		if err != nil {
			return errors.Wrap(err, "selecting handler")
		}

		h, err = handler.New(c, h)
		if err != nil {
			return errors.Wrap(err, "initializing handler")
		}

		log.WithField("address", *addr).Info("listening")
		if err := http.ListenAndServe(*addr, h); err != nil {
			return errors.Wrap(err, "binding")
		}

		return nil
	})
}
