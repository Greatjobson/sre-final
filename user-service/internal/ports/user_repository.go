package ports

import "context"

type User struct {
	ID       string
	FullName string
	Email    string
	Role     string
}

type UserRepository interface {
	List(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Count(ctx context.Context) (int64, error)
}
