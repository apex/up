// Package discard provides a reporter for discarding events.
package discard

import "github.com/apex/up/platform/event"

// Report events.
func Report(events <-chan *event.Event) {
	for range events {
		// :)
	}
}
