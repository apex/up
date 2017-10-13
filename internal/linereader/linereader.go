// Package linereader reads lines with indentation suppoort.
package linereader

import (
	"bufio"
	"io"
	"strings"
	"time"
)

// New line reader.
func New(r io.Reader) io.ReadCloser {
	pr, pw := io.Pipe()
	go read(r, pw)
	return pr
}

// read from r and write lines to w.
func read(r io.Reader, w *io.PipeWriter) {
	lines := make(chan string)
	logs := make(chan string)
	done := make(chan struct{})
	var buf []string

	// flush the buffer as a distinct log message
	flush := func() {
		if len(buf) == 0 {
			return
		}

		logs <- strings.Join(buf, "\n")
		buf = nil
	}

	// buffer line with indentation look-ahead
	buffer := func(s string) {
		if !indented(s) {
			flush()
		}
		buf = append(buf, s)
	}

	// copy logs to the writer as distinct Write calls
	go func() {
		for l := range logs {
			w.Write([]byte(l))
		}
		close(done)
	}()

	// buffer lines with indentation support, flushing
	// periodically as look-ahead is required for indents
	go func() {
		t := time.NewTicker(500 * time.Millisecond)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				flush()
			case l, ok := <-lines:
				if ok {
					buffer(l)
				} else {
					flush()
					close(logs)
					return
				}
			}
		}
	}()

	// scan lines from the reader
	s := bufio.NewScanner(r)
	for s.Scan() {
		lines <- s.Text()
	}
	close(lines)

	// wait for final writes
	<-done
	w.CloseWithError(s.Err())
}

// indented returns true if the line is indented.
func indented(s string) bool {
	return strings.HasPrefix(s, "  ") || strings.HasPrefix(s, "\t")
}
