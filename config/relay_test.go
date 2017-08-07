package config

import (
	"testing"
	"time"

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
