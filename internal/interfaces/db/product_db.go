package db

import (
	"context"

	"github.com/zhunismp/intent-products-api/internal/applications/repositories"
	"github.com/zhunismp/intent-products-api/internal/core/entities"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductRepositoryImpl struct {
	collection *mongo.Collection
}

func NewProductRepositoryImpl(db *mongo.Database) repositories.ProductRepository {
	return &ProductRepositoryImpl{
		collection: db.Collection("products"),
	}
}

func (r *ProductRepositoryImpl) CreateProduct(ctx context.Context, product entities.Product) error {
	_, err := r.collection.InsertOne(ctx, product)
	return err
}