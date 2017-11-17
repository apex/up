package docs

import (
	"github.com/pkg/browser"
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
)

var url = "https://up.docs.apex.sh"

func init() {
	cmd := root.Command("docs", "Open documentation website in the browser.")

	cmd.Example(`up docs`, "Open the documentation site.")

	cmd.Action(func(_ *kingpin.ParseContext) error {
		return browser.OpenURL(url)
	})
}
