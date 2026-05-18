package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Tedra-ez/AdvancedProgramming_Final/chat-service/internal/config"
	"github.com/Tedra-ez/AdvancedProgramming_Final/chat-service/internal/middleware"
	"github.com/Tedra-ez/AdvancedProgramming_Final/chat-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type chatMessage struct {
	ID         string    `json:"id"`
	FromUserID string    `json:"from_user_id"`
	ToUserID   string    `json:"to_user_id"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}

var (
	chatMu       sync.RWMutex
	chatMessages []chatMessage
)

func main() {
	godotenv.Load()
	cfg := config.Load()

	authService := services.NewAuthService(nil)
	server := gin.New()
	server.Use(
		gin.Recovery(),
		middleware.CORS(cfg.FrontendOrigin),
		middleware.Auth(authService),
		middleware.Metrics(),
		middleware.Logger(),
	)

	server.GET("/metrics", gin.WrapH(promhttp.Handler()))
	server.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "pong"}) })
	server.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "service": "chat-service"}) })

	server.GET("/chat/messages", listMessages)
	server.POST("/chat/messages", createMessage)

	addr := ":" + cfg.Port
	if err := server.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func createMessage(c *gin.Context) {
	var req struct {
		FromUserID string `json:"from_user_id"`
		ToUserID   string `json:"to_user_id" binding:"required"`
		Message    string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fromUserID := req.FromUserID
	if fromUserID == "" {
		if userID, ok := c.Get("user_id"); ok {
			fromUserID, _ = userID.(string)
		}
	}
	if fromUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from_user_id required or authenticate via token"})
		return
	}

	msg := chatMessage{
		ID:         strconv.FormatInt(time.Now().UnixNano(), 10),
		FromUserID: fromUserID,
		ToUserID:   req.ToUserID,
		Message:    req.Message,
		CreatedAt:  time.Now().UTC(),
	}

	chatMu.Lock()
	chatMessages = append(chatMessages, msg)
	chatMu.Unlock()

	c.JSON(http.StatusCreated, msg)
}

func listMessages(c *gin.Context) {
	userID := c.Query("user_id")
	peerID := c.Query("peer_id")
	if userID == "" {
		if tokenUserID, ok := c.Get("user_id"); ok {
			userID, _ = tokenUserID.(string)
		}
	}
	if userID == "" || peerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and peer_id are required"})
		return
	}

	chatMu.RLock()
	filtered := make([]chatMessage, 0)
	for _, m := range chatMessages {
		if (m.FromUserID == userID && m.ToUserID == peerID) || (m.FromUserID == peerID && m.ToUserID == userID) {
			filtered = append(filtered, m)
		}
	}
	chatMu.RUnlock()

	c.JSON(http.StatusOK, filtered)
}
