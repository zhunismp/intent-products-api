package domainerrors

type DomainError struct {
	Code       string // internal error code
	Message    string
	StatusCode int
}

func (e *DomainError) Error() string {
	return e.Message
}

func NewDomainError(code, message string, statusCode int) *DomainError {
	return &DomainError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

var (
	ErrorDuplicateProduct = NewDomainError(
		"e00001",
		"duplicate product",
		400,
	)

	ErrorProductNotFound = NewDomainError(
		"e00002",
		"product not found",
		404,
	)

	ErrorProductInput = NewDomainError(
		"e00003",
		"product input is invalid",
		400,
	)

	ErrorIllegalExecution = NewDomainError(
		"e00004",
		"user are not allowed to execute this action",
		400,
	)
)
