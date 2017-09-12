package linereader

import (
	"io"
	"strings"
	"testing"

	"github.com/tj/assert"
)

type writer struct {
	writes []string
}

// Write implementation.
func (w *writer) Write(b []byte) (int, error) {
	w.writes = append(w.writes, string(b))
	return len(b), nil
}

func TestReader_flatLines(t *testing.T) {
	input := `GET /
GET /account
GET /signup
GET /login
POST /login
`

	w := &writer{}
	r := New(strings.NewReader(input))

	_, err := io.Copy(w, r)
	assert.NoError(t, err)

	expected := []string{
		"GET /",
		"GET /account",
		"GET /signup",
		"GET /login",
		"POST /login",
	}

	assert.Equal(t, expected, w.writes)
}

func TestReader_indentedLines(t *testing.T) {
	input := `GET /
GET /account
GET /signup
POST /login
	user: Tobi
	referrer: something.com
GET /login
UncaughtException: Something exploded
  at foo
  at bar
  at baz
GET /
GET /
`

	w := &writer{}
	r := New(strings.NewReader(input))

	_, err := io.Copy(w, r)
	assert.NoError(t, err)

	expected := []string{
		"GET /",
		"GET /account",
		"GET /signup",
		"POST /login\n\tuser: Tobi\n\treferrer: something.com",
		"GET /login",
		"UncaughtException: Something exploded\n  at foo\n  at bar\n  at baz",
		"GET /",
		"GET /",
	}

	assert.Equal(t, expected, w.writes)
}
