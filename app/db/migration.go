package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RunMigration() error {
	ctx := context.Background()

	_, err := Collections().Releases().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "version_name", Value: 1},
			{Key: "platform", Value: 1},
			{Key: "app_id", Value: 1},
			{Key: "version_code", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	return nil
}
