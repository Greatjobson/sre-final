package ports

import "context"

type Product struct {
	ID          string
	Name        string
	Description string
	Category    string
	Price       float64
}

type ProductRepository interface {
	List(ctx context.Context) ([]Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	Create(ctx context.Context, product Product) (*Product, error)
}
