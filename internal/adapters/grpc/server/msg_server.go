package server

import (
	"context"
	ports "gitlab.com/g6834/team17/analytics-service/internal/ports/input"
	pb "gitlab.com/g6834/team17/grpc/analytics_messaging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type GrpcConsumer struct {
	server      *grpc.Server
	msgConsumer *MsgGrpcHandler
	logger      *zap.Logger
}

const gRPCaddr = `:50051`

func New(service ports.EventService, logger *zap.Logger) GrpcConsumer {
	grpcServer := grpc.NewServer()
	consumer := newGrpcConsumer(service, logger)

	return GrpcConsumer{server: grpcServer, msgConsumer: &consumer}
}

func (c GrpcConsumer) StartConsume(ctx context.Context) error {
	conf := net.ListenConfig{
		Control:   nil,
		KeepAlive: 0,
	}
	l, err := conf.Listen(ctx, "tcp", gRPCaddr)
	if err != nil {
		return err
	}

	c.server = grpc.NewServer()
	pb.RegisterAnalyticsMsgServiceServer(c.server, c.msgConsumer)

	return c.server.Serve(l)
}

func (c GrpcConsumer) StopConsume(ctx context.Context) error {
	gracefulDone := make(chan struct{})

	go func() {
		c.server.GracefulStop()
		gracefulDone <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-gracefulDone:
		return nil
	}
}
