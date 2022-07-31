package app

import (
	grpcConsumer "gitlab.com/g6834/team17/analytics-service/internal/adapters/grpc/server"
	i "gitlab.com/g6834/team17/analytics-service/internal/adapters/interfaces"
	kafkaConsumer "gitlab.com/g6834/team17/analytics-service/internal/adapters/kafka"
	ports "gitlab.com/g6834/team17/analytics-service/internal/ports/input"
	"go.uber.org/zap"
)

func NewConsumer(consumerType string, s ports.EventService, l *zap.Logger) i.MessageConsumer {
	switch consumerType {
	case "kafka":
		return kafkaConsumer.New(s, l)
	case "gRPC":
		return grpcConsumer.New(s, l)
	default:
		return kafkaConsumer.New(s, l)
	}
}
