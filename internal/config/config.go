package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP
		Log
	}

	// HTTP represents http server settings
	HTTP struct {
		Port string `env-required:"false" env:"HTTP_PORT"`
	}

	// Log represents a logger settings
	Log struct {
		Level string `env-required:"true" env:"LOG_LEVEL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig("./.env", cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
