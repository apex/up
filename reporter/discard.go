package reporter

import "github.com/apex/up/platform/event"

// Discard events.
func Discard(events <-chan *event.Event) {
	for range events {
		// :)
	}
}
