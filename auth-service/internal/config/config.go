package config

import "os"

type Config struct {
	MongoURI       string
	Port           string
	JWTSecret      string
	FrontendOrigin string
	APIBaseURL     string
	NATSURL        string
}

func Load() *Config {
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret_change_me"
	}

	return &Config{
		MongoURI:       os.Getenv("MONGODB_URI"),
		Port:           port,
		JWTSecret:      secret,
		FrontendOrigin: os.Getenv("FRONTEND_ORIGIN"),
		APIBaseURL:     os.Getenv("API_BASE_URL"),
		NATSURL:        os.Getenv("NATS_URL"),
	}
}
