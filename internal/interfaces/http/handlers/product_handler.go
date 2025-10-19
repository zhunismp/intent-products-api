package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/zhunismp/intent-products-api/internal/core/domainerrors"
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
		writeError(w, domainerrors.ErrorProductInput)
		return
	}

	// For debug purpose
	b, _ := json.Marshal(req)
	log.Println(string(b))

	input, err := transformer.ToCreateProductInput(req)
	if err != nil {
		writeError(w, err)
		return
	}

	createdProduct, err := h.productUsecase.CreateProduct(ctx, input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeResponse(w, createdProduct, "product created", http.StatusCreated, nil)
}

func (h *ProductHttpHandler) QueryProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer func() {
		cancel()
	}()

	input, err := transformer.ToQueryProductInput(r.URL.Query())
	if err != nil {
		writeError(w, err)
		return
	}

	resp, err := h.productUsecase.QueryProduct(ctx, *input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeResponse(w, resp, "query product successfully", http.StatusOK, transformer.ToPagination(input.Pagination))
}

func (h *ProductHttpHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer func() {
		cancel()
	}()

	var req transport.DeleteProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// TODO: add logging
		writeError(w, domainerrors.ErrorProductInput)
		return
	}

	input, err := transformer.ToDeleteProductInput(req)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := h.productUsecase.DeleteProduct(ctx, input); err != nil {
		writeError(w, err)
		return
	}

	writeResponse(w, nil, "product deleted successfully", 200, nil)
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	// default 500 error
	resp := transport.ErrorResponse{
		StatusCode:   http.StatusInternalServerError,
		ErrorMessage: "Something went wrong",
	}

	// enrich domain error
	var derr *domainerrors.DomainError
	if errors.As(err, &derr) {
		resp.StatusCode = derr.StatusCode
		resp.ErrorMessage = derr.Message
	}

	// write error
	w.WriteHeader(resp.StatusCode)
	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		http.Error(w, `{"error":"failed to encode error response"}`, http.StatusInternalServerError)
	}
}

func writeResponse(w http.ResponseWriter, data any, message string, status int, pagination *transport.Pagination) {
	resp := transport.SuccessResponse{
		StatusCode: status,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		http.Error(w, `{"error":"failed to encode error response"}`, http.StatusInternalServerError)
	}

}
