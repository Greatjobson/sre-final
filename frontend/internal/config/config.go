package config

import "os"

type Config struct {
	Port              string
	APIBaseURL        string
	AuthServiceURL    string
	ProductServiceURL string
	OrderServiceURL   string
	UserServiceURL    string
	ChatServiceURL    string
	AnalyticsURL      string
}

func Load() *Config {
	port := os.Getenv("FRONTEND_PORT")
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8081"
	}

	apiBaseURL := os.Getenv("API_BASE_URL")
	return &Config{
		Port:              port,
		APIBaseURL:        apiBaseURL,
		AuthServiceURL:    serviceURL("AUTH_SERVICE_URL", apiBaseURL),
		ProductServiceURL: serviceURL("PRODUCT_SERVICE_URL", apiBaseURL),
		OrderServiceURL:   serviceURL("ORDER_SERVICE_URL", apiBaseURL),
		UserServiceURL:    serviceURL("USER_SERVICE_URL", apiBaseURL),
		ChatServiceURL:    serviceURL("CHAT_SERVICE_URL", apiBaseURL),
		AnalyticsURL:      serviceURL("ANALYTICS_SERVICE_URL", serviceURL("ORDER_SERVICE_URL", apiBaseURL)),
	}
}

func serviceURL(envName, fallback string) string {
	if value := os.Getenv(envName); value != "" {
		return value
	}
	if fallback != "" {
		return fallback
	}
	return "http://localhost:8080"
}
