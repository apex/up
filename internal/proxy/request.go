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

	// custom fields
	req.Header.Set("X-Request-Id", e.RequestContext.RequestID)
	req.Header.Set("X-Stage", e.RequestContext.Stage)

	// host
	req.URL.Host = req.Header.Get("Host")
	req.Host = req.URL.Host

	// user
	auth := req.Header.Get("Authorization")
	if auth != "" {
		authType, creds, err := parseAuthHeader(auth)
		if err != nil {
			return nil, errors.Wrap(err, "parsing auth header")
		}

		// Specifically allow Basic auth to be parsed
		// Purposefully ignore other kinds of auth, to allow them to trivially pass through
		// by simply copying their headers
		if strings.ToLower(authType) == "basic" {
			userinfo, err := basic(creds)
			if err != nil {
				return nil, errors.Wrap(err, "parsing basic auth")
			}

			req.URL.User = userinfo
		}
	}

	// TODO: pass the original json input
	return req, nil
}

// parseAuthHeader will parse the authorization header to determine the type of auth used
// in this request.
func parseAuthHeader(s string) (authType, credentials string, err error) {
	p := strings.SplitN(s, " ", 2)
	if len(p) != 2 {
		return "", "", errors.New("malformed")
	}

	return p[0], p[1], nil
}

// basic parses credentials of an Authorization header, this assumes it is basic auth
func basic(creds string) (userinfo *url.Userinfo, err error) {
	b, err := base64.StdEncoding.DecodeString(creds)
	if err != nil {
		return nil, errors.Wrap(err, "decoding")
	}

	p := strings.SplitN(string(b), ":", 2)
	return url.UserPassword(p[0], p[1]), nil
}
