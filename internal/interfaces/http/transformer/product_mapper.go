package transformer

import (
	"github.com/zhunismp/intent-products-api/internal/core/dtos"
	"github.com/zhunismp/intent-products-api/internal/interfaces/http/transport"
)

func ToCreateProductInput(req transport.CreateProductRequest) (dtos.CreateProductInput, error) {
	return dtos.CreateProductInput{
		OwnerID: "101234567890123456789", // hardcoded google oauth's sub
		Title:   req.Title,
		Price:   req.Price,
		Link:    req.Link,
		Reasons: req.Reasons,
	}, nil
}
