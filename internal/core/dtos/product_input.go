package dtos

type CreateProductInput struct {
	UserID int64  `validate:"required,gt=0"`
	ReqID  string `validate:"required"`

	Title   string   `validate:"required"`
	Price   float64  `validate:"required,gt=0"`
	Link    *string  `validate:"omitempty,url"`
	Reasons []string `validate:"omitempty,dive"`
}
