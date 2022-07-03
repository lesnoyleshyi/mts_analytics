package main

import (
	"context"
	"gitlab.com/g6834/team17/analytics-service/internal/app"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	rootContext := context.Background()
	ctx, stop := signal.NotifyContext(rootContext,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer func() {
		// debug message. if we call os.Exit() in Start(), it's a goroutine leak
		log.Println("stop() was called")
		stop()
	}()

	errCh := make(chan error)
	go app.Start(ctx, errCh)

	select {
	case <-ctx.Done():
		app.Stop()
	case <-errCh:
		app.Stop()
	}
}
