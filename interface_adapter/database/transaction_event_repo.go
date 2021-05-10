package database

import (
	"github.com/hdlproject/es-transaction-service/entity"
	"github.com/hdlproject/es-transaction-service/helper"
	"github.com/hdlproject/es-transaction-service/use_case/output_port"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	TransactionEvent struct {
		ID     primitive.ObjectID `bson:"_id,omitempty"`
		Type   string             `bson:"type,omitempty"`
		Params interface{}        `bson:"params,omitempty"`
	}

	TopUp struct {
		UserID uint   `bson:"user_id,omitempty"`
		Amount uint64 `bson:"amount,omitempty"`
	}
)

type (
	transactionEventRepo struct {
		mongoClient                *MongoClient
		transactionEventCollection *mongo.Collection
	}
)

func NewTransactionEventRepo(mongoClient *MongoClient) output_port.TransactionEventRepo {
	return &transactionEventRepo{
		mongoClient:                mongoClient,
		transactionEventCollection: mongoClient.DB.Collection("transaction_events"),
	}
}

func (instance *transactionEventRepo) Insert(event entity.TransactionEvent) (string, error) {
	data, _ := TransactionEvent{}.getData(event)
	result, err := instance.transactionEventCollection.InsertOne(mongoClient.Context, data)
	if err != nil {
		return "", helper.WrapError(err)
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (TransactionEvent) getData(transactionEventEntity entity.TransactionEvent) (TransactionEvent, error) {
	var id primitive.ObjectID
	var err error
	if transactionEventEntity.ID != "" {
		id, err = primitive.ObjectIDFromHex(transactionEventEntity.ID)
		if err != nil {
			return TransactionEvent{}, helper.WrapError(err)
		}
	}

	var params interface{}
	switch v := transactionEventEntity.Params.(type) {
	case entity.TopUp:
		params = TopUp{}.getData(v)
	}

	return TransactionEvent{
		ID:     id,
		Type:   string(transactionEventEntity.Type),
		Params: params,
	}, nil
}

func (TopUp) getData(topUpEntity entity.TopUp) TopUp {
	return TopUp{
		UserID: topUpEntity.UserID,
		Amount: topUpEntity.Amount,
	}
}
