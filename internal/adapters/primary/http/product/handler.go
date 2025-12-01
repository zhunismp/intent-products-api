package product

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	dto "github.com/zhunismp/intent-products-api/internal/adapters/primary/http/shared/dto"
	core "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	"go.opentelemetry.io/otel"
)

// for development purpose.
// user id will extract directly from Access token.
// Hence, normal request body or param will not contain user id
const OwnerID = 11122

type ProductHttpHandler struct {
	productSvc   core.ProductUsecase
	reqValidator *validator.Validate
	logger       *otelzap.Logger
}

func NewProductHttpHandler(productSvc core.ProductUsecase, logger *otelzap.Logger) *ProductHttpHandler {
	// setup validator
	validate := validator.New()
	validate.RegisterValidation("date_after_opt", IsDateAfter)

	return &ProductHttpHandler{
		productSvc:   productSvc,
		reqValidator: validate,
		logger:       logger,
	}
}

func (h *ProductHttpHandler) CreateProduct(c fiber.Ctx) error {
	// tracer
	tr := otel.Tracer("product-handler")
	ctx, span := tr.Start(c.Context(), "CreateProduct")
	defer span.End()

	// h.logger.Ctx(c.Context()).Info("create product request received")

	req := new(CreateProductRequest)

	// parse request body
	if err := c.Bind().Body(&req); err != nil {
		// h.logger.Ctx(c.Context()).Warn("failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{ErrorMessage: "can not parse request body"})
	}

	// validate req
	if err := h.reqValidator.Struct(req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			errMap := GenerateErrorMap(errs)
			// h.logger.Ctx(c.Context()).Warn("request validation failed", zap.Any("errors", errMap))
			return c.Status(fiber.StatusBadRequest).JSON(dto.ValidationErrorResponse{
				ErrorMessage: "invalid request",
				ErrorFields:  errMap,
			})
		}

		// h.logger.Ctx(c.Context()).Error("unexpected validation error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{ErrorMessage: "something went wrong"})
	}

	// calling svc
	if err := h.productSvc.CreateProduct(ctx, OwnerID, req.Title, req.Price, req.Link, req.Reasons); err != nil {
		return dto.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Message: "product succesfully created",
		Data:    nil,
	})
}

func (h *ProductHttpHandler) GetProduct(c fiber.Ctx) error {
	// tracer
	tr := otel.Tracer("product-handler")
	ctx, span := tr.Start(c.Context(), "GetProduct")
	defer span.End()

	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// h.logger.Ctx(c.Context()).Warn("invalid product id parameter", zap.String("id", idStr), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{ErrorMessage: "can not parse id"})
	}

	productID := uint(id)

	// calling svc
	product, err := h.productSvc.GetProduct(ctx, OwnerID, productID)
	if err != nil {
		return dto.HandleError(c, err)
	}

	return dto.HandleResponse(c, fiber.StatusOK, "get product successfully", product)
}

func (h *ProductHttpHandler) GetProductByStatus(c fiber.Ctx) error {
	// tracer
	tr := otel.Tracer("product-handler")
	ctx, span := tr.Start(c.Context(), "GetProductByStatus")
	defer span.End()

	status := c.Params("status")

	// calling svc
	products, err := h.productSvc.GetProductByStatus(ctx, OwnerID, status)
	if err != nil {
		return dto.HandleError(c, err)
	}

	return dto.HandleResponse(c, fiber.StatusOK, "get product successfully", products)
}

func (h *ProductHttpHandler) MoveProductPosition(c fiber.Ctx) error {
	// tracer
	tr := otel.Tracer("product-handler")
	ctx, span := tr.Start(c.Context(), "UpdatePriority")
	defer span.End()

	// h.logger.Ctx(c.Context()).Info("update priority request received")

	req := new(UpdatePriorityRequest)

	// parse request body
	if err := c.Bind().Body(&req); err != nil {
		// h.logger.Ctx(c.Context()).Warn("failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{ErrorMessage: "can not parse request body"})
	}

	// validate req
	if err := h.reqValidator.Struct(req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			errMap := GenerateErrorMap(errs)
			// h.logger.Ctx(c.Context()).Warn("request validation failed", zap.Any("errors", errMap))
			return c.Status(fiber.StatusBadRequest).JSON(dto.ValidationErrorResponse{
				ErrorMessage: "invalid request",
				ErrorFields:  errMap,
			})
		}

		// h.logger.Ctx(c.Context()).Error("unexpected validation error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{ErrorMessage: "something went wrong"})
	}

	if err := h.productSvc.Move(ctx, OwnerID, req.ProductID, req.ProductIDAfter); err != nil {
		return dto.HandleError(c, err)
	}

	return dto.HandleResponse(c, fiber.StatusOK, "priority was updated successfully", nil)
}

func (h *ProductHttpHandler) DeleteProduct(c fiber.Ctx) error {
	// tracer
	tr := otel.Tracer("product-service")
	ctx, span := tr.Start(c.Context(), "DeleteProduct")
	defer span.End()

	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// h.logger.Ctx(c.Context()).Warn("invalid product id parameter", zap.String("id", idStr), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{ErrorMessage: "can not parse id"})
	}

	productID := uint(id)

	// calling svc
	if err := h.productSvc.DeleteProduct(ctx, OwnerID, productID); err != nil {
		return dto.HandleError(c, err)
	}

	return dto.HandleResponse(c, fiber.StatusOK, "product was deleted", id)
}
