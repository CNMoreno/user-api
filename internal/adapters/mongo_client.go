package adapters

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoClient(uri, dbName string) (*MongoClient, error) {
	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

func (mc *MongoClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return mc.client.Disconnect(ctx)
}

func (mc *MongoClient) GetDatabase() *mongo.Database {
	return mc.db
}
