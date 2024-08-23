package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RunMigration() error {
	ctx := context.Background()

	_, err := Collections().Releases().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "version_name", Value: 1},
			{Key: "platform", Value: 1},
			{Key: "app_id", Value: 1},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
