package messaging

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/hdlproject/es-transaction-service/config"
	"github.com/hdlproject/es-transaction-service/helper"
)

type (
	KafkaProducer struct {
		producer *kafka.Producer
	}

	ProducerCallback func(string, error)
)

var kafkaProducer *KafkaProducer

func GetKafkaProducer(config config.Kafka) (*KafkaProducer, error) {
	if kafkaProducer == nil {
		producer, err := newKafkaProducer(config)
		if err != nil {
			return nil, helper.WrapError(err)
		}

		kafkaProducer = producer
	}

	return kafkaProducer, nil
}

func newKafkaProducer(config config.Kafka) (*KafkaProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", config.Host, config.Port),
		"client.id":         "localhost",
		"acks":              "all",
	})
	if err != nil {
		return nil, helper.WrapError(err)
	}

	return &KafkaProducer{
		producer: producer,
	}, nil
}

func (instance *KafkaProducer) produce(topic, message string) (chan kafka.Event, error) {
	deliveryChan := make(chan kafka.Event, 10000)

	err := instance.producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(message),
		},
		deliveryChan,
	)
	if err != nil {
		return nil, helper.WrapError(err)
	}

	return deliveryChan, nil
}

func (instance *KafkaProducer) Produce(topic, message string, callback ProducerCallback) error {
	deliveryChan, err := instance.produce(topic, message)
	if err != nil {
		return helper.WrapError(err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		callback("", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message %v %s\n", m.String(), m.Value)
		callback(string(m.Value), nil)
	}
	close(deliveryChan)

	return nil
}

func (instance *KafkaProducer) AsyncProduce(topic, message string, callback ProducerCallback) error {
	_, err := instance.produce(topic, message)
	if err != nil {
		return helper.WrapError(err)
	}

	go func() {
		for e := range instance.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
					callback("", ev.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message %v %s\n", ev.String(), ev.Value)
					callback(string(ev.Value), nil)
				}
			}
		}
	}()

	return nil
}
