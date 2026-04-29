package api

import (
	handlers2 "github.com/Tedra-ez/AdvancedProgramming_Final/internal/handlers"
	middleware2 "github.com/Tedra-ez/AdvancedProgramming_Final/internal/middleware"
	"github.com/Tedra-ez/AdvancedProgramming_Final/internal/services"
	"github.com/gin-gonic/gin"
)

func SetUpWebRoutes(r *gin.Engine, pageHandler *handlers2.PageHandler, authSvc *services.AuthService) {

	admin := r.Group("/admin")
	admin.Use(middleware2.RequireAuth, middleware2.RequireAdmin)
	{
		admin.GET("", pageHandler.AdminDashboard)
		admin.GET("/orders", pageHandler.AdminOrders)
		admin.GET("/products", pageHandler.AdminProducts)
		admin.GET("/users", pageHandler.AdminUsers)
		admin.GET("/users/:userId/orders", pageHandler.AdminUserOrders)
		admin.GET("/analytics", pageHandler.AdminAnalytics)
	}
}

func SetUpAPIRoutes(r *gin.Engine, orderHandler *handlers2.OrderHandler, productHandler *handlers2.ProductHandler, authHandler *handlers2.AuthHandler, analyticsHandler *handlers2.AnalyticsHandler, authSvc *services.AuthService) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/logout", authHandler.Logout)
		auth.POST("/refresh", authHandler.Refresh)
	}

	orders := r.Group("/orders")
	{
		orders.GET("", orderHandler.ListOrdersByUser)
		orders.POST("", orderHandler.CreateOrder)
		orders.GET("/:id", orderHandler.GetOrderStatus)
		orders.PATCH("/:id/status", orderHandler.UpdateOrderStatus)
	}

	api := r.Group("/api")
	{
		api.GET("/product", productHandler.GetProducts)
		api.GET("/product/:id", productHandler.GetProductByID)

		analytics := api.Group("/analytics")
		analytics.Use(middleware2.RequireAuth, middleware2.RequireAdmin)
		{
			analytics.GET("/stats", analyticsHandler.DashboardStatsHandler())
			analytics.GET("/top-products", analyticsHandler.TopProductsHandler())
			analytics.GET("/revenue", analyticsHandler.RevenueHandler())
			analytics.GET("/orders-status", analyticsHandler.OrdersByStatusHandler())
		}

		adminAPI := api.Group("")
		adminAPI.Use(middleware2.RequireAuth, middleware2.RequireAdmin)
		adminAPI.POST("/product", productHandler.CreateProduct)
		adminAPI.PUT("/product/:id", productHandler.UpdateProduct)
		adminAPI.DELETE("/product/:id", productHandler.DeleteProduct)
	}
}
