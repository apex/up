// Package text implements a development-friendly textual handler.
package text

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/dustin/go-humanize"

	"github.com/apex/up/internal/colors"
	"github.com/apex/up/internal/util"
)

// TODO: rename since it's specific to log querying ATM
// TODO: output larger timestamp when older
// TODO: option to output UTC
// TODO: option to truncate
// TODO: move to apex/log?

// color function.
type colorFunc func(string) string

// omit fields.
var omit = map[string]bool{
	"app":    true,
	"region": true,
	"plugin": true,
}

// Colors mapping.
var Colors = [...]colorFunc{
	log.DebugLevel: colors.Gray,
	log.InfoLevel:  colors.Blue,
	log.WarnLevel:  colors.Yellow,
	log.ErrorLevel: colors.Red,
	log.FatalLevel: colors.Red,
}

// Strings mapping.
var Strings = [...]string{
	log.DebugLevel: "DEBU",
	log.InfoLevel:  "INFO",
	log.WarnLevel:  "WARN",
	log.ErrorLevel: "ERRO",
	log.FatalLevel: "FATA",
}

// Handler implementation.
type Handler struct {
	mu     sync.Mutex
	Writer io.Writer
	expand bool
	layout string
}

// New handler.
func New(w io.Writer) *Handler {
	return &Handler{
		Writer: w,
	}
}

// WithExpandedFields sets the expanded field state.
func (h *Handler) WithExpandedFields(v bool) *Handler {
	h.expand = v
	return h
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	switch {
	case h.expand:
		return h.handleExpanded(e)
	default:
		return h.handleInline(e)
	}
}

// handleExpanded fields.
func (h *Handler) handleExpanded(e *log.Entry) error {
	color := Colors[e.Level]
	level := Strings[e.Level]
	names := e.Fields.Names()

	h.mu.Lock()
	defer h.mu.Unlock()

	ts := formatDate(e.Timestamp.Local())
	fmt.Fprintf(h.Writer, "  %s %s %s\n", colors.Gray(ts), color(level), e.Message)

	for _, name := range names {
		if omit[name] {
			continue
		}

		v := e.Fields.Get(name)

		if v == "" {
			continue
		}

		fmt.Fprintf(h.Writer, "  %30s%s%v\n", color(name), colors.Gray(": "), value(name, v))
	}

	return nil
}

// handleInline fields.
func (h *Handler) handleInline(e *log.Entry) error {
	color := Colors[e.Level]
	level := Strings[e.Level]
	names := e.Fields.Names()

	h.mu.Lock()
	defer h.mu.Unlock()

	ts := formatDate(e.Timestamp.Local())
	fmt.Fprintf(h.Writer, "  %s %s %s", colors.Gray(ts), color(level), e.Message)

	for _, name := range names {
		if omit[name] {
			continue
		}

		v := e.Fields.Get(name)

		if v == "" {
			continue
		}

		fmt.Fprintf(h.Writer, " %s%s%v", color(name), colors.Gray(": "), value(name, v))
	}

	fmt.Fprintln(h.Writer)

	return nil
}

// value returns the formatted value.
func value(name string, v interface{}) interface{} {
	switch name {
	case "size":
		return humanize.Bytes(uint64(util.ToFloat(v)))
	case "duration":
		return time.Millisecond * time.Duration(util.ToFloat(v))
	default:
		return v
	}
}

// day duration.
var day = time.Hour * 24

// formatDate formats t relative to now.
func formatDate(t time.Time) string {
	switch d := time.Now().Sub(t); {
	case d >= day*7:
		return t.Format(`Jan 2` + dateSuffix(t) + ` 3:04:05pm`)
	case d >= day:
		return t.Format(`2` + dateSuffix(t) + ` 3:04:05pm`)
	default:
		return t.Format(`3:04:05pm`)
	}
}

// dateSuffix returns the date suffix for t.
func dateSuffix(t time.Time) string {
	switch t.Day() {
	case 1, 21, 31:
		return "st"
	case 2, 22:
		return "nd"
	case 3, 23:
		return "rd"
	default:
		return "th"
	}
}
