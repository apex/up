package errorpage

import (
	"path/filepath"
	"testing"

	"github.com/tj/assert"
)

// load pages from dir.
func load(t testing.TB, dir string) Pages {
	dir = filepath.Join("testdata", dir)
	pages, err := Load(dir)
	assert.NoError(t, err, "load")
	return pages
}

func TestPages_precedence(t *testing.T) {
	pages := load(t, ".")

	t.Run("code 500 match exact", func(t *testing.T) {
		p := pages.Match(500)
		assert.NotNil(t, p, "no match")

		html, err := p.Render(nil)
		assert.NoError(t, err)

		assert.Equal(t, "500 page.\n", html)
	})

	t.Run("code 404 match exact", func(t *testing.T) {
		p := pages.Match(404)
		assert.NotNil(t, p, "no match")

		html, err := p.Render(nil)
		assert.NoError(t, err)

		assert.Equal(t, "404 page.\n", html)
	})

	t.Run("code 200 match exact", func(t *testing.T) {
		p := pages.Match(200)
		assert.NotNil(t, p, "no match")

		html, err := p.Render(nil)
		assert.NoError(t, err)

		assert.Equal(t, "200 page.\n", html)
	})

	t.Run("code 403 match range", func(t *testing.T) {
		p := pages.Match(403)
		assert.NotNil(t, p, "no match")

		html, err := p.Render(nil)
		assert.NoError(t, err)

		assert.Equal(t, "4xx page.\n", html)
	})

	t.Run("502 match global", func(t *testing.T) {
		p := pages.Match(502)
		assert.NotNil(t, p, "no match")

		data := struct {
			StatusText string
			StatusCode int
		}{"Bad Gateway", 502}

		html, err := p.Render(data)
		assert.NoError(t, err)

		assert.Equal(t, "Bad Gateway - 502.\n", html)
	})
}
