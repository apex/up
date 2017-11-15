package deploy

import (
	"time"

	"github.com/pkg/errors"
	"github.com/tj/kingpin"

	"github.com/apex/log"
	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/internal/validate"
)

func init() {
	cmd := root.Command("deploy", "Deploy the project.").Default()
	stage := cmd.Arg("stage", "Target stage name.").Default("development").String()
	cmd.Example(`up deploy`, "Deploy the project to the development stage.")
	cmd.Example(`up deploy staging`, "Deploy the project to the staging stage.")
	cmd.Example(`up deploy production`, "Deploy the project to the production stage.")

	cmd.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		start := time.Now()

		stats.Track("Deploy", map[string]interface{}{
			"duration":             util.MillisecondsSince(start),
			"type":                 c.Type,
			"regions":              c.Regions,
			"stage":                *stage,
			"proxy_timeout":        c.Proxy.Timeout,
			"header_rules_count":   len(c.Headers),
			"redirect_rules_count": len(c.Redirects),
			"inject_rules_count":   len(c.Inject),
			"environment_count":    len(c.Environment),
			"dns_zone_count":       len(c.DNS.Zones),
			"stage_count":          len(c.Stages.List()),
			"stage_domain_count":   len(c.Stages.Domains()),
			"has_cors":             c.CORS != nil,
			"has_logs":             !c.Logs.Disable,
			"has_profile":          c.Profile != "",
			"has_error_pages":      c.ErrorPages.Enable,
			"app_name_hash":        util.Md5(c.Name),
		})

		done := make(chan bool)

		go func() {
			defer close(done)
			if err := stats.Client.Flush(); err != nil {
				log.WithError(err).Warn("flushing analytics")
			}
		}()

		if err := validate.Stage(*stage); err != nil {
			return err
		}

		defer util.Pad()()

		if err := p.Deploy(*stage); err != nil {
			return err
		}

		<-done

		return nil
	})
}
