package swagger_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type SwaggerAdapter struct {
	server *http.Server
}

const host = `localhost`
const port = `:9090`
const gracefulShutdownDelaySec = 30

func New() SwaggerAdapter {
	var adapter SwaggerAdapter

	s := http.Server{
		Addr:    port,
		Handler: adapter.routes(),
	}
	adapter.server = &s

	return adapter
}

func (a SwaggerAdapter) Start(ctx context.Context) error {
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

func (a SwaggerAdapter) Stop(ctx context.Context) error {
	if a.server == nil {
		return nil
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*gracefulShutdownDelaySec)
	defer cancel()

	err := a.server.Shutdown(timeoutCtx)
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
