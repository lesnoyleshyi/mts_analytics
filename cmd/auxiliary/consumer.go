package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"gopkg.in/Shopify/sarama.v1"
	"mts_analytics/pkg/kafka"
	"sync"
)

var brokers = []string{"127.0.0.1:9095", "127.0.0.1:9096", "127.0.0.1:9097"}
var topics = []string{`task`}

const app_name = `mok_task_service`
const host_ip = `lol_hz`

func main() {
	log.SetLevel(log.DebugLevel)
	consumer := kafka.Consumer{Id: 1488}
	ctx := context.Background()

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_0_0_0
	log.WithFields(log.Fields{
		"app_name":    app_name,
		"host_ip":     host_ip,
		"logger_name": "main.sarama.NewConsumerGroup",
	}).Debugf("Kafka version: %s", cfg.Version.String())
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	CG, err := sarama.NewConsumerGroup(brokers, "analytics", cfg)
	if err != nil {
		log.WithFields(log.Fields{
			"app_name":    app_name,
			"host_ip":     host_ip,
			"logger_name": "main.sarama.NewConsumerGroup",
		}).Fatalf("unable create consumer group 'analytics': %s", err)
	}
	log.WithFields(log.Fields{
		"app_name":    app_name,
		"host_ip":     host_ip,
		"logger_name": "main.sarama.NewConsumerGroup",
	}).Debug("New consumer group 'analytics' created")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			if err := CG.Consume(ctx, topics, &consumer); err != nil {
				log.WithFields(log.Fields{
					"app_name":    app_name,
					"host_ip":     host_ip,
					"logger_name": "main.CG.Consume_loop",
				}).Errorf("error consuming: %s", err)
			}
			if err := ctx.Err(); err != nil {
				log.WithFields(log.Fields{
					"app_name":    app_name,
					"host_ip":     host_ip,
					"logger_name": "main.CG.Consume_loop",
				}).Errorf("context catch error: %s. Stop consuming", err)
				break
			}
		}
		wg.Done()
	}()
	wg.Wait()
	if err := CG.Close(); err != nil {
		log.WithFields(log.Fields{
			"app_name":    app_name,
			"host_ip":     host_ip,
			"logger_name": "main.CG.Close",
		}).Errorf("error closing consumer group: %s", err)
	}
}
