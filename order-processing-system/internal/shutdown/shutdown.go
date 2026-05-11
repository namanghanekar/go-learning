package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func Context() context.Context {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ctx.Done()
		stop()
	}()
	return ctx
}
