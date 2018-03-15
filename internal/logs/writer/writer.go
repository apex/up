// Package writer provides an io.Writer for capturing
// process output as logs, so that stdout may become
// INFO, and stderr ERROR.
package writer

import (
	"bufio"
	"bytes"
	"encoding/json"

	"github.com/apex/log"
	"github.com/apex/up/internal/util"
	"github.com/pkg/errors"
)

// Writer struct.
type Writer struct {
	log   log.Interface
	level log.Level
}

// New writer with the given log level.
func New(l log.Level, ctx log.Interface) *Writer {
	return &Writer{
		log:   ctx,
		level: l,
	}
}

// Write implementation.
func (w *Writer) Write(b []byte) (int, error) {
	s := bufio.NewScanner(bytes.NewReader(b))

	for s.Scan() {
		if err := w.write(s.Text()); err != nil {
			return 0, err
		}
	}

	if err := s.Err(); err != nil {
		return 0, err
	}

	return len(b), nil
}

// write the line.
func (w *Writer) write(s string) error {
	if util.IsJSONLog(s) {
		return w.writeJSON(s)
	}

	return w.writeText(s)
}

// writeJSON writes a json log, interpreting it as a log.Entry.
func (w *Writer) writeJSON(s string) error {
	// TODO: make this less ugly in apex/log,
	// you should be able to write an arbitrary Entry.
	var e log.Entry

	if err := json.Unmarshal([]byte(s), &e); err != nil {
		return errors.Wrap(err, "unmarshaling")
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

	return nil
}

// writeText writes plain text.
func (w *Writer) writeText(s string) error {
	switch w.level {
	case log.InfoLevel:
		w.log.Info(s)
	case log.ErrorLevel:
		w.log.Error(s)
	}
	return nil
}
