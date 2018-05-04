package prune

import (
	"github.com/pkg/errors"
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
)

func init() {
	cmd := root.Command("prune", "Prune old S3 deployments of a stage.")

	cmd.Example(`up prune`, "Prune and retain the most recent 30 staging versions.")
	cmd.Example(`up prune -s production`, "Prune and retain the most recent 30 production versions.")
	cmd.Example(`up prune -s production -r 15`, "Prune and retain the most recent 15 production versions.")

	stage := cmd.Flag("stage", "Target stage name.").Short('s').Default("staging").String()
	versions := cmd.Flag("retain", "Number of versions to retain.").Short('r').Default("30").Int()

	cmd.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		region := c.Regions[0]

		stats.Track("Prune", map[string]interface{}{
			"versions": *versions,
			"stage":    *stage,
		})

		return p.Prune(region, *stage, *versions)
	})
}
