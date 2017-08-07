package zip

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/tj/go-archive"
)

// Build the given `dir`.
func Build(dir string) (io.ReadCloser, *archive.Stats, error) {
	gitignore, err := read(".gitignore")
	if err != nil {
		return nil, nil, errors.Wrap(err, "reading .gitignore")
	}
	defer gitignore.Close()

	npmignore, err := read(".npmignore")
	if err != nil {
		return nil, nil, errors.Wrap(err, "reading .npmignore")
	}
	defer npmignore.Close()

	upignore, err := read(".upignore")
	if err != nil {
		return nil, nil, errors.Wrap(err, "reading .upignore")
	}
	defer upignore.Close()

	r := io.MultiReader(
		gitignore,
		strings.NewReader("\n"),
		npmignore,
		strings.NewReader("\n!node_modules\n"),
		upignore)

	filter, err := archive.FilterPatterns(r)
	if err != nil {
		return nil, nil, errors.Wrap(err, "parsing ignore patterns")
	}

	buf := new(bytes.Buffer)
	zip := archive.NewZip(buf).WithFilter(filter)

	if err := zip.Open(); err != nil {
		return nil, nil, errors.Wrap(err, "opening")
	}

	if err := zip.AddDir(dir); err != nil {
		return nil, nil, errors.Wrap(err, "adding dir")
	}

	if err := zip.Close(); err != nil {
		return nil, nil, errors.Wrap(err, "closing")
	}

	return ioutil.NopCloser(buf), zip.Stats(), nil
}

// read file.
func read(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)

	if os.IsNotExist(err) {
		return ioutil.NopCloser(bytes.NewReader(nil)), nil
	}

	if err != nil {
		return nil, err
	}

	return f, nil
}
