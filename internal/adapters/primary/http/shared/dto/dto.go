package dto

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type ValidationErrorResponse struct {
	ErrorMessage string            `json:"errorMessage"`
	ErrorFields  map[string]string `json:"errorFields"`
}