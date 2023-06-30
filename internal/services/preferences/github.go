package preferences

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/go-github/v53/github"
	"github.com/samber/lo"
	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"golang.org/x/oauth2"
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

// FetchUserToken fetches user oauth code using verification code and persists it.
func (svc GitHubService) FetchUserToken(ctx context.Context, owner bot.ChatID, verificationCode string) error {
	cfg := svc.cfg.NewOAuthConfig()
	t, err := cfg.Exchange(ctx, verificationCode)
	if err != nil {
		return fmt.Errorf("failed to obtain OAuth token: %w", err)
	}

	if err := svc.store.SetGitHubToken(ctx, owner, t.AccessToken); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	return nil
}

func (svc GitHubService) getOAuthToken(ctx context.Context, uid bot.ChatID) (*oauth2.Token, error) {
	accessToken, err := svc.store.GetGitHubToken(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get github auth code: %w", err)
	}

	return &oauth2.Token{
		AccessToken: accessToken,
	}, nil
}

func (svc GitHubService) GetUntrackedRepositories(ctx context.Context, uid bot.ChatID) ([]string, error) {
	token, err := svc.getOAuthToken(ctx, uid)
	if err != nil {
		return nil, err
	}

	tokenSource := oauth2.StaticTokenSource(token)
	oauthClient := oauth2.NewClient(ctx, tokenSource)
	client := github.NewClient(oauthClient)
	repos, _, err := client.Repositories.List(ctx, "", &github.RepositoryListOptions{})
	if err != nil {
		//github.ErrorResponse
		return nil, err
	}

	return lo.Map(repos, func(r *github.Repository, _ int) string {
		return *r.FullName
	}), nil
}
