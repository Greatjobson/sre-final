package main

import (
	"context"
	"log"
	"time"

	"github.com/Tedra-ez/AdvancedProgramming_Final/auth-service/internal/config"
	"github.com/Tedra-ez/AdvancedProgramming_Final/auth-service/internal/db"
	backendhandlers "github.com/Tedra-ez/AdvancedProgramming_Final/auth-service/internal/handlers"
	"github.com/Tedra-ez/AdvancedProgramming_Final/auth-service/internal/middleware"
	"github.com/Tedra-ez/AdvancedProgramming_Final/auth-service/internal/repository"
	"github.com/Tedra-ez/AdvancedProgramming_Final/auth-service/internal/services"
	"github.com/Tedra-ez/AdvancedProgramming_Final/pkg/events"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	godotenv.Load()
	cfg := config.Load()

	if cfg.MongoURI == "" {
		log.Fatal("MONGODB_URI is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient, err := db.NewMongoDBClient(ctx, cfg.MongoURI)
	cancel()
	if err != nil {
		log.Fatalf("MongoDB: %v", err)
	}
	defer mongoClient.Close(context.Background())

	userRepo := repository.NewUserRepository(mongoClient.Collection("users"))
	authService := services.NewAuthService(userRepo)
	eventPublisher := events.NewPublisher(cfg.NATSURL, "auth-service")
	defer eventPublisher.Close()
	authHandler := backendhandlers.NewAuthHandler(authService, eventPublisher)

	server := gin.New()
	server.Use(gin.Recovery(), middleware.CORS(cfg.FrontendOrigin), middleware.Metrics(), middleware.Logger())
	server.GET("/metrics", gin.WrapH(promhttp.Handler()))
	server.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "pong"}) })
	server.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "service": "auth-service"}) })

	auth := server.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/logout", authHandler.Logout)
		auth.POST("/refresh", authHandler.Refresh)
	}

	addr := ":" + cfg.Port
	if err := server.Run(addr); err != nil {
		log.Fatal(err)
	}
}
