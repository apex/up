package docs

import (
	"github.com/pkg/browser"
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
)

// URL for documentation.
var url = `https://github.com/apex/up/tree/master/docs`

func init() {
	cmd := root.Command("docs", "Show docs in the browser.")
	cmd.Action(func(_ *kingpin.ParseContext) error {
		return browser.OpenURL(url)
	})
}
