package version

import (
	"fmt"

	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
)

func init() {
	cmd := root.Command("version", "Show version.")
	cmd.Action(func(_ *kingpin.ParseContext) error {
		fmt.Println(root.Cmd.GetVersion())
		return nil
	})
}
