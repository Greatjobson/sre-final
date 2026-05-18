package events

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher(url, serviceName string) *Publisher {
	if url == "" {
		log.Printf("NATS_URL is empty; event publishing disabled")
		return &Publisher{}
	}

	conn, err := nats.Connect(
		url,
		nats.Name(serviceName),
		nats.Timeout(2*time.Second),
		nats.ReconnectWait(2*time.Second),
		nats.MaxReconnects(-1),
	)
	if err != nil {
		log.Printf("NATS connection failed; event publishing disabled: %v", err)
		return &Publisher{}
	}

	return &Publisher{conn: conn}
}

func (p *Publisher) Close() {
	if p == nil || p.conn == nil {
		return
	}
	p.conn.Drain()
	p.conn.Close()
}

func (p *Publisher) PublishUserLogin(event UserLoginEvent) {
	p.publish(UserLoginSubject, event)
}

func (p *Publisher) PublishOrderCreated(event OrderCreatedEvent) {
	p.publish(OrderCreatedSubject, event)
}

func (p *Publisher) publish(subject string, event any) {
	if p == nil || p.conn == nil {
		return
	}

	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to encode %s event: %v", subject, err)
		return
	}
	if err := p.conn.Publish(subject, payload); err != nil {
		log.Printf("failed to publish %s event: %v", subject, err)
	}
}
