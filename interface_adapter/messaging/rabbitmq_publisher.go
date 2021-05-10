package messaging

import (
	"github.com/hdlproject/es-transaction-service/helper"
	"github.com/streadway/amqp"
)

type (
	RabbitMQPublisher struct {
		client  *RabbitMQClient
		channel *amqp.Channel
	}
)

var (
	rabbitMQPublisher *RabbitMQPublisher
)

func GetRabbitMQPublisher(client *RabbitMQClient) (*RabbitMQPublisher, error) {
	if rabbitMQPublisher == nil {
		publisher, err := newRabbitMQPublisher(client)
		if err != nil {
			return nil, helper.WrapError(err)
		}

		rabbitMQPublisher = publisher
	}

	return rabbitMQPublisher, nil
}

func newRabbitMQPublisher(client *RabbitMQClient) (*RabbitMQPublisher, error) {
	ch, err := client.Connection.Channel()
	if err != nil {
		return nil, helper.WrapError(err)
	}

	return &RabbitMQPublisher{
		channel: ch,
	}, nil
}
