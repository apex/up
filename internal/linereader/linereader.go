// Package linereader reads lines with indentation suppoort.
package linereader

import (
	"bufio"
	"io"
	"strings"
)

// New line reader.
func New(r io.Reader) io.ReadCloser {
	pr, pw := io.Pipe()
	go read(r, pw)
	return pr
}

// read from r and write lines to w.
func read(r io.Reader, w *io.PipeWriter) {
	s := bufio.NewScanner(r)
	var lines []string

	flush := func() {
		io.WriteString(w, strings.Join(lines, "\n"))
		lines = nil
	}

	for s.Scan() {
		line := s.Text()
		if !indented(line) {
			if len(lines) > 0 {
				flush()
			}
		}
		lines = append(lines, line)
	}

	if len(lines) > 0 {
		flush()
	}

	w.CloseWithError(s.Err())
}

// indented returns true if the line is indented.
func indented(s string) bool {
	return strings.HasPrefix(s, "  ") || strings.HasPrefix(s, "\t")
}
