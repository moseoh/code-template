package db

import (
	"auth-service/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connection wraps the database connection and queries
type Connection struct {
	Pool    *pgxpool.Pool
	Queries *Queries
}

// NewConnection creates a new PostgreSQL connection using pgx
func NewConnection(cfg config.DatabaseConfig) (*Connection, error) {
	// Build connection string
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	// Create connection pool
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create queries instance
	queries := New(pool)

	return &Connection{
		Pool:    pool,
		Queries: queries,
	}, nil
}

// Close closes the database connection
func (c *Connection) Close() {
	if c.Pool != nil {
		c.Pool.Close()
	}
}
