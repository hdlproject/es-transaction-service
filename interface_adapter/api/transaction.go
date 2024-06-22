package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hdlproject/es-transaction-service/use_case/input_port"

	"github.com/hdlproject/es-transaction-service/config"
	"github.com/hdlproject/es-transaction-service/interface_adapter/database"
	"github.com/hdlproject/es-transaction-service/interface_adapter/messaging"
	"github.com/hdlproject/es-transaction-service/use_case/interactor"
)

type (
	TransactionController struct {
		transactionService *transactionService
	}
)

func RegisterTransactionAPI(router *gin.RouterGroup) {
	configInstance, _ := config.GetInstance()

	mongoClient, _ := database.GetMongoDB(configInstance.EventStorage)

	rabbitMQClient, _ := messaging.GetRabbitMQClient(configInstance.EventBus)
	eventPublisher, _ := messaging.GetRabbitMQPublisher(rabbitMQClient)

	transactionPublisher, err := messaging.NewTransactionPublisher(eventPublisher)
	if err != nil {
		panic(err)
	}

	transactionController := NewTransactionController(
		interactor.NewTopUpUseCase(
			database.NewTransactionEventRepo(mongoClient),
			transactionPublisher,
		),
	)

	transactionRouter := router.Group("/transaction")
	transactionRouter.POST("/topup", transactionController.TopUp)
}

func NewTransactionController(topUpUseCase input_port.TopUpUseCase) *TransactionController {
	return &TransactionController{
		transactionService: newTransactionService(
			topUpUseCase,
		),
	}
}

func (instance *TransactionController) TopUp(ctx *gin.Context) {
	request, err := topUpRequest{}.parse(ctx)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, topUpResponse{Ok: false, Message: parseRequestFailure})
		return
	}

	response, err := instance.transactionService.topUp(ctx, request)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, topUpResponse{Ok: false, Message: defaultProcessError})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
