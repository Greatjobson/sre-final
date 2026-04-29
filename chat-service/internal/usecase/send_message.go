package usecase

import (
	"context"

	"github.com/Tedra-ez/AdvancedProgramming_Final/chat-service/internal/ports"
)

type SendMessage struct {
	repo ports.MessageRepository
}

func NewSendMessage(repo ports.MessageRepository) *SendMessage {
	return &SendMessage{repo: repo}
}

func (uc *SendMessage) Execute(ctx context.Context, message ports.Message) error {
	return uc.repo.Save(ctx, message)
}
