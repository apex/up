package proxy

import (
	"encoding/base64"
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

	// url.UserPassword(username, password)

	// base64 encoded body
	body := e.Body
	if e.IsBase64Encoded {
		b, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return nil, errors.Wrap(err, "decoding base64 body")
		}
		body = string(b)
	}

	// new request
	req, err := http.NewRequest(e.HTTPMethod, u.String(), strings.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	// remote addr
	req.RemoteAddr = e.RequestContext.Identity.SourceIP

	// header fields
	for k, v := range e.Headers {
		req.Header.Set(k, v)
	}

	// host
	req.URL.Host = req.Header.Get("Host")
	req.Host = req.URL.Host

	// user
	auth := req.Header.Get("Authorization")
	if auth != "" {
		user, pass, err := basic(auth)
		if err != nil {
			return nil, errors.Wrap(err, "parsing basic auth")
		}
		req.URL.User = url.UserPassword(user, pass)
	}

	// TODO: pass the original json input
	return req, nil
}

// basic auth parser.
func basic(s string) (user, pass string, err error) {
	p := strings.SplitN(s, " ", 2)

	if len(p) != 2 || p[0] != "Basic" {
		return "", "", errors.New("malformed")
	}

	b, err := base64.StdEncoding.DecodeString(p[1])
	if err != nil {
		return "", "", errors.Wrap(err, "decoding")
	}

	pair := strings.SplitN(string(b), ":", 2)
	return pair[0], pair[1], nil
}
