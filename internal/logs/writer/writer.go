// Package writer provides an io.Writer for capturing
// process output as logs, so that stdout may become
// INFO, and stderr ERROR.
package writer

import (
	"io"
	"strings"

	"github.com/apex/log"
)

// TODO: capture indented output as single log call
// TODO: json support?

// writer struct.
type writer struct {
	Level log.Level
	Log   log.Interface
}

// Write implementation.
func (w *writer) Write(b []byte) (int, error) {
	switch w.Level {
	case log.InfoLevel:
		w.Log.Info(strings.TrimSpace(string(b)))
	case log.ErrorLevel:
		w.Log.Error(strings.TrimSpace(string(b)))
	}

	return len(b), nil
}

// New writer with the given log level.
func New(l log.Level) io.Writer {
	return &writer{
		Level: l,
		Log:   log.WithField("app", true), // TODO: rename?
	}
}
