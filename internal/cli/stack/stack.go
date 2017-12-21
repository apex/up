package stack

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tj/kingpin"
	"github.com/tj/survey"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
)

func init() {
	cmd := root.Command("stack", "Stack resource management.")

	cmd.Example(`up stack`, "Show status of the stack resources.")
	cmd.Example(`up stack plan`, "Show resource changes.")
	cmd.Example(`up stack apply`, "Apply resource changes.")
	cmd.Example(`up stack delete`, "Delete the stack resources.")

	plan(cmd)
	apply(cmd)
	delete(cmd)
	status(cmd)
}

// plan changes.
func plan(cmd *kingpin.Cmd) {
	c := cmd.Command("plan", "Plan configuration changes.")
	c.Example(`up stack plan`, "Show changes planned.")

	c.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stats.Track("Plan Stack", nil)

		// TODO: multi-region
		return p.PlanStack(c.Regions[0])
	})
}

// apply changes.
func apply(cmd *kingpin.Cmd) {
	c := cmd.Command("apply", "Apply configuration changes.")
	c.Example(`up stack apply`, "Apply the changes of the previous plan.")

	c.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stats.Track("Apply Stack", map[string]interface{}{
			"dns_zone_count":     len(c.DNS.Zones),
			"stage_count":        len(c.Stages.List()),
			"stage_domain_count": len(c.Stages.Domains()),
		})

		// TODO: multi-region
		return p.ApplyStack(c.Regions[0])
	})
}

// delete resources.
func delete(cmd *kingpin.Cmd) {
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

		if *force {
			// TODO: multi-region
			return p.DeleteStack(c.Regions[0], wait)
		}

		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Really destroy stack %q?", c.Name),
		}

		var ok bool
		if err := survey.AskOne(prompt, &ok, nil); err != nil {
			return err
		}

		if !ok {
			util.LogPad("Aborted")
			return nil
		}

		return p.DeleteStack(c.Regions[0], wait)
	})
}

// status of the stack.
func status(cmd *kingpin.Cmd) {
	c := cmd.Command("status", "Show status of resources.").Default()

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
