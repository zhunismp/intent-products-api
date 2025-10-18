package transport

// TODO: support image upload
type CreateProductRequest struct {
	UserID string `json:"user_id"`
	ReqID  string `json:"request_id"`

	Title   string   `json:"title"`
	Price   float64  `json:"price"`
	Link    *string  `json:"link"`
	Reasons []string `json:"reasons"`
}
