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

	productRepo := repository.NewProductRepositoryMongo(mongoClient.Collection("products"))
	userRepo := repository.NewUserRepository(mongoClient.Collection("users"))
	orderItemRepo := repository.NewOrderItemRepositoryMongo(mongoClient.Collection("order_items"))
	orderRepo := repository.NewOrderRepositoryMongo(mongoClient.Collection("orders"), orderItemRepo)

	indexCtx, indexCancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := repository.EnsureMongoIndexes(indexCtx, mongoClient.Collection("orders"), mongoClient.Collection("order_items")); err != nil {
		indexCancel()
		log.Fatalf("MongoDB indexes: %v", err)
	}
	indexCancel()

	orderService := services.NewOrderService(orderRepo, productRepo, userRepo)
	orderHandler := handlers.NewOrderHandler(orderService)
	analyticsService := services.NewAnalyticsService(orderRepo, productRepo, userRepo)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
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

	orders := server.Group("/orders")
	{
		orders.GET("", orderHandler.ListOrdersByUser)
		orders.POST("", orderHandler.CreateOrder)
		orders.GET("/:id", orderHandler.GetOrderStatus)
		orders.PATCH("/:id/status", orderHandler.UpdateOrderStatus)
	}

	analytics := server.Group("/api/analytics")
	analytics.Use(middleware.RequireAdmin)
	{
		analytics.GET("/stats", analyticsHandler.DashboardStatsHandler())
		analytics.GET("/top-products", analyticsHandler.TopProductsHandler())
		analytics.GET("/revenue", analyticsHandler.RevenueHandler())
		analytics.GET("/orders-status", analyticsHandler.OrdersByStatusHandler())
	}

	addr := ":" + cfg.Port
	if err := server.Run(addr); err != nil {
		log.Fatal(err)
	}
}
