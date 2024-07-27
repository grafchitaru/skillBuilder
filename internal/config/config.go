package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	HTTPServerAddress   string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	PostgresDatabaseDsn string `env:"DATABASE_DSN" envDefault:"postgres://root:root@localhost:54322/app"`
	SecretKey           string `env:"SECRET_KEY" envDefault:"your_secret_key"`
	ClientServer        string `env:"CLIENT_SERVER" envDefault:"http://127.0.0.1:8081"`
}

type Configs interface {
	NewConfig() *Config
}

func NewConfig() *Config {
	var cfg Config

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = env.Parse(&cfg)
	if err != nil {
		fmt.Println("Can't parse  config: %w", err)
	}

	flag.StringVar(&cfg.HTTPServerAddress, "a", cfg.HTTPServerAddress, "HTTP server address")
	flag.StringVar(&cfg.PostgresDatabaseDsn, "d", cfg.PostgresDatabaseDsn, "PostgreSql database dsn")
	flag.StringVar(&cfg.ClientServer, "c", cfg.ClientServer, "Client Server address")

	flag.Parse()

	return &cfg
}
