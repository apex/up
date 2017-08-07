// Package text implements a development-friendly textual handler.
package text

import (
	"fmt"
	"io"
	"sync"

	"github.com/apex/log"

	"github.com/apex/up/internal/colors"
)

// TODO: rename since it's specific to log querying ATM
// TODO: output larger timestamp when older
// TODO: option to output UTC
// TODO: option to output expanded fields
// TODO: option to truncate

// color function.
type colorFunc func(string) string

// omit fields.
var omit = map[string]bool{
	"app":        true,
	"app_name":   true,
	"app_region": true,
	"plugin":     true,
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
	layout string
}

// New handler.
func New(w io.Writer) *Handler {
	return &Handler{
		Writer: w,
		layout: "3:04:05pm",
	}
}

// WithFormat sets the date format.
func (h *Handler) WithFormat(s string) *Handler {
	h.layout = s
	return h
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	color := Colors[e.Level]
	level := Strings[e.Level]
	names := e.Fields.Names()

	h.mu.Lock()
	defer h.mu.Unlock()

	ts := e.Timestamp.Local().Format(h.layout)
	fmt.Fprintf(h.Writer, "  %s %s %s", colors.Gray(ts), color(level), e.Message)

	for _, name := range names {
		if omit[name] {
			continue
		}

		v := e.Fields.Get(name)

		if v == "" {
			continue
		}

		fmt.Fprintf(h.Writer, " %s%s%v", color(name), colors.Gray(": "), v)
	}

	fmt.Fprintln(h.Writer)

	return nil
}
