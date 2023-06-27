package config

type GitHubConfig struct {
	ClientID string `envconfig:"GITHUB_CLIENT_ID" required:"true"`
}
