package main

import (
	"github.com/hdlproject/es-transaction-service/config"
	"github.com/hdlproject/es-transaction-service/interface_adapter/api"
	"github.com/hdlproject/es-transaction-service/interface_adapter/database"
	"github.com/hdlproject/es-transaction-service/interface_adapter/messaging"
)

func init() {
	configInstance, err := config.GetInstance()
	if err != nil {
		panic(err)
	}

	_, err = database.GetMongoDB(configInstance.EventStorage)
	if err != nil {
		panic(err)
	}

	rabbitMQClient, err := messaging.GetRabbitMQClient(configInstance.EventBus)
	if err != nil {
		panic(err)
	}

	_, err = messaging.GetRabbitMQPublisher(rabbitMQClient)
	if err != nil {
		panic(err)
	}

	_, err = messaging.GetRabbitMQSubscriber(rabbitMQClient)
	if err != nil {
		panic(err)
	}
}

func main() {
	configInstance, _ := config.GetInstance()

	api.ActivateAPI(configInstance.Port)
}
