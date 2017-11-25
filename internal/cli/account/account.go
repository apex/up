package account

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/token"
	"github.com/tj/go/clipboard"
	"github.com/tj/go/env"
	"github.com/tj/go/http/request"
	"github.com/tj/kingpin"
	"github.com/tj/survey"

	"github.com/apex/log"
	"github.com/apex/up/internal/account"
	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/userconfig"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/platform/event"
	"github.com/apex/up/reporter"
)

var (
	api = env.GetDefault("APEX_TEAMS_API", "https://teams.apex.sh")
	a   = account.New(api)
)

func init() {
	cmd := root.Command("account", "Manage account, plans, and billing.")
	cmd.Example(`up account login`, "Sign in or create account with interactive prompt.")
	cmd.Example(`up account login --email tj@example.com`, "Sign in or create account.")
	cmd.Example(`up account login --email tj@example.com --team apex-software`, "Sign in to a team.")
	cmd.Example(`up account cards`, "List credit cards.")
	cmd.Example(`up account cards add`, "Add credit card to Stripe.")
	cmd.Example(`up account cards rm ID`, "Remove credit card from Stripe.")
	cmd.Example(`up account subscribe`, "Subscribe to the Pro plan.")
	cmd.Example(`up account invite --email asya@example.com`, "Invite a team member to your active team.")
	cmd.Example(`up account invite --email asya@example.com --team apex-inc`, "Invite a team member to a specific team.")
	status(cmd)
	switchTeam(cmd)
	invite(cmd)
	login(cmd)
	logout(cmd)
	cards(cmd)
	subscribe(cmd)
	unsubscribe(cmd)
	copy(cmd)
}

// copy commands.
func copy(cmd *kingpin.CmdClause) {
	c := cmd.Command("ci", "Credentials for CI.")
	copy := c.Flag("copy", "Credentials to the clipboard.").Short('c').Bool()

	c.Action(func(_ *kingpin.ParseContext) error {
		var config userconfig.Config
		if err := config.Load(); err != nil {
			return errors.Wrap(err, "loading")
		}

		stats.Track("Copy Credentials", map[string]interface{}{
			"copy": *copy,
		})

		b, err := json.Marshal(config)
		if err != nil {
			return errors.Wrap(err, "marshaling")
		}

		if *copy {
			clipboard.Write(string(b))
			fmt.Println("Copied to clipboard!")
			return nil
		}

		fmt.Printf("%s\n", string(b))

		return nil
	})
}

// card commands.
func cards(cmd *kingpin.CmdClause) {
	c := cmd.Command("cards", "Card management.")
	addCard(c)
	removeCard(c)
	listCards(c)
}

// status of account.
func status(cmd *kingpin.CmdClause) {
	c := cmd.Command("status", "Status of your account.").Default()

	c.Action(func(_ *kingpin.ParseContext) error {
		var config userconfig.Config
		if err := config.Load(); err != nil {
			return errors.Wrap(err, "loading config")
		}

		defer util.Pad()()
		stats.Track("Account Status", nil)

		if !config.Authenticated() {
			util.LogName("status", "Signed out")
			return nil
		}

		t := config.GetActiveTeam()

		util.LogName("status", "Signed in")
		util.LogName("active team", t.ID)

		// TODO: list teams

		plans, err := a.GetPlans(t.Token)
		if err != nil {
			return errors.Wrap(err, "listing plans")
		}

		if len(plans) == 0 {
			util.LogName("subscription", "none")
			return nil
		}

		defer util.Pad()()

		// TODO: amount should reflect any coupon discount present
		for _, p := range plans {
			util.LogName("subscription", p.PlanName)
			util.LogName("amount", "$%0.2f/mo USD", float64(p.Amount)/100)
			util.LogName("created", p.CreatedAt.Format("January 2, 2006"))
			if p.Canceled {
				util.LogName("canceled", p.CanceledAt.Format("January 2, 2006"))
			}
		}

		return nil
	})
}

// remove card.
func removeCard(cmd *kingpin.CmdClause) {
	c := cmd.Command("rm", "Remove credit card.").Alias("remove")
	id := c.Arg("id", "Card ID.").Required().String()

	c.Action(func(_ *kingpin.ParseContext) error {
		t, err := userconfig.Require()
		if err != nil {
			return err
		}

		stats.Track("Remove Card", nil)

		if err := a.RemoveCard(t.Token, "card_"+*id); err != nil {
			return errors.Wrap(err, "removing card")
		}

		util.LogPad("Card removed")

		return nil
	})
}

