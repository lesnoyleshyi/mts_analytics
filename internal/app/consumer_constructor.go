package app

import (
	i "gitlab.com/g6834/team17/analytics-service/internal/adapters/interfaces"
	kafkaConsumer "gitlab.com/g6834/team17/analytics-service/internal/adapters/kafka"
	ports "gitlab.com/g6834/team17/analytics-service/internal/ports/input"
	"go.uber.org/zap"
)

func NewConsumer(consumerType string, s ports.EventService, l *zap.Logger) i.MessageConsumer {
	switch consumerType {
	case "kafka":
		return kafkaConsumer.New(s, l)
	}

	return nil
}
