package main

import (
	"context"
	"log"
	"time"

	"github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/config"
	"github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/db"
	"github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/handlers"
	"github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/middleware"
	"github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/repository"
	"github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/services"
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
	userHandler := handlers.NewUserHandler(authService)

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
	server.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "service": "user-service"}) })

	api := server.Group("/api/users")
	api.Use(middleware.RequireAdmin)
	{
		api.GET("", userHandler.GetUsers)
		api.GET("/count", userHandler.GetUserCount)
		api.GET("/:id", userHandler.GetUserByID)
	}

	addr := ":" + cfg.Port
	if err := server.Run(addr); err != nil {
		log.Fatal(err)
	}
}
