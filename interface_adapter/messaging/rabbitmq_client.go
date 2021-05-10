package messaging

import (
	"fmt"
	"github.com/hdlproject/es-transaction-service/config"
	"github.com/hdlproject/es-transaction-service/helper"
	"github.com/streadway/amqp"
)

type (
	RabbitMQClient struct {
		Connection *amqp.Connection
	}
)

const (
	rabbitMQUriTemplate = "amqp://%s:%s@%s:%s/"
)

var rabbitMQClient *RabbitMQClient

func GetRabbitMQClient(rabbitMQConfig config.EventBus) (*RabbitMQClient, error) {
	if rabbitMQClient == nil {
		client, err := newRabbitMQClient(rabbitMQConfig)
		if err != nil {
			return nil, helper.WrapError(err)
		}

		rabbitMQClient = client
	}

	return rabbitMQClient, nil
}

func newRabbitMQClient(rabbitMQConfig config.EventBus) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(fmt.Sprintf(rabbitMQUriTemplate,
		rabbitMQConfig.Username,
		rabbitMQConfig.Password,
		rabbitMQConfig.Host,
		rabbitMQConfig.Port,
	))
	if err != nil {
		return nil, helper.WrapError(err)
	}

	return &RabbitMQClient{
		Connection: conn,
	}, nil
}
