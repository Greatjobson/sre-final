package ports

import "context"

type Message struct {
	ID         string
	FromUserID string
	ToUserID   string
	Message    string
}

type MessageRepository interface {
	Save(ctx context.Context, message Message) error
	ListBetween(ctx context.Context, userID, peerID string) ([]Message, error)
}
