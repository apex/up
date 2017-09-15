// Package util haters gonna hate.
package util

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/apex/up/internal/colors"
	"github.com/pascaldekloe/name"
	"github.com/pkg/errors"
	"github.com/tj/backoff"
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

// ManagedByUp appends "Managed by Up".
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

// Camelcase string with optional args.
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

// NewInlineProgressInt with the given total.
func NewInlineProgressInt(total int) *progress.Bar {
	b := progress.NewInt(total)
	b.Template(`{{.Bar}} {{.Percent | printf "%0.0f"}}% {{.Text}}`)
	b.Width = 20
	b.StartDelimiter = colors.Gray("|")
	b.EndDelimiter = colors.Gray("|")
	b.Filled = colors.Purple("█")
	b.Empty = colors.Gray(" ")
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

// IsNotFound returns true if err is not nil and represents a missing resource.
func IsNotFound(err error) bool {
	switch {
	case err == nil:
		return false
	case strings.Contains(err.Error(), "does not exist"):
		return true
	case strings.Contains(err.Error(), "not found"):
		return true
	default:
		return false
	}
}

// IsThrottled returns true if err is not nil and represents a throttled request.
func IsThrottled(err error) bool {
	switch {
	case err == nil:
		return false
	case strings.Contains(err.Error(), "Throttling: Rate exceeded"):
		return true
	default:
		return false
	}
}

// Env returns a slice from environment variable map.
func Env(m map[string]string) (env []string) {
	for k, v := range m {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return
}

// PrefixLines prefixes the lines in s with prefix.
func PrefixLines(s string, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = prefix + l
	}
	return strings.Join(lines, "\n")
}

// Indent the given string.
func Indent(s string) string {
	return PrefixLines(s, "  ")
}

// WaitForListen blocks until `u` is listening with timeout.
func WaitForListen(u *url.URL, timeout time.Duration) error {
	timedout := time.After(timeout)

	b := backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    time.Second,
		Factor: 1.5,
	}

	for {
		select {
		case <-timedout:
			return errors.Errorf("timed out after %s", timeout)
		case <-time.After(b.Duration()):
			if IsListening(u) {
				return nil
			}
		}
	}
}

// IsListening returns true if there's a server listening on `u`.
func IsListening(u *url.URL) bool {
	conn, err := net.Dial("tcp", u.Host)
	if err != nil {
		return false
	}

	conn.Close()
	return true
}

// ExitStatus returns the exit status of cmd.
func ExitStatus(cmd *exec.Cmd, err error) string {
	ps := cmd.ProcessState

	if e, ok := err.(*exec.ExitError); ok {
		ps = e.ProcessState
	}

	if ps != nil {
		s, ok := ps.Sys().(syscall.WaitStatus)
		if ok {
			return fmt.Sprintf("%d", s.ExitStatus())
		}
	}

	return "?"
}

// MaybeClose closes v if it is an io.Closer.
func MaybeClose(v interface{}) error {
	if c, ok := v.(io.Closer); ok {
		return c.Close()
	}

	return nil
}

// IsSubdomain returns true if s is a subdomain.
func IsSubdomain(s string) bool {
	return len(strings.Split(s, ".")) > 2
}

// Domain returns the domain devoid of any subdomain(s).
func Domain(s string) string {
	p := strings.Split(s, ".")
	return strings.Join(p[len(p)-2:], ".")
}
