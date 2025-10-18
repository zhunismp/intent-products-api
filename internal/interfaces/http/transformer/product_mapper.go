package transformer

import (
	"github.com/zhunismp/intent-products-api/internal/core/dtos"
	"github.com/zhunismp/intent-products-api/internal/interfaces/http/transport"
)

func ToCreateProductInput(req transport.CreateProductRequest) dtos.CreateProductInput {
	return dtos.CreateProductInput{
		Title:   req.Title,
		Price:   req.Price,
		Link:    req.Link,
		Reasons: req.Reasons,
	}
}
