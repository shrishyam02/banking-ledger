package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ConnectMongo establishes a connection to the MongoDB database.
func ConnectMongo(connStr string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(connStr)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}

// func EnsureIndexes(ctx context.Context, collection *mongo.Collection, models []interface{}) error {

// 	for _, model := range models {
// 		indexes, ok := model.(Indexable)
// 		if !ok {
// 			continue
// 		}

// 		for _, index := range indexes.Indexes() {
// 			opts := options.Index()
// 			if index.Unique {
// 				opts.SetUnique(true)
// 			}
// 			if index.Sparse {
// 				opts.SetSparse(true)
// 			}

// 			keys := bson.M{}
// 			for k, v := range index.Keys {
// 				keys[k] = v
// 			}

// 			_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
// 				Keys:    keys,
// 				Options: opts,
// 			})
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// type Indexable interface {
// 	Indexes() []Index
// }

// type Index struct {
// 	Keys   map[string]interface{}
// 	Unique bool
// 	Sparse bool
// }
