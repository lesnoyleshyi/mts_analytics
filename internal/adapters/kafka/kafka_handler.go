package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"gitlab.com/g6834/team17/analytics-service/internal/config"
	ports "gitlab.com/g6834/team17/analytics-service/internal/ports/input"
	"go.uber.org/zap"
)

type MsgKafkaConsumer struct {
	client  *sarama.ConsumerGroup
	logger  *zap.Logger
	handler *consumerGroupHandler
}

func New(service ports.EventService, logger *zap.Logger) *MsgKafkaConsumer {
	var consumer MsgKafkaConsumer

	handler := new(consumerGroupHandler)
	handler.service = service
	handler.logger = logger
	handler.ready = make(chan struct{})

	consumer.handler = handler
	consumer.logger = logger

	consumer.client = nil

	return &consumer
}

func (c *MsgKafkaConsumer) StartConsume(ctx context.Context) error {
	var err error
	conf := config.GetConfig()

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V3_0_0_0
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(conf.Consumer.Brokers, conf.Consumer.ConsGroupID, cfg)
	if err != nil {
		return err
	}
	c.client = &client

	go func() {
		for {
			if err := client.Consume(ctx, conf.Consumer.Topics, c.handler); err != nil {
				c.logger.Warn("error consuming message from kafka", zap.Error(err))
				return
			}
			if ctx.Err() != nil {
				return
			}
			c.handler.ready = make(chan struct{})
		}
	}()

	<-c.handler.ready
	c.logger.Info("Kafka consumer starts successfully")

	return nil
}

func (c *MsgKafkaConsumer) StopConsume(ctx context.Context) error {
	if c.client == nil {
		return nil
	}
	return (*c.client).Close()
}
