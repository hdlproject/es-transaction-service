package api

import (
	"github.com/hdlproject/es-transaction-service/use_case/input_port"
)

type (
	topUpResponse struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	}
)

func (topUpResponse) fromUseCase(response input_port.TopUpResponse) topUpResponse {
	return topUpResponse{
		Ok:      response.Ok,
		Message: response.Message,
	}
}
