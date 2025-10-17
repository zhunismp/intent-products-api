package transport

// TODO: support image upload
type CreateProductRequest struct {
	Title   string   `json:"title" binding:"required" validate:"required"`
	Price   float64  `json:"price" binding:"required" validate:"required,gt=0"`
	Link    *string  `json:"link" validate:"omitempty,url"`
	Reasons []string `json:"reasons" validate:"omitempty,dive"`
}
