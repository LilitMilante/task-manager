package app

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string `env:"PORT"`

	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     int    `env:"POSTGRES_PORT"`
	PostgresName     string `env:"POSTGRES_NAME"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("load env: %w", err)
	}

	cfg := &Config{}

	err = env.ParseWithOptions(cfg, env.Options{RequiredIfNoDef: true})
	if err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	return cfg, nil
}
