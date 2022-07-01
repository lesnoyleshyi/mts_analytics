package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	ports "gitlab.com/g6834/team17/analytics-service/internal/ports/input"
	"go.uber.org/zap"
	"net/http"
)

type AdapterHTTP struct {
	events ports.EventService
	logger *zap.Logger
	server *http.Server
}

const httpAddr = `:80`

func New(eventService ports.EventService, logger *zap.Logger) AdapterHTTP {
	var adapter AdapterHTTP

	adapter.events = eventService
	adapter.logger = logger
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

func (a AdapterHTTP) Start() error {
	if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("couldn't start server: %w", err)
	}
	return nil
}

func (a AdapterHTTP) Stop(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

func (a AdapterHTTP) routes() http.Handler {
	r := chi.NewRouter()

	//r.Use(a.validate.Validate)
	r.Use(middleware.RequestID)

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf("request %s. ok", middleware.GetReqID(r.Context()))))
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
