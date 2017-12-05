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

// New writer with the given log level.
func New(l log.Level, ctx log.Interface) *Writer {
	pr, pw := io.Pipe()

	w := &Writer{
		PipeWriter: pw,
		LineReader: linereader.New(pr),
		done:       make(chan struct{}),
	}

	lw := &logWriter{
		log:   ctx,
		level: l,
	}

	go func() {
		defer close(w.done)
		io.Copy(lw, w.LineReader)
	}()

	return w
}

// Writer is a writer which copies lines with
// indentation support to a logWriter.
type Writer struct {
	*io.PipeWriter
	*linereader.LineReader
	done chan struct{}
}

// Close implementation.
func (w *Writer) Close() error {
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
	if util.IsJSONLog(string(b)) {
		return w.writeJSON(b)
	}

	return w.writeText(b)
}

// writeJSON writes a json log, interpreting it as a log.Entry.
func (w *logWriter) writeJSON(b []byte) (int, error) {
	// TODO: make this less ugly in apex/log,
	// you should be able to write an arbitrary Entry.
	var e log.Entry

	if err := json.Unmarshal(b, &e); err != nil {
		return 0, errors.Wrap(err, "unmarshaling")
	}

	switch e.Level {
	case log.DebugLevel:
		w.log.WithFields(e.Fields).Debug(e.Message)
	case log.InfoLevel:
		w.log.WithFields(e.Fields).Info(e.Message)
	case log.WarnLevel:
		w.log.WithFields(e.Fields).Warn(e.Message)
	case log.ErrorLevel:
		w.log.WithFields(e.Fields).Error(e.Message)
	case log.FatalLevel:
		// TODO: FATAL without exit...
		w.log.WithFields(e.Fields).Error(e.Message)
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
