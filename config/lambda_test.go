package config

import (
	"testing"

	"github.com/tj/assert"
)

func TestLambda(t *testing.T) {
	c := &Lambda{}
	assert.NoError(t, c.Default(), "default")
	assert.Equal(t, 5, c.Timeout, "timeout")
	assert.Equal(t, 128, c.Memory, "timeout")
}
