package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Tedra-ez/AdvancedProgramming_Final/auth-service/internal/ports"
)

var ErrNotImplemented = errors.New("postgres adapter is not implemented yet")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user ports.User) error {
	_ = ctx
	_ = user
	_ = r.db
	return ErrNotImplemented
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*ports.User, error) {
	_ = ctx
	_ = email
	_ = r.db
	return nil, ErrNotImplemented
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*ports.User, error) {
	_ = ctx
	_ = id
	_ = r.db
	return nil, ErrNotImplemented
}
