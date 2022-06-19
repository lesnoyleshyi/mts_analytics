package kafka

import (
	"context"
	log "github.com/sirupsen/logrus"
	"gopkg.in/Shopify/sarama.v1"
	"time"
)

const appName = `mok_task_service`
const hostIP = `lol_hz`

type Consumer struct {
	ID int
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	log.WithFields(log.Fields{
		"appName":     appName,
		"hostIP":      hostIP,
		"logger_name": "consumer",
	}).Debug("Setup() hook is called. Hello!")

	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	log.WithFields(log.Fields{
		"appName":     appName,
		"hostIP":      hostIP,
		"logger_name": "consumer",
	}).Debug("ConsumeClaim() loops have exited. Hello from Cleanup()")

	return nil
}

func (c *Consumer) ConsumeClaim(ses sarama.ConsumerGroupSession, cl sarama.ConsumerGroupClaim) error {
	log.WithFields(log.Fields{
		"appName":     appName,
		"hostIP":      hostIP,
		"logger_name": "consumer",
	}).Debug("Some claim was assigned. Hello from ConsumeClaim()")

	for msg := range cl.Messages() {
		log.WithFields(log.Fields{
			"appName":     appName,
			"hostIP":      hostIP,
			"logger_name": "consumer",
		}).Debugf("message got: %s. Partition: %d, offset: %d, topic: %s, ts: %s",
			msg.Value, msg.Partition, msg.Offset, msg.Topic, msg.Timestamp.Format(time.RFC3339))
		ses.MarkMessage(msg, "analytics")
	}

	return nil
}

func subscribe(ctx context.Context, topic string, consGr sarama.ConsumerGroup) error {
	consumer := Consumer{ID: 1}

	go func() {
		for {
			if err := consGr.Consume(ctx, []string{topic}, &consumer); err != nil {
				log.WithField("a", "che").Debug(err)
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	return nil
}

func StartConsuming(ctx context.Context) error {
	cfg := sarama.NewConfig()

	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	consGr, err := sarama.NewConsumerGroup(brokers, "analytics", cfg)
	if err != nil {
		log.WithField("a", "che").Warn(err)
	}

	return subscribe(ctx, "lol", consGr)
}
