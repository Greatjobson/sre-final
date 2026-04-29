package usecase

import (
	"context"

	"github.com/Tedra-ez/AdvancedProgramming_Final/order-service/internal/ports"
)

type CreateOrder struct {
	repo ports.OrderRepository
}

func NewCreateOrder(repo ports.OrderRepository) *CreateOrder {
	return &CreateOrder{repo: repo}
}

func (uc *CreateOrder) Execute(ctx context.Context, order ports.Order) error {
	if order.Status == "" {
		order.Status = "pending"
	}
	return uc.repo.Create(ctx, order)
}
