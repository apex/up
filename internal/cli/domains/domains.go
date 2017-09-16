package domains

import (
	"fmt"
	"time"

	"github.com/Bowery/prompt"
	"github.com/pkg/errors"
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/colors"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/platform"
	"github.com/apex/up/platform/lambda/cost"
)

// TODO: make list the default
// TODO: find/write prompt with more options
// TODO: add colors to prompts
// TODO: add validation for emails, phone numbers, postal codes etc
// TODO: add ability to move up/down lines more like a form
// TODO: add polling of registration status (it's async)
// TODO: auto-fill these details?

func init() {
	cmd := root.Command("domains", "Manage domain names.")
	cmd.Example(`up domains`, "List purchased domains.")
	cmd.Example(`up domains check example.com`, "Check availability of a domain.")
	cmd.Example(`up domains buy example.com`, "Purchase domain if it's available.")
	list(cmd)
	check(cmd)
	buy(cmd)
}

// USD.
var usd = colors.Gray("USD")

// Anwsers.
var (
	firstName   = new(string)
	lastName    = new(string)
	email       = new(string)
	phone       = new(string)
	address     = new(string)
	city        = new(string)
	state       = new(string)
	countryCode = new(string)
	zipCode     = new(string)
)

// Questions.
var questions = []*struct {
	Prompt string
	Value  *string
}{
	{"First name", firstName},
	{"Last name", lastName},
	{"Email", email},
	{"Phone", phone},
	{"Country code", countryCode},
	{"City", city},
	{"State or province", state},
	{"Zip code", zipCode},
	{"Address", address},
}

func buy(cmd *kingpin.CmdClause) {
	c := cmd.Command("buy", "Purchase a domain.")
	domain := c.Arg("domain", "Domain name.").Required().String()

	c.Action(func(_ *kingpin.ParseContext) error {
		defer util.Pad()()

		_, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		confirm, err := prompt.Basic("  Confirm domain:", true)
		if err != nil {
			return err
		}

		if confirm != *domain {
			return errors.New("domains do not match")
		}

		for _, q := range questions {
			s := fmt.Sprintf("  %s:", q.Prompt)
			v, err := prompt.Basic(s, true)
			if err != nil {
				return err
			}
			*q.Value = v
		}

		stats.Track("Register Domain", nil)

		contact := platform.DomainContact{
			Email:       *email,
			FirstName:   *firstName,
			LastName:    *lastName,
			CountryCode: *countryCode,
			City:        *city,
			Address:     *address,
			PhoneNumber: *phone,
			State:       *state,
			ZipCode:     *zipCode,
		}

		domains := p.Domains()
		if err := domains.Purchase(*domain, contact); err != nil {
			return errors.Wrap(err, "purshasing")
		}

		return nil
	})
}

func check(cmd *kingpin.CmdClause) {
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

func list(cmd *kingpin.CmdClause) {
	c := cmd.Command("list", "List purchased domains.").Default()

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
				s = " renews"
			}
			fmt.Printf("  %-40s %s %s\n", colors.Purple(d.Name), s, d.Expiry.Format(time.Stamp))
		}

		return nil
	})
}
