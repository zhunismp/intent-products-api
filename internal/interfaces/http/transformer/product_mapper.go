package transformer

import (
	"strconv"

	"github.com/zhunismp/intent-products-api/internal/core/domainerrors"
	"github.com/zhunismp/intent-products-api/internal/core/dtos"
	"github.com/zhunismp/intent-products-api/internal/interfaces/http/transport"
)

func ToCreateProductInput(req transport.CreateProductRequest) (dtos.CreateProductInput, error) {
	parsedID, err := strconv.ParseInt(req.UserID, 10, 64)
	if err != nil {
		// TODO: add log when parsing failed.
		return dtos.CreateProductInput{}, domainerrors.ErrorProducInput
	}

	return dtos.CreateProductInput{
		UserID:  parsedID,
		ReqID:   req.ReqID,
		Title:   req.Title,
		Price:   req.Price,
		Link:    req.Link,
		Reasons: req.Reasons,
	}, nil
}
