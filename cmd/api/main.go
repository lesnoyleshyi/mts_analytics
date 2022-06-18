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
const httpProfilePort = `:8085`
const kafkaTopic = "task"
const app_name = `mok_task_service`
const host_ip = `lol_hz`

func main() {
	log.SetLevel(log.DebugLevel)
	ctx := context.Background()

	repo := repository.New()
	defer repo.Pool.Close()
	service := services.New(repo)
	handler := handlers.New(service)
	server := http.Server{
		Addr:    httpPort,
		Handler: handler.NewMux(),
	}

	profileServer := http.Server{
		Addr:    httpProfilePort,
		Handler: handlers.NewProfiler(),
	}

	kafkaHandler, err := handlers.NewKafkaHandler(service)
	if err != nil {
		log.Fatalf("structure sucks: %s", err)
	}
	defer func() { _ = (*kafkaHandler.ConsumerGroup).Close() }()
	defer func() { _ = (*kafkaHandler.Client).Close() }()
	done := make(chan struct{})

	go start(&server, httpPort)
	go start(&profileServer, httpProfilePort)
	go consume(ctx, kafkaHandler, kafkaTopic)
	shutdown(&server, done)
	shutdown(&profileServer, done)

	<-done
	log.WithField("func:", "main").Println("Server stopped gracefully")
}

func start(srv *http.Server, port string) {
	log.WithField("func", "start").Debugf("Server started listening on port %s", port)
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

func consume(ctx context.Context, h *handlers.KafkaHandler, topics ...string) {
	var errCnt int64
	for {
		if err := (*h.ConsumerGroup).Consume(ctx, topics, h.Consumer); err != nil {
			log.WithFields(log.Fields{
				"app_name":    app_name,
				"host_ip":     host_ip,
				"logger_name": "main.CG.Consume_loop",
			}).Errorf("error consuming: %s", err)
			errCnt++
			if errCnt >= 300 {
				log.Error("TOO MUCH ERRORS")
				break
			}
		}
		if err := ctx.Err(); err != nil {
			log.WithFields(log.Fields{
				"app_name":    app_name,
				"host_ip":     host_ip,
				"logger_name": "main.CG.Consume_loop",
			}).Errorf("context catch error: %s. Stop consuming", err)
			break
		}
	}
}
