package preferences

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/v53/github"
	"github.com/samber/lo"
	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"github.com/x1unix/tg-stargazers-bot/internal/services/auth"
	"github.com/x1unix/tg-stargazers-bot/internal/util/collections"
	"golang.org/x/oauth2"
)

var (
	ErrMissingToken  = errors.New("missing GitHub token")
	ErrHookNotExists = errors.New("hook not exists")
)

type Hooks map[string]int64

type HookStore interface {
	// GetHooks returns a pair of all registered repositories and hooks.
	GetHooks(ctx context.Context, uid auth.UserID) (Hooks, error)

	// GetHook returns hook for associated repository.
	GetHook(ctx context.Context, uid auth.UserID, repo string) (int64, error)

	// AddHook stores a new pair of repository and hook
	AddHook(ctx context.Context, uid auth.UserID, repo string, hookID int64) error

	// RemoveHook removes one hook for a specified repository.
	RemoveHook(ctx context.Context, uid auth.UserID, repo string) error

	// TruncateHooks deletes all user hooks.
	TruncateHooks(ctx context.Context, uid auth.UserID) error

	// GetHookRepositories returns all repositories that have registered hooks.
	GetHookRepositories(ctx context.Context, uid auth.UserID) ([]string, error)
}

type GitHubTokenStore interface {
	GetGitHubToken(ctx context.Context, owner auth.UserID) (string, error)
	SetGitHubToken(ctx context.Context, owner auth.UserID, token string) error
	RemoveGitHubToken(ctx context.Context, owner auth.UserID) error
}

type WebhookURLBuilder interface {
	BuildWebhookURL(token string) *url.URL
}

type TokenProvider interface {
	// ProvideUserToken provides user auth token.
	ProvideUserToken(ctx context.Context, subject auth.UserID) (string, error)
}

type GitHubService struct {
	cfg           config.GitHubConfig
	urlBuilder    WebhookURLBuilder
	tokenStore    GitHubTokenStore
	hookStore     HookStore
	tokenProvider TokenProvider
}

