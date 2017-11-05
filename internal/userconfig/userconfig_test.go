package userconfig

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/tj/assert"
)

func init() {
	configDir = ".up-test"
}

func TestConfig(t *testing.T) {
	t.Run("load when missing", func(t *testing.T) {
		dir, _ := homedir.Dir()
		os.RemoveAll(filepath.Join(dir, configDir))

		c := Config{}
		assert.NoError(t, c.Load(), "load")
	})

	t.Run("save", func(t *testing.T) {
		c := Config{}
		assert.NoError(t, c.Load(), "load")
		assert.Equal(t, "", c.Token)

		c.Token = "foo-bar-baz"
		assert.NoError(t, c.Save(), "save")
	})

	t.Run("load after save", func(t *testing.T) {
		c := Config{}
		assert.NoError(t, c.Load(), "save")
		assert.Equal(t, "foo-bar-baz", c.Token)
	})
}
