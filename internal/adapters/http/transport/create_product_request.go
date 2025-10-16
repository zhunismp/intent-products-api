package transport

// TODO: support image upload
type CreateProductRequest struct {
	Title   string   `json:"title" binding:"required"`
	Price   float64  `json:"price" binding:"required"`
	Link    *string  `json:"link"`
	Reasons []string `json:"reasons"`
}
