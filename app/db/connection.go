package db

import (
	"context"
	"time"

	"github.com/tanapoln/capgo-server/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	conn *mongo.Client
	db   *mongo.Database
)

func InitDB(ctx context.Context) error {
	ctx, cancelFn := context.WithTimeout(ctx, time.Second*10)
	defer cancelFn()

	cfg := config.Get()
	_conn, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoConnectionString))
	if err != nil {
		return err
	}
	if err := _conn.Ping(ctx, nil); err != nil {
		return err
	}

	conn = _conn

	db = conn.Database(cfg.MongoDatabase)
	return nil
}

func Disconnect() error {
	if conn == nil {
		return nil
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFn()
	return conn.Disconnect(ctx)
}

func Connection() *mongo.Client {
	return conn
}

func Database() *mongo.Database {
	return db
}
