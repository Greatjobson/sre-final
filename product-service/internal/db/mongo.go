package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBClient struct {
	client   *mongo.Client
	database string
}

func NewMongoDBClient(ctx context.Context, uri string) (*MongoDBClient, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return &MongoDBClient{client: client, database: "clothes_store"}, nil
}

func (c *MongoDBClient) Collection(name string) *mongo.Collection {
	return c.client.Database(c.database).Collection(name)
}

func (c *MongoDBClient) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}
