package usecase

import (
	"context"
	"strings"

	"github.com/Tedra-ez/AdvancedProgramming_Final/auth-service/internal/ports"
)

type RegisterUser struct {
	users ports.UserRepository
}

func NewRegisterUser(users ports.UserRepository) *RegisterUser {
	return &RegisterUser{users: users}
}

func (uc *RegisterUser) Execute(ctx context.Context, fullName, email, passwordHash string) error {
	user := ports.User{
		FullName:     strings.TrimSpace(fullName),
		Email:        strings.ToLower(strings.TrimSpace(email)),
		PasswordHash: passwordHash,
		Role:         "customer",
	}
	return uc.users.Create(ctx, user)
}
