package deploys

import (
	"github.com/pkg/errors"
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
)

func init() {
	cmd := root.Command("deploys", "Show deployment history.")
	cmd.Example(`up deploys`, "Show all deployment history.")

	cmd.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stats.Track("Deploys", nil)

		region := c.Regions[0]
		return p.ShowDeploys(region)
	})
}
