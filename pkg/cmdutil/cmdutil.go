package cmdutil

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// InterruptChan returns a channel that is closed when the interrupt signal is received.
func InterruptChan() <-chan any {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ret := make(chan any, 1)
	go func() {
		s := <-c
		ret <- s
		close(ret)
	}()

	return ret
}

// NewInterruptContext returns a context that is cancelled when the interrupt signal is received.
func NewInterruptContext() (context.Context, context.CancelFunc) {
	interruptChan := InterruptChan()

	return InterruptContextFromChan(interruptChan)
}

// InterruptContextFromChan returns a context that is cancelled when the interruptChan is closed.
func InterruptContextFromChan(interruptChan <-chan any) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-interruptChan
		cancel()
	}()

	return ctx, cancel
}
