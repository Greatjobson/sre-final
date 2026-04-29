package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Tedra-ez/AdvancedProgramming_Final/user-service/internal/ports"
)

var ErrNotImplemented = errors.New("postgres adapter is not implemented yet")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) List(ctx context.Context) ([]ports.User, error) {
	_ = ctx
	_ = r.db
	return nil, ErrNotImplemented
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*ports.User, error) {
	_ = ctx
	_ = id
	_ = r.db
	return nil, ErrNotImplemented
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	_ = ctx
	_ = r.db
	return 0, ErrNotImplemented
}
