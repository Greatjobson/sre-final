package ports

import "context"

type Order struct {
	ID     string
	UserID string
	Status string
	Total  float64
}

type OrderRepository interface {
	Create(ctx context.Context, order Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	ListByUser(ctx context.Context, userID string) ([]Order, error)
}
