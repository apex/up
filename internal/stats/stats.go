// Package stats provides CLI analytics.
package stats

import (
	"github.com/apex/log"
	"github.com/tj/go-cli-analytics"
)

// p merged with track calls.
var p = map[string]interface{}{}

// Client for Segment analytics.
var Client = analytics.New(&analytics.Config{
	WriteKey: "qnvYCHktBBgACBkQ6V4dzh7aFCe8LF8u",
	Dir:      ".up",
})

// Track event `name` with optional `props`.
func Track(name string, props map[string]interface{}) {
	if props == nil {
		props = map[string]interface{}{}
	}

	for k, v := range p {
		props[k] = v
	}

	log.Debugf("track %q %v", name, props)
	Client.Track(name, props)
}

// SetProperties sets global properties.
func SetProperties(props map[string]interface{}) {
	p = props
}
