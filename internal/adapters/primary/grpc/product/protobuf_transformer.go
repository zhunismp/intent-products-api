package product

import (
	"time"

	core "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	productv1 "github.com/zhunismp/intent-proto/product/gen/go/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoProduct(product *core.Product) *productv1.Product {
	if product == nil {
		return nil
	}

	return &productv1.Product{
		Id:        product.ID,
		OwnerId:   product.OwnerID,
		Name:      product.Name,
		ImageUrl:  product.ImageUrl,
		Link:      product.Link,
		Price:     product.Price,
		Status:    product.Status,
		CreatedAt: timestampFromTime(product.CreatedAt),
		UpdatedAt: timestampFromTime(product.UpdatedAt),
	}
}

func timestampFromTime(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}
