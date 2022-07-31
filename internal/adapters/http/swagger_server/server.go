package swagger_server

import (
	"context"
	"errors"
	"fmt"
	"gitlab.com/g6834/team17/analytics-service/internal/config"
	"net"
	"net/http"
	"time"
)

type SwaggerAdapter struct {
	server *http.Server
}

func New() SwaggerAdapter {
	var adapter SwaggerAdapter
	cfg := config.GetConfig()
	docAddr := net.JoinHostPort(cfg.Rest.Host, cfg.Rest.DocPort)

	s := http.Server{ //nolint:exhaustruct
		Addr:    docAddr,
		Handler: adapter.routes(),
	}
	adapter.server = &s

	return adapter
}

func (a SwaggerAdapter) Start(ctx context.Context) error {
	srvErrChan := make(chan error)

	go func() {
		if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			srvErrChan <- fmt.Errorf("couldn't start documentation server: %w", err)
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

func (a SwaggerAdapter) Stop(ctx context.Context) error {
	if a.server == nil {
		return nil
	}

	cfg := config.GetConfig()
	gracefulShutdownDelaySec := time.Duration(cfg.Rest.GracefulTimeout)

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*gracefulShutdownDelaySec)
	defer cancel()

	err := a.server.Shutdown(timeoutCtx)
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
