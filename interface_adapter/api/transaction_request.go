package api

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/hdlproject/es-transaction-service/helper"
	"github.com/hdlproject/es-transaction-service/use_case/input_port"
)

type (
	topUpRequest struct {
		UserID uint   `json:"user_id"`
		Amount uint64 `json:"amount"`
	}
)

func (topUpRequest) parse(ctx *gin.Context) (request topUpRequest, err error) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return topUpRequest{}, helper.WrapError(err)
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		return topUpRequest{}, helper.WrapError(err)
	}

	return request, nil
}

func (instance topUpRequest) getUseCase() input_port.TopUpRequest {
	return input_port.TopUpRequest{
		UserID: instance.UserID,
		Amount: instance.Amount,
	}
}
