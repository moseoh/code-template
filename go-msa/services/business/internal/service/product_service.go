package service

import (
	"context"
	"fmt"

	"business-service/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProductService struct {
	queries *db.Queries
}

func NewProductService(queries *db.Queries) *ProductService {
	return &ProductService{
		queries: queries,
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, name, description string, price float64, stock int32) (*db.Product, error) {
	params := db.CreateProductParams{
		Name:        name,
		Description: pgtype.Text{String: description, Valid: description != ""},
		Price:       pgtype.Numeric{Int: nil, Exp: 0, NaN: false, InfinityModifier: 0, Valid: true}, // Will need proper conversion
		Stock:       stock,
	}

	// Convert float64 to pgtype.Numeric properly
	if err := params.Price.Scan(fmt.Sprintf("%.2f", price)); err != nil {
		return nil, fmt.Errorf("failed to convert price: %w", err)
	}

	product, err := s.queries.CreateProduct(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (*db.Product, error) {
	product, err := s.queries.GetProduct(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// GetProductByStringID retrieves a product by string ID (converts to UUID)
func (s *ProductService) GetProductByStringID(ctx context.Context, idStr string) (*db.Product, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	return s.GetProduct(ctx, id)
}

// ListProducts retrieves all products
func (s *ProductService) ListProducts(ctx context.Context) ([]*db.Product, error) {
	products, err := s.queries.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	return products, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, id uuid.UUID, name, description string, price float64, stock int32) (*db.Product, error) {
	params := db.UpdateProductParams{
		ID:          id,
		Name:        name,
		Description: pgtype.Text{String: description, Valid: description != ""},
		Stock:       stock,
	}

	// Convert float64 to pgtype.Numeric properly
	if err := params.Price.Scan(fmt.Sprintf("%.2f", price)); err != nil {
		return nil, fmt.Errorf("failed to convert price: %w", err)
	}

	product, err := s.queries.UpdateProduct(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

// UpdateProductByStringID updates a product by string ID (converts to UUID)
func (s *ProductService) UpdateProductByStringID(ctx context.Context, idStr, name, description string, price float64, stock int32) (*db.Product, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	return s.UpdateProduct(ctx, id, name, description, price, stock)
}

// DeleteProduct deletes a product by ID
func (s *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if err := s.queries.DeleteProduct(ctx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

// DeleteProductByStringID deletes a product by string ID (converts to UUID)
func (s *ProductService) DeleteProductByStringID(ctx context.Context, idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	return s.DeleteProduct(ctx, id)
}

// GetProductsByPriceRange retrieves products within a price range
func (s *ProductService) GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*db.Product, error) {
	var minPriceNumeric, maxPriceNumeric pgtype.Numeric

	if err := minPriceNumeric.Scan(fmt.Sprintf("%.2f", minPrice)); err != nil {
		return nil, fmt.Errorf("failed to convert min price: %w", err)
	}

	if err := maxPriceNumeric.Scan(fmt.Sprintf("%.2f", maxPrice)); err != nil {
		return nil, fmt.Errorf("failed to convert max price: %w", err)
	}

	params := db.GetProductsByPriceRangeParams{
		Price:   minPriceNumeric,
		Price_2: maxPriceNumeric,
	}

	products, err := s.queries.GetProductsByPriceRange(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by price range: %w", err)
	}

	return products, nil
}

// UpdateProductStock updates only the stock of a product
func (s *ProductService) UpdateProductStock(ctx context.Context, id uuid.UUID, stock int32) (*db.Product, error) {
	params := db.UpdateProductStockParams{
		ID:    id,
		Stock: stock,
	}

	product, err := s.queries.UpdateProductStock(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update product stock: %w", err)
	}

	return product, nil
}
