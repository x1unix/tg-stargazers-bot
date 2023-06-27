package services

import (
	"fmt"
	"github.com/x1unix/tg-stargazers-bot/internal/config"
)

type GitHubService struct {
	cfg config.GitHubConfig
}

func NewGitHubService(cfg config.GitHubConfig) *GitHubService {
	return &GitHubService{cfg: cfg}
}

func (svc GitHubService) GetAuthURL() string {
	return fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=repo",
		svc.cfg.ClientID)
}
