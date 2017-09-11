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
