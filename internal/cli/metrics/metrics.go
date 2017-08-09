package metrics

import (
	"time"

	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
)

func init() {
	cmd := root.Command("metrics", "Show project metrics.")
	cmd.Example(`up metrics`, "Show metrics for development stage.")
	cmd.Example(`up metrics production`, "Show metrics for production stage.")

	stage := cmd.Arg("stage", "Name of the stage.").Default("development").String()
	since := cmd.Flag("since", "Show logs since duration (30s, 5m, 2h, 1h30m).").Short('s').Default("24h").Duration()

	cmd.Action(func(_ *kingpin.ParseContext) error {
		region := root.Config.Regions[0]

		stats.Track("Metrics", map[string]interface{}{
			"stage": *stage,
			"since": since.Round(time.Second),
		})

		start := time.Now().Add(-*since)
		return root.Project.ShowMetrics(region, *stage, start)
	})
}
