package run

import (
	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
	"github.com/tj/kingpin"
)

func init() {
	cmd := root.Command("run", "Run a hook.")
	cmd.Example(`up run build`, "Run build hook.")
	cmd.Example(`up run clean`, "Run clean hook.")

	hook := cmd.Arg("hook", "Name of the hook to run.").Required().String()

	cmd.Action(func(_ *kingpin.ParseContext) error {
		defer util.Pad()()

		stats.Track("Hook", map[string]interface{}{
			"name": *hook,
		})

		return root.Project.RunHook(*hook)
	})
}
