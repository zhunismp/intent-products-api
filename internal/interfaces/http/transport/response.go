package transport

type SuccessResponse struct {
	StatusCode int32  `json:"statusCode"`
	Message    string `json:"message"`
	Data       []any  `json:"data"`
}

type ErrorResponse struct {
	StatusCode   int32  `json:"statusCode"`
	ErrorMessage string `json:"errorMessage"`
}
