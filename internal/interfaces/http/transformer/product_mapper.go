package transformer

import (
	"net/url"

	"github.com/zhunismp/intent-products-api/internal/core/domainerrors"
	"github.com/zhunismp/intent-products-api/internal/core/dtos"
	"github.com/zhunismp/intent-products-api/internal/interfaces/http/transport"
)

const hardcodedUserID = "101234567890123456789"

func ToCreateProductInput(req transport.CreateProductRequest) (dtos.CreateProductInput, error) {
	return dtos.CreateProductInput{
		OwnerID: hardcodedUserID, // hardcoded google oauth's sub
		Title:   req.Title,
		Price:   req.Price,
		Link:    req.Link,
		Reasons: req.Reasons,
	}, nil
}

func ToQueryProductInput(q url.Values) (*dtos.QueryProductInput, error) {
	input, err := NewQueryProductInput(
		hardcodedUserID,
		WithStart(q.Get("start")),
		WithEnd(q.Get("end")),
		WithStatus(q.Get("status")),
		WithPagination(q.Get("page"), q.Get("size")),
	)
	if err != nil {
		return nil, domainerrors.ErrorProductInput
	}

	return input, nil
}

func ToDeleteProductInput(req transport.DeleteProductRequest) (dtos.DeleteProductInput, error) {
	return dtos.DeleteProductInput{
		OwnerID:   hardcodedUserID,
		ProductID: req.ProductID,
	}, nil
}

func ToPagination(pagination *dtos.Pagination) *transport.Pagination {
	if pagination == nil {
		return nil
	}

	return &transport.Pagination{
		Page: pagination.Page,
		Size: pagination.Size,
	}
}
