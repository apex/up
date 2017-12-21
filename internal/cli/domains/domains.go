package domains

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tj/kingpin"
	"github.com/tj/survey"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/colors"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/platform"
	"github.com/apex/up/platform/lambda/cost"
)

// TODO: add ability to move up/down lines more like a form
// TODO: add polling of registration status (it's async)
// TODO: auto-fill these details from AWS account?

func init() {
	cmd := root.Command("domains", "Manage domain names.")
	cmd.Example(`up domains`, "List purchased domains.")
	cmd.Example(`up domains check example.com`, "Check availability of a domain.")
	cmd.Example(`up domains buy`, "Purchase a domain.")
	list(cmd)
	check(cmd)
	buy(cmd)
}

// USD.
var usd = colors.Gray("USD")

// Questions.
var questions = []*survey.Question{
	{
		Name:     "email",
		Prompt:   &survey.Input{Message: "Email:"},
		Validate: validateEmail,
	},
	{
		Name:     "firstname",
		Prompt:   &survey.Input{Message: "First name:"},
		Validate: survey.Required,
	},
	{
		Name:     "lastname",
		Prompt:   &survey.Input{Message: "Last name:"},
		Validate: survey.Required,
	},
	{
		Name:     "countrycode",
		Prompt:   &survey.Input{Message: "Country code:"},
		Validate: validateCountryCode,
	},
	{
		Name:     "city",
		Prompt:   &survey.Input{Message: "City:"},
		Validate: survey.Required,
	},
	{
		Name:     "address",
		Prompt:   &survey.Input{Message: "Address:"},
		Validate: survey.Required,
	},
	{
		Name:     "phonenumber",
		Prompt:   &survey.Input{Message: "Phone:"},
		Validate: validatePhoneNumber,
	},
	{
		Name:     "state",
		Prompt:   &survey.Input{Message: "State:"},
		Validate: survey.Required,
	},
	{
		Name:     "zipcode",
		Prompt:   &survey.Input{Message: "Zip code:"},
		Validate: survey.Required,
	},
}

// buy a domain.
func buy(cmd *kingpin.Cmd) {
	c := cmd.Command("buy", "Purchase a domain.")

	c.Action(func(_ *kingpin.ParseContext) error {
		defer util.Pad()()

		_, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		var domain string
		survey.AskOne(&survey.Input{
			Message: "Domain:",
		}, &domain, survey.Required)

		var contact platform.DomainContact

		if err := survey.Ask(questions, &contact); err != nil {
			return errors.Wrap(err, "prompting")
		}

		domains := p.Domains()
		if err := domains.Purchase(domain, contact); err != nil {
			return errors.Wrap(err, "purshasing")
		}

		return nil
	})
}

// check domain availability.
func check(cmd *kingpin.Cmd) {
	c := cmd.Command("check", "Check availability of a domain.")
	domain := c.Arg("domain", "Domain name.").Required().String()

	c.Action(func(_ *kingpin.ParseContext) error {
		defer util.Pad()()

		_, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stats.Track("Check Domain Availability", nil)

		domains := p.Domains()
		d, err := domains.Availability(*domain)
		if err != nil {
			return errors.Wrap(err, "fetching availability")
		}

		state := fmt.Sprintf("Domain %s is unavailable", d.Name)
		if d.Available {
			state = fmt.Sprintf("Domain %s is available for %s %s", d.Name, cost.Domain(d.Name), usd)
		}

		fmt.Printf("  %s\n", colors.Bool(d.Available)(state))

		if !d.Available {
			fmt.Printf("\n  Suggestions:\n")

			suggestions, err := domains.Suggestions(*domain)
			if err != nil {
				return errors.Wrap(err, "fetching suggestions")
			}

			fmt.Printf("\n")
			for _, d := range suggestions {
				price := cost.Domain(d.Name)
				fmt.Printf("  %-40s %s %s\n", colors.Purple(d.Name), price, usd)
			}
		}

		return nil
	})
}

// list domains purchased.
func list(cmd *kingpin.Cmd) {
	c := cmd.Command("ls", "List purchased domains.").Alias("list").Default()

	c.Action(func(_ *kingpin.ParseContext) error {
		defer util.Pad()()

		_, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stats.Track("List Domains", nil)

		domains, err := p.Domains().List()
		if err != nil {
			return errors.Wrap(err, "listing domains")
		}

		for _, d := range domains {
			s := "expires"
			if d.AutoRenew {
				s = "renews"
			}
			util.LogName(d.Name, "%s %s", s, d.Expiry.Format(time.Stamp))
		}

		return nil
	})
}

// validateEmail returns an error if the input does not look like an email.
func validateEmail(v interface{}) error {
	s := v.(string)
	i := strings.LastIndex(s, "@")

	if s == "" {
		return errors.New("Email is required.")
	}

	if i == -1 {
		return errors.New("Email is missing '@'.")
	}

	if i == len(s)-1 {
		return errors.New("Email is missing domain.")
	}

	return nil
}

// validateCountryCode returns an error if the input does not look like a valid country code.
func validateCountryCode(v interface{}) error {
	s := v.(string)

	if s == "" {
		return errors.New("Country code is required.")
	}

	if len(s) != 2 {
		return errors.New("Country codes must consist of two uppercase letters, such as CA or AU.")
	}

	return nil
}

// validatePhoneNumber returns an error if the input does not look like a valid phone number.
func validatePhoneNumber(v interface{}) error {
	s := v.(string)

	if s == "" {
		return errors.New("Phone number is required.")
	}

	if !strings.HasPrefix(s, "+") {
		return errors.New("Phone number must contain the country code, for example +1.2223334444 for Canada.")
	}

	return nil
}
