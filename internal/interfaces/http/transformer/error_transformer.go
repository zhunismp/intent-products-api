package transformer

import (
	"errors"

	"github.com/zhunismp/intent-products-api/internal/core/domainerrors"
	"github.com/zhunismp/intent-products-api/internal/interfaces/http/transport"
)

var badRequestErrors = []error{
    domainerrors.ErrorDuplicateProduct,
}

var notFoundErrors = []error{
    domainerrors.ErrorProductNotFound,
}

func ToErrorResponse(err error) transport.ErrorResponse {
	statusCode := mapStatusCode(err)

	return transport.ErrorResponse{
		StatusCode: statusCode,
		ErrorMessage: err.Error(),
	}
}

func mapStatusCode(err error) int {
	// Check for 400 status
    for _, e := range badRequestErrors {
        if errors.Is(err, e) {
            return 400
        }
    }

	// check for 404 status
    for _, e := range notFoundErrors {
        if errors.Is(err, e) {
            return 404
        }
    }

    return 500
}