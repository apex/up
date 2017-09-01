package cost

import (
	"testing"

	"github.com/tj/assert"
)

func TestRequests(t *testing.T) {
	table := []struct {
		count    int
		expected float64
	}{
		{0, 0.0},
		{1000, 0.001},
		{1000000, 1.0},
	}

	for _, row := range table {
		assert.Equal(t, row.expected, Requests(row.count))
	}
}
