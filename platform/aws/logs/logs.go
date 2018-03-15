// Package logs provides log management for AWS platforms.
package logs

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"

	"github.com/apex/log"
	jsonlog "github.com/apex/log/handlers/json"
	"github.com/apex/up"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/tj/aws/logs"

	"github.com/apex/up/internal/logs/parser"
	"github.com/apex/up/internal/logs/text"
	"github.com/apex/up/internal/util"
)

// Logs implementation.
type Logs struct {
	up.LogsConfig
	group string
	query string
	w     io.WriteCloser
	io.Reader
}

// New log tailer.
func New(group string, c up.LogsConfig) up.Logs {
	r, w := io.Pipe()

	query, err := parseQuery(c.Query)
	if err != nil {
		w.CloseWithError(err)
	}
	log.Debugf("query %q", query)

	l := &Logs{
		LogsConfig: c,
		query:      query,
		group:      group,
		Reader:     r,
		w:          w,
	}

	go l.start()

	return l
}

// start fetching logs.
func (l *Logs) start() {
	tailer := logs.New(logs.Config{
		Service:       cloudwatchlogs.New(session.New(aws.NewConfig().WithRegion(l.Region))),
		StartTime:     l.Since,
		PollInterval:  2 * time.Second,
		Follow:        l.Follow,
		FilterPattern: l.query,
		GroupNames:    []string{l.group},
	})

	var handler log.Handler

	if l.OutputJSON {
		handler = jsonlog.New(os.Stdout)
	} else {
		handler = text.New(os.Stdout).WithExpandedFields(l.Expand)
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

		// skip START / END logs since they are redundant
		if skippable(l.Message) {
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

// skippable returns true if the message is skippable.
func skippable(s string) bool {
	return strings.Contains(s, "END RequestId") ||
		strings.Contains(s, "START RequestId")
}
