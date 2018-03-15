// Package text implements a development-friendly textual handler.
package text

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/dustin/go-humanize"

	"github.com/apex/up/internal/colors"
	"github.com/apex/up/internal/util"
)

var (
	spacerPlaceholderBytes = []byte("{{spacer}}")
	spacerBytes            = []byte(colors.Gray(":"))
	newlineBytes           = []byte("\n")
	emptyBytes             = []byte("")
)

// color function.
type colorFunc func(string) string

// omit fields.
var omit = map[string]bool{
	"app":     true,
	"stage":   true,
	"region":  true,
	"plugin":  true,
	"commit":  true,
	"version": true,
}

// Colors mapping.
var Colors = [...]colorFunc{
	log.DebugLevel: colors.Gray,
	log.InfoLevel:  colors.Purple,
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
	fmt.Fprintf(h.Writer, "  %s %s %s\n", colors.Gray(ts), bold(color(level)), colors.Purple(e.Message))

	for _, name := range names {
		v := e.Fields.Get(name)

		if v == "" {
			continue
		}

		fmt.Fprintf(h.Writer, "    %s%s%v\n", color(name), colors.Gray(": "), value(name, v))
	}

	if len(names) > 0 {
		fmt.Fprintf(h.Writer, "\n")
	}

	return nil
}

// handleInline fields.
func (h *Handler) handleInline(e *log.Entry) error {
	var buf bytes.Buffer
	var fields int

	color := Colors[e.Level]
	level := Strings[e.Level]
	names := e.Fields.Names()
	ts := formatDate(e.Timestamp.Local())

	if stage, ok := e.Fields.Get("stage").(string); ok && stage != "" {
		fmt.Fprintf(&buf, "  %s %s %s %s %s{{spacer}}", colors.Gray(ts), bold(color(level)), colors.Gray(stage), colors.Gray(version(e)), colors.Purple(e.Message))
	} else {
		fmt.Fprintf(&buf, "  %s %s %s{{spacer}}", colors.Gray(ts), bold(color(level)), colors.Purple(e.Message))
	}

	for _, name := range names {
		if omit[name] {
			continue
		}

		v := e.Fields.Get(name)

		if v == "" {
			continue
		}

		fields++
		fmt.Fprintf(&buf, " %s%s%v", color(name), colors.Gray("="), value(name, v))
	}

	b := buf.Bytes()

	if fields > 0 {
		b = bytes.Replace(b, spacerPlaceholderBytes, spacerBytes, 1)
	} else {
		b = bytes.Replace(b, spacerPlaceholderBytes, emptyBytes, 1)
	}

	h.mu.Lock()
	h.Writer.Write(b)
	h.Writer.Write(newlineBytes)
	h.mu.Unlock()

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
	return t.Format(`Jan 2` + util.DateSuffix(t) + ` 03:04:05pm`)
}

// version returns the entry version via GIT commit or lambda version.
func version(e *log.Entry) string {
	if s, ok := e.Fields.Get("commit").(string); ok && s != "" {
		return s
	}

	if s, ok := e.Fields.Get("version").(string); ok && s != "" {
		return s
	}

	return ""
}

// bold string.
func bold(s string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", s)
}
