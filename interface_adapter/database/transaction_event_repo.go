package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/hdlproject/es-transaction-service/entity"
	"github.com/hdlproject/es-transaction-service/helper"
	"github.com/hdlproject/es-transaction-service/use_case/output_port"
)

type (
	transactionEvent struct {
		ID     primitive.ObjectID `bson:"_id,omitempty"`
		Type   string             `bson:"type,omitempty"`
		Params interface{}        `bson:"params,omitempty"`
	}

	topUp struct {
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

func (instance *transactionEventRepo) Insert(ctx context.Context, event entity.TransactionEvent) (string, error) {
	data, err := transactionEvent{}.getData(event)
	if err != nil {
		return "", helper.WrapError(err)
	}

	result, err := instance.transactionEventCollection.InsertOne(ctx, data)
	if err != nil {
		return "", helper.WrapError(err)
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (instance *transactionEventRepo) GetTotalBalanceByUserID(ctx context.Context) (map[uint]uint64, error) {
	cursor, err := instance.transactionEventCollection.Aggregate(ctx, mongo.Pipeline{
		{
			{"$group", bson.D{
				{"_id", "$params.user_id"},
				{"total_balance", bson.D{
					{"$sum", "$params.amount"},
				}},
			}},
		},
	})
	if err != nil {
		return nil, helper.WrapError(err)
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, helper.WrapError(err)
	}

	aggregatedValue := make(map[uint]uint64)
	for _, result := range results {
		aggregatedValue[uint(result["_id"].(int64))] = uint64(result["total_balance"].(int64))
	}

	return aggregatedValue, nil
}

func (instance *transactionEventRepo) DeleteAll(ctx context.Context) error {
	_, err := instance.transactionEventCollection.DeleteMany(ctx, bson.D{})
	if err != nil {
		return helper.WrapError(err)
	}

	return nil
}

func (transactionEvent) getData(transactionEventEntity entity.TransactionEvent) (transactionEvent, error) {
	var params interface{}
	switch v := transactionEventEntity.Params.(type) {
	case entity.TopUp:
		params = topUp{}.getData(v)
	}

	return transactionEvent{
		Type:   string(transactionEventEntity.Type),
		Params: params,
	}, nil
}

func (topUp) getData(topUpEntity entity.TopUp) topUp {
	return topUp{
		UserID: topUpEntity.UserID,
		Amount: topUpEntity.Amount,
	}
}
