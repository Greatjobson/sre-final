package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func EnsureMongoIndexes(ctx context.Context, orderCol, orderItemCol *mongo.Collection) error {
	if _, err := orderCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{"userId", 1}, {"createdAt", -1}}},
	}); err != nil {
		return err
	}

	if _, err := orderItemCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{"orderId", 1}, {"productId", 1}}},
	}); err != nil {
		return err
	}
	return nil
}
