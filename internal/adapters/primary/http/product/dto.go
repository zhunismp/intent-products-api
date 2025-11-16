package product

// model
type Sorting struct {
	Field     string `json:"field" validate:"required,oneof=title price created_at"`
	Direction string `json:"direction" validate:"required,oneof=asc desc"`
}

// request
type CreateProductRequest struct {
	Title   string   `json:"title" validate:"required"`
	Price   float64  `json:"price" validate:"min=1"`
	Link    *string  `json:"link" validate:"omitempty,url"`
	Reasons []string `json:"reasons" validate:"omitempty,dive,required"`
}

type UpdateCauseStatusRequest struct {
	ProductID string `json:"productId" validate:"required"`
	CauseID   string `json:"causeId" validate:"required"`
	Status    bool   `json:"status" validate:"required"`
}
