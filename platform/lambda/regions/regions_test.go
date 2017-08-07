package regions

import (
	"testing"

	"github.com/tj/assert"
)

func TestMatch(t *testing.T) {
	t.Run("explicit", func(t *testing.T) {
		v := Match([]string{"us-west-2", "us-east-1"})
		assert.Equal(t, []string{"us-west-2", "us-east-1"}, v)
	})

	t.Run("glob all", func(t *testing.T) {
		v := Match([]string{"*"})
		assert.Equal(t, All, v)
	})

	t.Run("glob some", func(t *testing.T) {
		v := Match([]string{"us-west-*", "ca-*"})
		e := []string{"us-west-2", "us-west-1", "ca-central-1"}
		assert.Equal(t, e, v)
	})
}
