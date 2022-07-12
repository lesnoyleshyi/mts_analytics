package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/http/interfaces"
	ports "gitlab.com/g6834/team17/analytics-service/internal/ports/input"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type AdapterHTTP struct {
	events    ports.EventService
	validator interfaces.MiddlewareValidator
	logger    *zap.Logger
	server    *http.Server
}

const httpAddr = `:80`
const gracefulShutdownDelaySec = 30

func New(s ports.EventService, l *zap.Logger, v interfaces.MiddlewareValidator) AdapterHTTP {
	var adapter AdapterHTTP

	adapter.events = s
	adapter.validator = v
	adapter.logger = l
	server := http.Server{
		Addr:    httpAddr,
		Handler: adapter.routes(),
		// we could wrap *zap.Logger in adapter to pass here
		ErrorLog: nil,
		// maybe we should pass context from main.go here
		BaseContext: nil,
		// or here
		ConnContext: nil,
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

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*gracefulShutdownDelaySec)
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

	r.Use(middleware.RequestID)
	r.Use(a.validator.Validate)

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
	a.logger.Info("error serving request", zap.Error(err))
	// http.Error requires response be plain text
	//w.Header().Set("Content-Type", "application/json")
	http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", msg), status)
}
