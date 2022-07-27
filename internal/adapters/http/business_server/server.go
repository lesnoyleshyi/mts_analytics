package business_server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/http/interfaces"
	ports "gitlab.com/g6834/team17/analytics-service/internal/ports/input"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type AdapterHTTP struct {
	events    ports.EventService
	validator interfaces.MiddlewareValidator
	responder interfaces.Responder
	logger    *zap.Logger
	server    *http.Server
}

const httpAddr = `:8282`
const gracefulShutdownDelaySec = 30

func New(s ports.EventService,
	l *zap.Logger,
	v interfaces.MiddlewareValidator,
	r interfaces.Responder) AdapterHTTP {
	var adapter AdapterHTTP

	adapter.events = s
	adapter.validator = v
	adapter.logger = l
	adapter.responder = r
	server := http.Server{ //nolint:exhaustruct
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
			srvErrChan <- fmt.Errorf("couldn't start business server: %w", err)
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

	r.Use(cors.Default().Handler)
	r.Use(a.validator.Validate)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	r.Mount("/", a.routeEvents())

	return r
}
