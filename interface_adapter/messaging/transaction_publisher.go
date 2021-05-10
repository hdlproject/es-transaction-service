package messaging

import (
	"encoding/json"
	"github.com/hdlproject/es-transaction-service/entity"
	"github.com/hdlproject/es-transaction-service/helper"
	"github.com/hdlproject/es-transaction-service/use_case/output_port"
	"github.com/streadway/amqp"
)

type (
	transactionPublisher struct {
		publisher *RabbitMQPublisher
	}
)

func NewTransactionPublisher(publisher *RabbitMQPublisher) (output_port.TransactionPublisher, error) {
	err := publisher.channel.ExchangeDeclare(
		"transactions_topic",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, helper.WrapError(err)
	}

	return &transactionPublisher{
		publisher: publisher,
	}, nil
}

func (instance *transactionPublisher) Publish(event entity.TransactionEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return helper.WrapError(err)
	}

	err = instance.publisher.channel.Publish(
		"transactions_topic",
		string(event.Type),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return helper.WrapError(err)
	}

	return nil
}
