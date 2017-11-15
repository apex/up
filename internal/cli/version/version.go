package version

import (
	"fmt"

	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
)

func init() {
	cmd := root.Command("version", "Show version.")
	cmd.Action(func(_ *kingpin.ParseContext) error {
		stats.Track("Show Version", nil)
		fmt.Printf("%s\n", root.Cmd.GetVersion())
		return nil
	})
}
