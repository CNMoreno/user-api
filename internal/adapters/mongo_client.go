package adapters

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TIMEOUT time for timeout connection.
const (
	TIMEOUT = 10
)

// MongoClient encapsulates the mongo connection.
type MongoClient struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoClient create Client for Mongo.
func NewMongoClient(uri, dbName string) (*MongoClient, error) {
	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	return &MongoClient{
		client: client,
		db:     client.Database(dbName),
	}, nil
}

// Close disconnect connection to Mongo.
func (mc *MongoClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT*time.Second)
	defer cancel()

	return mc.client.Disconnect(ctx)
}

// GetDatabase handles to get DB in connection to Mongo.
func (mc *MongoClient) GetDatabase() *mongo.Database {
	return mc.db
}
