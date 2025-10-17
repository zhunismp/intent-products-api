package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/zhunismp/intent-products-api/internal/core/usecases"
	"github.com/zhunismp/intent-products-api/internal/interfaces/http/mapper"
	"github.com/zhunismp/intent-products-api/internal/interfaces/http/transport"
)

type ProductHttpHandler struct {
	productUsecase usecases.ProductUsecase
}

func NewProductHttpHandler(productUsecase usecases.ProductUsecase) *ProductHttpHandler {
	return &ProductHttpHandler{
		productUsecase: productUsecase,
	}
}

func (h *ProductHttpHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer func() {
		cancel()
	}()

	var req transport.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// TODO: add logging
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)

		json.NewEncoder(w).Encode(
			transport.ErrorResponse{
				StatusCode:   400,
				ErrorMessage: err.Error(),
			},
		)

		return
	}

	b, _ := json.Marshal(req)
	log.Println(string(b))

	if err := h.productUsecase.CreateProduct(ctx, mapper.ToCreateProductInput(req)); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(int(err.Status)) // TODO: remove typecasting

		json.NewEncoder(w).Encode(
			transport.ErrorResponse{
				StatusCode:   int32(err.Status),
				ErrorMessage: err.Message,
			},
		)

		return
	}
}
