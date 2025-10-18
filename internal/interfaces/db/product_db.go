package db

import (
	"context"
	"log"

	"github.com/zhunismp/intent-products-api/internal/applications/repositories"
	"github.com/zhunismp/intent-products-api/internal/core/entities"
	"github.com/zhunismp/intent-products-api/internal/core/domainerrors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ProductRepositoryImpl struct {
	collection *mongo.Collection
}

func NewProductRepositoryImpl(ctx context.Context, db *mongo.Database) repositories.ProductRepository {
	collection := db.Collection("products")

	// apply index constraint
	if err := applyIndex(ctx, collection); err != nil {
		log.Fatalf("Error while apply index constraint %v", err)
	}

	return &ProductRepositoryImpl{
		collection: collection,
	}
}

func (r *ProductRepositoryImpl) CreateProduct(ctx context.Context, product entities.Product) error {
	_, err := r.collection.InsertOne(ctx, product)
	if mongo.IsDuplicateKeyError(err) {
		return domainerrors.ErrorDuplicateProduct
	}
	return err
}

func applyIndex(ctx context.Context, collection *mongo.Collection) error {
	model := mongo.IndexModel{
		Keys: bson.D{
			{Key: "product_id", Value: 1},
			{Key: "owner_id", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("unique_user_product"),
	}

	_, err := collection.Indexes().CreateOne(ctx, model)
	return err
}
