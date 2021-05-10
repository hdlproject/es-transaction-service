package output_port

import "github.com/hdlproject/es-transaction-service/entity"

type (
	TransactionPublisher interface {
		Publish(event entity.TransactionEvent) error
	}
)
