package main

import (
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Tedra-ez/AdvancedProgramming_Final/frontend/internal/config"
	"github.com/Tedra-ez/AdvancedProgramming_Final/internal/middleware"
	"github.com/Tedra-ez/AdvancedProgramming_Final/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type pageServer struct {
	templates map[string]*template.Template
	client    *serviceClient
}

func main() {
	godotenv.Load()

	cfg := config.Load()

	staticDir, err := mustResolveDir("static", filepath.Join("frontend", "static"))
	if err != nil {
		log.Fatalf("static dir: %v", err)
	}

	server := gin.Default()
	authService := services.NewAuthService(nil)
	server.Use(middleware.Auth(authService))
	registerAPIProxies(server, cfg)
	server.GET("/metrics", gin.WrapH(promhttp.Handler()))
	server.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"msg": "pong"}) })
	server.Static("/static", staticDir)

	templateDir, err := mustResolveDir("templates", filepath.Join("frontend", "templates"))
	if err != nil {
		log.Fatalf("template dir: %v", err)
	}
	pages, err := newPageServer(templateDir, newServiceClient(cfg))
	if err != nil {
		log.Fatalf("templates: %v", err)
	}
	pages.register(server)

	addr := ":" + cfg.Port
	if err := server.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func newPageServer(templateDir string, client *serviceClient) (*pageServer, error) {
	basePath := filepath.Join(templateDir, "base.html")
	pages := []string{
		"shop", "index", "account", "login", "register",
		"admin_orders", "admin_products", "admin_dashboard",
		"admin_users", "admin_analytics", "account_orders",
		"product", "wishlist", "cart", "checkout",
	}

	parsed := make(map[string]*template.Template, len(pages))
	for _, page := range pages {
		path := filepath.Join(templateDir, page+".html")
		tmpl, err := template.ParseFiles(basePath, path)
		if err != nil {
			return nil, err
		}
		parsed[page] = tmpl
	}

	return &pageServer{templates: parsed, client: client}, nil
}

func (p *pageServer) register(server *gin.Engine) {
	server.GET("/", p.Index)
	server.GET("/shop", p.Shop)
	server.GET("/product/:id", p.Product)
	server.GET("/account", p.Account)
	server.GET("/wishlist", p.Wishlist)
	server.GET("/cart", p.Cart)
	server.GET("/checkout", p.Checkout)
	server.GET("/login", p.Login)
	server.GET("/register", p.Register)
	server.GET("/account/orders", middleware.RequireAuth, p.AccountOrders)

	admin := server.Group("/admin")
	admin.Use(middleware.RequireAuth, middleware.RequireAdmin)
	{
		admin.GET("", p.AdminDashboard)
		admin.GET("/dashboard", p.AdminDashboard)
		admin.GET("/orders", p.AdminOrders)
		admin.GET("/products", p.AdminProducts)
		admin.GET("/users", p.AdminUsers)
		admin.GET("/users/:userId/orders", p.AdminUserOrders)
		admin.GET("/analytics", p.AdminAnalytics)
	}
}

func registerAPIProxies(server *gin.Engine, cfg *config.Config) {
	registerReverseProxy(server, "/auth", cfg.AuthServiceURL)
	registerReverseProxy(server, "/api/product", cfg.ProductServiceURL)
	registerReverseProxy(server, "/orders", cfg.OrderServiceURL)
	registerReverseProxy(server, "/api/users", cfg.UserServiceURL)
	registerReverseProxy(server, "/api/analytics", cfg.AnalyticsURL)
	registerReverseProxy(server, "/chat", cfg.ChatServiceURL)
}

func registerReverseProxy(server *gin.Engine, routePrefix, baseURL string) {
	target, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("invalid service URL for %s: %v", routePrefix, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxyHandler := gin.WrapH(proxy)
	server.Any(routePrefix, proxyHandler)
	server.Any(routePrefix+"/*proxyPath", proxyHandler)
}

func mustResolveDir(candidates ...string) (string, error) {
	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return candidate, nil
		}
	}
	return "", os.ErrNotExist
}
