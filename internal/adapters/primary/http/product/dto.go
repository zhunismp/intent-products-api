package product

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type QueryProductRequest struct {
	Start  *time.Time `json:"start"`
	End    *time.Time `json:"end"`
	Status *string    `json:"status"`
}

func (r QueryProductRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.End, validation.When(r.End != nil && r.Start != nil, validation.Min(r.Start))),
		validation.Field(&r.Status, validation.When(r.Status != nil, validation.In("staging", "valid", "bought"))),
	)
}

type CreateProductRequest struct {
	Title   string   `json:"title"`
	Price   float64  `json:"price"`
	Link    *string  `json:"link"`
	Reasons []string `json:"reasons"`
}

func (r CreateProductRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Title, validation.Required),
		validation.Field(&r.Price, validation.Min(1)),
	)
}

type DeleteProductRequest struct {
	ProductID string `params:"id"`
}

type UpdateProductRequest struct{}

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
