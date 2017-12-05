package lambda

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"

	"github.com/apex/log"
	jsonlog "github.com/apex/log/handlers/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/tj/go/term"

	"github.com/apex/up/internal/logs/parser"
	"github.com/apex/up/internal/logs/text"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/platform"
	"github.com/apex/up/platform/lambda/logs"
)

// TODO: move formatting logic outside of platform, reader interface
// TODO: optionally expand fields

// Logs implementation.
type Logs struct {
	platform *Platform
	region   string
	query    string
	follow   bool
	expand   bool
	since    time.Time
	w        io.WriteCloser
	io.Reader
}

// NewLogs returns a new logs tailer.
func NewLogs(p *Platform, region, query string) platform.Logs {
	r, w := io.Pipe()

	query, err := parseQuery(query)
	if err != nil {
		w.CloseWithError(err)
	}
	log.Debugf("query %q", query)

	l := &Logs{
		platform: p,
		region:   region,
		query:    query,
		Reader:   r,
		w:        w,
	}

	go l.start()

	return l
}

// Since implementation.
func (l *Logs) Since(t time.Time) {
	l.since = t
}

// Follow implementation.
func (l *Logs) Follow() {
	log.Debug("follow")
	l.follow = true
}

// Expand implementation.
func (l *Logs) Expand() {
	log.Debug("expand")
	l.expand = true
}

// start fetching logs.
func (l *Logs) start() {
	// TODO: flag to override and allow querying other groups
	// TODO: apply backoff instead of PollInterval
	group := "/aws/lambda/" + l.platform.config.Name

	config := logs.Config{
		Service:       cloudwatchlogs.New(session.New(aws.NewConfig().WithRegion(l.region))),
		StartTime:     l.since,
		PollInterval:  2 * time.Second,
		Follow:        l.follow,
		FilterPattern: l.query,
	}

	tailer := &logs.Logs{
		Config:     config,
		GroupNames: []string{group},
	}

	// TODO: delegate isatty stuff...

	var handler log.Handler

	if term.IsTerminal() {
		handler = text.New(os.Stdout).WithExpandedFields(l.expand)
	} else {
		handler = jsonlog.New(os.Stdout)
	}

	// TODO: transform to reader of nl-delimited json, move to apex/log?
	// TODO: marshal/unmarshal as JSON so that numeric values are always float64... remove util.ToFloat()
	for l := range tailer.Start() {
		line := strings.TrimSpace(l.Message)

		// json log
		if util.IsJSONLog(line) {
			var e log.Entry
			err := json.Unmarshal([]byte(line), &e)
			if err != nil {
				log.Fatalf("error parsing json: %s", err)
			}

			handler.HandleLog(&e)
			continue
		}

		// lambda textual logs
		handler.HandleLog(&log.Entry{
			Timestamp: l.Timestamp,
			Level:     log.InfoLevel,
			Message:   strings.TrimRight(l.Message, " \n"),
		})
	}

	// TODO: refactor interface to delegate
	if err := tailer.Err(); err != nil {
		panic(err)
	}

	l.w.Close()
}

// parseQuery parses and converts the query to a CW friendly syntax.
func parseQuery(s string) (string, error) {
	if s == "" {
		return s, nil
	}

	n, err := parser.Parse(s)
	if err != nil {
		return "", err
	}

	return n.String(), nil
}
