package db

import (
	"context"
	"fmt"
	"log"

	"github.com/zhunismp/intent-products-api/internal/core/domainerrors"
	"github.com/zhunismp/intent-products-api/internal/core/dtos"
	"github.com/zhunismp/intent-products-api/internal/core/entities"
	"github.com/zhunismp/intent-products-api/internal/core/repositories"
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

func (r *ProductRepositoryImpl) CreateProduct(ctx context.Context, product entities.Product) (*entities.Product, error) {
	_, err := r.collection.InsertOne(ctx, product)
	if mongo.IsDuplicateKeyError(err) {
		return nil, domainerrors.ErrorDuplicateProduct
	}
	return &product, err
}

func (r *ProductRepositoryImpl) QueryProduct(ctx context.Context, query dtos.QueryProductInput) ([]entities.Product, error) {
	filter := generateFilters(query)
	findOptions := options.Find()
	if query.Pagination != nil {
		findOptions.SetSkip(int64((query.Pagination.Page - 1) * query.Pagination.Size))
		findOptions.SetLimit(int64(query.Pagination.Size))
		findOptions.SetSort(bson.D{{Key: "added_at", Value: -1}})
	}

	// filter
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer cursor.Close(ctx)

	// parsing
	var products []entities.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	return products, nil
}

func (r *ProductRepositoryImpl) DeleteProduct(ctx context.Context, input dtos.DeleteProductInput) error {
	// filter user id and product id
	filter := bson.M{
		"id":       input.ProductID,
		"owner_id": input.OwnerID,
	}

	// delete
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	// user delete product which is not belong to user.
	if result.DeletedCount == 0 {
		return domainerrors.ErrorIllegalExecution
	}

	return nil
}

func applyIndex(ctx context.Context, collection *mongo.Collection) error {
	model := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "owner_id", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("unique_user_product"),
	}

	_, err := collection.Indexes().CreateOne(ctx, model)
	return err
}

func generateFilters(query dtos.QueryProductInput) bson.M {
	filter := bson.M{
		"owner_id": query.OwnerID,
	}

	// Create at filter
	addedAt := bson.M{}
	if query.Filters.Start != nil {
		addedAt["$gte"] = *query.Filters.Start
	}
	if query.Filters.End != nil {
		addedAt["$lte"] = *query.Filters.End
	}
	if len(addedAt) > 0 {
		filter["added_at"] = addedAt
	}

	// Status filter
	if query.Filters.Status != nil {
		filter["status"] = *query.Filters.Status
	}

	return filter
}
