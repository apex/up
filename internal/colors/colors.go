// Package colors provides colors used by the CLI.
package colors

import (
	color "github.com/aybabtme/rgbterm"
)

// Func is a color function.
type Func func(string) string

// Gray string.
func Gray(s string) string {
	return color.FgString(s, 150, 150, 150)
}

// Blue string.
func Blue(s string) string {
	return color.FgString(s, 77, 173, 247)
}

// Cyan string.
func Cyan(s string) string {
	return color.FgString(s, 34, 184, 207)
}

// Green string.
func Green(s string) string {
	return color.FgString(s, 0, 200, 255)
}

// Red string.
func Red(s string) string {
	return color.FgString(s, 194, 37, 92)
}

// Yellow string.
func Yellow(s string) string {
	return color.FgString(s, 252, 196, 25)
}

// Purple string.
func Purple(s string) string {
	return color.FgString(s, 96, 97, 190)
}

// Bool returns a color func based on the state.
func Bool(ok bool) Func {
	if ok {
		return Purple
	}

	return Red
}
