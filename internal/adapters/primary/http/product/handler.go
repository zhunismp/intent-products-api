package product

import (
	"log/slog"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	dto "github.com/zhunismp/intent-products-api/internal/adapters/primary/http/shared/dto"
	core "github.com/zhunismp/intent-products-api/internal/core/domain/product"
)

// for development purpose.
// user id will extract directly from Access token.
// Hence, normal request body or param will not contain user id
const OwnerID = 11122

type ProductHttpHandler struct {
	productSvc   core.ProductUsecase
	reqValidator *validator.Validate
	logger       *slog.Logger
}

func NewProductHttpHandler(productSvc core.ProductUsecase, logger *slog.Logger) *ProductHttpHandler {
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

	// calling svc
	if err := h.productSvc.CreateProduct(c.Context(), OwnerID, req.Title, req.Price, req.Link, req.Reasons); err != nil {
		return dto.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Message: "product succesfully created",
		Data:    nil,
	})
}

func (h *ProductHttpHandler) GetProduct(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{ErrorMessage: "can not parse id"})
	}

	productID := uint(id)

	// calling svc
	product, err := h.productSvc.GetProduct(c.Context(), OwnerID, productID)
	if err != nil {
		return dto.HandleError(c, err)
	}

	return dto.HandleResponse(c, fiber.StatusOK, "get product successfully", product)
}

func (h *ProductHttpHandler) GetAllProducts(c fiber.Ctx) error {
	status := c.Query("status")
	page := dto.QueryInt(c, "page", 1)
	size := dto.QueryInt(c, "size", 20)

	filter := &core.Filter{
		Status: status,
		Page:   page,
		Size:   size,
	}

	// calling svc
	products, err := h.productSvc.GetAllProducts(c.Context(), OwnerID, filter)
	if err != nil {
		return dto.HandleError(c, err)
	}

	return dto.HandleResponse(c, fiber.StatusOK, "get product successfully", products)
}

func (h *ProductHttpHandler) MoveProductPosition(c fiber.Ctx) error {
	req := new(UpdatePriorityRequest)

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

	if err := h.productSvc.Move(c.Context(), OwnerID, req.ProductID, req.ProductIDAfter); err != nil {
		return dto.HandleError(c, err)
	}

	return dto.HandleResponse(c, fiber.StatusOK, "priority was updated successfully", nil)
}

func (h *ProductHttpHandler) DeleteProduct(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{ErrorMessage: "can not parse id"})
	}

	productID := uint(id)

	// calling svc
	if err := h.productSvc.DeleteProduct(c.Context(), OwnerID, productID); err != nil {
		return dto.HandleError(c, err)
	}

	return dto.HandleResponse(c, fiber.StatusOK, "product was deleted", id)
}
