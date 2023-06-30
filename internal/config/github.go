package config

import "net/url"

type GitHubConfig struct {
	BaseURL  *url.URL `envconfig:"GITHUB_BASE_URL" default:"https://github.com"`
	ClientID string   `envconfig:"GITHUB_CLIENT_ID" required:"true"`
}
