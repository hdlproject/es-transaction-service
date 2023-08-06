package output_port

import (
	"context"

	"github.com/hdlproject/es-transaction-service/entity"
)

type (
	TransactionEventRepo interface {
		Insert(ctx context.Context, event entity.TransactionEvent) (string, error)
		GetTotalBalanceByUserID(ctx context.Context) (map[uint]uint64, error)
		DeleteAll(ctx context.Context) error
	}
)
