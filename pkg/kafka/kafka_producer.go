package kafka

import (
	"gopkg.in/Shopify/sarama.v1"
	"strings"
	"time"
)

type Msg struct {
	TaskUUID  string
	EventType string
	UserUUID  string
	Timestamp time.Time
}

func (m Msg) Encode() ([]byte, error) {
	b := strings.Builder{}
	b.WriteString(`{"task_uuid": "`)
	b.WriteString(m.TaskUUID)
	b.WriteString(`", "event": "`)
	b.WriteString(m.EventType)
	b.WriteString(`", "user_uuid": "`)
	b.WriteString(m.UserUUID)
	b.WriteString(`", "timestamp": "`)
	b.WriteString(m.Timestamp.Format(time.RFC3339))
	b.WriteString(`"}`)

	return []byte(b.String()), nil
}

func (m Msg) Length() int {
	return len(m.TaskUUID) + len(m.EventType) + len(m.UserUUID) + len(m.Timestamp.Format(time.RFC3339)) + 64
}

var brokers = []string{"127.0.0.1:9095", "127.0.0.1:9096", "127.0.0.1:9097"}

func NewSyncProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.ClientID = "mok_for_task_service"
	producer, err := sarama.NewSyncProducer(brokers, config)

	//prt, offs, err := producer.SendMessage()

	return producer, err
}

func NewAsyncProducer() (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.ClientID = "mok_for_task_service"
	producer, err := sarama.NewAsyncProducer(brokers, config)

	return producer, err
}
