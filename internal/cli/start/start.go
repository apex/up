package start

import (
	"net/http"

	"github.com/apex/log"
	"github.com/pkg/errors"
	"github.com/tj/kingpin"

	"github.com/apex/up/handler"
	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
)

func init() {
	cmd := root.Command("start", "Start development server.")
	cmd.Example(`up start`, "Start development server on port 3000.")
	cmd.Example(`up start --address :5000`, "Start development server on port 5000.")

	addr := cmd.Flag("address", "Address for server.").Default(":3000").String()

	cmd.Action(func(_ *kingpin.ParseContext) error {
		_, _, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stats.Track("Start", map[string]interface{}{
			"address": *addr,
		})

		h, err := handler.New()
		if err != nil {
			return errors.Wrap(err, "initializing handler")
		}

		log.WithField("address", *addr).Infof("listening")
		if err := http.ListenAndServe(*addr, h); err != nil {
			return errors.Wrap(err, "binding")
		}

		return nil
	})
}
