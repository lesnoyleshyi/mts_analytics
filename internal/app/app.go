package app

import (
	"context"
	httpAdapter "gitlab.com/g6834/team17/analytics-service/internal/adapters/http"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/postgres"
	"gitlab.com/g6834/team17/analytics-service/internal/domain/usecases"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"os"
)

//var httpServer *httpAdapter.AdapterHTTP

//const PGConnStr = `postgres://team17:mNgd3ETbhVGd@91.185.93.23:5432/events`

func Start(ctx context.Context) {
	PGConnStr := os.Getenv("PG_CONNSTR")

	logger, _ := zap.NewProduction()
	// seems not ok: New is constructor - it's for construction. not running.
	// looks like it'd be better to run postgres connection within errgroup,
	// the same way as starting httpAdapter etc.
	storage, err := postgres.New(ctx, PGConnStr)
	// is it ok to Close pool in this function?
	defer storage.Pool.Close()
	if err != nil {
		logger.Fatal("error initialising db", zap.Error(err))
	}
	eventService := usecases.New(storage)
	httpServer := httpAdapter.New(eventService, logger)

	group, _ := errgroup.WithContext(ctx)
	group.Go(httpServer.Start)

	logger.Info("application is starting")

	if err := group.Wait(); err != nil {
		// may be should panic instead of fatal-ing. Is it necessary to call stop in main.go?
		logger.Fatal("application start fail", zap.Error(err))
	}
}

func Stop(ctx context.Context) error {
	//return httpServer.Stop(ctx)
	return nil
}
