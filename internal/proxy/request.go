package proxy

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// NewRequest returns a new http.Request from the given Lambda event.
func NewRequest(e *Input) (*http.Request, error) {
	// path
	u, err := url.Parse(e.Path)
	if err != nil {
		return nil, errors.Wrap(err, "parsing path")
	}

	// querystring
	q := u.Query()
	for k, v := range e.QueryStringParameters {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	// base64 encoded body
	body := e.Body
	if e.IsBase64Encoded {
		b, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return nil, errors.Wrap(err, "decoding base64 body")
		}
		body = string(b)
	}

	getBody := func() (io.ReadCloser, error) {
		return ioutil.NopCloser(strings.NewReader(body)), nil
	}

	// even though this can't fail still handle the error to defend against future changes
	bodyReader, err := getBody()
	if err != nil {
		return nil, err
	}

	// new request
	req, err := http.NewRequest(e.HTTPMethod, u.String(), bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	req.GetBody = getBody

	// remote addr
	req.RemoteAddr = e.RequestContext.Identity.SourceIP

	// header fields
	for k, v := range e.Headers {
		req.Header.Set(k, v)
	}

	// custom fields
	req.Header.Set("X-Request-Id", e.RequestContext.RequestID)
	req.Header.Set("X-Stage", e.RequestContext.Stage)

	// host
	req.URL.Host = req.Header.Get("Host")
	req.Host = req.URL.Host

	return req, nil
}
