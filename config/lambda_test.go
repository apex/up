package config

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/tj/assert"
)

func TestLambda(t *testing.T) {
	c := &Lambda{}
	assert.NoError(t, c.Default(), "default")
	assert.Equal(t, 0, c.Timeout, "timeout")
	assert.Equal(t, 512, c.Memory, "timeout")
}

func TestLambda_Override(t *testing.T) {
	c := &Config{}

	l := &Lambda{
		Warm:      aws.Bool(true),
		WarmCount: 30,
	}

	l.Override(c)

	assert.Equal(t, true, *c.Lambda.Warm)
	assert.Equal(t, 30, c.Lambda.WarmCount)
}
