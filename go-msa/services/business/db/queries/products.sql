-- name: GetProduct :one
SELECT id, name, description, price, stock, created_at, updated_at
FROM products
WHERE id = $1;

-- name: ListProducts :many
SELECT id, name, description, price, stock, created_at, updated_at
FROM products
ORDER BY created_at DESC;

-- name: CreateProduct :one
INSERT INTO products (name, description, price, stock)
VALUES ($1, $2, $3, $4)
RETURNING id, name, description, price, stock, created_at, updated_at;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, description = $3, price = $4, stock = $5
WHERE id = $1
RETURNING id, name, description, price, stock, created_at, updated_at;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: GetProductsByPriceRange :many
SELECT id, name, description, price, stock, created_at, updated_at
FROM products
WHERE price BETWEEN $1 AND $2
ORDER BY price ASC;

-- name: UpdateProductStock :one
UPDATE products
SET stock = $2
WHERE id = $1
RETURNING id, name, description, price, stock, created_at, updated_at;