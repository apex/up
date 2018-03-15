package config

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/tj/assert"
)

func TestLambda(t *testing.T) {
	c := &Lambda{}
	assert.NoError(t, c.Default(), "default")
	assert.Equal(t, 60, c.Timeout, "timeout")
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

func TestLambda_Policy(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		c := &Lambda{}
		assert.NoError(t, c.Default(), "default")
		assert.Len(t, c.Policy, 1)
		assert.Equal(t, defaultPolicy, c.Policy[0])
	})

	t.Run("specified", func(t *testing.T) {
		c := &Lambda{
			Policy: []IAMPolicyStatement{
				{
					"Effect":   "Allow",
					"Resource": "*",
					"Action": []string{
						"s3:List*",
						"s3:Get*",
					},
				},
			},
		}

		assert.NoError(t, c.Default(), "default")
		assert.Len(t, c.Policy, 2)
		assert.Equal(t, defaultPolicy, c.Policy[1])
	})
}
