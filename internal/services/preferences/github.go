package preferences

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
)

var ErrMissingToken = errors.New("missing GitHub token")

type Store interface {
	GetGitHubToken(ctx context.Context, owner bot.ChatID) (string, error)
	SetGitHubToken(ctx context.Context, owner bot.ChatID, token string) error
	GetRepositories(ctx context.Context, owner bot.ChatID) ([]string, error)
	AddRepository(ctx context.Context, owner bot.ChatID, repo string) error
	RemoveRepository(ctx context.Context, owner bot.ChatID, repo string) error
}

type GitHubService struct {
	cfg   config.GitHubConfig
	store Store
}

func NewGitHubService(cfg config.GitHubConfig, store Store) *GitHubService {
	return &GitHubService{cfg: cfg, store: store}
}

func (svc GitHubService) BuildAuthURL(redirectUri *url.URL) string {
	params := url.Values{
		"client_id":    []string{svc.cfg.ClientID},
		"redirect_uri": []string{redirectUri.String()},
		"scope":        []string{"repo"},
	}

	newUrl := svc.cfg.BaseURL.JoinPath("/login/oauth/authorize")
	newUrl.RawQuery = params.Encode()
	return newUrl.String()
}

func (svc GitHubService) SetGitHubToken(ctx context.Context, owner bot.ChatID, token string) error {
	if err := svc.store.SetGitHubToken(ctx, owner, token); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	return nil
}
