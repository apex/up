// Package relay provides a reverse proxy which
// relays requests to your "vanilla" HTTP server,
// and supports crash recovery.
package relay

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/facebookgo/freeport"
	"github.com/jpillora/backoff"
	"github.com/pkg/errors"

	"github.com/apex/up"
	"github.com/apex/up/internal/logs"
	"github.com/apex/up/internal/logs/writer"
)

// TODO: Wait() and handle error
// TODO: add timeout
// TODO: scope all plugin logs to their plugin name
// TODO: utilize BufferPool
// TODO: if the first Start() fails then bail

// log context.
var ctx = logs.Plugin("relay")

// DefaultTransport used by relay.
var DefaultTransport http.RoundTripper = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   2 * time.Second,
		KeepAlive: 2 * time.Second,
		DualStack: true,
	}).DialContext,
	MaxIdleConns:          0,
	MaxIdleConnsPerHost:   10,
	IdleConnTimeout:       5 * time.Minute,
	TLSHandshakeTimeout:   2 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

// Proxy is a reverse proxy and sub-process monitor
// for ensuring your web server is running.
type Proxy struct {
	config *up.Config

	mu       sync.Mutex
	restarts int
	port     int
	target   *url.URL
	*httputil.ReverseProxy
}

// New proxy.
func New(c *up.Config) (http.Handler, error) {
	p := &Proxy{
		config: c,
	}

	if err := p.Start(); err != nil {
		return nil, err
	}

	return p, nil
}

// Start the server.
func (p *Proxy) Start() error {
	if err := p.start(); err != nil {
		return err
	}

	p.ReverseProxy = httputil.NewSingleHostReverseProxy(p.target)
	p.ReverseProxy.Transport = p

	// TODO: configurable timeout
	ctx.Infof("waiting for %s", p.target.String())

	timeout := time.Duration(p.config.Proxy.Timeout) * time.Second
	if err := waitForListen(p.target, timeout); err != nil {
		return errors.Wrapf(err, "waiting for %s to be in listening state", p.target.String())
	}

	return nil
}

// Restart the server.
func (p *Proxy) Restart() error {
	ctx.Warn("restarting")
	p.restarts++

	if err := p.Start(); err != nil {
		return err
	}

	ctx.WithField("restarts", p.restarts).Warn("restarted")
	return nil
}

// ServeHTTP implementation.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ReverseProxy.ServeHTTP(w, r)
}

// RoundTrip implementation.
func (p *Proxy) RoundTrip(r *http.Request) (*http.Response, error) {
	// TODO: give up after N attempts

	b := p.config.Proxy.Backoff.Backoff()

retry:
	// replace host as it will change on restart
	r.URL.Host = p.target.Host
	res, err := DefaultTransport.RoundTrip(r)

	// everything is fine
	if err == nil {
		return res, nil
	}

	// temporary error, try again
	if e, ok := err.(net.Error); ok && e.Temporary() {
		ctx.WithError(err).Warn("temporary error")
		time.Sleep(b.Duration())
		goto retry
	}

	// timeout error, try again
	if e, ok := err.(net.Error); ok && e.Timeout() {
		ctx.WithError(err).Warn("timed out")
		time.Sleep(b.Duration())
		goto retry
	}

	// restart the server, try again
	ctx.WithError(err).Error("network error")
	if err := p.Restart(); err != nil {
		return nil, errors.Wrap(err, "restarting")
	}

	goto retry
}

// environment returns the server env variables.
func (p *Proxy) environment() []string {
	return []string{
		env("PORT", p.port),
		env("UP_RESTARTS", p.restarts),
	}
}

// start the server on a free port.
func (p *Proxy) start() error {
	port, err := freeport.Get()
	if err != nil {
		return errors.Wrap(err, "getting free port")
	}
	p.port = port
	ctx.Infof("found free port %d", port)

	ctx.Infof("executing %q", p.config.Proxy.Command)
	cmd := exec.Command("sh", "-c", p.config.Proxy.Command)
	cmd.Stdout = writer.New(log.InfoLevel)
	cmd.Stderr = writer.New(log.ErrorLevel)
	env := append(p.environment(), "PATH=node_modules/.bin:"+os.Getenv("PATH"))
	cmd.Env = append(os.Environ(), env...)

	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "running command")
	}

	target, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", port))
	if err != nil {
		return errors.Wrap(err, "parsing url")
	}
	p.target = target

	return nil
}

// env returns an environment variable.
func env(name string, val interface{}) string {
	return fmt.Sprintf("%s=%v", name, val)
}

// waitForListen blocks until `u` is listening with timeout.
func waitForListen(u *url.URL, timeout time.Duration) error {
	timedout := time.After(timeout)

	b := backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    time.Second,
		Factor: 1.5,
	}

	for {
		select {
		case <-timedout:
			return errors.Errorf("timed out after %s", timeout)
		case <-time.After(b.Duration()):
			if isListening(u) {
				return nil
			}
		}
	}
}

// isListening returns true if there's a server listening on `u`.
func isListening(u *url.URL) bool {
	conn, err := net.Dial("tcp", u.Host)
	if err != nil {
		return false
	}

	conn.Close()
	return true
}
