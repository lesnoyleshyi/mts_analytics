package app

import (
	"context"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/grpc/client"
	httpAdapter "gitlab.com/g6834/team17/analytics-service/internal/adapters/http"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/http/business_server"
	httpInterfaces "gitlab.com/g6834/team17/analytics-service/internal/adapters/http/interfaces"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/http/profile_server"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/http/swagger_server"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/interfaces"
	"gitlab.com/g6834/team17/analytics-service/internal/config"
	"gitlab.com/g6834/team17/analytics-service/internal/domain/usecases"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"log"
	"sync"
)

var err error
var logger *zap.Logger
var storage interfaces.Storage
var responder httpInterfaces.Responder
var gRPCValidator client.AuthClient
var messageConsumer interfaces.MessageConsumer
var businessServer business_server.AdapterHTTP
var profileServer profile_server.ProfileAdapter
var documentationServer swagger_server.SwaggerAdapter

const pathToConfigFile = `./config.yaml`

func Start(ctx context.Context, errChannel chan<- error) {
	if err := config.ReadConfigYML(pathToConfigFile); err != nil {
		log.Fatalf("error reading config file %s: %s", pathToConfigFile, err)
	}
	cfg := config.GetConfig()

	logger = NewLogger()

	responder = httpAdapter.NewJSONResponder(logger)

	// perhaps could be hide in NewValidator same way as NewStorage()
	gRPCValidator = client.NewGrpcAuth()
	validator := httpAdapter.NewJWTValidator(&gRPCValidator, responder, logger)

	storage = NewStorage(cfg.DB.Type)

	eventService := usecases.NewEventService(storage)

	messageConsumer = NewConsumer(cfg.Consumer.Type, eventService, logger)

	businessServer = business_server.New(eventService, logger, &validator, responder)
	profileServer = profile_server.NewProfileServer(logger)
	documentationServer = swagger_server.New()

	group, gctx := errgroup.WithContext(ctx)
	group.Go(func() error { return storage.Connect(gctx) })
	group.Go(func() error { return gRPCValidator.Connect(gctx) })
	group.Go(func() error { return businessServer.Start(gctx) })
	group.Go(func() error { return profileServer.Start(gctx) })
	group.Go(func() error { return documentationServer.Start(gctx) })
	group.Go(func() error { return messageConsumer.StartConsume(gctx) })

	logger.Info("application is starting")

	if err = group.Wait(); err != nil {
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
		if err := businessServer.Stop(ctx); err != nil {
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
