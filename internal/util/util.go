// Package util haters gonna hate.
package util

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/apex/up/internal/colors"
	humanize "github.com/dustin/go-humanize"
	"github.com/pascaldekloe/name"
	"github.com/pkg/errors"
	"github.com/tj/backoff"
	"github.com/tj/go-progress"
	"github.com/tj/go/term"
	"golang.org/x/net/publicsuffix"
)

// ClearHeader removes all content header fields.
func ClearHeader(h http.Header) {
	h.Del("Content-Type")
	h.Del("Content-Length")
	h.Del("Content-Encoding")
	h.Del("Content-Range")
	h.Del("Content-MD5")
	h.Del("Cache-Control")
	h.Del("ETag")
	h.Del("Last-Modified")
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
	case strings.Contains(err.Error(), "ResourceNotFoundException"):
		return true
	case strings.Contains(err.Error(), "does not exist"):
		return true
	case strings.Contains(err.Error(), "not found"):
		return true
	default:
		return false
	}
}

// IsBucketExists returns true if err is not nil and represents an existing bucket.
func IsBucketExists(err error) bool {
	switch {
	case err == nil:
		return false
	case strings.Contains(err.Error(), "BucketAlreadyOwnedByYou"):
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

// IsNoCredentials returns true if err is not nil and represents missing credentials.
func IsNoCredentials(err error) bool {
	switch {
	case err == nil:
		return false
	case strings.Contains(err.Error(), "NoCredentialProviders"):
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

// DefaultString returns d unless s is present.
func DefaultString(s *string, d string) string {
	if s == nil || *s == "" {
		return d
	}

	return *s
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

// LogClear clears the line and outputs a log message.
func LogClear(msg string, v ...interface{}) {
	term.MoveUp(1)
	term.ClearLine()
	fmt.Printf("\r     %s\n", colors.Purple(fmt.Sprintf(msg, v...)))
}

// LogTitle outputs a log title.
func LogTitle(msg string, v ...interface{}) {
	fmt.Printf("\n     \x1b[1m%s\x1b[m\n\n", fmt.Sprintf(msg, v...))
}

// LogName outputs a log message with name.
func LogName(name, msg string, v ...interface{}) {
	fmt.Printf("     %s %s\n", colors.Purple(name+":"), fmt.Sprintf(msg, v...))
}

// LogListItem outputs a list item.
func LogListItem(msg string, v ...interface{}) {
	fmt.Printf("      • %s\n", fmt.Sprintf(msg, v...))
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

// Domain returns the effective domain (TLD plus one).
func Domain(s string) string {
	d, err := publicsuffix.EffectiveTLDPlusOne(s)
	if err != nil {
		panic(errors.Wrap(err, "effective domain"))
	}

	return d
}

// CertDomainNames returns the certificate domain name
// and alternative names for a requested domain.
func CertDomainNames(s string) []string {
	// effective domain
	if Domain(s) == s {
		return []string{s, "*." + s}
	}

	// subdomain
	return []string{RemoveSubdomains(s, 1), "*." + RemoveSubdomains(s, 1)}
}

// IsWildcardDomain returns true if the domain is a wildcard.
func IsWildcardDomain(s string) bool {
	return strings.HasPrefix(s, "*.")
}

// WildcardMatches returns true if wildcard is a wildcard domain
// and it satisfies the given domain.
func WildcardMatches(wildcard, domain string) bool {
	if !IsWildcardDomain(wildcard) {
		return false
	}

	w := RemoveSubdomains(wildcard, 1)
	d := RemoveSubdomains(domain, 1)
	return w == d
}

// RemoveSubdomains returns the domain without the n left-most subdomain(s).
func RemoveSubdomains(s string, n int) string {
	domains := strings.Split(s, ".")
	return strings.Join(domains[n:], ".")
}

// ParseSections returns INI style sections from r.
func ParseSections(r io.Reader) (sections []string, err error) {
	s := bufio.NewScanner(r)

	for s.Scan() {
		t := s.Text()
		if strings.HasPrefix(t, "[") {
			sections = append(sections, strings.Trim(t, "[]"))
		}
	}

	err = s.Err()
	return
}

// StringMapKeys returns keys for m.
func StringMapKeys(m map[string]string) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return
}

// UniqueStrings returns a string slice of unique values.
func UniqueStrings(s []string) (v []string) {
	m := make(map[string]struct{})
	for _, val := range s {
		_, ok := m[val]
		if !ok {
			v = append(v, val)
			m[val] = struct{}{}
		}
	}
	return
}

// IsCI returns true if the env looks like it's a CI platform.
func IsCI() bool {
	return os.Getenv("CI") == "true"
}

// EnvironMap returns environment as a map.
func EnvironMap() map[string]string {
	m, _ := ParseEnviron(os.Environ())
	return m
}

// ParseEnviron returns environment as a map from the given env slice.
func ParseEnviron(env []string) (map[string]string, error) {
	m := make(map[string]string)

	for i := 0; i < len(env); i++ {
		s := env[i]

		if strings.ContainsRune(s, '=') {
			p := strings.SplitN(s, "=", 2)
			m[p[0]] = p[1]
			continue
		}

		if i == len(env)-1 {
			return nil, errors.Errorf("%q is missing a value", s)
		}

		m[s] = env[i+1]
		i++
	}

	return m, nil
}

// EncodeAlias encodes an alias string so that it conforms to the
// requirement of matching (?!^[0-9]+$)([a-zA-Z0-9-_]+).
func EncodeAlias(s string) string {
	return "commit-" + strings.Replace(s, ".", "_", -1)
}

// DecodeAlias decodes an alias string which was encoded by
// the EncodeAlias function.
func DecodeAlias(s string) string {
	s = strings.Replace(s, "_", ".", -1)
	s = strings.Replace(s, "commit-", "", 1)
	return s
}

// DateSuffix returns the date suffix for t.
func DateSuffix(t time.Time) string {
	switch t.Day() {
	case 1, 21, 31:
		return "st"
	case 2, 22:
		return "nd"
	case 3, 23:
		return "rd"
	default:
		return "th"
	}
}

// StripLerna strips the owner portion of a Lerna-based tag. See #670 for
// details. They are in the form of "@owner/repo@0.5.0".
func StripLerna(s string) string {
	if strings.HasPrefix(s, "@") {
		p := strings.Split(s, "@")
		return p[len(p)-1]
	}

	return s
}

// FixMultipleSetCookie staggers the casing of each set-cookie
// value to trick API Gateway into setting multiple in the response.
func FixMultipleSetCookie(h http.Header) {
	cookies := h["Set-Cookie"]

	if len(cookies) == 0 {
		return
	}

	h.Del("Set-Cookie")

	for i, v := range cookies {
		h[BinaryCase("set-cookie", i)] = []string{v}
	}
}

// BinaryCase ported from https://github.com/Gi60s/binary-case/blob/master/index.js#L86.
func BinaryCase(s string, n int) string {
	var res []rune

	for _, c := range s {
		if c >= 65 && c <= 90 {
			if n&1 > 0 {
				c += 32
			}
			res = append(res, c)
			n >>= 1
		} else if c >= 97 && c <= 122 {
			if n&1 > 0 {
				c -= 32
			}
			res = append(res, c)
			n >>= 1
		} else {
			res = append(res, c)
		}
	}

	return string(res)
}

// RelativeDate returns a date formatted relative to now.
func RelativeDate(t time.Time) string {
	switch d := time.Since(t); {
	case d <= 12*time.Hour:
		return humanize.RelTime(time.Now(), t, "from now", "ago")
	case d <= 24*time.Hour:
		return t.Format(`Today at 03:04:05pm`)
	case d <= 24*time.Hour*2:
		return t.Format(`Yesterday at 03:04:05pm`)
	default:
		return t.Format(`Jan 2` + DateSuffix(t) + ` 03:04:05pm`)
	}
}

var numericRe = regexp.MustCompile(`^[0-9]+$`)

// IsNumeric returns true if s is numeric.
func IsNumeric(s string) bool {
	return numericRe.MatchString(s)
}
