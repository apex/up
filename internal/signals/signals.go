package signals

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type closer func(os.Signal)

var (
	signals  = make(chan os.Signal, 1)
	exit     = make(chan bool, 1)
	closers  = make([]closer, 0)
	incoming os.Signal
)

// Init signals channel
func root() {
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
}

// Add closer
func AddCloser(fn closer) {
	closers = append([]closer{fn}, closers...)
}

// Capture all signals
func Capture() {
	go func() {
		<-signals
		for _, fn := range closers {
			fmt.Println("Executing")
			fn(incoming)
		}
	}()
}
