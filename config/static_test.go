package config

import (
	"os"
	"testing"

	"github.com/tj/assert"
)

func TestStatic(t *testing.T) {
	cwd, _ := os.Getwd()

	table := []struct {
		Static
		valid bool
	}{
		{Static{Dir: cwd}, true},
		{Static{Dir: cwd + "/static_test.go"}, false},
	}

	for _, row := range table {
		if row.valid {
			assert.NoError(t, row.Validate())
		} else {
			assert.Error(t, row.Validate())
		}
	}
}
