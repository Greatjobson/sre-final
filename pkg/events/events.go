package events

import "time"

const (
	UserLoginSubject    = "auth.login"
	OrderCreatedSubject = "orders.created"
)

type UserLoginEvent struct {
	EventID    string    `json:"event_id"`
	UserID     string    `json:"user_id,omitempty"`
	Email      string    `json:"email"`
	Name       string    `json:"name,omitempty"`
	Role       string    `json:"role,omitempty"`
	RemoteAddr string    `json:"remote_addr,omitempty"`
	UserAgent  string    `json:"user_agent,omitempty"`
	LoggedAt   time.Time `json:"logged_at"`
}

type OrderCreatedEvent struct {
	EventID   string    `json:"event_id"`
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	Status    string    `json:"status"`
	Total     float64   `json:"total"`
	ItemCount int       `json:"item_count"`
	CreatedAt time.Time `json:"created_at"`
}
