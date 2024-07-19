package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// ContextWithSignal sets up a signal listener, which will cancel the returned context when
// an interupt (SIGINT OR SIGTERM) is received.
func ContextWithSignal(ctx context.Context) context.Context {
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		<-quitCh
	}()
	return ctx

}
