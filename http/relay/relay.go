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
	"github.com/pkg/errors"

	"github.com/apex/up"
	"github.com/apex/up/internal/logs"
	"github.com/apex/up/internal/logs/writer"
	"github.com/apex/up/internal/util"
)

// log context.
var ctx = logs.Plugin("relay")

// Proxy is a reverse proxy and sub-process monitor
// for ensuring your web server is running.
type Proxy struct {
	config *up.Config

	// transport used for the reverse proxy.
	transport http.RoundTripper

	// stdout is the log writer for structured logging output.
	stdout *writer.Writer

	// stderr is the log writer for structured logging output.
	stderr *writer.Writer

	mu sync.Mutex

	// restarts is the restart count.
	restarts int

	// url is the active application url.
	url *url.URL

	// ReverseProxy is the reverse proxy making the requests to the app.
	*httputil.ReverseProxy

	// cmd is the current child process of the app.
	cmd *exec.Cmd
}

// New proxy.
//
// We want to buffer the cleanup channel so that we can bound the
// number of concurrent processes executing, and prevent exhausting
// the ulimits of the host OS.
func New(c *up.Config) (http.Handler, error) {
	stdout, err := log.ParseLevel(c.Logs.Stdout)
	if err != nil {
		return nil, errors.Wrap(err, "invalid stdout log level")
	}

	stderr, err := log.ParseLevel(c.Logs.Stderr)
	if err != nil {
		return nil, errors.Wrap(err, "invalid stderr log level")
	}

	timeout := time.Duration(c.Proxy.Timeout) * time.Second

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 2 * time.Second,
			DualStack: true,
		}).DialContext,
		ResponseHeaderTimeout: timeout,
		DisableKeepAlives:     true,
	}

	p := &Proxy{
		config:    c,
		stdout:    writer.New(stdout, ctx),
		stderr:    writer.New(stderr, ctx),
		transport: transport,
	}

	if err := p.Start(); err != nil {
		return nil, err
	}

	return p, nil
}

// Start the server.
func (p *Proxy) Start() error {
	if err := p.startServer(); err != nil {
		return err
	}

	p.ReverseProxy = httputil.NewSingleHostReverseProxy(p.url)
	p.ReverseProxy.Transport = p

	start := time.Now()
	timeout := time.Duration(p.config.Proxy.ListenTimeout) * time.Second
	ctx.Info("waiting for app to listen on PORT")

	if err := util.WaitForListen(p.url, timeout); err != nil {
		return errors.Wrapf(err, "waiting for %s to be in listening state", p.url.String())
	}

	ctx.WithField("duration", util.MillisecondsSince(start)).Info("app listening")
	return nil
}

// Restart the server.
func (p *Proxy) Restart() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	ctx.Warn("restarting")
	p.restarts++

	if p.cmd != nil {
		if err := p.cmd.Process.Kill(); err != nil {
			ctx.WithError(err).Error("killing application process")
		}
	}

	if err := p.Start(); err != nil {
		return err
	}

	ctx.WithField("restarts", p.restarts).Warn("restarted")
	return nil
}

// RoundTrip implementation.
func (p *Proxy) RoundTrip(r *http.Request) (*http.Response, error) {
	res, err := p.transport.RoundTrip(r)

	// temporary error
	if e, ok := err.(net.Error); ok && e.Temporary() {
		ctx.WithError(err).Warn("request temporary error")
		return res, err
	}

	// timeout error
	if e, ok := err.(net.Error); ok && e.Timeout() {
		ctx.WithError(err).Warn("request timeout")
		return res, err
	}

	// network error
	if err != nil {
		ctx.WithError(err).Error("request network error")
		if err := p.Restart(); err != nil {
			ctx.WithError(err).Error("restarting")
		}
	}

	return res, err
}

// environment returns the server env variables.
func (p *Proxy) environment() []string {
	return []string{
		env("PORT", p.url.Port()),
		env("UP_RESTARTS", p.restarts),
	}
}

// startServer the server on a free port.
func (p *Proxy) startServer() error {
	port, err := freeport.Get()
	if err != nil {
		return errors.Wrap(err, "getting free port")
	}

	target, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", port))
	if err != nil {
		return errors.Wrap(err, "parsing url")
	}

	p.url = target

	ctx.WithField("command", p.config.Proxy.Command).WithField("PORT", port).Info("starting app")
	p.cmd = p.command(p.config.Proxy.Command, p.environment())

	if err := p.cmd.Start(); err != nil {
		return errors.Wrap(err, "running command")
	}

	ctx.Info("started app")
	return nil
}

// command returns the command for spawning a server.
func (p *Proxy) command(s string, env []string) *exec.Cmd {
	cmd := exec.Command("sh", "-c", s)
	cmd.Stdout = p.stdout
	cmd.Stderr = p.stderr
	cmd.Env = append(os.Environ(), append(env, "PATH=node_modules/.bin:"+os.Getenv("PATH"))...)
	return cmd
}

// env returns an environment variable.
func env(name string, val interface{}) string {
	return fmt.Sprintf("%s=%v", name, val)
}
