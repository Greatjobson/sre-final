package usecase

import (
	"context"

	"github.com/Tedra-ez/AdvancedProgramming_Final/user-service/internal/ports"
)

type GetUser struct {
	repo ports.UserRepository
}

func NewGetUser(repo ports.UserRepository) *GetUser {
	return &GetUser{repo: repo}
}

func (uc *GetUser) Execute(ctx context.Context, id string) (*ports.User, error) {
	return uc.repo.GetByID(ctx, id)
}
