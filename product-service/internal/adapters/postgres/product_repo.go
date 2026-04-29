package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Tedra-ez/AdvancedProgramming_Final/product-service/internal/ports"
)

var ErrNotImplemented = errors.New("postgres adapter is not implemented yet")

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) List(ctx context.Context) ([]ports.Product, error) {
	_ = ctx
	_ = r.db
	return nil, ErrNotImplemented
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*ports.Product, error) {
	_ = ctx
	_ = id
	_ = r.db
	return nil, ErrNotImplemented
}

func (r *ProductRepository) Create(ctx context.Context, product ports.Product) (*ports.Product, error) {
	_ = ctx
	_ = product
	_ = r.db
	return nil, ErrNotImplemented
}
