package messaging

import (
	"github.com/hdlproject/es-transaction-service/helper"
	"github.com/streadway/amqp"
)

type (
	RabbitMQSubscriber struct {
		client  *RabbitMQClient
		channel *amqp.Channel
	}
)

var (
	rabbitMQSubscriber *RabbitMQSubscriber
)

func GetRabbitMQSubscriber(client *RabbitMQClient) (*RabbitMQSubscriber, error) {
	if rabbitMQSubscriber == nil {
		subscriber, err := newRabbitMQSubscriber(client)
		if err != nil {
			return nil, helper.WrapError(err)
		}

		rabbitMQSubscriber = subscriber
	}

	return rabbitMQSubscriber, nil
}

func newRabbitMQSubscriber(client *RabbitMQClient) (*RabbitMQSubscriber, error) {
	ch, err := client.Connection.Channel()
	if err != nil {
		return nil, helper.WrapError(err)
	}

	return &RabbitMQSubscriber{
		channel: ch,
	}, nil
}
