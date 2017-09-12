// Package writer provides an io.Writer for capturing
// process output as logs, so that stdout may become
// INFO, and stderr ERROR.
package writer

import (
	"encoding/json"
	"io"

	"github.com/apex/log"
	"github.com/apex/up/internal/linereader"
	"github.com/apex/up/internal/util"
	"github.com/pkg/errors"
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

// Close implementation.
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
	if util.IsJSON(string(b)) {
		return w.writeJSON(b)
	}

	return w.writeText(b)
}

// writeJSON writes a json log,
// interpreting it as a log.Entry.
func (w *logWriter) writeJSON(b []byte) (int, error) {
	var m map[string]interface{}

	// parse
	if err := json.Unmarshal(b, &m); err != nil {
		return 0, errors.Wrap(err, "unmarshaling")
	}

	// see if it looks like a log
	level, ok := m["level"].(string)
	if !ok {
		return w.writeText(b)
	}

	// parse log level
	lvl, err := log.ParseLevel(level)
	if err != nil {
		return 0, errors.Wrap(err, "parsing level")
	}

	// message
	var msg string
	if s, ok := m["message"].(string); ok {
		msg = s
	}

	// reserved fields
	delete(m, "level")
	delete(m, "timestamp")
	delete(m, "message")
	delete(m, "fields")

	// fields
	fields := log.Fields{}
	for k, v := range m {
		fields[k] = v
	}

	ctx := w.log.WithFields(fields)

	// log
	switch lvl {
	case log.DebugLevel:
		ctx.Debug(msg)
	case log.InfoLevel:
		ctx.Info(msg)
	case log.WarnLevel:
		ctx.Warn(msg)
	case log.ErrorLevel:
		ctx.Error(msg)
	case log.FatalLevel:
		ctx.Fatal(msg)
	}

	return len(b), nil
}

// writeText writes plain text.
func (w *logWriter) writeText(b []byte) (int, error) {
	switch w.level {
	case log.InfoLevel:
		w.log.Info(string(b))
	case log.ErrorLevel:
		w.log.Error(string(b))
	}

	return len(b), nil
}
