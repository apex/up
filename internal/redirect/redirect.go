// Package redirect provides compiling and matching
// redirect and rewrite rules.
package redirect

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fanyang01/radix"
	"github.com/pkg/errors"
)

// placeholders regexp.
var placeholders = regexp.MustCompile(`:(\w+)`)

// Rule is a single redirect rule.
type Rule struct {
	Path     string `json:"path"`
	Location string `json:"location"`
	Status   int    `json:"status"`
	names    map[string]bool
	dynamic  bool
	sub      string
	path     *regexp.Regexp
}

// URL returns the final destination after substitutions from path.
func (r Rule) URL(path string) string {
	return r.path.ReplaceAllString(path, r.sub)
}

// IsDynamic returns true if a splat or placeholder is used.
func (r *Rule) IsDynamic() bool {
	return r.dynamic
}

// IsRewrite returns true if the rule represents a rewrite.
func (r *Rule) IsRewrite() bool {
	return r.Status == 200 || r.Status == 0
}

// Compile the rule.
func (r *Rule) Compile() {
	r.path, r.names = compilePath(r.Path)
	r.sub = compileSub(r.Path, r.Location, r.names)
	r.dynamic = isDynamic(r.Path)
}

// Rules map of paths to redirects.
type Rules map[string]Rule

// Matcher for header lookup.
type Matcher struct {
	t *radix.PatternTrie
}

// Lookup returns fields for the given path.
func (m *Matcher) Lookup(path string) *Rule {
	v, ok := m.t.Lookup(path)
	if !ok {
		return nil
	}

	r := v.(Rule)
	return &r
}

// Compile the given rules to a trie.
func Compile(rules Rules) (*Matcher, error) {
	t := radix.NewPatternTrie()
	m := &Matcher{t}

	for path, rule := range rules {
		rule.Path = path
		rule.Compile()
		t.Add(compilePattern(path), rule)
		t.Add(compilePattern(path)+"/", rule)
	}

	return m, nil
}

// compileSub returns a substitution string.
func compileSub(path, s string, names map[string]bool) string {
	// splat
	s = strings.Replace(s, `:splat`, `${splat}`, -1)

	// placeholders
	s = placeholders.ReplaceAllStringFunc(s, func(v string) string {
		name := v[1:]

		// TODO: refactor to not panic
		if !names[name] {
			panic(errors.Errorf("placeholder %q is not present in the path pattern %q", v, path))
		}

		return fmt.Sprintf("${%s}", name)
	})

	return s
}

// compilePath returns a regexp for substitutions and return
// a map of placeholder names for validation.
func compilePath(s string) (*regexp.Regexp, map[string]bool) {
	names := make(map[string]bool)

	// escape
	s = regexp.QuoteMeta(s)

	// splat
	s = strings.Replace(s, `\*`, `(?P<splat>.*?)`, -1)

	// placeholders
	s = placeholders.ReplaceAllStringFunc(s, func(v string) string {
		name := v[1:]
		names[name] = true
		return fmt.Sprintf(`(?P<%s>[^/]+)`, name)
	})

	// trailing slash
	s += `\/?`

	s = fmt.Sprintf(`^%s$`, s)
	return regexp.MustCompile(s), names
}

// compilePattern to a syntax usable by the trie.
func compilePattern(s string) string {
	return placeholders.ReplaceAllString(s, "*")
}

// isDynamic returns true for splats or placeholders.
func isDynamic(s string) bool {
	return hasPlaceholder(s) || hasSplat(s)
}

// hasPlaceholder returns true for placeholders
func hasPlaceholder(s string) bool {
	return strings.ContainsRune(s, ':')
}

// hasSplat returns true for splats.
func hasSplat(s string) bool {
	return strings.ContainsRune(s, '*')
}
