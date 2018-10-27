package metrics

import (
	"time"

	"github.com/pkg/errors"
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
)

func init() {
	cmd := root.Command("metrics", "Show project metrics.")
	cmd.Example(`up metrics`, "Show metrics for staging environment.")
	cmd.Example(`up metrics -s production`, "Show metrics for production environment.")

	stage := cmd.Flag("stage", "Target stage name.").Short('s').Default("staging").String()
	since := cmd.Flag("since", "Show metrics since duration (30s, 5m, 2h, 1h30m, 3d, 1M).").Short('S').Default("1M").String()

	cmd.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		s, err := util.ParseDuration(*since)
		if err != nil {
			return errors.Wrap(err, "parsing --since duration")
		}

		region := c.Regions[0]

		stats.Track("Metrics", map[string]interface{}{
			"stage": *stage,
			"since": s.Round(time.Second),
		})

		start := time.Now().UTC().Add(-s)
		return p.ShowMetrics(region, *stage, start)
	})
}
