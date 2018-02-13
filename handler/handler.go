// Package handler provides what is essentially the core of Up's
// reverse proxy, complete with all middleware for handling
// logging, redirectcs, static file serving and so on.
package handler

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/apex/up"
	"github.com/apex/up/http/cors"
	"github.com/apex/up/http/errorpages"
	"github.com/apex/up/http/gzip"
	"github.com/apex/up/http/headers"
	"github.com/apex/up/http/inject"
	"github.com/apex/up/http/logs"
	"github.com/apex/up/http/ping"
	"github.com/apex/up/http/poweredby"
	"github.com/apex/up/http/redirects"
	"github.com/apex/up/http/relay"
	"github.com/apex/up/http/static"
)

// FromConfig returns the handler based on user config.
func FromConfig(c *up.Config) (http.Handler, error) {
	switch c.Type {
	case "server":
		return relay.New(c)
	case "static":
		return static.New(c), nil
	default:
		return nil, errors.Errorf("unknown .type %q", c.Type)
	}
}

// New handler complete with all Up middleware.
func New(c *up.Config, h http.Handler) (http.Handler, error) {
	h = poweredby.New("up", h)

	h, err := headers.New(c, h)
	if err != nil {
		return nil, errors.Wrap(err, "headers")
	}

	h, err = errorpages.New(c, h)
	if err != nil {
		return nil, errors.Wrap(err, "error pages")
	}

	h, err = inject.New(c, h)
	if err != nil {
		return nil, errors.Wrap(err, "inject")
	}

	h = cors.New(c, h)

	h, err = redirects.New(c, h)
	if err != nil {
		return nil, errors.Wrap(err, "redirects")
	}

	h = gzip.New(c, h)

	h, err = logs.New(c, h)
	if err != nil {
		return nil, errors.Wrap(err, "logs")
	}

	h = ping.New(c, h)

	return h, nil
}
