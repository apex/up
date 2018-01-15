package cost

import (
	"testing"

	"github.com/tj/assert"
)

func TestRequests(t *testing.T) {
	table := []struct {
		requests int
		expected float64
	}{
		{0, 0.0},
		{1000, 0.001},
		{1000000, 1.0},
	}

	for _, row := range table {
		assert.Equal(t, row.expected, Requests(row.requests))
	}
}

func TestRate(t *testing.T) {
	table := []struct {
		memory   int
		expected float64
	}{
		{-1, 0.0},
		{0, 0.0},
		{128, 2.08e-7},
		{156, 0.0},
	}

	for _, row := range table {
		assert.Equal(t, row.expected, Rate(row.memory))
	}
}

func TestInvocations(t *testing.T) {
	table := []struct {
		invocations int
		expected    float64
	}{
		{0, 0.0},
		{1, 2.0e-7},
		{1.0e7, 2.0},
	}

	for _, row := range table {
		assert.Equal(t, row.expected, Invocations(row.invocations))
	}
}

func TestDuration(t *testing.T) {
	table := []struct {
		duration int
		memory   int
		expected float64
	}{
		{0, 128, 0},
		{100000, 256, 4.17e-4},
		{1e8, 1536, 2.501},
	}

	for _, row := range table {
		assert.Equal(t, row.expected, Duration(row.duration, row.memory))
	}
}
