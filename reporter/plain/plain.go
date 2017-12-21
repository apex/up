// Package plain provides plain-text reporting for CI.
package plain

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/apex/up/platform/event"
)

// Report on events.
func Report(events <-chan *event.Event) {
	r := reporter{
		events: events,
	}

	r.Start()
}

// reporter struct.
type reporter struct {
	events <-chan *event.Event
}

// complete log with duration.
func (r *reporter) complete(name, value string, d time.Duration) {
	duration := fmt.Sprintf("(%s)", d.Round(time.Millisecond))
	fmt.Printf("     %s %s %s\n", name+":", value, duration)
}

// log line.
func (r *reporter) log(name, value string) {
	fmt.Printf("     %s %s\n", name+":", value)
}

// error line.
func (r *reporter) error(name, value string) {
	fmt.Printf("     %s %s\n", name+":", value)
}

// Start handling events.
func (r *reporter) Start() {
	for e := range r.events {
		switch e.Name {
		case "account.login.verify":
			r.log("verify", "Check your email for a confirmation link")
		case "account.login.verified":
			r.log("verify", "complete")
		case "hook":
			r.log("hook", e.String("name"))
		case "hook.complete":
			r.complete("hook", e.String("name"), e.Duration("duration"))
		case "platform.build.zip":
			s := fmt.Sprintf("%s files, %s", humanize.Comma(e.Int64("files")), humanize.Bytes(uint64(e.Int("size_compressed"))))
			r.complete("build", s, e.Duration("duration"))
		case "platform.deploy.complete":
			s := "complete"
			if v := e.String("version"); v != "" {
				s = "version " + v
			}
			r.complete("deploy", s, e.Duration("duration"))
		}
	}
}
