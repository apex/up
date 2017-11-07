package env

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/tj/kingpin"
	"github.com/tj/prompt/colors"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/secret"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/internal/validate"
	"github.com/apex/up/platform"
)

// TODO: logging utils
// TODO: differentiate between ADD and SET ... one overrides
// TODO: rename --desc? util...
// TODO: add prompt for remove and --force
// TODO: add option for viewing secret history
// TODO: date format util for domains too
// TODO: prefix for the project or document the lack of prefix?
// TODO: add docs
// TODO: optional '=' ?
// TODO: restart without redeploy?
// TODO: do this automatically for env var changes?
// TODO: multiple env vars with '=' ? rm too?
// TODO: better errors when not found....
// TODO: option to require them for an env before deploying / starting

func init() {
	cmd := root.Command("env", "Manage encrypted env variables.")
	cmd.Example(`up env`, "List variables available to all stages.")
	cmd.Example(`up env add MONGO_URL "mongodb://db1.example.net:2500/" -s production`, "Add a production env variable.")
	cmd.Example(`up env add MONGO_URL "mongodb://db2.example.net:2500/" -s staging`, "Add a staging env variable.")
	cmd.Example(`up env add S3_KEY xxxxxxx`, "Add add a global env variable for all stages.")
	cmd.Example(`up env add S3_KEY xxxxxxx -s production`, "Add a stage specific env var to override the previous.")
	cmd.Example(`up env rm S3_KEY`, "Remove a variable.")
	cmd.Example(`up env rm S3_KEY -s production`, "Remove a production variable.")
	list(cmd)
	add(cmd)
	remove(cmd)
}

// list variables.
func list(cmd *kingpin.CmdClause) {
	c := cmd.Command("ls", "List variables.").Alias("list").Default()

	c.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stats.Track("List Secrets", nil)

		secrets, err := p.Secrets("").List()
		if err != nil {
			return errors.Wrap(err, "listing secrets")
		}

		if len(secrets) == 0 {
			return nil
		}

		grouped := secret.GroupByStage(secret.FilterByApp(secrets, c.Name))

		for _, name := range []string{"all", "staging", "production"} {
			s, ok := grouped[name]
			if !ok {
				continue
			}

			util.LogTitle(name)
			outputSecrets(s)
		}

		println()

		return nil
	})
}

func outputSecrets(secrets []*platform.Secret) {
	for _, s := range secrets {
		mod := fmt.Sprintf("Modified %s by %s", humanize.Time(s.LastModified), s.LastModifiedUser)
		desc := colors.Gray(util.DefaultString(&s.Description, "No description"))
		if s.Value != "" {
			desc = s.Value
		}
		name := colors.Purple(s.Name)
		_ = desc
		util.LogName(name, mod)
	}
}

// add variables.
func add(cmd *kingpin.CmdClause) {
	c := cmd.Command("add", "Add a variable.").Alias("set")
	key := c.Arg("name", "Variable name.").Required().String()
	val := c.Arg("value", "Variable value.").Required().String()
	stage := c.Flag("stage", "Stage name.").Short('s').String()
	desc := c.Flag("desc", "Variable description message.").Short('d').String()
	plain := c.Flag("plain", "Store as cleartext (unencrypted).").Short('c').Bool()

	c.Action(func(_ *kingpin.ParseContext) error {
		if err := validate.OptionalStage(*stage); err != nil {
			return err
		}

		_, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stats.Track("Add Secret", map[string]interface{}{
			"cleartext": *plain,
			"stage":     *stage,
			"has_desc":  *desc != "",
		})

		if err := p.Secrets(*stage).Add(*key, *val, *desc, *plain); err != nil {
			return errors.Wrap(err, "adding secret")
		}

		util.LogPad("Added " + *key)

		return nil
	})
}

// remove variables.
func remove(cmd *kingpin.CmdClause) {
	c := cmd.Command("rm", "Remove a variable.").Alias("remove")
	stage := c.Flag("stage", "Stage name.").Short('s').String()
	key := c.Arg("name", "Variable name.").Required().String()

	c.Action(func(_ *kingpin.ParseContext) error {
		if err := validate.OptionalStage(*stage); err != nil {
			return err
		}

		defer util.Pad()()

		_, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stats.Track("Remove Secret", nil)

		if err := p.Secrets(*stage).Remove(*key); err != nil {
			return errors.Wrap(err, "removing secret")
		}

		util.LogPad("Removed " + *key)

		return nil
	})
}
