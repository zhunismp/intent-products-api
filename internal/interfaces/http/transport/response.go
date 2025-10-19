package transport

type SuccessResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Pagination *Pagination `json:"pagination"`
	Data       any         `json:"data"`
}

type Pagination struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type ErrorResponse struct {
	StatusCode   int    `json:"statusCode"`
	ErrorMessage string `json:"errorMessage"`
}
