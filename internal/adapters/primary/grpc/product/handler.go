package product

import (
	"context"
	"errors"

	core "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/apperrors"
	productv1 "github.com/zhunismp/intent-proto/product/gen/go/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductGrpcHandler struct {
	productv1.UnimplementedProductServiceServer
	productSvc core.ProductUsecase
}

func NewProductGrpcHandler(productSvc core.ProductUsecase) *ProductGrpcHandler {
	return &ProductGrpcHandler{productSvc: productSvc}
}

func (h *ProductGrpcHandler) BatchGetProduct(ctx context.Context, in *productv1.BatchGetProductRequest) (*productv1.BatchGetProdcutResponse, error) {
	// Validate input
	if in.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	if len(in.ProductIdList) == 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id_list cannot be empty")
	}

	// Transform to command
	cmd := core.BatchGetProductCmd{
		OwnerID:    in.UserId,
		ProductIDs: in.ProductIdList,
	}

	// calling svc
	products, err := h.productSvc.BatchGetProduct(ctx, cmd)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			return nil, status.Error(apperrors.MapToGrpcStatus(appErr.Code), appErr.Message)
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	// Transform domain products to proto
	protoProducts := make([]*productv1.Product, 0, len(products))
	for i := range products {
		protoProducts = append(protoProducts, toProtoProduct(&products[i]))
	}

	return &productv1.BatchGetProdcutResponse{
		Products: protoProducts,
	}, nil
}
