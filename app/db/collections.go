package db

import "go.mongodb.org/mongo-driver/mongo"

type collections struct {
}

func (c collections) Bundles() *mongo.Collection {
	return Database().Collection("bundles")
}

func (c collections) Releases() *mongo.Collection {
	return Database().Collection("releases")
}

func Collections() collections {
	return collections{}
}
