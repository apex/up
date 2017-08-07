package redirect

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func rule(from, to string) Rule {
	r := Rule{
		Path:     from,
		Location: to,
	}

	r.Compile()
	return r
}

func TestRule_URL(t *testing.T) {
	t.Run("exact", func(t *testing.T) {
		s := rule("/docs", "/help").URL("/docs")
		assert.Equal(t, "/help", s)
	})

	t.Run("splat one segment", func(t *testing.T) {
		r := rule("/docs/*", "/help/:splat")
		assert.Equal(t, "/help/foo", r.URL("/docs/foo"))
	})

	t.Run("splat many segments", func(t *testing.T) {
		r := rule("/docs/*", "/help/:splat")
		assert.Equal(t, "/help/foo/bar/baz", r.URL("/docs/foo/bar/baz"))
	})

	t.Run("placeholder", func(t *testing.T) {
		r := rule("/shop/:brand", "/store/:brand")
		assert.Equal(t, "/store/apple", r.URL("/shop/apple"))
	})

	t.Run("placeholders", func(t *testing.T) {
		r := rule("/shop/:brand/category/:cat", "/products/:brand/:cat")
		assert.Equal(t, "/products/apple/laptops", r.URL("/shop/apple/category/laptops"))
	})

	t.Run("placeholders trailing slash", func(t *testing.T) {
		r := rule("/docs/:product/guides/:guide", "/help/:product/:guide")
		assert.Equal(t, "/help/ping/alerting", r.URL("/docs/ping/guides/alerting/"))
	})

	t.Run("placeholders rearranged", func(t *testing.T) {
		r := rule("/shop/:brand/category/:cat", "/products/:cat/:brand")
		assert.Equal(t, "/products/laptops/apple", r.URL("/shop/apple/category/laptops"))
	})

	t.Run("placeholders mismatch", func(t *testing.T) {
		// TODO: sorry :D
		err := func() (err error) {
			defer func() {
				err = recover().(error)
			}()

			rule("/shop/:brand/category/:category", "/products/:cat/:brand")
			return nil
		}()

		assert.EqualError(t, err, `placeholder ":cat" is not present in the path pattern "/shop/:brand/category/:category"`)
	})
}

func TestMatcher_Lookup(t *testing.T) {
	rules := Rules{
		"/docs/:product/guides/:guide": Rule{
			Location: "/help/:product/:guide",
			Status:   301,
		},
		"/blog": Rule{
			Location: "https://blog.apex.sh",
			Status:   302,
		},
		"/articles/*": Rule{
			Location: "/guides/:splat",
		},
	}

	m, err := Compile(rules)
	assert.NoError(t, err, "compile")

	t.Run("exact", func(t *testing.T) {
		assert.NotNil(t, m.Lookup("/blog"))
	})

	t.Run("exact trailing slash", func(t *testing.T) {
		assert.NotNil(t, m.Lookup("/blog/"))
	})

	t.Run("placeholders", func(t *testing.T) {
		assert.NotNil(t, m.Lookup("/docs/ping/guides/alerts"))
	})

	// TODO: need to fork the trie to be less greedy
	// t.Run("mismatch", func(t *testing.T) {
	// 	r := m.Lookup("/docs/ping/another/guides/alerts")
	// 	assert.NotNil(t, r)
	// })

	t.Run("splat one segment", func(t *testing.T) {
		assert.NotNil(t, m.Lookup("/articles/alerting"))
	})

	t.Run("splat many segments", func(t *testing.T) {
		assert.NotNil(t, m.Lookup("/articles/alerting/pagerduty"))
		assert.NotNil(t, m.Lookup("/articles/alerting/pagerduty/"))
	})
}

func BenchmarkMatcher_Lookup(b *testing.B) {
	rules := Rules{
		"/docs/:product/guides/:guide": Rule{
			Location: "/help/:product/:guide",
			Status:   301,
		},
	}

	m, err := Compile(rules)
	assert.NoError(b, err, "compile")

	b.ResetTimer()

	b.Run("match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fmt.Printf("%#v\n", m.Lookup("/docs/ping/guides/alerts"))
		}
	})

	b.Run("mismatch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.Lookup("/some/other/page")
		}
	})
}
