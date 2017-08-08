// Package progressreader provides a reader progress bar.
package progressreader

import (
	"io"
	"os"

	"github.com/apex/up/internal/util"
	progress "github.com/tj/go-progress"
	"github.com/tj/go/term"
)

// reader wrapping a progress bar.
type reader struct {
	io.ReadCloser
	p       *progress.Bar
	written int
}

// Read implementation.
func (r *reader) Read(b []byte) (int, error) {
	n, err := r.ReadCloser.Read(b)
	r.written += n
	r.p.ValueInt(r.written)
	io.WriteString(os.Stdout, term.CenterLine(r.p.String()))
	return n, err
}

// New returns a progress bar reader.
func New(size int, r io.ReadCloser) io.ReadCloser {
	return &reader{
		ReadCloser: r,
		p:          util.NewProgressInt(size),
	}
}