func NewGitHubService(
	cfg config.GitHubConfig,
	urlBuilder WebhookURLBuilder,
	tokenProvider TokenProvider,
	hookStore HookStore,
	tokenStore GitHubTokenStore,
) *GitHubService {
	return &GitHubService{
		cfg:           cfg,
		urlBuilder:    urlBuilder,
		tokenProvider: tokenProvider,
		hookStore:     hookStore,
		tokenStore:    tokenStore,
	}
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
func (svc GitHubService) FetchUserToken(ctx context.Context, owner auth.UserID, verificationCode string) error {
	cfg := svc.cfg.NewOAuthConfig()
	t, err := cfg.Exchange(ctx, verificationCode)
	if err != nil {
		return fmt.Errorf("failed to obtain OAuth token: %w", err)
	}

	if err := svc.tokenStore.SetGitHubToken(ctx, owner, t.AccessToken); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	return nil
}

func (svc GitHubService) getOAuthToken(ctx context.Context, uid auth.UserID) (*oauth2.Token, error) {
	accessToken, err := svc.tokenStore.GetGitHubToken(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get github auth code: %w", err)
	}

	return &oauth2.Token{
		AccessToken: accessToken,
	}, nil
}

// GetUntrackedRepositories returns a list of available untracked repositories.
func (svc GitHubService) GetUntrackedRepositories(ctx context.Context, uid auth.UserID) ([]string, error) {
	client, err := svc.getClient(ctx, uid)
	if err != nil {
		return nil, err
	}

	repos, _, err := client.Repositories.List(ctx, "", &github.RepositoryListOptions{})
	if err != nil {
		//github.ErrorResponse
		return nil, err
	}

	trackedRepos, err := collections.SetFromResult(svc.hookStore.GetHookRepositories(ctx, uid))
	if err != nil {
		return nil, fmt.Errorf("failed to get a list of tracked repos: %w", err)
	}

	return lo.Map(
		lo.Filter(repos, func(item *github.Repository, _ int) bool {
			return !trackedRepos.Has(*item.FullName)
		}),
		func(r *github.Repository, _ int) string {
			return *r.FullName
		},
	), nil
}

// TrackRepository installs webhook on GitHub repository.
func (svc GitHubService) TrackRepository(ctx context.Context, uid auth.UserID, repo string) error {
	owner, repoName, err := splitOwnerAndRepo(repo)
	if err != nil {
		return err
	}

	webhookToken, err := svc.tokenProvider.ProvideUserToken(ctx, uid)
	if err != nil {
		return fmt.Errorf("failed to obtain callback token: %w", err)
	}

	client, err := svc.getClient(ctx, uid)
	if err != nil {
		return err
	}

	hookUrl := svc.urlBuilder.BuildWebhookURL(webhookToken).String()
	hookCfg := &github.Hook{
		Events: []string{"star"},
		Config: map[string]any{
			"url":          hookUrl,
			"content_type": "json",
		},
	}

	hook, _, err := client.Repositories.CreateHook(ctx, owner, repoName, hookCfg)
	if err != nil {
		return fmt.Errorf("failed to create hook: %w", err)
	}

	err = svc.hookStore.AddHook(ctx, uid, repo, *hook.ID)
	return nil
}

// UntrackRepository removes webhook from GitHub repo and unregisters repository.
func (svc GitHubService) UntrackRepository(ctx context.Context, uid auth.UserID, repo string) error {
	owner, repoName, err := splitOwnerAndRepo(repo)
	if err != nil {
		return err
	}

	hook, err := svc.hookStore.GetHook(ctx, uid, repo)
	if err != nil {
		return err
	}

	client, err := svc.getClient(ctx, uid)
	if err != nil {
		return err
	}

	_, err = client.Repositories.DeleteHook(ctx, owner, repoName, hook)
	if err != nil {
		return fmt.Errorf("failed to remove hook from GitHub: %w", err)
	}

	if err := svc.hookStore.RemoveHook(ctx, uid, repo); err != nil {
		return fmt.Errorf("failed to update hooks list: %w", err)
	}

	return nil
}

// TruncateUserData removes all webhooks from user repos and truncates any stored user information (repos list, tokens).
func (svc GitHubService) TruncateUserData(ctx context.Context, uid auth.UserID) error {
	client, err := svc.getClient(ctx, uid)
	if err != nil {
		return err
	}

	hooks, err := svc.hookStore.GetHooks(ctx, uid)
	if err != nil {
		return err
	}

	for fullRepoName, hookID := range hooks {
		owner, repo, err := splitOwnerAndRepo(fullRepoName)
		if err != nil {
			return err
		}

		if _, err := client.Repositories.DeleteHook(ctx, owner, repo, hookID); err != nil {
			return fmt.Errorf("failed to remove hook from repo %q: %w", fullRepoName, err)
		}
	}

	if err := svc.hookStore.TruncateHooks(ctx, uid); err != nil {
		return fmt.Errorf("failed to remove all stored hooks: %w", err)
	}

	if err := svc.tokenStore.RemoveGitHubToken(ctx, uid); err != nil {
		return fmt.Errorf("failed to remove stored GitHub token: %w", err)
	}

	return nil
}

func (svc GitHubService) getClient(ctx context.Context, uid auth.UserID) (*github.Client, error) {
	token, err := svc.getOAuthToken(ctx, uid)
	if err != nil {
		return nil, err
	}

	tokenSource := oauth2.StaticTokenSource(token)
	oauthClient := oauth2.NewClient(ctx, tokenSource)
	client := github.NewClient(oauthClient)
	return client, nil
}

func splitOwnerAndRepo(str string) (string, string, error) {
	chunks := strings.Split(str, "/")
	if len(chunks) < 2 {
		return "", "", fmt.Errorf("invalid repository name format: %q", str)
	}

	return chunks[0], chunks[1], nil
}
