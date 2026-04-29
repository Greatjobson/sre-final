package ports

import "context"

type User struct {
	ID           string
	FullName     string
	Email        string
	PasswordHash string
	Role         string
}

type UserRepository interface {
	Create(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
}
