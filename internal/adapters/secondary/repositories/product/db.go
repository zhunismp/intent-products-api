package product

import (
	"context"
	"fmt"
	"log"

	domain "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/apperrors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type productRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(ctx context.Context, db *mongo.Database) domain.ProductRepository {
	collection, err := NewProductCollectionBuilder(db).
		SetIndex(ctx).
		Build()
	if err != nil {
		log.Fatalf("failed to initialize product repository %v", err)
	}

	return &productRepository{collection: collection}
}

func (p *productRepository) CreateProduct(ctx context.Context, product domain.Product) (*domain.Product, error) {
	_, err := p.collection.InsertOne(ctx, product)
	if mongo.IsDuplicateKeyError(err) {
		return nil, apperrors.New("FORBIDDEN", fmt.Sprintf("owner id %s attempted to create existing product name '%s'", product.OwnerID, product.Name), err)
	}

	return &product, err
}

func (p *productRepository) QueryProduct(ctx context.Context, spec domain.QueryProductSpec) ([]domain.Product, error) {
	filter := generateFilters(spec)
	options := generateOptions(spec)

	cursor, err := p.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "failed to query products: %w", err)
	}
	defer cursor.Close(ctx)

	// parsing
	var products []domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "failed to decode products", err)
	}

	return products, nil
}

func (p *productRepository) DeleteProduct(ctx context.Context, ownerID string, productID string) error {
	filter := bson.M{
		"id":       productID,
		"owner_id": ownerID,
	}

	result, err := p.collection.DeleteOne(ctx, filter)
	if err != nil {
		return apperrors.New("INTERNAL_ERROR", "failed to delete product: %w", err)
	}

	// user delete product which is not belong to user.
	if result.DeletedCount == 0 {
		return apperrors.New("NOT_FOUND", fmt.Sprintf("owner id %s do not own product id %s", ownerID, productID), err)
	}

	return nil
}

func generateFilters(spec domain.QueryProductSpec) bson.M {
	filter := bson.M{
		"owner_id": spec.OwnerID,
	}

	if spec.Status != nil {
		filter["status"] = spec.Status
	}

	// Build date range filter if at least one of start or end is provided
	if spec.Start != nil || spec.End != nil {
		dateFilter := bson.M{}
		if spec.Start != nil {
			dateFilter["$gte"] = spec.Start
		}
		if spec.End != nil {
			dateFilter["$lte"] = spec.End
		}
		filter["added_at"] = dateFilter
	}

	return filter
}

func generateOptions(spec domain.QueryProductSpec) *options.FindOptionsBuilder {

	if spec.Sort == nil {
		return nil
	}

	return options.Find().
		SetSort(bson.D{{Key: spec.Sort.SortField, Value: spec.Sort.SortDirection}})

}
