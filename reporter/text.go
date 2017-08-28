// Package reporter provides event-based reporting for the CLI,
// aka this is what the user sees.
package reporter

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/dustin/go-humanize"
	"github.com/tj/go-progress"
	"github.com/tj/go-spin"
	"github.com/tj/go/term"

	"github.com/apex/up/internal/colors"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/platform/event"
	"github.com/apex/up/platform/lambda/cost"
	"github.com/apex/up/platform/lambda/stack"
)

// TODO: platform-specific reporting should live in the platform
// TODO: typed events would be nicer.. refactor event names
// TODO: refactor, this is a hot mess :D

// Text outputs human friendly textual reporting.
func Text(events <-chan *event.Event) {
	r := reporter{
		events:  events,
		spinner: spin.New(),
	}

	r.Start()
}

// reporter struct.
type reporter struct {
	events         <-chan *event.Event
	spinner        *spin.Spinner
	prevTime       time.Time
	bar            *progress.Bar
	inlineProgress bool
	pendingName    string
	pendingValue   string
}

// spin the spinner by moving to the start of the line and re-printing.
func (r *reporter) spin() {
	if r.pendingName != "" {
		r.pending(r.pendingName, r.pendingValue)
	}
}

// pending log with spinner.
func (r *reporter) pending(name, value string) {
	r.pendingName = name
	r.pendingValue = value
	term.ClearLine()
	fmt.Printf("\r %30s %s", colors.Purple(r.spinner.Next()+" "+name+":"), value)
}

// complete log with duration.
func (r *reporter) complete(name, value string, d time.Duration) {
	r.pendingName = ""
	r.pendingValue = ""
	term.ClearLine()
	duration := fmt.Sprintf("(%s)", d.Round(time.Millisecond))
	fmt.Printf("\r %30s %s %s\n", colors.Purple(name+":"), value, colors.Gray(duration))
}

// log line
func (r *reporter) log(name, value string) {
	fmt.Printf("\r %35s %s\n", colors.Purple(name+":"), value)
}

