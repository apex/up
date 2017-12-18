package upgrade

import (
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tj/go-update"
	"github.com/tj/go-update/stores/apex"
	"github.com/tj/go-update/stores/github"
	"github.com/tj/go/env"
	"github.com/tj/go/http/request"
	"github.com/tj/go/term"
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/progressreader"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/userconfig"
	"github.com/apex/up/internal/util"
)

var releasesAPI = env.GetDefault("APEX_RELEASES_API", "https://releases.apex.sh")

func init() {
	cmd := root.Command("upgrade", "Install the latest release of Up.")
	cmd.Action(func(_ *kingpin.ParseContext) error {
		version := root.Cmd.GetVersion()
		start := time.Now()

		term.HideCursor()
		defer term.ShowCursor()

		var config userconfig.Config
		if err := config.Load(); err != nil {
			return errors.Wrap(err, "loading user config")
		}

		// open-source edition
		p := &update.Manager{
			Command: "up",
			Store: &github.Store{
				Owner:   "apex",
				Repo:    "up",
				Version: version,
			},
		}

		// commercial edition
		if t := config.GetActiveTeam(); t != nil {
			// we pass 0.0.0 here beause the OSS
			// binary should always upgrade to Pro
			// regardless of versions matching.
			p.Store = &apex.Store{
				URL:       releasesAPI,
				Product:   "up",
				Version:   "0.0.0",
				Plan:      "pro",
				AccessKey: t.Token,
			}
		}

		// fetch the new releases
		releases, err := p.LatestReleases()

		if request.IsClient(err) {
			return errors.Wrap(err, "You're not subscribed to Up Pro")
		}

		if err != nil {
			return errors.Wrap(err, "fetching releases")
		}

		// no updates
		if len(releases) == 0 {
			util.LogPad("No updates available, you're good :)")
			return nil
		}

		// latest release
		latest := releases[0]

		// find the tarball for this system
		a := latest.FindTarball(runtime.GOOS, runtime.GOARCH)
		if a == nil {
			return errors.Errorf("failed to find a binary for %s %s", runtime.GOOS, runtime.GOARCH)
		}

		// download tarball to a tmp dir
		tarball, err := a.DownloadProxy(progressreader.New)
		if err != nil {
			return errors.Wrap(err, "downloading tarball")
		}

		// install it
		if err := p.Install(tarball); err != nil {
			return errors.Wrap(err, "installing")
		}

		term.ClearAll()

		if strings.Contains(a.URL, "up/pro") {
			util.LogPad("Updated %s to %s Pro :)", version, latest.Version)
		} else {
			util.LogPad("Updated %s to %s OSS :)", version, latest.Version)
		}

		stats.Track("Upgrade", map[string]interface{}{
			"from":     version,
			"to":       latest.Version,
			"duration": time.Since(start).Round(time.Millisecond),
		})

		return nil
	})
}
