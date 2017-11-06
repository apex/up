package account

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/token"
	"github.com/tj/go/env"
	"github.com/tj/kingpin"
	"github.com/tj/survey"

	"github.com/apex/up/internal/account"
	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/userconfig"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/platform/event"
	"github.com/apex/up/reporter"
)

// TODO: input masking (Go seems to be lacking some term utilities)
// TODO: add nicer error if they try to subscribe without a card
// TODO: prompt user for next step after login (would you like to add a card... etc)

var (
	accountAPI = env.GetDefault("APEX_USERS_API", "https://users.apex.sh")
	a          = account.New(accountAPI)
)

func init() {
	cmd := root.Command("account", "Manage account, plans, and billing.")
	cmd.Example(`up account login`, "Sign in or create account.")
	cmd.Example(`up account cards`, "List credit cards.")
	cmd.Example(`up account cards add`, "Add credit card to Stripe.")
	cmd.Example(`up account cards rm ID`, "Remove credit card from Stripe.")
	cmd.Example(`up account subscribe`, "Subscribe to the Pro plan.")
	status(cmd)
	login(cmd)
	logout(cmd)
	cards(cmd)
	subscribe(cmd)
	unsubscribe(cmd)
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

		if config.Token == "" {
			util.LogName("status", "Signed out")
			return nil
		}

		util.LogName("status", "Signed in")
		util.LogName("email", config.Email)

		plans, err := a.GetPlans(config.Token)
		if err != nil {
			return errors.Wrap(err, "listing plans")
		}

		if len(plans) == 0 {
			return nil
		}

		defer util.Pad()()

		for _, p := range plans {
			util.LogName("name", p.PlanName)
			util.LogName("amount", "$%0.2f/mo USD", float64(p.Amount)/100)
			util.LogName("created", p.CreatedAt.Format("January 2, 2006"))
		}

		return nil
	})
}

// remove card.
func removeCard(cmd *kingpin.CmdClause) {
	c := cmd.Command("rm", "Remove credit card.").Alias("remove")
	id := c.Arg("id", "Card ID.").Required().String()

	c.Action(func(_ *kingpin.ParseContext) error {
		config, err := userconfig.Require()
		if err != nil {
			return err
		}

		stats.Track("Remove Card", nil)

		if err := a.RemoveCard(config.Token, "card_"+*id); err != nil {
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
		config, err := userconfig.Require()
		if err != nil {
			return err
		}

		stats.Track("List Cards", nil)

		cards, err := a.GetCards(config.Token)
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
		config, err := userconfig.Require()
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

		if err := a.AddCard(config.Token, tok.ID); err != nil {
			return errors.Wrap(err, "adding card")
		}

		return nil
	})
}

// login user.
func login(cmd *kingpin.CmdClause) {
	c := cmd.Command("login", "Sign in to your account.")

	c.Action(func(_ *kingpin.ParseContext) error {
		_, _, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		defer util.Pad()()
		stats.Track("Login", nil)

		// prompt
		var email string
		prompt := &survey.Input{Message: "email:"}
		survey.AskOne(prompt, &email, survey.Required)

		events := make(event.Events)
		go reporter.Text(events)
		events.Emit("account.login.verify", nil)

		// send email
		code, err := a.Login(email)
		if err != nil {
			return errors.Wrap(err, "login")
		}

		// access key
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		key, err := a.PollAccessKey(ctx, email, code)
		if err != nil {
			return errors.Wrap(err, "getting access key")
		}

		events.Emit("account.login.verified", nil)
		err = userconfig.Alter(func(c *userconfig.Config) {
			c.Email = email
			c.Token = key
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

		// confirm
		var ok bool
		err = survey.AskOne(&survey.Confirm{
			Message: "Subscribe to Up Pro",
		}, &ok, nil)

		if err != nil {
			return err
		}

		if !ok {
			util.LogPad("Aborted")
			return nil
		}

		// coupon
		var coupon string
		err = survey.AskOne(&survey.Input{
			Message: "Coupon (optional):",
		}, &coupon, nil)

		if err != nil {
			return err
		}

		stats.Track("Subscribe", map[string]interface{}{
			"coupon": coupon,
		})

		if err := a.AddPlan(config.Token, "up", "pro", coupon); err != nil {
			return errors.Wrap(err, "subscribing")
		}

		err = userconfig.Alter(func(c *userconfig.Config) {
			c.Plan = "pro"
		})

		if err != nil {
			return errors.Wrap(err, "saving config")
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

		stats.Track("Unsubscribe", nil)

		if err := a.RemovePlan(config.Token, "up", "pro"); err != nil {
			return errors.Wrap(err, "unsubscribing")
		}

		util.LogPad("Unsubscribed")

		return nil
	})
}
