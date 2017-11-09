package util

import (
	"os/exec"
	"testing"
	"time"

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

func TestParseDuration(t *testing.T) {
	t.Run("day", func(t *testing.T) {
		v, err := ParseDuration("1d")
		assert.NoError(t, err, "parsing")
		assert.Equal(t, time.Hour*24, v)
	})

	t.Run("day with faction", func(t *testing.T) {
		v, err := ParseDuration("1.5d")
		assert.NoError(t, err, "parsing")
		assert.Equal(t, time.Duration(float64(time.Hour*24)*1.5), v)
	})

	t.Run("month", func(t *testing.T) {
		v, err := ParseDuration("1mo")
		assert.NoError(t, err, "parsing")
		assert.Equal(t, time.Hour*24*30, v)

		v, err = ParseDuration("1M")
		assert.NoError(t, err, "parsing")
		assert.Equal(t, time.Hour*24*30, v)
	})

	t.Run("month with faction", func(t *testing.T) {
		v, err := ParseDuration("1.5mo")
		assert.NoError(t, err, "parsing")
		assert.Equal(t, time.Duration(float64(time.Hour*24*30)*1.5), v)
	})

	t.Run("default", func(t *testing.T) {
		v, err := ParseDuration("15m")
		assert.NoError(t, err, "parsing")
		assert.Equal(t, 15*time.Minute, v)
	})
}
