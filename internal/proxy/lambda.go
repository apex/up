// Package proxy provides API Gateway and Lambda interoperability.
package proxy

import (
	"encoding/json"
	"net/http"

	"github.com/apex/go-apex"
	"github.com/pkg/errors"
)

// NewHandler returns an apex.Handler.
func NewHandler(h http.Handler) apex.Handler {
	return apex.HandlerFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		e := new(Input)

		err := json.Unmarshal(event, e)
		if err != nil {
			return nil, errors.Wrap(err, "parsing proxy event")
		}

		req, err := NewRequest(e)
		if err != nil {
			return nil, errors.Wrap(err, "creating new request from event")
		}

		res := NewResponse()
		h.ServeHTTP(res, req)
		return res.End(), nil
	})
}