// list cards.
func listCards(cmd *kingpin.CmdClause) {
	c := cmd.Command("ls", "List credit cards.").Alias("list").Default()

	c.Action(func(_ *kingpin.ParseContext) error {
		t, err := userconfig.Require()
		if err != nil {
			return err
		}

		stats.Track("List Cards", nil)

		cards, err := a.GetCards(t.Token)
		if err != nil {
			return errors.Wrap(err, "listing cards")
		}

		defer util.Pad()()
		for _, c := range cards {
			id := strings.Replace(c.ID, "card_", "", 1)
			util.LogName(id, "%s ending in %s", c.Brand, c.LastFour)
		}

		return nil
	})
}

// add card.
func addCard(cmd *kingpin.CmdClause) {
	c := cmd.Command("add", "Add credit card.")
	c.Action(func(_ *kingpin.ParseContext) error {
		t, err := userconfig.Require()
		if err != nil {
			return err
		}

		stats.Track("Add Card", nil)

		defer util.Pad()()

		card, err := account.PromptForCard()
		if err != nil {
			return errors.Wrap(err, "prompting for card")
		}

		tok, err := token.New(&stripe.TokenParams{
			Card: &card,
		})

		if err != nil {
			return errors.Wrap(err, "requesting card token")
		}

		if err := a.AddCard(t.Token, tok.ID); err != nil {
			return errors.Wrap(err, "adding card")
		}

		return nil
	})
}

// invite user.
func invite(cmd *kingpin.CmdClause) {
	c := cmd.Command("invite", "Invite a team member.")
	c.Example(`up account invite --email asya@example.com`, "Invite a team member to your active team.")
	c.Example(`up account invite --email asya@example.com --team apex-inc`, "Invite a team member to a specific team.")
	email := c.Flag("email", "Email address.").String()
	team := c.Flag("team", "Team id (or current team).").String()

	c.Action(func(_ *kingpin.ParseContext) error {
		t, err := userconfig.Require()
		if err != nil {
			return err
		}

		if *team == "" {
			*team = t.ID
		}

		stats.Track("Invite", map[string]interface{}{
			"team":  *team,
			"email": *email,
		})

		if err := a.AddInvite(t.Token, *email); err != nil {
			return errors.Wrap(err, "adding invite")
		}

		util.LogPad("Invited %s to team %s", *email, *team)

		return nil
	})
}

// switchTeam team.
func switchTeam(cmd *kingpin.CmdClause) {
	c := cmd.Command("switch", "Switch active team.")
	c.Example(`up account switch`, "Switch teams interactively.")

	c.Action(func(_ *kingpin.ParseContext) error {
		defer util.Pad()()

		var config userconfig.Config
		if err := config.Load(); err != nil {
			return errors.Wrap(err, "loading user config")
		}

		var options []string
		for _, t := range config.GetTeams() {
			options = append(options, t.ID)
		}
		sort.Strings(options)

		var team string
		prompt := &survey.Select{
			Message: "",
			Options: options,
			Default: config.Team,
		}

		if err := survey.AskOne(prompt, &team, survey.Required); err != nil {
			return err
		}

		stats.Track("Switch Team", nil)

		err := userconfig.Alter(func(c *userconfig.Config) {
			c.Team = team
		})

		if err != nil {
			return errors.Wrap(err, "saving config")
		}

		return nil
	})
}

