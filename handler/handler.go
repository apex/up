// Package handler provides what is essentially the core of Up's
// reverse proxy, complete with all middleware for handling
// logging, redirectcs, static file serving and so on.
package handler

import (
	"net/http"

	"github.com/apex/log"
	"github.com/pkg/errors"

	"github.com/apex/up"
	"github.com/apex/up/http/cors"
	"github.com/apex/up/http/errorpages"
	"github.com/apex/up/http/headers"
	"github.com/apex/up/http/inject"
	"github.com/apex/up/http/logs"
	"github.com/apex/up/http/poweredby"
	"github.com/apex/up/http/redirects"
	"github.com/apex/up/http/relay"
	"github.com/apex/up/http/static"
)

// New reads up.json to configure and initialize
// the http handler chain for serving an Up application.
func New() (http.Handler, error) {
	c, err := up.ReadConfig("up.json")
	if err != nil {
		return nil, errors.Wrap(err, "reading config")
	}

	log.WithFields(log.Fields{
		"name": c.Name,
		"type": c.Type,
	}).Info("starting")

	var h http.Handler

	switch c.Type {
	case "server":
		h, err = relay.New(c)
		if err != nil {
			return nil, errors.Wrap(err, "initializing relay")
		}
	case "static":
		h = static.New(c)
	}

	h = poweredby.New("up", h)

	h, err = headers.New(c, h)
	if err != nil {
		return nil, errors.Wrap(err, "initializing headers")
	}

	h, err = errorpages.New(c, h)
	if err != nil {
		return nil, errors.Wrap(err, "initializing error pages")
	}

	h, err = inject.New(c, h)
	if err != nil {
		return nil, errors.Wrap(err, "initializing inject")
	}

	h = cors.New(c, h)

	h, err = redirects.New(c, h)
	if err != nil {
		return nil, errors.Wrap(err, "initializing redirects")
	}

	h, err = logs.New(c, h)
	if err != nil {
		return nil, errors.Wrap(err, "initializing logs")
	}

	return h, nil
}
