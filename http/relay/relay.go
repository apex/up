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
	"syscall"
	"time"

	"github.com/apex/log"
	"github.com/facebookgo/freeport"
	"github.com/jpillora/backoff"
	"github.com/pkg/errors"

	"github.com/apex/up"
	"github.com/apex/up/internal/logs"
	"github.com/apex/up/internal/logs/writer"
)

// TODO: add timeout
// TODO: scope all plugin logs to their plugin name
// TODO: utilize BufferPool

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

	// cmdCleanup is a channel that queues abandoned commands so they can be cleaned up
	// and resources reclaimed
	cmdCleanup chan *exec.Cmd

	// maxRetries is the number of times to retry a single request before failing alltogether
	maxRetries int

	*httputil.ReverseProxy
}

// New proxy.
func New(c *up.Config) (http.Handler, error) {
	p := &Proxy{
		config: c,
		// We want to buffer this channel so that we can bound the number of concurrent processes
		// currently executing, and prevent exhausting the ulimits of the host OS
		cmdCleanup: make(chan *exec.Cmd, 3),
		maxRetries: 3,
	}

	if err := p.Start(); err != nil {
		return nil, err
	}

	// Launch a goroutine for cleaning up the old commands as they are abandoned
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

	// Only remember this command it if was successfully started
	p.cmd = cmd
	ctx.Infof("proxy (pid=%d) started", cmd.Process.Pid)

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

// cleanupAbandoned consumes the cmdCleanup channel and signals abandoned processes to shut down
// and release their resources.
func (p *Proxy) cleanupAbandoned() {
	for cmd := range p.cmdCleanup {
		if cmd == nil {
			continue
		}

		// Set up a channel to wait for this process to complete processing cleanly
		waitDone := make(chan bool, 1)
		go func() {
			err := cmd.Wait()
			ps := cmd.ProcessState
			if e, ok := err.(*exec.ExitError); ok && e != nil {
				ps = e.ProcessState
			}

			exitCode := "?"
			if ps != nil {
				sys := ps.Sys()
				if x, ok := sys.(syscall.WaitStatus); ok {
					exitCode = fmt.Sprintf("%d", x.ExitStatus())
				}
			}

			ctx.Infof("proxy (pid=%d) exited with code=%s", cmd.Process.Pid, exitCode)
			waitDone <- true
		}()

		// We have deemed this command suitable for cleanup, but we aren't positive the reason
		// was because of an actual process shutdown.  First try to nicely send a SIGINT.
		cmd.Process.Signal(os.Interrupt)

		select {
		case <-waitDone:
			continue
		case <-time.After(10 * time.Second):
			ctx.Warnf("proxy (pid=%d) sending SIGKILL", cmd.Process.Pid)
			cmd.Process.Kill()
			<-waitDone
		}
	}
}
