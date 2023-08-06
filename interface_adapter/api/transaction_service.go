package api

import (
	"context"

	"github.com/hdlproject/es-transaction-service/helper"
	"github.com/hdlproject/es-transaction-service/use_case/interactor"
)

type (
	transactionService struct {
		topUpUseCase *interactor.TopUp
	}
)

func newTransactionService(topUpUseCase *interactor.TopUp) *transactionService {
	return &transactionService{
		topUpUseCase: topUpUseCase,
	}
}

func (instance *transactionService) topUp(ctx context.Context, request topUpRequest) (topUpResponse, error) {
	useCaseRequest := request.getUseCase()

	useCaseResponse, err := instance.topUpUseCase.TopUp(ctx, useCaseRequest)
	if err != nil {
		return topUpResponse{}, helper.WrapError(err)
	}

	return topUpResponse{}.fromUseCase(useCaseResponse), nil
}
