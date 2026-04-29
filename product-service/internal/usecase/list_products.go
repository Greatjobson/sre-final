package usecase

import (
	"context"

	"github.com/Tedra-ez/AdvancedProgramming_Final/product-service/internal/ports"
)

type ListProducts struct {
	repo ports.ProductRepository
}

func NewListProducts(repo ports.ProductRepository) *ListProducts {
	return &ListProducts{repo: repo}
}

func (uc *ListProducts) Execute(ctx context.Context) ([]ports.Product, error) {
	return uc.repo.List(ctx)
}
