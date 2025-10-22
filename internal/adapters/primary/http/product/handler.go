package product

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	core "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/apperrors"
	"go.uber.org/zap"
)

// for development purpose.
const OwnerID = "101234567890123456789"

type ProductHttpHandler struct {
	productSvc core.ProductUsecase
}

func NewProductHttpHandler(productSvc core.ProductUsecase, logger *zap.Logger) *ProductHttpHandler {
	return &ProductHttpHandler{
		productSvc: productSvc,
		logger:     logger,
	}
}

func (h *ProductHttpHandler) CreateProduct(c fiber.Ctx) error {
	req := new(CreateProductRequest)

	// parse request body
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{ErrorMessage: "can not parse request body"})
	}

	// transform cmd
	cmd := core.CreateProductCmd{
		OwnerID: OwnerID,
		Title:   req.Title,
		Price:   req.Price,
		Link:    req.Link,
		Reasons: req.Reasons,
	}

	// calling svc
	product, err := h.productSvc.CreateProduct(c.Context(), cmd)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			return c.Status(apperrors.MapToHttpCode(appErr.Code)).JSON(ErrorResponse{ErrorMessage: appErr.Message})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{ErrorMessage: "something went wrong"})
	}

	return c.Status(200).JSON(SuccessResponse{
		Message: "product succesfully created",
		Data:    product,
	})
}

func (h *ProductHttpHandler) QueryProduct(c fiber.Ctx) error {
	req := new(QueryProductRequest)

	// Parse request body
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{ErrorMessage: "can not parse request body"})
	}

	// transform cmd
	cmd := core.QueryProductCmd{
		OwnerID: OwnerID,
		Filters: core.QueryProductFilter{
			Start:  req.Start,
			End:    req.End,
			Status: req.Status,
		},
	}

	// calling svc
	products, err := h.productSvc.QueryProduct(c.Context(), cmd)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			return c.Status(apperrors.MapToHttpCode(appErr.Code)).JSON(ErrorResponse{ErrorMessage: appErr.Message})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{ErrorMessage: "something went wrong"})
	}

	return c.Status(fiber.StatusOK).JSON(
		SuccessResponse{
			Message: "query product success",
			Data:    products,
		},
	)
}

func (h *ProductHttpHandler) DeleteProduct(c fiber.Ctx) error {
	id := c.Params("id")

	// transform cmd
	cmd := core.DeleteProductCmd{
		OwnerID:   OwnerID,
		ProductID: id,
	}

	// calling svc
	if err := h.productSvc.DeleteProduct(c.Context(), cmd); err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			return c.Status(apperrors.MapToHttpCode(appErr.Code)).JSON(ErrorResponse{ErrorMessage: appErr.Message})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{ErrorMessage: "something went wrong"})
	}

	return c.Status(fiber.StatusOK).JSON(
		SuccessResponse{
			Message: "product was deleted",
			Data:    id,
		},
	)
}
