package metrics

import (
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
)

// TODO: add --since

func init() {
	cmd := root.Command("metrics", "Show project metrics.")
	cmd.Example(`up metrics`, "Show metrics for development stage.")
	cmd.Example(`up metrics production`, "Show metrics for production stage.")

	stage := cmd.Arg("stage", "Name of the stage.").Default("development").String()

	cmd.Action(func(_ *kingpin.ParseContext) error {
		region := root.Config.Regions[0]

		stats.Track("Metrics", map[string]interface{}{
			"stage": *stage,
		})

		return root.Project.ShowMetrics(region, *stage)
	})
}
