package messaging

import (
	"encoding/json"
	"github.com/hdlproject/es-transaction-service/entity"
	"github.com/hdlproject/es-transaction-service/helper"
	"github.com/hdlproject/es-transaction-service/use_case/output_port"
	"github.com/streadway/amqp"
)

type (
	transactionEvent struct {
		ID     string      `json:"id"`
		Type   string      `json:"type"`
		Params interface{} `json:"params"`
	}

	topUp struct {
		UserID uint   `json:"user_id"`
		Amount uint64 `json:"amount"`
	}

	transactionPublisher struct {
		publisher    *RabbitMQPublisher
		exchangeName string
	}
)

func NewTransactionPublisher(publisher *RabbitMQPublisher) (output_port.TransactionPublisher, error) {
	exchangeName := "transactions_direct"

	err := publisher.channel.ExchangeDeclare(
		exchangeName,
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
		publisher:    publisher,
		exchangeName: exchangeName,
	}, nil
}

func (instance *transactionPublisher) Publish(event entity.TransactionEvent) error {
	publishedEvent := transactionEvent{}.fromEntity(event)

	body, err := json.Marshal(publishedEvent)
	if err != nil {
		return helper.WrapError(err)
	}

	err = instance.publisher.channel.Publish(
		instance.exchangeName,
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

func (instance transactionEvent) fromEntity(transactionEventEntity entity.TransactionEvent) transactionEvent {
	var params interface{}
	switch v := transactionEventEntity.Params.(type) {
	case entity.TopUp:
		params = topUp{}.fromEntity(v)
	}

	return transactionEvent{
		ID:     transactionEventEntity.ID,
		Type:   string(transactionEventEntity.Type),
		Params: params,
	}
}

func (instance topUp) fromEntity(topUpEntity entity.TopUp) topUp {
	return topUp{
		UserID: topUpEntity.UserID,
		Amount: topUpEntity.Amount,
	}
}
