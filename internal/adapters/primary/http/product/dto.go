package product

// request
type CreateProductRequest struct {
	Title   string   `json:"title" validate:"required"`
	Price   float64  `json:"price" validate:"min=1"`
	Link    *string  `json:"link" validate:"omitempty,url"`
	Reasons []string `json:"reasons" validate:"omitempty,dive,required"`
}

type UpdatePriorityRequest struct {
	ProductID       uint `json:"productId" validate:"required"`
	ProductIDBefore *uint `json:"productIdBefore"`
	ProductIDAfter  *uint `json:"productIdAfter"`
}