// login user.
func login(cmd *kingpin.CmdClause) {
	c := cmd.Command("login", "Sign in to your account.")
	c.Example(`up account login`, "Sign in or create account with interactive prompt.")
	c.Example(`up account login --email tj@example.com`, "Sign in or create account.")
	c.Example(`up account login --email tj@example.com --team apex-software`, "Sign in to a team.")
	email := c.Flag("email", "Email address.").String()
	team := c.Flag("team", "Team id.").String()

	c.Action(func(_ *kingpin.ParseContext) error {
		defer util.Pad()()

		var config userconfig.Config
		if err := config.Load(); err != nil {
			return errors.Wrap(err, "loading user config")
		}

		stats.Track("Login", map[string]interface{}{
			"team_count": len(config.GetTeams()),
		})

		// email from config
		if t := config.GetActiveTeam(); *email == "" && t != nil {
			*email = t.Email
		}

		// email prompt
		if *email == "" {
			var s string
			prompt := &survey.Input{Message: "email:"}
			survey.AskOne(prompt, &s, survey.Required)
			*email = s
		}

		// events
		events := make(event.Events)
		go reporter.Text(events)
		events.Emit("account.login.verify", nil)

		// log context
		l := log.WithFields(log.Fields{
			"email": *email,
			"team":  *team,
		})

		// authenticate
		var code string
		var err error
		if t := config.GetActiveTeam(); t != nil {
			l.Debug("login with token")
			code, err = a.LoginWithToken(t.Token, *email, *team)
		} else {
			l.Debug("login without token")
			code, err = a.Login(*email, *team)
		}

		if err != nil {
			return errors.Wrap(err, "login")
		}

		// personal team
		if *team == "" {
			team = email
		}

		// access key
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		l.WithField("team", *team).Debug("poll for access token")
		token, err := a.PollAccessToken(ctx, *email, *team, code)
		if err != nil {
			return errors.Wrap(err, "getting access token")
		}

		events.Emit("account.login.verified", nil)
		err = userconfig.Alter(func(c *userconfig.Config) {
			c.Team = *team
			c.AddTeam(&userconfig.Team{
				Token: token,
				ID:    *team,
				Email: *email,
			})
		})

		if err != nil {
			return errors.Wrap(err, "config")
		}

		return nil
	})
}

// logout user.
func logout(cmd *kingpin.CmdClause) {
	c := cmd.Command("logout", "Sign out of your account.")

	c.Action(func(_ *kingpin.ParseContext) error {
		stats.Track("Logout", nil)

		var config userconfig.Config
		if err := config.Save(); err != nil {
			return errors.Wrap(err, "saving")
		}

		util.LogPad("Signed out")

		return nil
	})
}

// subscribe to plan.
func subscribe(cmd *kingpin.CmdClause) {
	c := cmd.Command("subscribe", "Subscribe to the Pro plan.")

	c.Action(func(_ *kingpin.ParseContext) error {
		config, err := userconfig.Require()
		if err != nil {
			return err
		}

		defer util.Pad()()

		// TODO: fetch from plan
		amount := 2000

		// coupon
		var couponID string
		err = survey.AskOne(&survey.Input{
			Message: "Coupon (optional):",
		}, &couponID, nil)

		if err != nil {
			return err
		}

		// coupon provided
		if strings.TrimSpace(couponID) != "" {
			coupon, err := a.GetCoupon(couponID)
			if err != nil && !request.IsNotFound(err) {
				return errors.Wrap(err, "fetching coupon")
			}

			if coupon == nil {
				util.Log("Coupon is invalid")
			} else {
				amount = coupon.Discount(amount)
				util.Log("Coupon savings: %s", coupon.Description())
			}
		}

		// confirm
		var ok bool
		total := fmt.Sprintf("%0.2f", float64(amount)/100)
		err = survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("Subscribe to Up Pro for $%s/mo USD?", total),
		}, &ok, nil)

		if err != nil {
			return err
		}

		if !ok {
			util.LogPad("Aborted")
			stats.Track("Abort Subscription", nil)
			return nil
		}

		stats.Track("Subscribe", map[string]interface{}{
			"coupon": couponID,
		})

		if err := a.AddPlan(config.Token, "up", "pro", couponID); err != nil {
			return errors.Wrap(err, "subscribing")
		}

		util.Log("Subscribed")

		return nil
	})
}

// unsubscribe from plan.
func unsubscribe(cmd *kingpin.CmdClause) {
	c := cmd.Command("unsubscribe", "Unsubscribe from the Pro plan.")

	c.Action(func(_ *kingpin.ParseContext) error {
		config, err := userconfig.Require()
		if err != nil {
			return err
		}

		defer util.Pad()()

		// confirm
		var ok bool
		err = survey.AskOne(&survey.Confirm{
			Message: "Are you sure you want to unsubscribe?",
		}, &ok, nil)

		if err != nil {
			return err
		}

		if !ok {
			util.LogPad("Aborted")
			return nil
		}

		stats.Track("Unsubscribe", nil)

		if err := a.RemovePlan(config.Token, "up", "pro"); err != nil {
			return errors.Wrap(err, "unsubscribing")
		}

		util.Log("Unsubscribed")

		return nil
	})
}
