package transport

type SuccessResponse struct {
	StatusCode int  `json:"statusCode"`
	Message    string `json:"message"`
	Data       any  `json:"data"`
}

type ErrorResponse struct {
	StatusCode   int  `json:"statusCode"`
	ErrorMessage string `json:"errorMessage"`
}
