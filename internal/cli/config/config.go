package open

import (
	"encoding/json"
	"os"

	"github.com/tj/kingpin"
	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
)

func init() {
	cmd := root.Command("config", "Show configuration after defaults and validation.")
	cmd.Example(`up config`, "Show the config.")

	cmd.Action(func(_ *kingpin.ParseContext) error {
		stats.Track("Show Config", nil)

		// note that config is already read in root.go
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(root.Config)

		return nil
	})
}
