package deploy

import (
	"time"

	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
)

func init() {
	cmd := root.Command("deploy", "Deploy the project.").Default()
	stage := cmd.Arg("stage", "Target stage name.").Default("development").String()
	cmd.Example(`up deploy`, "Deploy the project to the development stage.")
	cmd.Example(`up deploy production`, "Deploy the project to the production stage.")

	cmd.Action(func(_ *kingpin.ParseContext) error {
		start := time.Now()
		defer util.Pad()()

		c := root.Config

		stats.Track("Deploy", map[string]interface{}{
			"duration":             time.Since(start) / time.Millisecond,
			"type":                 c.Type,
			"regions":              c.Regions,
			"stage":                *stage,
			"header_rules_count":   len(c.Headers),
			"redirect_rules_count": len(c.Redirects),
			"inject_rules_count":   len(c.Inject),
			"has_cors":             c.CORS != nil,
			"has_logs":             !c.Logs.Disable,
		})

		go stats.Client.Flush()

		if err := root.Project.Deploy(*stage); err != nil {
			return err
		}

		return nil
	})
}
