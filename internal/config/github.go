package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"net/url"
)

type GitHubConfig struct {
	BaseURL      *url.URL `envconfig:"GITHUB_BASE_URL" default:"https://github.com"`
	ClientID     string   `envconfig:"GITHUB_CLIENT_ID" required:"true"`
	ClientSecret string   `envconfig:"GITHUB_CLIENT_SECRET" required:"true"`
}

func (cfg GitHubConfig) NewOAuthConfig() oauth2.Config {
	return oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     github.Endpoint,
	}
}
