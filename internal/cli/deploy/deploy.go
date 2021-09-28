package deploy

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/tj/go/git"
	"github.com/tj/go/term"
	"github.com/tj/kingpin"

	"github.com/apex/up"
	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/setup"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/internal/validate"
)

func init() {
	cmd := root.Command("deploy", "Deploy the project.").Default()
	stage := cmd.Arg("stage", "Target stage name.").Default("production").String()
	noBuild := cmd.Flag("no-build", "Disable build related hooks.").Bool()

	cmd.Example(`up deploy`, "Deploy to the staging environment.")
	cmd.Example(`up deploy production`, "Deploy to the production environment.")
	cmd.Example(`up deploy --no-build`, "Skip build hooks, useful in CI when a separate build step is used.")

	cmd.Action(func(_ *kingpin.ParseContext) error {
		return deploy(*stage, !*noBuild)
	})
}

func deploy(stage string, build bool) error {
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
			return errors.New("Cannot find credentials, visit https://apex.sh/docs/up/credentials/ for help.")
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

	// git information
	commit, err := getCommit()
	if err != nil {
		return errors.Wrap(err, "fetching git commit")
	}

	defer util.Pad()()
	start := time.Now()

	if err := p.Init(stage); err != nil {
		return errors.Wrap(err, "initializing")
	}

	if err := p.Deploy(up.Deploy{
		Stage:  stage,
		Commit: util.StripLerna(commit.Describe()),
		Author: commit.Author.Name,
		Build:  build,
	}); err != nil {
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
		"lambda_memory":        c.Lambda.Memory,
		"has_cors":             c.CORS != nil,
		"has_logs":             !c.Logs.Disable,
		"has_profile":          c.Profile != "",
		"has_error_pages":      c.ErrorPages.Enable,
		"app_name_hash":        util.Md5(c.Name),
		"is_git":               commit.Author.Name != "",
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

// getCommit returns the git information when available.
func getCommit() (git.Commit, error) {
	c, err := git.GetCommit(".", "HEAD")
	if err != nil && !isIgnorable(err) {
		return git.Commit{}, err
	}

	if c == nil {
		return git.Commit{}, nil
	}

	return *c, nil
}

// isIgnorable returns true if the GIT error is ignorable.
func isIgnorable(err error) bool {
	switch err {
	case git.ErrLookup, git.ErrNoRepo, git.ErrDirty:
		return true
	default:
		return false
	}
}
