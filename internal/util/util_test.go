package util

import (
	"os/exec"
	"testing"

	"github.com/tj/assert"
)

func TestExitStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cmd := exec.Command("echo", "hello", "world")
		code := ExitStatus(cmd, cmd.Run())
		assert.Equal(t, "0", code)
	})

	t.Run("missing", func(t *testing.T) {
		cmd := exec.Command("nope")
		code := ExitStatus(cmd, cmd.Run())
		assert.Equal(t, "?", code)
	})

	t.Run("failure", func(t *testing.T) {
		cmd := exec.Command("sh", "-c", `echo hello && exit 5`)
		code := ExitStatus(cmd, cmd.Run())
		assert.Equal(t, "5", code)
	})
}

func TestIsSubdomain(t *testing.T) {
	assert.Equal(t, false, IsSubdomain("gh-polls.com"))
	assert.Equal(t, true, IsSubdomain("staging.gh-polls.com"))
	assert.Equal(t, true, IsSubdomain("blog.gh-polls.com"))
	assert.Equal(t, true, IsSubdomain("foo.bar.gh-polls.com"))
}

func TestDomain(t *testing.T) {
	assert.Equal(t, "gh-polls.com", Domain("gh-polls.com"))
	assert.Equal(t, "gh-polls.com", Domain("staging.gh-polls.com"))
	assert.Equal(t, "gh-polls.com", Domain("blog.gh-polls.com"))
	assert.Equal(t, "gh-polls.com", Domain("foo.api.gh-polls.com"))
}
