package domain

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type ProductRepository interface {
	Create(product *Product) error
	FindByID(id string) (*Product, error)
	FindAll() ([]*Product, error)
	Update(product *Product) error
	Delete(id string) error
}

type PostgresProductRepository struct {
	db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) Create(product *Product) error {
	query := `
		INSERT INTO products (name, description, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRow(query, product.Name, product.Description, product.Price, product.Stock, product.CreatedAt, product.UpdatedAt).Scan(&product.ID)
	return err
}

func (r *PostgresProductRepository) FindByID(id string) (*Product, error) {
	query := `
		SELECT id, name, description, price, stock, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	product := &Product{}
	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *PostgresProductRepository) FindAll() ([]*Product, error) {
	query := `
		SELECT id, name, description, price, stock, created_at, updated_at
		FROM products
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *PostgresProductRepository) Update(product *Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, stock = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := r.db.Exec(query, product.Name, product.Description, product.Price, product.Stock, product.UpdatedAt, product.ID)
	return err
}

func (r *PostgresProductRepository) Delete(id string) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
