package interactor

import (
	"context"

	"github.com/hdlproject/es-transaction-service/entity"
	"github.com/hdlproject/es-transaction-service/helper"
	"github.com/hdlproject/es-transaction-service/use_case/input_port"
	"github.com/hdlproject/es-transaction-service/use_case/output_port"
)

type (
	TopUp struct {
		transactionEventRepo output_port.TransactionEventRepo
		transactionPublisher output_port.TransactionPublisher
	}
)

func NewTopUpUseCase(transactionEventRepo output_port.TransactionEventRepo,
	transactionPublisher output_port.TransactionPublisher) *TopUp {

	return &TopUp{
		transactionEventRepo: transactionEventRepo,
		transactionPublisher: transactionPublisher,
	}
}

func (instance *TopUp) TopUp(ctx context.Context, request input_port.TopUpRequest) (response input_port.TopUpResponse, err error) {
	transactionEvent := entity.TransactionEvent{
		Type: entity.TransactionTypeTopUp,
		Params: entity.TopUp{
			UserID: request.UserID,
			Amount: request.Amount,
		},
	}
	eventID, err := instance.transactionEventRepo.Insert(ctx, transactionEvent)
	if err != nil {
		return input_port.TopUpResponse{}, helper.WrapError(err)
	}
	transactionEvent.ID = eventID

	err = instance.transactionPublisher.Publish(transactionEvent)
	if err != nil {
		return input_port.TopUpResponse{}, helper.WrapError(err)
	}

	return input_port.TopUpResponse{
		Ok:      true,
		Message: success,
	}, nil
}
