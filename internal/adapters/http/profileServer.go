package http

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ProfileAdapter struct {
	logger *zap.Logger
	server *http.Server
}

const httpProfileAddr = `:8080`

func NewProfileServer(logger *zap.Logger) ProfileAdapter {
	var adapter ProfileAdapter

	adapter.logger = logger
	s := http.Server{ //nolint:exhaustruct
		Addr:    httpProfileAddr,
		Handler: adapter.routeProfiles(),
	}
	adapter.server = &s

	return adapter
}

func (a ProfileAdapter) Start(ctx context.Context) error {
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

func (a ProfileAdapter) Stop(ctx context.Context) error {
	if a.server == nil {
		a.logger.Info("profile server wasn't initialised, stop() is no-op")
		return nil
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*gracefulShutdownDelaySec)
	defer cancel()

	err := a.server.Shutdown(timeoutCtx)
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return err
	}

	a.logger.Info("profile server stopped gracefully")

	return nil
}
