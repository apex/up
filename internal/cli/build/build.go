package build

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"

	humanize "github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/tj/go/git"
	"github.com/tj/go/term"
	"github.com/tj/kingpin"

	"github.com/apex/up"
	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/colors"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
)

func init() {
	cmd := root.Command("build", "Build zip file.")
	cmd.Example(`up build`, "Build archive and save to ./out.zip")
	cmd.Example(`up build > /tmp/out.zip`, "Build archive and output to file via stdout.")
	cmd.Example(`up build --size`, "Build archive and list files by size.")

	stage := cmd.Flag("stage", "Target stage name.").Short('s').Default("staging").String()
	size := cmd.Flag("size", "Show zip contents size information.").Bool()

	cmd.Action(func(_ *kingpin.ParseContext) error {
		defer util.Pad()()

		_, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stats.Track("Build", nil)

		// git information
		commit, err := getCommit()
		if err != nil {
			return errors.Wrap(err, "fetching git commit")
		}

		if err := p.Init(*stage); err != nil {
			return errors.Wrap(err, "initializing")
		}

		if err := p.Build(up.Build{
			Stage:  *stage,
			Commit: util.StripLerna(commit.Describe()),
			Author: commit.Author.Name,
			Hooks:  true,
		}); err != nil {
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

// getCommit returns the git information when available.
func getCommit() (git.Commit, error) {
	c, err := git.GetCommit(".", "HEAD")
	if err != nil && !isIgnorable(err) {
		return git.Commit{}, err
	}

	if c == nil {
		return git.Commit{}, nil
	}

	return *c, nil
}

// isIgnorable returns true if the GIT error is ignorable.
func isIgnorable(err error) bool {
	switch err {
	case git.ErrLookup, git.ErrNoRepo, git.ErrDirty:
		return true
	default:
		return false
	}
}
