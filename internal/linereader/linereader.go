// Package linereader reads lines with indentation suppoort.
package linereader

import (
	"bufio"
	"io"
	"strings"
	"time"
)

// LineReader is a line reader with indentation support.
type LineReader struct {
	io.ReadCloser
	lines chan string
	logs  chan string
	flush chan struct{}
	done  chan struct{}
}

// Flush the logs.
func (l *LineReader) Flush() {
	l.flush <- struct{}{}
}

// New line reader.
func New(r io.Reader) *LineReader {
	pr, pw := io.Pipe()

	lr := &LineReader{
		ReadCloser: pr,
		lines:      make(chan string),
		logs:       make(chan string),
		done:       make(chan struct{}),
		flush:      make(chan struct{}),
	}

	go lr.read(r, pw)

	return lr
}

// read from r and write lines to w.
func (l *LineReader) read(r io.Reader, w *io.PipeWriter) {
	var buf []string

	// flush the buffer as a distinct log message
	flush := func() {
		if len(buf) == 0 {
			return
		}

		l.logs <- strings.Join(buf, "\n")
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
		for s := range l.logs {
			w.Write([]byte(s))
		}
		close(l.done)
	}()

	// buffer lines with indentation support, flushing
	// periodically as look-ahead is required for indents
	go func() {
		for {
			select {
			case <-time.After(500 * time.Millisecond):
				flush()
			case <-l.flush:
				flush()
			case s, ok := <-l.lines:
				if ok {
					buffer(s)
				} else {
					flush()
					close(l.logs)
					return
				}
			}
		}
	}()

	// scan lines from the reader
	s := bufio.NewScanner(r)
	for s.Scan() {
		l.lines <- s.Text()
	}
	close(l.lines)

	// wait for final writes
	<-l.done
	w.CloseWithError(s.Err())
}

// indented returns true if the line is indented.
func indented(s string) bool {
	return strings.HasPrefix(s, "  ") || strings.HasPrefix(s, "\t")
}
