package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func LoadEnvFile(envFile string) error {
	if err := godotenv.Load(envFile); err != nil {
		return fmt.Errorf("failed to load environment variables file %q: %w", envFile, err)
	}

	return nil
}

func FromEnv() (*Config, error) {
	cfg := new(Config)
	err := envconfig.Process("", cfg)
	return cfg, err
}
