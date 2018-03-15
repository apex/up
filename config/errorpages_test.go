package config

import (
	"testing"

	"github.com/tj/assert"
)

func TestErrorPages(t *testing.T) {
	c := &ErrorPages{}
	assert.NoError(t, c.Default(), "default")
	assert.Equal(t, ".", c.Dir, "dir")
}
