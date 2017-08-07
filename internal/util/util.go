// Package util haters gonna hate.
package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/apex/up/internal/colors"
	"github.com/pascaldekloe/name"
	"github.com/pkg/errors"
	"github.com/tj/go-progress"
)

// Fields retained when clearing.
var keepFields = map[string]bool{
	"X-Powered-By": true,
}

// ClearHeader removes all header fields.
func ClearHeader(h http.Header) {
	for k := range h {
		if keepFields[k] {
			continue
		}

		h.Del(k)
	}
}

// ManagedByUp string.
func ManagedByUp(s string) string {
	if s == "" {
		return "Managed by Up."
	}

	return s + " (Managed by Up)."
}

// Exists returns true if the file exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ReadFileJSON reads json from the given path.
func ReadFileJSON(path string, v interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "reading")
	}

	if err := json.Unmarshal(b, &v); err != nil {
		return errors.Wrap(err, "unmarshaling")
	}

	return nil
}

// Camelcase .....
func Camelcase(s string, v ...interface{}) string {
	return name.CamelCase(fmt.Sprintf(s, v...), true)
}

// NewProgressInt with the given total.
func NewProgressInt(total int) *progress.Bar {
	b := progress.NewInt(total)
	b.Template(`{{.Bar}} {{.Percent | printf "%0.0f"}}% {{.Text}}`)
	b.Width = 35
	b.StartDelimiter = colors.Gray("|")
	b.EndDelimiter = colors.Gray("|")
	b.Filled = colors.Purple("█")
	b.Empty = colors.Gray("░")
	return b
}

// Pad helper.
func Pad() func() {
	println()
	return func() {
		println()
	}
}

// Fatal error.
func Fatal(err error) {
	fmt.Fprintf(os.Stderr, "\n  %s %s\n\n", colors.Red("Error:"), err)
	os.Exit(1)
}

// IsJSON returns true if the msg looks like json.
func IsJSON(s string) bool {
	return len(s) > 1 && s[0] == '{' && s[len(s)-1] == '}'
}
