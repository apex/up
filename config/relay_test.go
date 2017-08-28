package config

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/tj/assert"
)

func TestBackoff_Default(t *testing.T) {
	a := &Backoff{}
	assert.NoError(t, a.Default(), "default")

	b := &Backoff{
		Min:      100,
		Max:      500,
		Factor:   2,
		Attempts: 3,
	}

	assert.Equal(t, b, a)
}

func TestBackoff_Backoff(t *testing.T) {
	a := &Backoff{}
	assert.NoError(t, a.Default(), "default")

	b := a.Backoff()
	assert.Equal(t, time.Millisecond*100, b.Min)
	assert.Equal(t, time.Millisecond*500, b.Max)
}

func TestRelay_Default(t *testing.T) {
	r := Relay{}
	b := Backoff{}
	b.Default()
	assert.NoError(t, r.Default(), "default")

	assert.Equal(t, r.Command, "./server")
	assert.Equal(t, r.Backoff, b)
	assert.True(t, *r.GzipCompression)
}

func TestRelay_Custom(t *testing.T) {
	r := Relay{}
	b := Backoff{}
	b.Default()
	b.Max = 100
	r.Backoff = b
	r.Command = "./test"
	r.GzipCompression = aws.Bool(false)
	assert.NoError(t, r.Default(), "default")

	assert.Equal(t, r.Command, "./test")
	assert.Equal(t, r.Backoff, b)
	assert.False(t, *r.GzipCompression)
}
