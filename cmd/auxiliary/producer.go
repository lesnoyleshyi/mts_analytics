package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/Shopify/sarama.v1"
	"mts_analytics/pkg/kafka"
	"time"
)

const app_name = `mok_task_service`
const host_ip = `lol_hz`

const TaskUUID1 = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01"
const TaskUUID2 = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02"
const TaskUUID3 = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03"
const TaskUUID4 = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a04"

const UserUUID1 = "0a75fc43-584a-46a6-8e3c-686edbea43f8"
const UserUUID2 = "97364b12-0dcf-487f-8a03-9051d5c1d8ba"
const UserUUID3 = "c158314b-a043-4e25-8bd0-4de99bfd4601"
const UserUUID4 = "469481c6-0376-4a9f-a2e0-8d72b99bde67"
const UserUUID5 = "eafe0883-5aa0-4f8a-850f-ed002b56c131"

const e_created = "created"
const e_sent_to = "sent_to"
const e_signed_by_user = "approved_by"
const e_rejected_by_user = "rejected_by"
const e_signed = "signed"
const e_sent = "sent"

var messages = []kafka.Msg{
	{
		TaskUUID:  TaskUUID1,
		EventType: e_created,
		UserUUID:  UserUUID1,
		Timestamp: time.Date(2022, 01, 10, 11, 30, 0, 0, time.UTC),
	},
	{
		TaskUUID:  TaskUUID1,
		EventType: e_sent_to,
		UserUUID:  UserUUID2,
		Timestamp: time.Date(2022, 01, 10, 11, 30, 10, 0, time.UTC),
	},
	{
		TaskUUID:  TaskUUID1,
		EventType: e_signed_by_user,
		UserUUID:  UserUUID2,
		Timestamp: time.Date(2022, 01, 10, 11, 43, 0, 0, time.UTC),
	},
	{
		TaskUUID:  TaskUUID1,
		EventType: e_sent_to,
		UserUUID:  UserUUID3,
		Timestamp: time.Date(2022, 01, 10, 11, 43, 5, 0, time.UTC),
	},
	{
		TaskUUID:  TaskUUID1,
		EventType: e_signed_by_user,
		UserUUID:  UserUUID3,
		Timestamp: time.Date(2022, 01, 10, 15, 10, 5, 0, time.UTC),
	},
	{
		TaskUUID:  TaskUUID1,
		EventType: e_signed,
		UserUUID:  "",
		Timestamp: time.Date(2022, 01, 10, 15, 10, 10, 0, time.UTC),
	},
	{
		TaskUUID:  TaskUUID1,
		EventType: e_sent,
		UserUUID:  "",
		Timestamp: time.Date(2022, 01, 10, 15, 10, 12, 0, time.UTC),
	},
}

func main() {
	log.SetLevel(log.DebugLevel)
	p, err := kafka.NewAsyncProducer()
	if err != nil {
		log.WithFields(log.Fields{
			"app_name":    app_name,
			"host_ip":     host_ip,
			"logger_name": "main.NewSyncProducer",
		}).Errorf("can't create async producer: %s", err)
	}
	log.WithFields(log.Fields{
		"app_name":    app_name,
		"host_ip":     host_ip,
		"logger_name": "main.NewSyncProducer",
	}).Debug("Async producer created!")
	errCh := p.Errors()
	sucCh := p.Successes()
	done := make(chan struct{})
	go readKafkaResponse(errCh, sucCh, done)

	msgOutputCh := p.Input()

	go sendMsgs(msgOutputCh, done)

	<-done
	p.AsyncClose()
}

func readKafkaResponse(
	errCh <-chan *sarama.ProducerError,
	successCh <-chan *sarama.ProducerMessage,
	done <-chan struct{}) {
	logErr := func(err error) {
		log.WithFields(log.Fields{
			"app_name":    app_name,
			"host_ip":     host_ip,
			"logger_name": "main.readKafkaResponse",
		}).Error(err)
	}
	logSuccess := func(msg *sarama.ProducerMessage) {
		val, err := msg.Value.Encode()
		if err != nil {
			log.WithFields(log.Fields{
				"app_name":    app_name,
				"host_ip":     host_ip,
				"logger_name": "main.readKafkaResponse",
			}).Errorf("Message sent successful, but encoding suck: %s", err)
		}
		log.WithFields(log.Fields{
			"app_name":    app_name,
			"host_ip":     host_ip,
			"logger_name": "main.NewSyncProducer",
		}).Debug(string(val))
	}

	for {
		select {
		case err := <-errCh:
			{
				logErr(err)
			}
		case msg := <-successCh:
			{
				logSuccess(msg)
			}
		case <-done:
			{
				return
			}
		}
	}
}

func sendMsgs(outputCh chan<- *sarama.ProducerMessage, done chan<- struct{}) {
	for _, m := range messages {
		sm := sarama.ProducerMessage{
			Topic:     "task",
			Key:       sarama.ByteEncoder("my_random_key"),
			Value:     m,
			Timestamp: time.Now(),
		}
		outputCh <- &sm
		time.Sleep(time.Second * 5)
	}
	log.WithFields(log.Fields{
		"app_name":    app_name,
		"host_ip":     host_ip,
		"logger_name": "main.sendMsgs",
	}).Debug("All messages are sent")
	close(done)
}
