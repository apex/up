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

	// cmd refers to the currently running (active) proxy subprocss
	cmd *exec.Cmd

	// cmdCleanup is a channel that queues abandoned commands
	// so they can be cleaned up and resources reclaimed.
	cmdCleanup chan *exec.Cmd

	// maxRetries is the number of times to retry a single
	// request before failing alltogether.
	maxRetries int

	// shutdownTimeout is the amount of time to wait between sending
	// a SIGINT and finally killing with a SIGKILL.
	shutdownTimeout time.Duration

	*httputil.ReverseProxy
}

// New proxy.
//
// We want to buffer the cleanup channel so that we can bound the
// number of concurrent processes executing, and prevent exhausting
// the ulimits of the host OS.
func New(c *up.Config) (http.Handler, error) {
	p := &Proxy{
		config:          c,
		cmdCleanup:      make(chan *exec.Cmd, 3),
		maxRetries:      c.Proxy.Backoff.Attempts,
		shutdownTimeout: time.Duration(c.Proxy.ShutdownTimeout) * time.Second,
	}

	if err := p.Start(); err != nil {
		return nil, err
	}

	go p.cleanupAbandoned()

	return p, nil
}

// Start the server.
func (p *Proxy) Start() error {
	if err := p.start(); err != nil {
		return err
	}

	p.ReverseProxy = httputil.NewSingleHostReverseProxy(p.target)
	p.ReverseProxy.Transport = p

	timeout := time.Duration(p.config.Proxy.ListenTimeout) * time.Second
	ctx.Infof("waiting for %s to listen (timeout %s)", p.target.String(), timeout)

	if err := util.WaitForListen(p.target, timeout); err != nil {
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
	b := p.config.Proxy.Backoff.Backoff()
	retries := 0

retry:
	// replace host as it will change on restart
	r.URL.Host = p.target.Host
	res, err := DefaultTransport.RoundTrip(r)

	// everything is fine
	if err == nil {
		return res, nil
	}

	// retries exceeded
	if retries >= p.maxRetries {
		return nil, err
	}

	retries++

	// temporary error, try again
	if e, ok := err.(net.Error); ok && e.Temporary() {
		ctx.WithError(err).Warn("temporary")
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
	ctx.WithError(err).Error("network")
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
	cmd := p.cmd
	p.cmd = nil

	// Send this previous command to be cleaned up (Waited on, killed if necessary)
	p.cmdCleanup <- cmd

	port, err := freeport.Get()
	if err != nil {
		return errors.Wrap(err, "getting free port")
	}

	ctx.Infof("found free port %d", port)
	target, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", port))
	if err != nil {
		return errors.Wrap(err, "parsing url")
	}

	p.port = port
	p.target = target

	ctx.Infof("executing %q", p.config.Proxy.Command)

	cmd = exec.Command("sh", "-c", p.config.Proxy.Command)
	cmd.Stdout = writer.New(log.InfoLevel)
	cmd.Stderr = writer.New(log.ErrorLevel)
	env := append(p.environment(), "PATH=node_modules/.bin:"+os.Getenv("PATH"))
	cmd.Env = append(os.Environ(), env...)
	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "running command")
	}

	// Only remember this if it was successfully started
	p.cmd = cmd
	ctx.Infof("proxy (pid=%d) started", cmd.Process.Pid)

	return nil
}

// cleanupAbandoned consumes the cmdCleanup channel and signals
// abandoned processes to shut down and release their resources.
func (p *Proxy) cleanupAbandoned() {
	for cmd := range p.cmdCleanup {
		if cmd == nil {
			continue
		}

		done := make(chan bool, 1)

		go func() {
			err := cmd.Wait()
			code := util.ExitStatus(cmd, err)
			ctx.Infof("proxy (pid=%d) exited with code=%s", cmd.Process.Pid, code)
			done <- true
		}()

		// We have deemed this command suitable for cleanup,
		// but we aren't positive the reason was because of an actual
		// process shutdown. First try to nicely send a SIGINT.
		cmd.Process.Signal(os.Interrupt)

		select {
		case <-done:
			continue
		case <-time.After(p.shutdownTimeout):
			ctx.Warnf("proxy (pid=%d) sending SIGKILL", cmd.Process.Pid)
			cmd.Process.Kill()
			<-done
		}
	}
}

// env returns an environment variable.
func env(name string, val interface{}) string {
	return fmt.Sprintf("%s=%v", name, val)
}
