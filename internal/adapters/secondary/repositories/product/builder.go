package product

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const productCollectionName = "product"

type ProductCollectionBuilder interface {
	SetIndex(context.Context) ProductCollectionBuilder
	Build() (*mongo.Collection, error)
}

type productCollectionBuilder struct {
	collection *mongo.Collection
	err        error
}

func NewProductCollectionBuilder(db *mongo.Database) ProductCollectionBuilder {
	return &productCollectionBuilder{
		collection: db.Collection(productCollectionName),
		err:        nil,
	}
}

func (p *productCollectionBuilder) SetIndex(ctx context.Context) ProductCollectionBuilder {
	index := mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: 1}, {Key: "owner_id", Value: 1}},
		Options: options.Index().
			SetUnique(true).
			SetName("unique_user_product"),
	}

	if _, err := p.collection.Indexes().CreateOne(ctx, index); err != nil {
		p.err = fmt.Errorf("error occur while setting index")
	}

	return p
}

func (p *productCollectionBuilder) Build() (*mongo.Collection, error) {
	return p.collection, p.err
}
