package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/Tedra-ez/AdvancedProgramming_Final/pkg/events"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type eventCounters struct {
	logins atomic.Uint64
	orders atomic.Uint64
}

func main() {
	godotenv.Load()

	port := getenv("BACKEND_PORT", getenv("PORT", "8087"))
	natsURL := getenv("NATS_URL", "nats://localhost:4222")

	conn, err := nats.Connect(
		natsURL,
		nats.Name("notification-service"),
		nats.Timeout(2*time.Second),
		nats.ReconnectWait(2*time.Second),
		nats.MaxReconnects(-1),
	)
	if err != nil {
		log.Fatalf("NATS connection failed: %v", err)
	}
	defer conn.Close()

	counters := &eventCounters{}
	mustSubscribe(conn, events.UserLoginSubject, func(msg *nats.Msg) {
		var event events.UserLoginEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("notification: invalid login event: %v", err)
			return
		}
		counters.logins.Add(1)
		log.Printf("notification: user login email=%s user_id=%s role=%s ip=%s at=%s",
			event.Email, event.UserID, event.Role, event.RemoteAddr, event.LoggedAt.Format(time.RFC3339))
	})
	mustSubscribe(conn, events.OrderCreatedSubject, func(msg *nats.Msg) {
		var event events.OrderCreatedEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("notification: invalid order event: %v", err)
			return
		}
		counters.orders.Add(1)
		log.Printf("notification: order created order_id=%s user_id=%s status=%s total=%.2f items=%d at=%s",
			event.OrderID, event.UserID, event.Status, event.Total, event.ItemCount, event.CreatedAt.Format(time.RFC3339))
	})

	server := gin.New()
	server.Use(gin.Recovery())
	server.GET("/metrics", gin.WrapH(promhttp.Handler()))
	server.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"msg": "pong"}) })
	server.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "notification-service",
			"nats":    conn.Status().String(),
		})
	})
	server.GET("/notifications/events", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_logins":    counters.logins.Load(),
			"orders_created": counters.orders.Load(),
		})
	})

	if err := server.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func mustSubscribe(conn *nats.Conn, subject string, handler nats.MsgHandler) {
	if _, err := conn.Subscribe(subject, handler); err != nil {
		log.Fatalf("NATS subscribe %s failed: %v", subject, err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
