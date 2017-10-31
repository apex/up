// Package progressreader provides an io.Reader progress bar.
package progressreader

import (
	"io"
	"sync"

	"github.com/apex/up/internal/util"
	"github.com/tj/go-progress"
	"github.com/tj/go/term"
)

// reader wrapping a progress bar.
type reader struct {
	io.ReadCloser
	p       *progress.Bar
	render  func(string)
	written int
	sync.Once
}

// Read implementation.
func (r *reader) Read(b []byte) (int, error) {
	r.Do(term.ClearAll)
	n, err := r.ReadCloser.Read(b)
	r.written += n
	r.p.ValueInt(r.written)
	r.render(term.CenterLine(r.p.String()))
	return n, err
}

// New returns a progress bar reader.
func New(size int, r io.ReadCloser) io.ReadCloser {
	return &reader{
		ReadCloser: r,
		p:          util.NewProgressInt(size),
		render:     term.Renderer(),
	}
}
