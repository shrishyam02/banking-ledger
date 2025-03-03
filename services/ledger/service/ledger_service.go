package service

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type LedgerService struct {
	collection *mongo.Collection
}

func NewLedgerService(db *mongo.Database) *LedgerService {
	return &LedgerService{
		collection: db.Collection("transactions"),
	}
}

func (s *LedgerService) HandleMessage(ctx context.Context, msg kafka.Message) error {
	var transaction map[string]interface{}
	if err := json.Unmarshal(msg.Value, &transaction); err != nil {
		return err
	}

	_, err := s.collection.InsertOne(ctx, transaction)
	if err != nil {
		return err
	}

	return nil
}

func (s *LedgerService) GetAccountTransactionHistory(ctx context.Context, accountID string) ([]map[string]interface{}, error) {
	filter := map[string]interface{}{
		"accountId": accountID,
	}
	cursor, err := s.collection.Find(ctx, filter, options.Find().SetSort(map[string]interface{}{"acceptedAt": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []map[string]interface{}
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *LedgerService) GetTransactionHistory(ctx context.Context, id string) ([]map[string]interface{}, error) {
	filter := map[string]interface{}{
		"id": id,
	}
	cursor, err := s.collection.Find(ctx, filter, options.Find().SetSort(map[string]interface{}{"acceptedAt": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []map[string]interface{}
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}
