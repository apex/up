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
		handler = text.New(os.Stdout)
	} else {
		handler = jsonlog.New(os.Stdout)
	}

	// TODO: transform to reader of nl-delimited json, move to apex/log?
	for l := range tailer.Start() {
		line := strings.TrimSpace(l.Message)

		if !util.IsJSON(line) {
			// fmt.Fprint(l.w, e.Message) // TODO: ignore? json-ify?
			continue
		}

		var e log.Entry
		err := json.Unmarshal([]byte(line), &e)
		if err != nil {
			log.Fatalf("error parsing json: %s", err)
		}

		handler.HandleLog(&e)
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
