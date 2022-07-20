package app

import (
	"context"
	"flag"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/grpc/client"
	grpcServer "gitlab.com/g6834/team17/analytics-service/internal/adapters/grpc/server"
	httpAdapter "gitlab.com/g6834/team17/analytics-service/internal/adapters/http"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/http/business_server"
	httpInterfaces "gitlab.com/g6834/team17/analytics-service/internal/adapters/http/interfaces"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/http/profile_server"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/http/swagger_server"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/interfaces"
	"gitlab.com/g6834/team17/analytics-service/internal/domain/usecases"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"sync"
)

var err error
var logger *zap.Logger
var storage Storage
var responder httpInterfaces.Responder
var gRPCValidator client.AuthClient
var httpServer business_server.AdapterHTTP
var profileServer profile_server.ProfileAdapter
var documentationServer swagger_server.SwaggerAdapter
var messageConsumer interfaces.MessageConsumer

func Start(ctx context.Context, errChannel chan<- error) {
	// should be hide in some config-initialising function
	storageType := flag.String("storage", "postgres",
		"defines storage type: postgres, mongo, cache, etc")
	flag.Parse()
	if storageType == nil {
		*storageType = "postgres"
	}

	logger, _ = zap.NewProduction()

	responder = httpAdapter.NewJSONResponder(logger)

	// perhaps could be hide in NewValidator same way as NewStorage()
	gRPCValidator = client.NewGrpcAuth()
	validator := httpAdapter.NewJWTValidator(&gRPCValidator, responder, logger)

	storage = NewStorage(*storageType)

	eventService := usecases.NewEventService(storage)
	httpServer = business_server.New(eventService, logger, &validator, responder)

	profileServer = profile_server.NewProfileServer(logger)

	documentationServer = swagger_server.New()

	messageConsumer = grpcServer.New(eventService, logger)

	group, gctx := errgroup.WithContext(ctx)
	group.Go(func() error { return storage.Connect(gctx) })
	group.Go(func() error { return gRPCValidator.Connect(gctx) })
	group.Go(func() error { return httpServer.Start(gctx) })
	group.Go(func() error { return profileServer.Start(gctx) })
	group.Go(func() error { return documentationServer.Start(gctx) })
	group.Go(func() error { return messageConsumer.StartConsume(gctx) })

	logger.Info("application is starting")

	if err = group.Wait(); err != nil {
		// may be should panic instead of fatal-ing. Is it necessary to call stop() in main.go?
		logger.Error("application start fail", zap.Error(err))
		errChannel <- err
	}
}

func Stop() {
	var wg sync.WaitGroup
	// TODO decide what kind of context should be passed in each case
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
		if err := documentationServer.Stop(ctx); err != nil {
			logger.Warn("swagger server shutdown error", zap.Error(err))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := gRPCValidator.Disconnect(ctx); err != nil {
			logger.Warn("validator disconnection error", zap.Error(err))
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := messageConsumer.StopConsume(ctx); err != nil {
			logger.Warn("error stopping message consuming gracefully")
		} else {
			logger.Info("message consumer stopped gracefully")
		}
	}()

	wg.Wait()
	logger.Info("application stopped successfully")
}
