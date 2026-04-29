package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Tedra-ez/AdvancedProgramming_Final/chat-service/internal/ports"
)

var ErrNotImplemented = errors.New("postgres adapter is not implemented yet")

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Save(ctx context.Context, message ports.Message) error {
	_ = ctx
	_ = message
	_ = r.db
	return ErrNotImplemented
}

func (r *MessageRepository) ListBetween(ctx context.Context, userID, peerID string) ([]ports.Message, error) {
	_ = ctx
	_ = userID
	_ = peerID
	_ = r.db
	return nil, ErrNotImplemented
}
