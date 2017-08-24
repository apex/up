package open

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/tj/go/term"
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
)

func init() {
	cmd := root.Command("build", "Build zip file.")

	cmd.Example(`up build`, "Build archive and save to ./out.zip")
	cmd.Example(`up build > /tmp/out.zip`, "Build archive and output to file via stdout.")

	cmd.Action(func(_ *kingpin.ParseContext) error {
		defer util.Pad()()

		_, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stats.Track("Build", nil)

		if err := p.Build(); err != nil {
			return errors.Wrap(err, "building")
		}

		r, err := p.Zip()
		if err != nil {
			return errors.Wrap(err, "zip")
		}

		out := os.Stdout

		if term.IsTerminal() {
			f, err := os.Create("out.zip")
			if err != nil {
				return errors.Wrap(err, "creating zip")
			}
			defer f.Close()
			out = f
		}

		_, err = io.Copy(out, r)
		return err
	})
}
