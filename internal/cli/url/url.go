package url

import (
	"fmt"

	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/tj/go/clipboard"
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/validate"
)

func init() {
	cmd := root.Command("url", "Show, open, or copy a stage endpoint.")

	cmd.Example(`up url`, "Show the development endpoint.")
	cmd.Example(`up url --open`, "Open the development endpoint in the browser.")
	cmd.Example(`up url --copy`, "Copy the development endpoint to the clipboard.")
	cmd.Example(`up url production`, "Show the production endpoint.")
	cmd.Example(`up url -o production`, "Open the production endpoint in the browser.")
	cmd.Example(`up url -c production`, "Copy the production endpoint to the clipboard.")

	stage := cmd.Arg("stage", "Name of the stage.").Default("development").String()
	open := cmd.Flag("open", "Open endpoint in the browser.").Short('o').Bool()
	copy := cmd.Flag("copy", "Copy endpoint to the clipboard.").Short('c').Bool()

	cmd.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		region := c.Regions[0]

		stats.Track("URL", map[string]interface{}{
			"region": region,
			"stage":  *stage,
			"open":   *open,
			"copy":   *copy,
		})

		if err := validate.Stage(*stage); err != nil {
			return err
		}

		url, err := p.URL(region, *stage)
		if err != nil {
			return err
		}

		switch {
		case *open:
			browser.OpenURL(url)
		case *copy:
			clipboard.Write(url)
			fmt.Println("Copied to clipboard!")
		default:
			fmt.Println(url)
		}

		return nil
	})
}
