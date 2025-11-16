package product

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	dto "github.com/zhunismp/intent-products-api/internal/adapters/primary/http/shared/dto"
	core "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/apperrors"
)

// for development purpose.
// user id will extract directly from Access token.
// Hence, normal request body or param will not contain user id
const OwnerID = "101234567890123456789"

type ProductHttpHandler struct {
	productSvc   core.ProductUsecase
	reqValidator *validator.Validate
}

func NewProductHttpHandler(productSvc core.ProductUsecase) *ProductHttpHandler {
	// setup validator
	validate := validator.New()
	validate.RegisterValidation("date_after_opt", IsDateAfter)

	return &ProductHttpHandler{productSvc: productSvc, reqValidator: validate}
}

func (h *ProductHttpHandler) CreateProduct(c fiber.Ctx) error {
	req := new(CreateProductRequest)

	// parse request body
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{ErrorMessage: "can not parse request body"})
	}

	// validate req
	if err := h.reqValidator.Struct(req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			errMap := GenerateErrorMap(errs)
			return c.Status(fiber.StatusBadRequest).JSON(dto.ValidationErrorResponse{
				ErrorMessage: "invalid request",
				ErrorFields:  errMap,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{ErrorMessage: "something went wrong"})
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
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Message: "product succesfully created",
		Data:    product,
	})
}

func (h *ProductHttpHandler) GetProduct(c fiber.Ctx) error {
	id := c.Params("id")

	// transform cmd
	cmd := core.GetProductCmd{
		OwnerID:   OwnerID,
		ProductID: id,
	}

	// calling svc
	product, err := h.productSvc.GetProduct(c.Context(), cmd)
	if err != nil {
		return handleError(c, err)
	}

	return handleResponse(c, fiber.StatusOK, "get product successfully", product)
}

func (h *ProductHttpHandler) GetProductByStatus(c fiber.Ctx) error {
	status := c.Params("status")

	// transform cmd
	cmd := core.GetProductByStatusCmd{
		OwnerID: OwnerID,
		Status:  status,
	}

	// calling svc
	products, err := h.productSvc.GetProductByStatus(c.Context(), cmd)
	if err != nil {
		return handleError(c, err)
	}

	return handleResponse(c, fiber.StatusOK, "get product successfully", products)
}

func (h *ProductHttpHandler) UpdateCauseStatus(c fiber.Ctx) error {
	req := new(UpdateCauseStatusRequest)

	// parse request body
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{ErrorMessage: "can not parse request body"})
	}

	// validate req
	if err := h.reqValidator.Struct(req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			errMap := GenerateErrorMap(errs)
			return c.Status(fiber.StatusBadRequest).JSON(dto.ValidationErrorResponse{
				ErrorMessage: "invalid request",
				ErrorFields:  errMap,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{ErrorMessage: "something went wrong"})
	}

	// transform cmd
	causeStatus := core.CauseStatus{
		CauseID: req.CauseID,
		Status:  req.Status,
	}

	cmd := core.UpdateCauseStatusCmd{
		OwnerID:     OwnerID,
		ProductID:   req.ProductID,
		CauseStatus: causeStatus,
	}

	// calling svc
	product, err := h.productSvc.UpdateCauseStatus(c, cmd)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Message: "cause status updated successfully",
		Data:    product,
	})
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
		return handleError(c, err)
	}

	return handleResponse(c, fiber.StatusOK, "product was deleted", id)
}

func handleError(c fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return c.Status(apperrors.MapToHttpCode(appErr.Code)).JSON(dto.ErrorResponse{ErrorMessage: appErr.Message})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{ErrorMessage: "something went wrong"})
}

func handleResponse(c fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(
		dto.SuccessResponse{
			Message: message,
			Data:    data,
		},
	)
}
