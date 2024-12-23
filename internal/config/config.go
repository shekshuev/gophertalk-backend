package config

import (
	"log"
	"time"

	"github.com/caarlos0/env/v6"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	ServerAddress       string        `env:"SERVER_ADDRESS"`
	DatabaseDSN         string        `env:"DATABASE_DSN"`
	AccessTokenExpires  time.Duration `env:"ACCESS_TOKEN_EXPIRES"`
	RefreshTokenExpires time.Duration `env:"REFRESH_TOKEN_EXPIRES"`
	AccessTokenSecret   string        `env:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret  string        `env:"REFRESH_TOKEN_SECRET"`
}

func GetConfig() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("Error starting server", err)
	}
	return cfg
}
