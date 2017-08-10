package deploy

import (
	"time"

	"github.com/pkg/errors"
	"github.com/tj/kingpin"

	"github.com/apex/log"
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
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		start := time.Now()
		defer util.Pad()()

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

		done := make(chan bool)

		go func() {
			defer close(done)
			if err := stats.Client.Flush(); err != nil {
				log.WithError(err).Warn("flushing analytics")
			}
		}()

		if err := p.Deploy(*stage); err != nil {
			return err
		}

		<-done

		return nil
	})
}
