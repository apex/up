// Package header provides path-matched header injection rules.
package header

import (
	"github.com/fanyang01/radix"
)

// Fields map.
type Fields map[string]string

// Rules map of paths to fields.
type Rules map[string]Fields

// Matcher for header lookup.
type Matcher struct {
	t *radix.PatternTrie
}

// Lookup returns fields for the given path.
func (m *Matcher) Lookup(path string) Fields {
	v, ok := m.t.Lookup(path)
	if !ok {
		return nil
	}

	return v.(Fields)
}

// Compile the given rules to a trie.
func Compile(rules Rules) (*Matcher, error) {
	t := radix.NewPatternTrie()
	m := &Matcher{t}

	for path, fields := range rules {
		t.Add(path, fields)
	}

	return m, nil
}

// Merge returns a new rules set giving precedence to `b`.
func Merge(a, b Rules) Rules {
	r := make(Rules)

	for path, fields := range a {
		r[path] = fields
	}

	for path, fields := range b {
		if _, ok := r[path]; !ok {
			r[path] = make(Fields)
		}

		for name, val := range fields {
			r[path][name] = val
		}
	}

	return r
}
