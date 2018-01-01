package signal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/apex/up/internal/util"
)

// close funcs.
var fns []Func

// Init signals channel
func init() {
	s := make(chan os.Signal, 1)
	go trap(s)
	signal.Notify(s, syscall.SIGINT)
}

// Func is a close function.
type Func func() error

// Add registers a close handler func.
func Add(fn Func) {
	fns = append(fns, fn)
}

// trap signals to invoke callbacks and exit.
func trap(ch chan os.Signal) {
	<-ch
	for _, fn := range fns {
		if err := fn(); err != nil {
			util.Fatal(err)
		}
	}
	os.Exit(1)
}
