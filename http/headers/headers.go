// Package headers provides header injection support.
package headers

import (
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/pkg/errors"
	hdr "github.com/tj/go-headers"

	"github.com/apex/up"
	"github.com/apex/up/internal/header"
)

// TODO: document precedence and/or add options
// TODO: maybe allow storing _headers in Static.Dir?

// filename of headers file.
var filename = "_headers"

// Headers handler injects headers.
type Headers struct {
	next  http.Handler
	rules *header.Matcher
}

// New headers handler.
func New(c *up.Config, next http.Handler) (http.Handler, error) {
	rulesFromFile, err := readFromFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "reading header file")
	}

	rules, err := header.Compile(header.Merge(rulesFromFile, c.Headers))
	if err != nil {
		return nil, errors.Wrap(err, "compiling header")
	}

	log.Debugf("header rules from _headers file: %d", len(rulesFromFile))
	log.Debugf("header rules from up.json: %d", len(c.Headers))

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fields := rules.Lookup(r.URL.Path)

		for k, v := range fields {
			w.Header().Set(k, v)
		}

		next.ServeHTTP(w, r)
	})

	return h, nil
}

// readFromFile reads from a Netlify style headers file.
func readFromFile(path string) (header.Rules, error) {
	rules := make(header.Rules)

	f, err := os.Open(path)

	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "opening headers file")
	}

	defer f.Close()

	h, err := hdr.Parse(f)
	if err != nil {
		return nil, errors.Wrap(err, "parsing")
	}

	for path, fields := range h {
		rules[path] = make(header.Fields)
		for name, vals := range fields {
			rules[path][name] = vals[0]
		}
	}

	return rules, nil
}
