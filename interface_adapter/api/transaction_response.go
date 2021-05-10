package api

import "github.com/hdlproject/es-transaction-service/use_case/interactor"

type (
	topUpResponse struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	}
)

func (topUpResponse) fromUseCase(response interactor.TopUpResponse) topUpResponse {
	return topUpResponse{
		Ok:      response.Ok,
		Message: response.Message,
	}
}
