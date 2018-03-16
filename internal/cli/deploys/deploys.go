package deploys

import (
	"time"

	"github.com/pkg/errors"
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
)

func init() {
	cmd := root.Command("deploys", "Show deployment history.")
	cmd.Example(`up deploys`, "Show all deployment history.")

	cmd.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		start := time.Now()

		region := c.Regions[0]
		if err := p.ShowDeploys(region); err != nil {
			return err
		}

		stats.Track("Deploys", map[string]interface{}{
			"duration": util.MillisecondsSince(start),
		})

		return nil
	})
}
