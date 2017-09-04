package header

import (
	"testing"

	"github.com/tj/assert"
)

func TestMatcher_Lookup(t *testing.T) {
	rules := Rules{
		"*": {
			"X-Type": "html",
		},
		"*.css": {
			"X-Type": "css",
		},
		"/docs/alerts": {
			"X-Type": "docs alerts",
		},
		"/docs/*": {
			"X-Type": "docs",
		},
	}

	m, err := Compile(rules)
	assert.NoError(t, err, "compile")

	assert.Equal(t, Fields{"X-Type": "html"}, m.Lookup("/something"))
	assert.Equal(t, Fields{"X-Type": "html"}, m.Lookup("/docs"))
	assert.Equal(t, Fields{"X-Type": "docs"}, m.Lookup("/docs/"))
	assert.Equal(t, Fields{"X-Type": "css"}, m.Lookup("/style.css"))
	assert.Equal(t, Fields{"X-Type": "css"}, m.Lookup("/public/css/style.css"))
	assert.Equal(t, Fields{"X-Type": "docs"}, m.Lookup("/docs/checks"))
	assert.Equal(t, Fields{"X-Type": "docs alerts"}, m.Lookup("/docs/alerts"))
}

func TestMerge(t *testing.T) {
	rules := Rules{
		"*": {
			"X-Type": "html",
			"X-Foo":  "bar",
		},
		"/login": {
			"X-Something": "here",
		},
	}

	rules = Merge(rules, Rules{
		"*": {
			"X-Type": "pdf",
		},
		"/admin": {
			"X-Something": "here",
		},
	})

	expected := Rules{
		"*": {
			"X-Type": "pdf",
			"X-Foo":  "bar",
		},
		"/login": {
			"X-Something": "here",
		},
		"/admin": {
			"X-Something": "here",
		},
	}

	assert.Equal(t, expected, rules)
}
