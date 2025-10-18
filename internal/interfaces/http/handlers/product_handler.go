package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/zhunismp/intent-products-api/internal/core/usecases"
	"github.com/zhunismp/intent-products-api/internal/interfaces/http/transformer"
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
		writeError(w, err)
		return
	}

	b, _ := json.Marshal(req)
	log.Println(string(b))

	createdProduct, err := h.productUsecase.CreateProduct(ctx, transformer.ToCreateProductInput(req))
	if err != nil {
		writeError(w, err)
		return
	}

	writeResponse(w, createdProduct, "product created")
}

func writeError(w http.ResponseWriter, err error) {
	errResponse := transformer.ToErrorResponse(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errResponse.StatusCode)
	json.NewEncoder(w).Encode(errResponse)
}

func writeResponse(w http.ResponseWriter, data any, message string) {
	response := transformer.ToResponse(data, message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)

}
