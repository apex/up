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
	cmd := root.Command("upgrade", "Install the latest or specified version of Up.")
	cmd.Example(`up upgrade`, "Upgrade to the latest version available.")
	cmd.Example(`up upgrade -t 0.4.4`, "Upgrade to the specified version.")
	target := cmd.Flag("target", "Target version for upgrade.").Short('t').String()

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

		// fetch latest or specified release
		r, err := getLatestOrSpecified(p, *target)
		if err != nil {
			return errors.Wrap(err, "fetching latest release")
		}

		// no updates
		if r == nil {
			util.LogPad("No updates available, you're good :)")
			return nil
		}

		// find the tarball for this system
		a := r.FindTarball(runtime.GOOS, runtime.GOARCH)
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
			util.LogPad("Updated %s to %s Pro", versionName(version), r.Version)
		} else {
			util.LogPad("Updated %s to %s OSS", versionName(version), r.Version)
		}

		stats.Track("Upgrade", map[string]interface{}{
			"from":     version,
			"to":       r.Version,
			"duration": time.Since(start).Round(time.Millisecond),
		})

		return nil
	})
}

// getLatestOrSpecified returns the latest or specified release.
func getLatestOrSpecified(s update.Store, version string) (*update.Release, error) {
	if version == "" {
		return getLatest(s)
	}

	return s.GetRelease(version)
}

// getLatest returns the latest release, error, or nil when there is none.
func getLatest(s update.Store) (*update.Release, error) {
	releases, err := s.LatestReleases()

	if request.IsClient(err) {
		return nil, errors.Wrap(err, "You're not subscribed to Up Pro")
	}

	if err != nil {
		return nil, errors.Wrap(err, "fetching releases")
	}

	if len(releases) == 0 {
		return nil, nil
	}

	return releases[0], nil
}

// versionName returns the humanized version name.
func versionName(s string) string {
	if strings.Contains(s, "-pro") {
		return strings.Replace(s, "-pro", "", 1) + " Pro"
	}

	return s + " OSS"
}
