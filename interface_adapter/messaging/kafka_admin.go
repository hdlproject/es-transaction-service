package messaging

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/hdlproject/es-transaction-service/helper"
)

type (
	KafkaAdmin struct {
		admin *kafka.AdminClient
	}
)

var kafkaAdmin *KafkaAdmin

func GetKafkaAdminFromProducer(producer *kafka.Producer) (*KafkaAdmin, error) {
	if kafkaAdmin == nil {
		admin, err := newKafkaAdmin(producer)
		if err != nil {
			return nil, helper.WrapError(err)
		}

		kafkaAdmin = admin
	}

	return kafkaAdmin, nil
}

func newKafkaAdmin(producer *kafka.Producer) (*KafkaAdmin, error) {
	admin, err := kafka.NewAdminClientFromProducer(producer)
	if err != nil {
		return nil, helper.WrapError(err)
	}

	return &KafkaAdmin{
		admin: admin,
	}, nil
}

func (instance *KafkaAdmin) CreateTopic(ctx context.Context, topic string) error {
	result, err := instance.admin.CreateTopics(ctx, []kafka.TopicSpecification{
		{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 1,
			Config: map[string]string{
				"confluent.value.schema.validation": "false",
			},
		},
	})
	if err != nil {
		return helper.WrapError(err)
	}

	for _, item := range result {
		if item.Error.Code() != kafka.ErrNoError && item.Error.Code() != kafka.ErrTopicAlreadyExists {
			return helper.WrapError(item.Error)
		}
	}

	return nil
}
