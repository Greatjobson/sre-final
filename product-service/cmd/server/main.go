package main

import (
	"context"
	"log"
	"time"

	"github.com/Tedra-ez/AdvancedProgramming_Final/internal/config"
	"github.com/Tedra-ez/AdvancedProgramming_Final/internal/db"
	backendhandlers "github.com/Tedra-ez/AdvancedProgramming_Final/internal/handlers"
	"github.com/Tedra-ez/AdvancedProgramming_Final/internal/middleware"
	"github.com/Tedra-ez/AdvancedProgramming_Final/internal/repository"
	"github.com/Tedra-ez/AdvancedProgramming_Final/internal/services"
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
	productService := services.NewProductService(productRepo)
	productHandler := backendhandlers.NewProductHandler(productService)
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
	server.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "service": "product-service"}) })

	api := server.Group("/api")
	{
		api.GET("/product", productHandler.GetProducts)
		api.GET("/product/:id", productHandler.GetProductByID)
	}

	adminAPI := api.Group("")
	adminAPI.Use(middleware.RequireAdmin)
	{
		adminAPI.POST("/product", productHandler.CreateProduct)
		adminAPI.PUT("/product/:id", productHandler.UpdateProduct)
		adminAPI.DELETE("/product/:id", productHandler.DeleteProduct)
	}

	addr := ":" + cfg.Port
	if err := server.Run(addr); err != nil {
		log.Fatal(err)
	}
}
