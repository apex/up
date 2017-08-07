package stack

import (
	"fmt"

	"github.com/tj/go-prompt"
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
)

func init() {
	cmd := root.Command("stack", "Stack resource management.")
	cmd.Example(`up stack`, "Create or update the configured resources.")
	delete(cmd)
	show(cmd)
}

func delete(cmd *kingpin.CmdClause) {
	c := cmd.Command("delete", "Delete configured resources.")
	c.Example(`up stack delete`, "Delete stack with confirmation prompt.")
	c.Example(`up stack delete --force`, "Delete stack without confirmation prompt.")
	c.Example(`up stack delete --wait`, "Wait for deletion to complete before exiting.")
	c.Example(`up stack delete -fw`, "Force and wait for deletion.")

	force := c.Flag("force", "Force deletion without prompt.").Short('f').Bool()
	wait := c.Flag("wait", "Wait for deletion to complete.").Short('w').Bool()

	c.Action(func(_ *kingpin.ParseContext) error {
		name := root.Config.Name

		stats.Track("Delete Stack", map[string]interface{}{
			"force": *force,
			"wait":  *wait,
		})

		if !*force && !prompt.Confirm("  Really destroy the stack %q?  ", name) {
			fmt.Println("aborting")
			return nil
		}

		// TODO: multi-region
		return root.Project.DeleteStack(root.Config.Regions[0], *wait)
	})
}

// TODO: rename? decide on conventions
// TODO: make the default?
func show(cmd *kingpin.CmdClause) {
	c := cmd.Command("show", "Show the status of the stack.").Hidden()

	c.Action(func(_ *kingpin.ParseContext) error {
		stats.Track("Show Stack", nil)

		// TODO: multi-region
		return root.Project.ShowStack(root.Config.Regions[0])
	})
}
