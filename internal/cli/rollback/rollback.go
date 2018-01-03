package rollback

import (
	"github.com/pkg/errors"
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
)

func init() {
	cmd := root.Command("rollback", "Rollback to a previous deployment.")
	cmd.Example(`up rollback`, "Rollback to the previous production version.")
	cmd.Example(`up rollback 50`, "Rollback to an explicit production version.")

	stage := cmd.Flag("stage", "Target stage name.").Short('s').Default("production").String()
	version := cmd.Arg("version", "Target version for rollback.").String()

	cmd.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		// TODO: multi-region
		r := c.Regions[0]
		v := *version

		if err := p.Rollback(r, *stage, v); err != nil {
			return errors.Wrap(err, "rollback")
		}

		stats.Track("Rollback", map[string]interface{}{
			"has_version": v != "",
			"stage":       *stage,
		})

		util.LogPad("Rollback complete")

		return nil
	})
}
