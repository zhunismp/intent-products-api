package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/zhunismp/intent-products-api/internal/adapters/http/transport"
	"github.com/zhunismp/intent-products-api/internal/core/usecases"
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
    	http.Error(w, "invalid JSON body", http.StatusBadRequest)
    	return
	}

	b, _ := json.Marshal(req)
	log.Println(string(b))

	if err := h.productUsecase.CreateProduct(ctx, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}