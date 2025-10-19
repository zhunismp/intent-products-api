package transport

// TODO: support image upload
type CreateProductRequest struct {
	Title   string   `json:"title"`
	Price   float64  `json:"price"`
	Link    *string  `json:"link"`
	Reasons []string `json:"reasons"`
}

type DeleteProductRequest struct {
	ProductID string `json:"product_id"`
}
