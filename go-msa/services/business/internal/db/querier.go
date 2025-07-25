// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateProduct(ctx context.Context, arg CreateProductParams) (*Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	GetProduct(ctx context.Context, id uuid.UUID) (*Product, error)
	GetProductsByPriceRange(ctx context.Context, arg GetProductsByPriceRangeParams) ([]*Product, error)
	ListProducts(ctx context.Context) ([]*Product, error)
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (*Product, error)
	UpdateProductStock(ctx context.Context, arg UpdateProductStockParams) (*Product, error)
}

var _ Querier = (*Queries)(nil)
