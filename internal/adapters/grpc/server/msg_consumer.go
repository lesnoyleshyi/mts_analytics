package server

import (
	"context"
	"errors"
	"gitlab.com/g6834/team17/analytics-service/internal/domain/entity"
	ports "gitlab.com/g6834/team17/analytics-service/internal/ports/input"
	pb "gitlab.com/g6834/team17/grpc/analytics_messaging"
	"go.uber.org/zap"
	"strings"
)

type MsgGrpcHandler struct {
	service ports.EventService
	logger  *zap.Logger
	pb.UnimplementedAnalyticsMsgServiceServer
}

// TODO remove this error. Domain layer should return its own errors
var ErrSaveMessageFailed = errors.New("can't handle message")

func newGrpcConsumer(s ports.EventService, l *zap.Logger) MsgGrpcHandler {
	return MsgGrpcHandler{service: s, logger: l}
}

func (c MsgGrpcHandler) SendMessage(ctx context.Context, msg *pb.EventMessage) (*pb.Response, error) {
	var event entity.Event

	event.TaskUUID = msg.TaskUuid
	event.EventType = strings.ToLower(msg.EventType.String())
	event.UserUUID = msg.UserUuid
	event.Timestamp = msg.Timestamp.AsTime()

	if err := c.service.Save(ctx, event); err != nil {
		c.logger.Warn("can't save event from message stream", zap.Error(err))
		// TODO return error from domain
		return &pb.Response{Success: false}, ErrSaveMessageFailed
	}

	return &pb.Response{Success: true}, nil
}
