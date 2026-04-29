package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/api"
	"github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/config"
	"github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/db"
	handlers2 "github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/handlers"
	"github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/middleware"
	repository2 "github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/repository"
	services2 "github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// .env is optional: loaded locally, ignored on Render (env vars are injected by the platform)
	godotenv.Load()

	cfg := config.Load()

	server := gin.Default()
	server.Use(middleware.CORS(cfg.FrontendOrigin))
	server.GET("/metrics", gin.WrapH(promhttp.Handler()))
	server.Static("/static", resolveExistingDir("static", filepath.Join("frontend", "static")))
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "pong"})
	})
	if cfg.MongoURI == "" {
		log.Fatalf("error when connecting to mongo, please specify MONGODB_URI in .env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient, err := db.NewMongoDBClient(ctx, cfg.MongoURI)
	cancel()
	if err != nil {
		log.Fatalf("MongoDB: %v", err)
	}

	defer func() {
		if err := mongoClient.Close(context.Background()); err != nil {
			log.Println("MongoDB close:", err)
		}
	}()

	productCol := mongoClient.Collection("products")
	productRepo := repository2.NewProductRepositoryMongo(productCol)
	productService := services2.NewProductService(productRepo)
	productHandler := handlers2.NewProductHandler(productService)

	userCol := mongoClient.Collection("users")
	userRepo := repository2.NewUserRepository(userCol)

	authService := services2.NewAuthService(userRepo)
	authHandler := handlers2.NewAuthHandler(authService)

	orderItemCol := mongoClient.Collection("order_items")
	orderCol := mongoClient.Collection("orders")
	indexCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := repository2.EnsureMongoIndexes(indexCtx, orderCol, orderItemCol); err != nil {
		cancel()
		log.Fatalf("MongoDB indexes: %v", err)
	}
	cancel()
	orderItemRepo := repository2.NewOrderItemRepositoryMongo(orderItemCol)
	orderRepo := repository2.NewOrderRepositoryMongo(orderCol, orderItemRepo)
	orderService := services2.NewOrderService(orderRepo, productRepo, userRepo)
	orderHandler := handlers2.NewOrderHandler(orderService)

	analyticsService := services2.NewAnalyticsService(orderRepo, productRepo, userRepo)
	analyticsHandler := handlers2.NewAnalyticsHandler(analyticsService)

	templateDir := resolveExistingDir("templates", filepath.Join("frontend", "templates"))
	pageHandler, err := handlers2.NewPageHandler(productService, orderService, authService, analyticsService, templateDir)
	if err != nil {
		log.Fatalf("templates: %v", err)
	}

	server.Use(middleware.Metrics(), middleware.Logger(), middleware.Auth(authService))
	api.SetUpWebRoutes(server, pageHandler, authService)
	api.SetUpAPIRoutes(server, orderHandler, productHandler, authHandler, analyticsHandler, authService)

	addr := ":" + cfg.Port
	if err := server.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func resolveExistingDir(candidates ...string) string {
	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return candidate
		}
	}
	return candidates[0]
}
