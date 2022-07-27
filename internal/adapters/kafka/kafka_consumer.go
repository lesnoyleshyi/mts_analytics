package kafka

import (
	"errors"
	s "github.com/Shopify/sarama"
	"github.com/mailru/easyjson"
	"gitlab.com/g6834/team17/analytics-service/internal/domain/entity"
	ports "gitlab.com/g6834/team17/analytics-service/internal/ports/input"
	"go.uber.org/zap"
)

// it's a stub! should be retuned from service
var ErrAlreadyExists = errors.New("event is already in database")

type consumerGroupHandler struct {
	service ports.EventService
	logger  *zap.Logger
	ready   chan struct{}
}

func (c *consumerGroupHandler) Setup(s.ConsumerGroupSession) error {
	c.logger.Debug("new kafka session started")
	close(c.ready)

	return nil
}

func (c *consumerGroupHandler) Cleanup(s s.ConsumerGroupSession) error {
	c.logger.Debug("kafka session finished", zap.String("memberID", s.MemberID()))

	return nil
}

func (c *consumerGroupHandler) ConsumeClaim(s s.ConsumerGroupSession, cl s.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-cl.Messages():
			c.logger.Debug("message consumed")

			event := new(entity.Event)
			err := easyjson.Unmarshal(message.Value, event)
			if err != nil {
				c.logger.Warn("unable unmarshall msg from kafka to Event struct",
					zap.Error(err))
				s.MarkMessage(message, "Unmarshall error")
				break
			}

			err = c.service.Save(s.Context(), *event)
			if err != nil {
				if !errors.Is(err, ErrAlreadyExists) {
					c.logger.Warn("unable save msg from Kafka", zap.Error(err))
					s.MarkMessage(message, "duplicated message")
					break
				}
				s.MarkMessage(message, "domain error")
			}
			s.MarkMessage(message, "")
		case <-s.Context().Done():
			return nil
		}
	}
}
