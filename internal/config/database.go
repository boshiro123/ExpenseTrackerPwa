package config

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(ctx context.Context, cfg Config) *mongo.Database {
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	return client.Database(cfg.DBName)
}
