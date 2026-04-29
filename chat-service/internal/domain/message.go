package domain

import "time"

type Message struct {
	ID         string
	FromUserID string
	ToUserID   string
	Message    string
	CreatedAt  time.Time
}
