package util

import (
	"os/exec"
	"strings"
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

	t.Run("week", func(t *testing.T) {
		v, err := ParseDuration("1w")
		assert.NoError(t, err, "parsing")
		assert.Equal(t, time.Hour*24*7, v)

		v, err = ParseDuration("2w")
		assert.NoError(t, err, "parsing")
		assert.Equal(t, time.Hour*24*7*2, v)
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

func TestDomain(t *testing.T) {
	assert.Equal(t, "example.com", Domain("example.com"))
	assert.Equal(t, "example.com", Domain("api.example.com"))
	assert.Equal(t, "example.com", Domain("v1.api.example.com"))

	assert.Equal(t, "example.co.uk", Domain("example.co.uk"))
	assert.Equal(t, "example.co.uk", Domain("api.example.co.uk"))
	assert.Equal(t, "example.co.uk", Domain("v1.api.example.co.uk"))
}

func TestParseSections(t *testing.T) {
	r := strings.NewReader(`[personal]
aws_access_key_id = personal_key
aws_secret_access_key = personal_secret
[app]
aws_access_key_id = app_key
aws_secret_access_key = app_secret
[foo_bar]
aws_access_key_id = foo_bar_key
aws_secret_access_key = foo_bar_secret
`)

	v, err := ParseSections(r)
	assert.NoError(t, err)

	assert.Equal(t, []string{"personal", "app", "foo_bar"}, v)
}
