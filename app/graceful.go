package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ContextWithSignal sets up a signal listener, which will cancel the returned context when
// an interrupt (SIGINT OR SIGTERM) is received.
func ContextWithSignal(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		<-quitCh
	}()
	return ctx

}

// SleepContext sleeps for the given duration or until the context is cancelled.
// If the context is cancelled, it returns the error from the context.
func SleepContext(ctx context.Context, duration time.Duration) error {
	select {
	case <-time.After(duration):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
