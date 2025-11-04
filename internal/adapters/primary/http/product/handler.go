package product

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	core "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/apperrors"
)

// for development purpose.
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
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{ErrorMessage: "can not parse request body"})
	}

	// validate req
	if err := h.reqValidator.Struct(req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			errMap := GenerateErrorMap(errs)
			return c.Status(fiber.StatusBadRequest).JSON(ValidationErrorResponse{
				ErrorMessage: "invalid request",
				ErrorFields:  errMap,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{ErrorMessage: "something went wrong"})
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

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
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
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			return c.Status(apperrors.MapToHttpCode(appErr.Code)).JSON(ErrorResponse{ErrorMessage: appErr.Message})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{ErrorMessage: "something went wrong"})
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Message: "get product successfully",
		Data:    product,
	})
}

func (h *ProductHttpHandler) UpdateCauseStatus(c fiber.Ctx) error {
	req := new(UpdateCauseStatusRequest)

	// parse request body
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{ErrorMessage: "can not parse request body"})
	}

	// validate req
	if err := h.reqValidator.Struct(req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			errMap := GenerateErrorMap(errs)
			return c.Status(fiber.StatusBadRequest).JSON(ValidationErrorResponse{
				ErrorMessage: "invalid request",
				ErrorFields:  errMap,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{ErrorMessage: "something went wrong"})
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
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			return c.Status(apperrors.MapToHttpCode(appErr.Code)).JSON(ErrorResponse{ErrorMessage: appErr.Message})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{ErrorMessage: "something went wrong"})
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Message: "cause status updated successfully",
		Data:    product,
	})
}

func (h *ProductHttpHandler) QueryProduct(c fiber.Ctx) error {
	req := new(QueryProductRequest)

	// Parse request body
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{ErrorMessage: "can not parse request body"})
	}

	// validate req
	if err := h.reqValidator.Struct(req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			errMap := GenerateErrorMap(errs)
			return c.Status(fiber.StatusBadRequest).JSON(ValidationErrorResponse{
				ErrorMessage: "invalid request",
				ErrorFields:  errMap,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{ErrorMessage: "something went wrong"})
	}

	// transform cmd
	var sort *core.Sort
	if req.Sort != nil {
		sort = &core.Sort{
			Field:     req.Sort.Field,
			Direction: req.Sort.Direction,
		}
	}

	cmd := core.QueryProductCmd{
		OwnerID: OwnerID,
		Start:   req.Start,
		End:     req.End,
		Status:  req.Status,
		Sort:    sort,
	}

	// calling svc
	products, err := h.productSvc.QueryProducts(c.Context(), cmd)
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
