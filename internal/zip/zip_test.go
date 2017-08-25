package zip

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"

	"github.com/tj/assert"
)

// TODO: better tests

func TestBuild(t *testing.T) {
	os.Chdir("testdata")
	defer os.Chdir("..")

	zip, _, err := Build(".")
	assert.NoError(t, err)

	out, err := ioutil.TempDir(os.TempDir(), "-up")
	assert.NoError(t, err, "tmpdir")
	dst := filepath.Join(out, "out.zip")

	f, err := os.Create(dst)
	assert.NoError(t, err, "create")

	_, err = io.Copy(f, zip)
	assert.NoError(t, err, "copy")

	assert.NoError(t, f.Close(), "close")

	cmd := exec.Command("unzip", "out.zip")
	cmd.Dir = out
	assert.NoError(t, cmd.Run(), "unzip")

	files, err := ioutil.ReadDir(out)
	assert.NoError(t, err, "readdir")

	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}
	sort.Strings(names)

	assert.Equal(t, []string{"bar.js", "foo.js", "index.js", "out.zip"}, names)
}
