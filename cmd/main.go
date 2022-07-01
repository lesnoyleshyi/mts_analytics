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

	go app.Start(ctx)
	<-ctx.Done()
	app.Stop(ctx)
}
