package messaging

import (
	"fmt"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/hdlproject/es-transaction-service/config"
	"github.com/hdlproject/es-transaction-service/helper"
)

type (
	KafkaConsumer struct {
		consumer       *kafka.Consumer
		processingSize int
		wg             sync.WaitGroup
	}

	ConsumerCallback func(string, error)
)

var kafkaConsumer *KafkaConsumer

func GetKafkaConsumer(config config.Kafka) (*KafkaConsumer, error) {
	if kafkaConsumer == nil {
		consumer, err := newKafkaConsumer(config)
		if err != nil {
			return nil, helper.WrapError(err)
		}

		kafkaConsumer = consumer
	}

	return kafkaConsumer, nil
}

func newKafkaConsumer(config config.Kafka) (*KafkaConsumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", config.Host, config.Port),
		"group.id":          "test",
		"auto.offset.reset": "smallest",
	})
	if err != nil {
		return nil, helper.WrapError(err)
	}

	return &KafkaConsumer{
		consumer: consumer,
	}, nil
}

func (instance *KafkaConsumer) Consume(topics []string, listenDuration time.Duration, callback ConsumerCallback) error {
	err := instance.consumer.SubscribeTopics(topics, nil)
	if err != nil {
		return helper.WrapError(err)
	}

	go func() {
		var processCount int

		for {
			ev := instance.consumer.Poll(int(listenDuration))
			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("Received message %v %s\n", e.String(), e.Value)

				processCount++
				instance.wg.Add(1)
				go func() {
					defer instance.wg.Done()

					callback(string(e.Value), nil)
				}()

				if processCount >= instance.processingSize {
					instance.wg.Wait()
					processCount = 0
					instance.consumer.CommitMessage(e)
					fmt.Printf("Committed message: %v\n", e.String())
				}
			case kafka.Error:
				fmt.Printf("Received error %v\n", e)
				callback("", e)
			default:
				//fmt.Printf("Ignored message %v\n", e)
			}
		}
	}()

	return nil
}
