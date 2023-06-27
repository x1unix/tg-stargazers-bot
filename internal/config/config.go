package config

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func NewApplicationContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
}

func loadEnvFile(envFile string) error {
	if err := godotenv.Load(envFile); err != nil {
		return fmt.Errorf("failed to load environment variables file %q: %w", envFile, err)
	}

	return nil
}

func FromEnv(flags CommandFlags) (*Config, error) {
	if flags.EnvFile != "" {
		if err := loadEnvFile(flags.EnvFile); err != nil {
			return nil, err
		}
	}

	cfg := new(Config)
	err := envconfig.Process("", cfg)
	return cfg, err
}
