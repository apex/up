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

// flusher is the interface used for flushing logs.
type flusher interface {
	Flush()
}

// DefaultTransport used by relay.
var DefaultTransport http.RoundTripper = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   2 * time.Second,
		KeepAlive: 2 * time.Second,
		DualStack: true,
	}).DialContext,
	DisableKeepAlives: true,
}

// Proxy is a reverse proxy and sub-process monitor
// for ensuring your web server is running.
type Proxy struct {
	config *up.Config

	mu       sync.Mutex
	restarts int
	port     int
	target   *url.URL

	// cmd is the active running user application sub-process
	cmd *exec.Cmd

	// stdout is the log writer for structured logging output
	stdout *writer.Writer

	// stderr is the log writer for structured logging output
	stderr *writer.Writer

	// cmdCleanup is a channel that queues abandoned commands
	// so they can be cleaned up and resources reclaimed.
	cmdCleanup chan *exec.Cmd

	// maxRetries is the number of times to retry a single
	// request before failing altogether.
	maxRetries int

	// shutdownTimeout is the amount of time to wait between sending
	// a SIGINT and finally killing with a SIGKILL.
	shutdownTimeout time.Duration

	// timeout is the amount of time that a response may take,
	// including any retry attempts made.
	timeout time.Duration

	*httputil.ReverseProxy
}

// New proxy.
//
// We want to buffer the cleanup channel so that we can bound the
// number of concurrent processes executing, and prevent exhausting
// the ulimits of the host OS.
func New(c *up.Config) (http.Handler, error) {
	stdout, err := log.ParseLevel(c.Logs.Stdout)
	if err != nil {
		return nil, errors.Wrap(err, "invalid stdout error level")
	}

	stderr, err := log.ParseLevel(c.Logs.Stderr)
	if err != nil {
		return nil, errors.Wrap(err, "invalid stdout error level")
	}

	p := &Proxy{
		config:          c,
		cmdCleanup:      make(chan *exec.Cmd, 3),
		maxRetries:      c.Proxy.Backoff.Attempts,
		timeout:         time.Duration(c.Proxy.Timeout) * time.Second,
		shutdownTimeout: time.Duration(c.Proxy.ShutdownTimeout) * time.Second,
		stdout:          writer.New(stdout, ctx),
		stderr:          writer.New(stderr, ctx),
	}

	defer p.flushLogs()
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

	start := time.Now()
	timeout := time.Duration(p.config.Proxy.ListenTimeout) * time.Second
	ctx.WithField("url", p.target.String()).Info("waiting for server to listen")

	if err := util.WaitForListen(p.target, timeout); err != nil {
		return errors.Wrapf(err, "waiting for %s to be in listening state", p.target.String())
	}

	ctx.WithField("duration", util.MillisecondsSince(start)).Info("server is listening")
	return nil
}

// Restart the server.
func (p *Proxy) Restart() error {
	defer p.flushLogs()

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
	p.flushLogs()
}

// RoundTrip implementation.
func (p *Proxy) RoundTrip(r *http.Request) (*http.Response, error) {
	b := p.config.Proxy.Backoff.Backoff()
	start := time.Now()
	attempts := -1

retry:
	attempts++

	// Starting on the second attempt, we need to rewind the body if we can
	// The DefaultTransport.RoundTrip will only rewind it for us in non-err scenarios
	if attempts > 0 && r.Body != http.NoBody && r.Body != nil && r.GetBody != nil {
		newBody, err := r.GetBody()
		if err != nil {
			return nil, err
		}
		r.Body = newBody
	}

	// replace host as it will change on restart
	r.URL.Host = p.target.Host
	res, err := DefaultTransport.RoundTrip(r)

	// retries disabled, don't create noise in the logs
	if p.maxRetries == 0 {
		return res, err
	}

	// attempts exceeded, respond as-is
	if attempts >= p.maxRetries {
		ctx.Warn("retry attempts exceeded")
		return res, err
	}

	// timeout exceeded, respond as-is
	if time.Since(start) >= p.timeout {
		// TODO: timeout in-flight as well
		ctx.Warn("retry timeout exceeded")
		return res, err
	}

	// we got an error response, retry if possible
	if err == nil && res.StatusCode >= 500 && isIdempotent(r) {
		ctx.WithField("status", res.StatusCode).Warn("retrying idempotent request")
		goto retry
	}

	// we got a response
	if err == nil {
		return res, nil
	}

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

	var restartErr error
	if restartErr = p.Restart(); restartErr != nil {
		// we want to restart, but not mask the error above
		ctx.WithError(restartErr).Error("restarting")
	}

	// retry idempotent requests
	if restartErr == nil && isIdempotent(r) {
		ctx.Info("retrying idempotent request")
		goto retry
	}

	return nil, errors.Wrap(err, "network")
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

	ctx.WithField("port", port).Info("found free port")
	target, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", port))
	if err != nil {
		return errors.Wrap(err, "parsing url")
	}

	p.port = port
	p.target = target

	ctx.WithField("command", p.config.Proxy.Command).Info("executing")
	cmd = p.command(p.config.Proxy.Command, p.environment())
	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "running command")
	}

	// Only remember this if it was successfully started
	p.cmd = cmd
	ctx.WithField("pid", cmd.Process.Pid).Info("proxy started")

	return nil
}

// cleanupAbandoned consumes the cmdCleanup channel and signals
// abandoned processes to shut down and release their resources.
func (p *Proxy) cleanupAbandoned() {
	for cmd := range p.cmdCleanup {
		if cmd == nil {
			continue
		}

		done := make(chan bool)

		go func() {
			defer close(done)
			code := util.ExitStatus(cmd, cmd.Wait())
			ctx.WithField("pid", cmd.Process.Pid).WithField("code", code).Info("proxy exited")
		}()

		// We have deemed this command suitable for cleanup,
		// but we aren't positive the reason was because of an actual
		// process shutdown. First try to nicely send a SIGINT.
		cmd.Process.Signal(os.Interrupt)

		select {
		case <-done:
			continue
		case <-time.After(p.shutdownTimeout):
			ctx.WithField("pid", cmd.Process.Pid).Warn("proxy sending SIGKILL")
			cmd.Process.Kill()
			<-done
		}
	}
}

// command returns the command for spawning a server.
func (p *Proxy) command(s string, env []string) *exec.Cmd {
	cmd := exec.Command("sh", "-c", s)
	cmd.Stdout = p.stdout
	cmd.Stderr = p.stderr
	cmd.Env = append(os.Environ(), append(env, "PATH=node_modules/.bin:"+os.Getenv("PATH"))...)
	return cmd
}

// flushLogs flushes any pending logs.
func (p *Proxy) flushLogs() {
	p.stdout.Flush()
	p.stderr.Flush()
}

// env returns an environment variable.
func env(name string, val interface{}) string {
	return fmt.Sprintf("%s=%v", name, val)
}

// isIdempotent returns true if the request is considered idempotent.
func isIdempotent(req *http.Request) bool {
	switch req.Method {
	case "GET", "HEAD", "OPTIONS":
		return true
	default:
		return false
	}
}
