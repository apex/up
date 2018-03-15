package run

import (
	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
	"github.com/pkg/errors"
	"github.com/tj/kingpin"
)

func init() {
	cmd := root.Command("run", "Run a hook.")
	hook := cmd.Arg("hook", "Name of the hook to run.").Required().String()
	stage := cmd.Arg("stage", "Target stage name.").Default("development").String()
	cmd.Example(`up run build`, "Run build hook.")
	cmd.Example(`up run clean`, "Run clean hook.")

	cmd.Action(func(_ *kingpin.ParseContext) error {
		_, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		defer util.Pad()()

		stats.Track("Hook", map[string]interface{}{
			"name": *hook,
		})

		if err := p.Init(*stage); err != nil {
			return errors.Wrap(err, "initializing")
		}

		return p.RunHook(*hook)
	})
}
