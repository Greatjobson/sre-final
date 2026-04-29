package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Tedra-ez/AdvancedProgramming_Final/order-service/internal/ports"
)

var ErrNotImplemented = errors.New("postgres adapter is not implemented yet")

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order ports.Order) error {
	_ = ctx
	_ = order
	_ = r.db
	return ErrNotImplemented
}

func (r *OrderRepository) GetByID(ctx context.Context, id string) (*ports.Order, error) {
	_ = ctx
	_ = id
	_ = r.db
	return nil, ErrNotImplemented
}

func (r *OrderRepository) ListByUser(ctx context.Context, userID string) ([]ports.Order, error) {
	_ = ctx
	_ = userID
	_ = r.db
	return nil, ErrNotImplemented
}
