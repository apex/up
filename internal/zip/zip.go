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

var transform = archive.TransformFunc(func(r io.Reader, i os.FileInfo) (io.Reader, os.FileInfo) {
	name := strings.Replace(i.Name(), "\\", "/", -1)

	i = archive.Info{
		Name:     name,
		Size:     i.Size(),
		Mode:     i.Mode() | 0555,
		Modified: i.ModTime(),
		Dir:      i.IsDir(),
	}.FileInfo()

	return r, i
})

// Build the given `dir`.
func Build(dir string) (io.ReadCloser, *archive.Stats, error) {
	upignore, err := read(".upignore")
	if err != nil {
		return nil, nil, errors.Wrap(err, "reading .upignore")
	}
	defer upignore.Close()

	r := io.MultiReader(
		strings.NewReader(".*\n"),
		strings.NewReader("\n!node_modules/**\n!.pypath/**\n"),
		upignore,
		strings.NewReader("\n!main\n!server\n!_proxy.js\n!byline.js\n!up.json\n!pom.xml\n!build.gradle\n!project.clj\n"))

	filter, err := archive.FilterPatterns(r)
	if err != nil {
		return nil, nil, errors.Wrap(err, "parsing ignore patterns")
	}

	buf := new(bytes.Buffer)
	zip := archive.NewZip(buf).
		WithFilter(filter).
		WithTransform(transform)

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
