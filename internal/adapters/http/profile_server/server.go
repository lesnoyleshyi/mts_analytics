package profile_server

import (
	"context"
	"errors"
	"fmt"
	"gitlab.com/g6834/team17/analytics-service/internal/config"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

type ProfileAdapter struct {
	logger *zap.Logger
	server *http.Server
}

func NewProfileServer(logger *zap.Logger) ProfileAdapter {
	var adapter ProfileAdapter
	var cfg = config.GetConfig()

	adapter.logger = logger.With(zap.String("host_port", cfg.Rest.DebugPort))
	s := http.Server{ //nolint:exhaustruct
		Addr:    net.JoinHostPort(cfg.Rest.Host, cfg.Rest.DebugPort),
		Handler: adapter.routeProfiles(),
	}
	adapter.server = &s

	return adapter
}

func (a ProfileAdapter) Start(ctx context.Context) error {
	srvErrChan := make(chan error)

	go func() {
		a.logger.Info("profile server starts listen", zap.String("port", a.server.Addr))

		if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			srvErrChan <- fmt.Errorf("couldn't start profile server: %w", err)
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

	cfg := config.GetConfig()
	timeoutCtx, cancel := context.WithTimeout(ctx,
		time.Second*time.Duration(cfg.Rest.GracefulTimeout))
	defer cancel()

	err := a.server.Shutdown(timeoutCtx)
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return err
	}

	a.logger.Info("profile server stopped gracefully")

	return nil
}
