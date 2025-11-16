package product

import (
	core "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	productv1 "github.com/zhunismp/intent-proto/product/gen/go/proto/v1"
)

type ProductGrpcHandler struct {
	productv1.UnimplementedProductServiceServer
	productSvc core.ProductUsecase
}

func NewProductGrpcHandler(productSvc core.ProductUsecase) *ProductGrpcHandler {
	return &ProductGrpcHandler{productSvc: productSvc}
}
