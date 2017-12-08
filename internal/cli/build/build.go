package build

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/tj/go/term"
	"github.com/tj/kingpin"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/colors"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/platform/lambda/runtime"
)

func init() {
	cmd := root.Command("build", "Build zip file.")
	stage := cmd.Arg("stage", "Target stage name.").Default("development").String()
	size := cmd.Flag("size", "Show zip contents size information.").Bool()
	cmd.Example(`up build`, "Build archive and save to ./out.zip")
	cmd.Example(`up build > /tmp/out.zip`, "Build archive and output to file via stdout.")
	cmd.Example(`up build --size`, "Build archive and list files by size.")

	cmd.Action(func(_ *kingpin.ParseContext) error {
		defer util.Pad()()

		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stats.Track("Build", nil)

		if err := runtime.New(c).Init(*stage); err != nil {
			log.Fatalf("error initializing: %s", err)
		}

		if err := p.Build(); err != nil {
			return errors.Wrap(err, "building")
		}

		r, err := p.Zip()
		if err != nil {
			return errors.Wrap(err, "zip")
		}

		var out io.Writer
		var buf bytes.Buffer

		switch {
		default:
			out = os.Stdout
		case *size:
			out = &buf
		case term.IsTerminal(os.Stdout.Fd()):
			f, err := os.Create("out.zip")
			if err != nil {
				return errors.Wrap(err, "creating zip")
			}
			defer f.Close()
			out = f
		}

		if _, err := io.Copy(out, r); err != nil {
			return errors.Wrap(err, "copying")
		}

		if *size {
			z, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
			if err != nil {
				return errors.Wrap(err, "opening zip")
			}

			files := z.File

			sort.Slice(files, func(i int, j int) bool {
				a := files[i]
				b := files[j]
				return a.UncompressedSize64 > b.UncompressedSize64
			})

			fmt.Printf("\n")
			for _, f := range files {
				size := humanize.Bytes(f.UncompressedSize64)
				fmt.Printf("  %10s %s\n", size, colors.Purple(f.Name))
			}
		}

		return err
	})
}
