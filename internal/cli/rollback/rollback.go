package rollback

import (
	"github.com/pkg/errors"
	"github.com/tj/go/term"
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
)

func init() {
	cmd := root.Command("rollback", "Rollback to a previous deployment.")
	cmd.Example(`up rollback`, "Rollback to the previous staging version.")
	cmd.Example(`up rollback 15c46ba`, "Rollback to a specific git commit")
	cmd.Example(`up rollback v1.7.2`, "Rollback to a specific git tag")
	cmd.Example(`up rollback -s production`, "Rollback to the previous production version.")
	cmd.Example(`up rollback -s production 15c46ba`, "Rollback to a specific git commit")
	cmd.Example(`up rollback -s production v1.7.2`, "Rollback to a specific git tag")

	stage := cmd.Flag("stage", "Target stage name.").Short('s').Default("staging").String()
	version := cmd.Arg("version", "Target version for rollback.").String()

	cmd.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		defer util.Pad()()

		// TODO: multi-region
		r := c.Regions[0]
		v := *version

		util.Log("Rolling back %s", *stage)

		stats.Track("Rollback", map[string]interface{}{
			"has_version": v != "",
			"stage":       *stage,
		})

		err = p.Rollback(r, *stage, v)
		if util.IsNotFound(err) {
			term.MoveUp(1)
			term.ClearLine()
			return errors.Errorf("Rollback failed – version %q not found", *version)
		}

		if err != nil {
			return errors.Wrap(err, "rollback")
		}

		util.LogClear("Rollback complete")

		return nil
	})
}
