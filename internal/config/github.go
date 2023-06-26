package config

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/x1unix/tg-stargazers-bot/internal/util/keyutil"
)

type GitHubConfig struct {
	AppSlug        string `envconfig:"GITHUB_APP_SLUG" required:"true"`
	AppSecret      string `envconfig:"GITHUB_APP_SECRET" required:"true"`
	WebHookSecret  string `envconfig:"GITHUB_WEBHOOK_SECRET" required:"true"`
	PrivateKeyFile string `envconfig:"GITHUB_APP_PRIVATE_KEY_FILE" required:"true"`
}

func (cfg GitHubConfig) PrivateKey() (*rsa.PrivateKey, error) {
	keyFile, err := os.ReadFile(cfg.PrivateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	key, err := keyutil.ParseRSAPrivateKeyFromPEM(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key file %q: %w", cfg.PrivateKeyFile, err)
	}

	return key, nil
}
