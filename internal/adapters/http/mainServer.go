package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"gitlab.com/g6834/team17/analytics-service/internal/config"
	ports "gitlab.com/g6834/team17/analytics-service/internal/ports/input"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

type AdapterHTTP struct {
	events ports.EventService
	logger *zap.Logger
	server *http.Server
}

func New(eventService ports.EventService, logger *zap.Logger) AdapterHTTP {
	var adapter AdapterHTTP
	var cfg = config.GetConfig()

	adapter.events = eventService
	adapter.logger = logger.With(zap.String("host_port", cfg.Rest.BusinessPort))
	server := http.Server{ //nolint:exhaustruct
		Addr:     net.JoinHostPort(cfg.Rest.Host, cfg.Rest.BusinessPort),
		Handler:  adapter.routes(),
		ErrorLog: zap.NewStdLog(logger),
	}
	adapter.server = &server

	return adapter
}

func (a AdapterHTTP) Start(ctx context.Context) error {
	srvErrChan := make(chan error)

	go func() {
		if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			srvErrChan <- fmt.Errorf("couldn't start server: %w", err)
		}
		srvErrChan <- nil
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-srvErrChan:
		return err
	}
}

func (a AdapterHTTP) Stop(ctx context.Context) error {
	if a.server == nil {
		a.logger.Info("main server wasn't initialised, stop() is no-op")
		return nil
	}

	cfg := config.GetConfig()
	timeoutCtx, cancel := context.WithTimeout(ctx,
		time.Second*time.Duration(cfg.Rest.GracefulTimeout))
	defer cancel()

	err := a.server.Shutdown(timeoutCtx)
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return err
	}

	a.logger.Info("main server stopped gracefully")

	return nil
}

func (a AdapterHTTP) routes() http.Handler {
	r := chi.NewRouter()

	//r.Use(a.validate.Validate)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	r.Mount("/", a.routeEvents())

	return r
}

func (a AdapterHTTP) respondSuccess(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := fmt.Fprintf(w, "{\"result\":\"%s\"}", msg); err != nil {
		a.logger.Warn("error writing response", zap.Error(err))
	}
}

func (a AdapterHTTP) respondError(w http.ResponseWriter, msg string, status int, err error) {
	a.logger.Info("error serving request",
		zap.Error(err),
	)
	// http.Error requires response be plain text
	//w.Header().Set("Content-Type", "application/json")
	http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", msg), status)
}
