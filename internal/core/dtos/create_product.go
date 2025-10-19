package dtos

type CreateProductInput struct {
	OwnerID string   `validate:"required"`
	Title   string   `validate:"required"`
	Price   float64  `validate:"required,gt=0"`
	Link    *string  `validate:"omitempty,url"`
	Reasons []string `validate:"omitempty,dive"`
}
