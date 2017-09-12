// Package writer provides an io.Writer for capturing
// process output as logs, so that stdout may become
// INFO, and stderr ERROR.
package writer

import (
	"io"

	"github.com/apex/log"
	"github.com/apex/up/internal/linereader"
)

// TODO: json support?
// TODO: logfmt support?
// TODO: rename "app" field

// New writer with the given log level.
func New(l log.Level) io.WriteCloser {
	pr, pw := io.Pipe()

	w := &writer{
		PipeWriter: pw,
		done:       make(chan bool),
	}

	lw := &logWriter{
		log:   log.WithField("app", true),
		level: l,
	}

	go func() {
		defer close(w.done)
		io.Copy(lw, linereader.New(pr))
	}()

	return w
}

// writer is a writer which copies lines with
// indentation support to a logWriter.
type writer struct {
	*io.PipeWriter
	done chan bool
}

// Close
func (w *writer) Close() error {
	if err := w.PipeWriter.Close(); err != nil {
		return err
	}

	<-w.done
	return nil
}

// logWriter is a write which logs distinct writes as log lines.
type logWriter struct {
	log   log.Interface
	level log.Level
}

// Write implementation.
func (w *logWriter) Write(b []byte) (int, error) {
	switch w.level {
	case log.InfoLevel:
		w.log.Info(string(b))
	case log.ErrorLevel:
		w.log.Error(string(b))
	}

	return len(b), nil
}
