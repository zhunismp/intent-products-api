package product

import (
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
	if err := h.productSvc.CreateProduct(c.Context(), cmd); err != nil {
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

	// transform cmd
	cmd := core.GetProductCmd{
		OwnerID:   OwnerID,
		ProductID: uint(id),
	}

	// calling svc
	product, err := h.productSvc.GetProduct(c.Context(), cmd)
	if err != nil {
		return dto.HandleError(c, err)
	}

	return dto.HandleResponse(c, fiber.StatusOK, "get product successfully", product)
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
		return dto.HandleError(c, err)
	}

	return dto.HandleResponse(c, fiber.StatusOK, "get product successfully", products)
}

func (h *ProductHttpHandler) UpdatePriority(c fiber.Ctx) error {
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

	cmd := core.UpdatePriorityCmd{
		OwnerID:         OwnerID,
		ProductID:       req.ProductID,
		ProductIDBefore: req.ProductIDBefore,
		ProductIDAfter:  req.ProductIDAfter,
	}

	if err := h.productSvc.UpdatePriority(c.Context(), cmd); err != nil {
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

	// transform cmd
	cmd := core.DeleteProductCmd{
		OwnerID:   OwnerID,
		ProductID: uint(id),
	}

	// calling svc
	if err := h.productSvc.DeleteProduct(c.Context(), cmd); err != nil {
		return dto.HandleError(c, err)
	}

	return dto.HandleResponse(c, fiber.StatusOK, "product was deleted", id)
}
