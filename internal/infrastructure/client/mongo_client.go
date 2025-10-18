package client

import (
	"context"
	"log"

	"github.com/zhunismp/intent-products-api/internal/infrastructure/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func NewMongoClient(ctx context.Context, cfg *config.Config) *mongo.Client {
	client, err := mongo.Connect(options.Client().ApplyURI(cfg.Mongo.URI))

	if err != nil {
		log.Fatal("❌ failed to connect to MongoDB: %w", err)
	}

	// ping to check connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("❌ failed to ping MongoDB: %w", err)
	}

	return client
}