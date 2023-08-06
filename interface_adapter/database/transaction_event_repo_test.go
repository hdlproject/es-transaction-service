package database

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hdlproject/es-transaction-service/config"
	"github.com/hdlproject/es-transaction-service/entity"
)

func TestTransactionEventRepo_Aggregate(t *testing.T) {
	ctx := context.Background()

	mongoClient, err := GetMongoDB(config.EventStorage{
		Host:     "localhost",
		Port:     "27017",
		Username: "root",
		Password: "example",
		Name:     "es-transaction-service",
	})
	if err != nil {
		t.Fatal(err)
	}

	transactionEventRepo := NewTransactionEventRepo(mongoClient)

	defer func() {
		_ = transactionEventRepo.DeleteAll(ctx)
	}()

	userID1 := uint(1)
	userID2 := uint(2)

	data := []entity.TransactionEvent{
		{
			Type: entity.TransactionTypeTopUp,
			Params: entity.TopUp{
				UserID: userID1,
				Amount: 10000,
			},
		},
		{
			Type: entity.TransactionTypeTopUp,
			Params: entity.TopUp{
				UserID: userID1,
				Amount: 10000,
			},
		},
		{
			Type: entity.TransactionTypeTopUp,
			Params: entity.TopUp{
				UserID: userID2,
				Amount: 30000,
			},
		},
	}

	for index, item := range data {
		id, err := transactionEventRepo.Insert(ctx, item)
		if err != nil {
			t.Fatal(err)
		}
		data[index].ID = id
	}

	result, err := transactionEventRepo.GetTotalBalanceByUserID(ctx)
	if err != nil {
		t.Fatal(err)
	}

	expectedResult := map[uint]uint64{
		1: 20000,
		2: 30000,
	}

	if diff := cmp.Diff(expectedResult, result); diff != "" {
		t.Fatalf("(-want/+got)\n%s", diff)
	}
}
