package prune

import (
	"github.com/pkg/errors"
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
)

func init() {
	cmd := root.Command("prune", "Prune old S3 deployments.")

	cmd.Example(`up prune`, "Prune and retain the most recent 60 versions.")
	cmd.Example(`up prune --retain 15`, "Prune and retain the most recent 15 versions.")

	versions := cmd.Flag("retain", "Number of versions to retain.").Short('r').Default("60").Int()

	cmd.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		region := c.Regions[0]

		stats.Track("Prune", map[string]interface{}{
			"versions": *versions,
		})

		return p.Prune(region, *versions)
	})
}
