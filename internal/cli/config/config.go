package config

import (
	"encoding/json"
	"os"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
	"github.com/pkg/errors"
	"github.com/tj/kingpin"
)

func init() {
	cmd := root.Command("config", "Show configuration after defaults and validation.")
	cmd.Example(`up config`, "Show the config.")

	cmd.Action(func(_ *kingpin.ParseContext) error {
		c, _, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stats.Track("Show Config", nil)

		// note that config is already read in root.go
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(c)

		return nil
	})
}
