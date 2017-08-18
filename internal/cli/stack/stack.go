package stack

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tj/go-prompt"
	"github.com/tj/kingpin"

	"github.com/apex/log"
	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
)

func init() {
	cmd := root.Command("stack", "Stack resource management.")
	cmd.Example(`up stack`, "Show status of the stack resources.")
	cmd.Example(`up stack delete`, "Delete the stack resources.")
	plan(cmd)
	apply(cmd)
	delete(cmd)
	show(cmd)
}

func plan(cmd *kingpin.CmdClause) {
	c := cmd.Command("plan", "Plan configuration changes.")
	c.Example(`up stack plan`, "Plan changes to configuration.")

	c.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		// stats.Track("Plan Stack", nil)

		// TODO: multi-region
		return p.PlanStack(c.Regions[0])
	})
}

func apply(cmd *kingpin.CmdClause) {
	c := cmd.Command("apply", "Apply configuration changes.")
	c.Example(`up stack apply`, "Apply changes to configuration.")

	c.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		// stats.Track("Plan Stack", nil)

		// TODO: multi-region
		return p.ApplyStack(c.Regions[0])
	})
}

func delete(cmd *kingpin.CmdClause) {
	c := cmd.Command("delete", "Delete configured resources.")
	c.Example(`up stack delete`, "Delete stack with confirmation prompt.")
	c.Example(`up stack delete --force`, "Delete stack without confirmation prompt.")
	c.Example(`up stack delete --async`, "Don't wait for deletion to complete.")
	c.Example(`up stack delete -fa`, "Force asynchronous deletion.")

	force := c.Flag("force", "Skip the confirmation prompt.").Short('f').Bool()
	async := c.Flag("async", "Perform deletion asynchronously.").Short('a').Bool()

	c.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		wait := !*async
		defer util.Pad()()

		stats.Track("Delete Stack", map[string]interface{}{
			"force": *force,
			"wait":  wait,
		})

		if !*force && !prompt.Confirm("  Really destroy the stack %q?  ", c.Name) {
			fmt.Printf("\n")
			log.Info("aborting")
			return nil
		}

		// TODO: multi-region
		return p.DeleteStack(c.Regions[0], wait)
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
