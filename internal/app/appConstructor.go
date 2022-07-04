package app

import (
	"context"
	httpAdapter "gitlab.com/g6834/team17/analytics-service/internal/adapters/http"
	"gitlab.com/g6834/team17/analytics-service/internal/config"
	"gitlab.com/g6834/team17/analytics-service/internal/domain/usecases"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"log"
	"sync"
)

var err error
var logger *zap.Logger
var storage Storage
var httpServer httpAdapter.AdapterHTTP
var profileServer httpAdapter.ProfileAdapter

const pathToConfigFile = `config.yaml`

//const PGConnStr = `postgres://team17:mNgd3ETbhVGd@91.185.93.23:5432/events`

func Start(ctx context.Context, errChannel chan<- error) {
	if err := config.ReadConfigYML(pathToConfigFile); err != nil {
		log.Fatalf("error reading config file %s: %s", pathToConfigFile, err)
	}
	cfg := config.GetConfig()

	logger = NewLogger()
	storage = NewStorage(cfg.DB.Type)
	eventService := usecases.NewEventService(storage)
	httpServer = httpAdapter.New(eventService, logger)
	profileServer = httpAdapter.NewProfileServer(logger)

	group, gctx := errgroup.WithContext(ctx)
	group.Go(func() error { return storage.Connect(gctx) })
	group.Go(func() error { return httpServer.Start(gctx) })
	group.Go(func() error { return profileServer.Start(gctx) })

	logger.Info("application is starting")

	if err = group.Wait(); err != nil {
		logger.Error("application start fail", zap.Error(err))
		errChannel <- err
	}
}

func Stop() {
	var wg sync.WaitGroup
	ctx := context.Background()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Stop(ctx); err != nil {
			logger.Warn("main server shutdown error", zap.Error(err))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := profileServer.Stop(ctx); err != nil {
			logger.Warn("profile server shutdown error", zap.Error(err))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := storage.Close(ctx); err != nil {
			logger.Warn("error on storage closing", zap.Error(err))
		} else {
			logger.Info("storage closed gracefully")
		}
	}()

	wg.Wait()
	logger.Info("application stopped successfully")
}
