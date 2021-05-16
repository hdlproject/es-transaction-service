package interactor

import (
	"github.com/hdlproject/es-transaction-service/entity"
	"github.com/hdlproject/es-transaction-service/helper"
	"github.com/hdlproject/es-transaction-service/use_case/output_port"
)

type (
	TopUpRequest struct {
		UserID uint
		Amount uint64
	}

	TopUpResponse struct {
		Ok      bool
		Message string
	}

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

func (instance *TopUp) TopUp(request TopUpRequest) (response TopUpResponse, err error) {
	transactionEvent := entity.TransactionEvent{
		Type: entity.TransactionTypeTopUp,
		Params: entity.TopUp{
			UserID: request.UserID,
			Amount: request.Amount,
		},
	}
	eventID, err := instance.transactionEventRepo.Insert(transactionEvent)
	if err != nil {
		return TopUpResponse{}, helper.WrapError(err)
	}
	transactionEvent.ID = eventID

	err = instance.transactionPublisher.Publish(transactionEvent)
	if err != nil {
		return TopUpResponse{}, helper.WrapError(err)
	}

	return TopUpResponse{
		Ok:      true,
		Message: success,
	}, nil
}
