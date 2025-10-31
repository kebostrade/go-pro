package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/internal/model"
	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/internal/service"
	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/pkg/response"
	"github.com/go-chi/chi/v5"
)

// ProductHandler handles product-related HTTP requests
type ProductHandler struct {
	service *service.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

// Create creates a new product
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	product, err := h.service.Create(r.Context(), &req)
	if err != nil {
		response.InternalServerError(w, "Failed to create product")
		return
	}

	response.Created(w, product)
}

// GetByID retrieves a product by ID
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid product ID")
		return
	}

	product, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Product not found")
		return
	}

	response.Success(w, product)
}

// List retrieves a list of products
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	category := r.URL.Query().Get("category")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	var products []*model.Product
	var err error

	if category != "" {
		products, err = h.service.ListByCategory(r.Context(), category, limit, offset)
	} else {
		products, err = h.service.List(r.Context(), limit, offset)
	}

	if err != nil {
		response.InternalServerError(w, "Failed to retrieve products")
		return
	}

	// Get total count for pagination
	total, _ := h.service.Count(r.Context())
	totalPages := (total + limit - 1) / limit

	meta := &response.Meta{
		Page:       offset/limit + 1,
		PerPage:    limit,
		TotalPages: totalPages,
		TotalCount: total,
	}

	response.WithMeta(w, http.StatusOK, products, meta)
}

// Update updates a product
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid product ID")
		return
	}

	var req model.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	product, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		response.InternalServerError(w, "Failed to update product")
		return
	}

	response.Success(w, product)
}

// Delete deletes a product
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid product ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		response.InternalServerError(w, "Failed to delete product")
		return
	}

	response.NoContent(w)
}

