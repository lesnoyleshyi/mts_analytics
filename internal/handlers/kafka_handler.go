package handlers

import (
	"encoding/json"
	"github.com/jackc/pgconn"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/Shopify/sarama.v1"
	"mts_analytics/internal/domain"
	"time"
)

var brokers = []string{"kafka-1:9092", "kafka-2:9092", "kafka-3:9092"}

const consumerGroupID string = `analytics`

type KafkaHandler struct {
	Client        *sarama.Client
	ConsumerGroup *sarama.ConsumerGroup
	Consumer      eventConsumer
}

type eventConsumer struct {
	service service
	id      int
}

func (c eventConsumer) Setup(s sarama.ConsumerGroupSession) error {
	log.WithFields(log.Fields{
		"app_name":    app_name,
		"host_ip":     host_ip,
		"logger_name": "kafka_handler.Setup",
	}).Debugf("Setup() hook is called. Consumer #%d with cluster member ID %s will consume messages",
		c.id, s.MemberID())

	return nil
}

func (c eventConsumer) ConsumeClaim(s sarama.ConsumerGroupSession, cl sarama.ConsumerGroupClaim) error {
	log.WithFields(log.Fields{
		"app_name":    app_name,
		"host_ip":     host_ip,
		"logger_name": "kafka_handler.ConsumeClaim",
	}).Debug("A claim was assigned. Start Processing messages")
	for msg := range cl.Messages() {
		var event domain.Event
		log.WithFields(log.Fields{
			"app_name":    app_name,
			"host_ip":     host_ip,
			"logger_name": "kafka_handler.ConsumeClaim",
		}).Debugf("message got: %s. Partition: %d, offset: %d, topic: %s, ts: %s",
			msg.Value, msg.Partition, msg.Offset, msg.Topic, msg.Timestamp.Format(time.RFC3339))
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.WithFields(log.Fields{
				"app_name":    app_name,
				"host_ip":     host_ip,
				"logger_name": "kafka_handler.ConsumeClaim.Unmarshal",
			}).Warn(err)
		}
		err := c.service.Save(event)
		if err != nil {
			log.WithFields(log.Fields{
				"app_name":    app_name,
				"host_ip":     host_ip,
				"logger_name": "kafka_handler.ConsumeClaim.Save",
			}).Warn(err)
			pgErr := &pgconn.PgError{}
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				log.Debug("Duplicates are not allowed!")
				s.MarkMessage(msg, "analytics")
			}
		} else {
			s.MarkMessage(msg, "analytics")
			log.WithFields(log.Fields{
				"app_name":    app_name,
				"host_ip":     host_ip,
				"logger_name": "kafka_handler.ConsumeClaim.Save",
			}).Debug("Message processed successfully")
		}
	}
	return nil
}

func (c eventConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	log.WithFields(log.Fields{
		"app_name":    app_name,
		"host_ip":     host_ip,
		"logger_name": "kafka_handler.Cleanup",
	}).Debug("ConsumeClaim() loop has exited. Parent context is cancelled or a server-side rebalance cycle is initiated")
	return nil
}

func NewKafkaHandler(service service) (*KafkaHandler, error) {
	conf := sarama.NewConfig()
	conf.Version = sarama.V2_0_0_0
	conf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewClient(brokers, conf)
	if err != nil {
		return nil, err
	}
	consGroup, err := sarama.NewConsumerGroupFromClient(consumerGroupID, client)
	if err != nil {
		return nil, err
	}
	handler := KafkaHandler{
		Client:        &client,
		ConsumerGroup: &consGroup,
		Consumer:      eventConsumer{id: 1, service: service},
	}
	return &handler, nil
}
