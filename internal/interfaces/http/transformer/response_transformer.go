package transformer

import "github.com/zhunismp/intent-products-api/internal/interfaces/http/transport"

func ToResponse(data any, message string) transport.SuccessResponse {
	return transport.SuccessResponse{
		StatusCode: 200,
		Message: message,
		Data: data,
	}
}