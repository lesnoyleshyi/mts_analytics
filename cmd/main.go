package main

import (
	"context"
	"gitlab.com/g6834/team17/analytics-service/internal/app"
	"os/signal"
	"syscall"
)

func main() {
	rootContext := context.Background()
	ctx, stop := signal.NotifyContext(rootContext,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	errCh := make(chan error)
	go app.Start(ctx, errCh)

	select {
	case <-ctx.Done():
		app.Stop()
	case <-errCh:
		app.Stop()
	}
}
