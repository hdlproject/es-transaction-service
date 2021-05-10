package output_port

import "github.com/hdlproject/es-transaction-service/entity"

type (
	TransactionEventRepo interface {
		Insert(event entity.TransactionEvent) (string, error)
	}
)
