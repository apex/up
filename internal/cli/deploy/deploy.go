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
	"github.com/apex/up/platform/lambda/runtime"
)

func init() {
	cmd := root.Command("deploy", "Deploy the project.").Default()
	stage := cmd.Arg("stage", "Target stage name.").Default("staging").String()
	cmd.Example(`up deploy`, "Deploy the project the staging environment.")
	cmd.Example(`up deploy production`, "Deploy the project to the production environment.")

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

	// validate stage name
	if err := validate.List(stage, c.Stages.RemoteNames()); err != nil {
		return err
	}

	// stage overrides
	if err := c.Override(stage); err != nil {
		return errors.Wrap(err, "overriding")
	}

	defer util.Pad()()
	start := time.Now()

<<<<<<< HEAD
	if err := p.Init(stage); err != nil {
=======
	r := runtime.New(c)
	if err := r.Init(stage); err != nil {
>>>>>>> add initialization of env vars for deployments. Closes #458
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
		"has_error_pages":      !c.ErrorPages.Disable,
		"app_name_hash":        util.Md5(c.Name),
		"alerts_count":         len(c.Alerts),
		"actions_count":        len(c.Actions),
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
