package main

import (
	"context"
	"os/signal"
	"syscall"
)

func main() {
	rootContext := context.Background()
	ctx, stop := signal.NotifyContext(rootContext,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	go app.Start(ctx)
	<- ctx.Done()
	app.Stop()
}
