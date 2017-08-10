package stack

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tj/go-prompt"
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
)

func init() {
	cmd := root.Command("stack", "Stack resource management.")
	cmd.Example(`up stack`, "Show status of the stack resources.")
	cmd.Example(`up stack delete`, "Delete the stack resources.")
	cmd.Example(`up stack delete -w`, "Delete resources and wait for completion.")
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
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stats.Track("Delete Stack", map[string]interface{}{
			"force": *force,
			"wait":  *wait,
		})

		if !*force && !prompt.Confirm("  Really destroy the stack %q?  ", c.Name) {
			fmt.Println("aborting")
			return nil
		}

		// TODO: multi-region
		return p.DeleteStack(c.Regions[0], *wait)
	})
}

// TODO: rename status, info, show? decide on conventions
func show(cmd *kingpin.CmdClause) {
	c := cmd.Command("show", "Show status of resources.").Default()

	c.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stats.Track("Show Stack", nil)

		// TODO: multi-region
		return p.ShowStack(c.Regions[0])
	})
}
