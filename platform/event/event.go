package event

import (
	"fmt"
	"strings"
	"time"

	"github.com/apex/log"
)

// Events channel.
type Events chan *Event

// Emit an event.
func (e Events) Emit(name string, fields Fields) {
	if !strings.Contains(name, ".event") {
		log.Debugf("event %s %v", name, fields)
	}

	e <- &Event{
		Name:   name,
		Fields: fields,
	}
}

// Time an event.
func (e Events) Time(name string, fields Fields) func() {
	start := time.Now()

	e.Emit(name, fields)

	return func() {
		if fields == nil {
			fields = make(Fields)
		}

		f := make(Fields)

		for k, v := range fields {
			f[k] = v
		}

		f["duration"] = time.Since(start)
		e.Emit(name+".complete", f)
	}
}

// Fields for an event.
type Fields map[string]interface{}

// Event is a representation of an operation performed
// by a platform, and is used for reporting.
type Event struct {
	Name   string
	Fields Fields
}

// String value.
func (e *Event) String(name string) string {
	v, ok := e.Fields[name].(string)
	if !ok {
		panic(fmt.Errorf("%#v field %s is not a string", e, name))
	}
	return v
}

// Duration value.
func (e *Event) Duration(name string) time.Duration {
	v, ok := e.Fields[name].(time.Duration)
	if !ok {
		panic(fmt.Errorf("%#v field %s is not a time.Duration", e, name))
	}
	return v
}

// Int64 value.
func (e *Event) Int64(name string) int64 {
	v, ok := e.Fields[name].(int64)
	if !ok {
		panic(fmt.Errorf("%#v field %s is not a int64", e, name))
	}
	return v
}

// Int value.
func (e *Event) Int(name string) int {
	v, ok := e.Fields[name].(int)
	if !ok {
		panic(fmt.Errorf("%#v field %s is not a int", e, name))
	}
	return v
}
