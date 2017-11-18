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

func TestConfig_file(t *testing.T) {
	t.Run("load when missing", func(t *testing.T) {
		dir, _ := homedir.Dir()
		os.RemoveAll(filepath.Join(dir, configDir))

		c := Config{}
		assert.NoError(t, c.Load(), "load")
	})

	t.Run("save", func(t *testing.T) {
		c := Config{}
		assert.NoError(t, c.Load(), "load")
		assert.Equal(t, "", c.Team)

		c.Team = "apex"
		assert.NoError(t, c.Save(), "save")
	})

	t.Run("load after save", func(t *testing.T) {
		c := Config{}
		assert.NoError(t, c.Load(), "save")
		assert.Equal(t, "apex", c.Team)
	})
}

func TestConfig_env(t *testing.T) {
	t.Run("load", func(t *testing.T) {
		os.Setenv("UP_CONFIG", `{ "team": "tj@apex.sh" }`)
		c := Config{}
		assert.NoError(t, c.Load(), "load")
		assert.Equal(t, "tj@apex.sh", c.Team)
	})
}
