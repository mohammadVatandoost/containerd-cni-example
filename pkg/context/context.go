package context

import (
	"context"
	"os"
	"os/signal"
)

func WithSignalCancellation(
	ctx context.Context,
	signals ...os.Signal,
) (context.Context, context.CancelFunc) {
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, signals...)

	newContext, cancel := context.WithCancel(ctx)
	go func() {
		<-gracefulStop
		cancel()
	}()
	return newContext, cancel
}
