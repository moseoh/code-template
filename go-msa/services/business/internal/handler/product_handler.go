package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"business-service/internal/db"
	"business-service/internal/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int32   `json:"stock"`
}

type ProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int32   `json:"stock"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// convertProductToResponse converts db.Product to ProductResponse
func convertProductToResponse(product *db.Product) *ProductResponse {
	description := ""
	if product.Description.Valid {
		description = product.Description.String
	}

	// Convert pgtype.Numeric to float64
	price := 0.0
	if product.Price.Valid {
		if val, err := product.Price.Float64Value(); err == nil {
			if val.Valid {
				price = val.Float64
			}
		}
	}

	return &ProductResponse{
		ID:          product.ID.String(),
		Name:        product.Name,
		Description: description,
		Price:       price,
		Stock:       product.Stock,
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Authenticated-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation
	if req.Name == "" {
		http.Error(w, "Product name is required", http.StatusBadRequest)
		return
	}
	if req.Price < 0 {
		http.Error(w, "Price must be non-negative", http.StatusBadRequest)
		return
	}
	if req.Stock < 0 {
		http.Error(w, "Stock must be non-negative", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	product, err := h.productService.CreateProduct(ctx, req.Name, req.Description, req.Price, req.Stock)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create product: %v", err), http.StatusInternalServerError)
		return
	}

	response := convertProductToResponse(product)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Authenticated-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	products, err := h.productService.ListProducts(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch products: %v", err), http.StatusInternalServerError)
		return
	}

	responses := make([]*ProductResponse, len(products))
	for i, product := range products {
		responses[i] = convertProductToResponse(product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Authenticated-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract ID from URL path (assuming /products/{id})
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}
	id := parts[len(parts)-1]

	if id == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	product, err := h.productService.GetProductByStringID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "invalid product ID") {
			http.Error(w, "Invalid product ID format", http.StatusBadRequest)
			return
		}
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	response := convertProductToResponse(product)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Authenticated-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract ID from URL path
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}
	id := parts[len(parts)-1]

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation
	if req.Name == "" {
		http.Error(w, "Product name is required", http.StatusBadRequest)
		return
	}
	if req.Price < 0 {
		http.Error(w, "Price must be non-negative", http.StatusBadRequest)
		return
	}
	if req.Stock < 0 {
		http.Error(w, "Stock must be non-negative", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	product, err := h.productService.UpdateProductByStringID(ctx, id, req.Name, req.Description, req.Price, req.Stock)
	if err != nil {
		if strings.Contains(err.Error(), "invalid product ID") {
			http.Error(w, "Invalid product ID format", http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to update product: %v", err), http.StatusInternalServerError)
		return
	}

	response := convertProductToResponse(product)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Authenticated-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract ID from URL path
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}
	id := parts[len(parts)-1]

	ctx := r.Context()
	if err := h.productService.DeleteProductByStringID(ctx, id); err != nil {
		if strings.Contains(err.Error(), "invalid product ID") {
			http.Error(w, "Invalid product ID format", http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to delete product: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) GetProductsByPriceRange(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Authenticated-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")

	if minPriceStr == "" || maxPriceStr == "" {
		http.Error(w, "min_price and max_price query parameters are required", http.StatusBadRequest)
		return
	}

	minPrice, err := strconv.ParseFloat(minPriceStr, 64)
	if err != nil {
		http.Error(w, "Invalid min_price format", http.StatusBadRequest)
		return
	}

	maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
	if err != nil {
		http.Error(w, "Invalid max_price format", http.StatusBadRequest)
		return
	}

	if minPrice < 0 || maxPrice < 0 || minPrice > maxPrice {
		http.Error(w, "Invalid price range", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	products, err := h.productService.GetProductsByPriceRange(ctx, minPrice, maxPrice)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch products: %v", err), http.StatusInternalServerError)
		return
	}

	responses := make([]*ProductResponse, len(products))
	for i, product := range products {
		responses[i] = convertProductToResponse(product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}
