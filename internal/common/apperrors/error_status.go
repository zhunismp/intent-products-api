package apperrors

type StatusCode int32

const (
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
	StatusNotFound            StatusCode = 404
)
