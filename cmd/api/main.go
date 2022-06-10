package main

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"mts_analytics/internal/handlers"
	"mts_analytics/internal/repository"
	"mts_analytics/internal/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const httpPort string = ":8080"

func main() {
	log.SetLevel(log.DebugLevel)

	repo := repository.New()
	service := services.New(repo)
	handler := handlers.New(service)
	server := http.Server{
		Addr:    httpPort,
		Handler: handler.NewMux(),
	}

	done := make(chan struct{})

	go start(&server)
	shutdown(&server, done)

	<-done
	log.WithField("func:", "main").Println("Server stopped gracefully")
}

func start(srv *http.Server) {
	log.WithField("func", "start").Debugf("Server started listening on port %s", httpPort)
	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.WithField("func", "start").Fatalf("Server failed unexpectidly with error: %s", err)
	}
}

func shutdown(srv *http.Server, done chan struct{}) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-ch
	log.WithField("func", "shutdown").Printf("%s signal caught, shutting down gracefully", sig)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer func() {
		cancel()
		close(done)
	}()
	if err := srv.Shutdown(ctx); err != nil {
		log.WithField("func", "shutdown").Fatalf("Server shutdown failed: %s", err)
	}
}
