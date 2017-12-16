package deploy

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/tj/go/term"
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/setup"
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
		return deploy(*stage)
	})
}

func deploy(stage string) error {
retry:
	c, p, err := root.Init()

	// missing up.json non-interactive
	if isMissingConfig(err) && !term.IsTerminal(os.Stdin.Fd()) {
		return errors.New("Cannot find ./up.json configuration file.")
	}

	// missing up.json interactive
	if isMissingConfig(err) {
		err := setup.Create()

		if err == setup.ErrNoCredentials {
			return errors.New("Cannot find credentials, visit https://up.docs.apex.sh/#aws_credentials for help.")
		}

		if err != nil {
			return errors.Wrap(err, "setup")
		}

		util.Log("Deploying the project and creating resources.")
		goto retry
	}

	// unrelated error
	if err != nil {
		return errors.Wrap(err, "initializing")
	}

	if err := validate.Stage(stage); err != nil {
		return err
	}

	defer util.Pad()()
	start := time.Now()

	if err := p.Init(stage); err != nil {
		return errors.Wrap(err, "initializing")
	}

	if err := p.Deploy(stage); err != nil {
		return err
	}

	stats.Track("Deploy", map[string]interface{}{
		"duration":             util.MillisecondsSince(start),
		"type":                 c.Type,
		"regions":              c.Regions,
		"stage":                stage,
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

	stats.Flush()
	return nil
}

// isMissingConfig returns true if the error represents a missing up.json.
func isMissingConfig(err error) bool {
	err = errors.Cause(err)
	e, ok := err.(*os.PathError)
	return ok && e.Path == "up.json"
}
