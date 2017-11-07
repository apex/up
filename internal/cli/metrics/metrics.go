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
	cmd.Example(`up metrics`, "Show metrics for development stage.")
	cmd.Example(`up metrics production`, "Show metrics for production stage.")

	stage := cmd.Arg("stage", "Name of the stage.").Default("development").String()
	since := cmd.Flag("since", "Show logs since duration (30s, 5m, 2h, 1h30m, 3d, 1M).").Short('s').Default("1d").String()

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
