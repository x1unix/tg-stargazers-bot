package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/x1unix/tg-stargazers-bot/internal/config"
)

type startupParams struct {
	envFile string
}

func main() {
	var params startupParams
	pflag.StringVarP(&params.envFile, "env-file", "f", "", "Environment file to load (.env)")
	pflag.Parse()

	if err := mainErr(params); err != nil {
		die("Error:", err)
	}
}

func mainErr(params startupParams) error {
	if params.envFile != "" {
		if err := config.LoadEnvFile(params.envFile); err != nil {
			return err
		}
	}

	cfg, err := config.FromEnv()
	if err != nil {
		return fmt.Errorf("failed to load config from environment: %w", err)
	}

	logger, err := cfg.Log.NewLogger()
	if err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}

	logger.Info("Hello")
	return nil
}

func die(args ...any) {
	_, _ = fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}
