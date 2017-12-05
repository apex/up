// Package util haters gonna hate.
package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
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
	"golang.org/x/net/publicsuffix"
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
	fmt.Fprintf(os.Stderr, "\n     %s %s\n\n", colors.Red("Error:"), err)
	os.Exit(1)
}

// IsJSON returns true if the string looks like json.
func IsJSON(s string) bool {
	return len(s) > 1 && s[0] == '{' && s[len(s)-1] == '}'
}

// IsJSONLog returns true if the string looks likes a json log.
func IsJSONLog(s string) bool {
	return IsJSON(s) && strings.Contains(s, `"level"`)
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

// StringsContains returns true if list contains s.
func StringsContains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// BasePath returns a normalized base path,
// stripping the leading '/' if present.
func BasePath(s string) string {
	return strings.TrimLeft(s, "/")
}

// LogPad outputs a log message with padding.
func LogPad(msg string, v ...interface{}) {
	defer Pad()()
	Log(msg, v...)
}

// Log outputs a log message.
func Log(msg string, v ...interface{}) {
	fmt.Printf("     %s\n", colors.Purple(fmt.Sprintf(msg, v...)))
}

// LogTitle outputs a log title.
func LogTitle(msg string, v ...interface{}) {
	fmt.Printf("\n     \x1b[1m%s\x1b[m\n\n", fmt.Sprintf(msg, v...))
}

// LogName outputs a log message with name.
func LogName(name, msg string, v ...interface{}) {
	fmt.Printf("     %s %s\n", colors.Purple(name+":"), fmt.Sprintf(msg, v...))
}

// ToFloat returns a float or NaN.
func ToFloat(v interface{}) float64 {
	switch n := v.(type) {
	case int:
		return float64(n)
	case int8:
		return float64(n)
	case int16:
		return float64(n)
	case int32:
		return float64(n)
	case int64:
		return float64(n)
	case uint:
		return float64(n)
	case uint8:
		return float64(n)
	case uint16:
		return float64(n)
	case uint32:
		return float64(n)
	case uint64:
		return float64(n)
	case float32:
		return float64(n)
	case float64:
		return n
	default:
		return math.NaN()
	}
}

// Milliseconds returns the duration as milliseconds.
func Milliseconds(d time.Duration) int {
	return int(d / time.Millisecond)
}

// MillisecondsSince returns the duration as milliseconds relative to time t.
func MillisecondsSince(t time.Time) int {
	return int(time.Since(t) / time.Millisecond)
}

// ParseDuration string with day and month approximation support.
func ParseDuration(s string) (d time.Duration, err error) {
	r := strings.NewReader(s)

	switch {
	case strings.HasSuffix(s, "d"):
		var v float64
		_, err = fmt.Fscanf(r, "%fd", &v)
		d = time.Duration(v * float64(24*time.Hour))
	case strings.HasSuffix(s, "w"):
		var v float64
		_, err = fmt.Fscanf(r, "%fw", &v)
		d = time.Duration(v * float64(24*time.Hour*7))
	case strings.HasSuffix(s, "mo"):
		var v float64
		_, err = fmt.Fscanf(r, "%fmo", &v)
		d = time.Duration(v * float64(30*24*time.Hour))
	case strings.HasSuffix(s, "M"):
		var v float64
		_, err = fmt.Fscanf(r, "%fM", &v)
		d = time.Duration(v * float64(30*24*time.Hour))
	default:
		d, err = time.ParseDuration(s)
	}

	return
}

// Md5 returns an md5 hash for s.
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Domain returns the effective domain. For example
// the string "api.example.com" becomes "example.com",
// while "api.example.co.uk" becomes "example.co.uk".
func Domain(s string) string {
	d, err := publicsuffix.EffectiveTLDPlusOne(s)
	if err != nil {
		panic(errors.Wrap(err, "effective domain"))
	}

	return d
}
