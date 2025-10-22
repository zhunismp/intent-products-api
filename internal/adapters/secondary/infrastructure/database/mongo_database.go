package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func NewMongoDatabase(ctx context.Context, host, user, password, dbname, port, sslmode, TimeZone string) (*mongo.Database, func(context.Context)) {
	mongoConnStr := buildMongoConnectionString(host, user, password, dbname, port, sslmode, TimeZone)
	client, err := mongo.Connect(options.Client().ApplyURI(mongoConnStr))

	if err != nil {
		log.Fatal("❌ failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("❌ failed to ping MongoDB: %w", err)
	}

	return client.Database(dbname), func(ctx context.Context) { client.Disconnect(ctx) }
}

func buildMongoConnectionString(host, user, password, dbname, port, sslmode, TimeZone string) string {
	// mongodb://username:password@host:port/dbname?sslmode=sslmode&authSource=admin&timezone=timezone
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?ssl=%s&authSource=admin&tz=%s",
		user,
		password,
		host,
		port,
		dbname,
		sslmode,
		TimeZone,
	)
}