// Start handling events.
func (r *reporter) Start() {
	tick := time.NewTicker(150 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			r.spin()
		case e := <-r.events:
			switch e.Name {
			case "hook":
				r.pending("hook", e.String("name"))
			case "hook.complete":
				r.complete("hook", e.String("name"), e.Duration("duration"))
			case "deploy", "stack.delete", "platform.stack.apply":
				term.HideCursor()
			case "deploy.complete", "stack.delete.complete", "platform.stack.apply.complete":
				term.ShowCursor()
			case "platform.build":
				r.pending("build", "")
			case "platform.build.zip":
				s := fmt.Sprintf("%s files, %s", humanize.Comma(e.Int64("files")), humanize.Bytes(uint64(e.Int("size_compressed"))))
				r.complete("build", s, e.Duration("duration"))
			case "platform.deploy":
				r.pending("deploy", "")
			case "platform.deploy.complete":
				r.complete("deploy", "complete", e.Duration("duration"))
			case "platform.function.create":
				r.inlineProgress = true
			case "stack.create":
				r.inlineProgress = true
			case "platform.stack.report":
				if r.inlineProgress {
					r.bar = util.NewInlineProgressInt(e.Int("total"))
					r.pending("stack", r.bar.String())
				} else {
					r.bar = util.NewProgressInt(e.Int("total"))
					io.WriteString(os.Stdout, term.CenterLine(r.bar.String()))
				}
			case "platform.stack.report.event":
				if r.inlineProgress {
					r.bar.ValueInt(e.Int("complete"))
					r.pending("stack", r.bar.String())
				} else {
					r.bar.ValueInt(e.Int("complete"))
					io.WriteString(os.Stdout, term.CenterLine(r.bar.String()))
				}
			case "platform.stack.report.complete":
				if r.inlineProgress {
					r.complete("stack", "complete", e.Duration("duration"))
				} else {
					term.ClearAll()
					term.ShowCursor()
				}
			case "platform.stack.show", "platform.stack.show.complete":
				fmt.Printf("\n")
			case "platform.stack.show.stack":
				s := e.Fields["stack"].(*cloudformation.Stack)
				fmt.Printf("  %s: %s\n", colors.Purple("status"), stack.Status(*s.StackStatus))
				if reason := s.StackStatusReason; reason != nil {
					fmt.Printf("  %s: %s\n", colors.Purple("reason"), *reason)
				}
				fmt.Printf("\n")
			case "platform.stack.show.event":
				event := e.Fields["event"].(*cloudformation.StackEvent)
				kind := *event.ResourceType
				status := stack.Status(*event.ResourceStatus)
				color := colors.Purple
				if status.State() == stack.Failure {
					color = colors.Red
				}
				fmt.Printf("  %s\n", color(kind))
				fmt.Printf("    %s: %v\n", color("id"), *event.LogicalResourceId)
				fmt.Printf("    %s: %s\n", color("status"), status)
				if reason := event.ResourceStatusReason; reason != nil {
					fmt.Printf("    %s: %s\n", color("reason"), *reason)
				}
				fmt.Printf("\n")
			case "stack.plan":
				fmt.Printf("\n")
			case "platform.stack.plan.change":
				c := e.Fields["change"].(*cloudformation.Change).ResourceChange
				color := actionColor(*c.Action)
				fmt.Printf("  %s %s\n", color(*c.Action), *c.ResourceType)
				fmt.Printf("    %s: %s\n", color("id"), *c.LogicalResourceId)
				if c.Replacement != nil {
					fmt.Printf("    %s: %s\n", color("replace"), *c.Replacement)
				}
				fmt.Printf("\n")
			case "metrics", "metrics.complete":
				fmt.Printf("\n")
			case "metrics.value":
				switch n := e.String("name"); n {
				case "Duration min", "Duration avg", "Duration max":
					r.log(n, fmt.Sprintf("%dms", e.Int("value")))
				case "Requests":
					v := humanize.Comma(int64(e.Int("value")))
					c := cost.Requests(e.Int("value"))
					r.log(n, fmt.Sprintf("%s %s", v, currency(c)))
				case "Duration sum":
					d := time.Millisecond * time.Duration(e.Int("value"))
					c := cost.Duration(e.Int("value"), e.Int("memory"))
					r.log(n, fmt.Sprintf("%s %s", d, currency(c)))
				case "Invocations":
					d := humanize.Comma(int64(e.Int("value")))
					c := cost.Invocations(e.Int("value"))
					r.log(n, fmt.Sprintf("%s %s", d, currency(c)))
				default:
					r.log(n, fmt.Sprintf("%s", humanize.Comma(int64(e.Int("value")))))
				}
			}

			r.prevTime = time.Now()
		}
	}
}

// currency format.
func currency(n float64) string {
	return colors.Gray(fmt.Sprintf("($%0.2f)", n))
}

// countEventsByStatus returns the number of events with the given state.
func countEventsByStatus(events []*cloudformation.StackEvent, desired stack.Status) (n int) {
	for _, e := range events {
		status := stack.Status(*e.ResourceStatus)

		if *e.ResourceType == "AWS::CloudFormation::Stack" {
			continue
		}

		if status == desired {
			n++
		}
	}

	return
}

// countEventsComplete returns the number of completed or failed events.
func countEventsComplete(events []*cloudformation.StackEvent) (n int) {
	for _, e := range events {
		status := stack.Status(*e.ResourceStatus)

		if *e.ResourceType == "AWS::CloudFormation::Stack" {
			continue
		}

		if status.IsDone() {
			n++
		}
	}

	return
}

// actionColor returns a color func by action.
func actionColor(s string) colors.Func {
	switch s {
	case "Add":
		return colors.Purple
	case "Remove":
		return colors.Red
	case "Modify":
		return colors.Blue
	default:
		return colors.Gray
	}
}
