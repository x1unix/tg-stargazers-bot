package repository

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/x1unix/tg-stargazers-bot/internal/services/auth"
	"github.com/x1unix/tg-stargazers-bot/internal/services/preferences"
)

const (
	githubTokenKeyPrefix = "github:token:"
)

var _ preferences.GitHubTokenStore = (*GitHubTokensRepository)(nil)

type GitHubTokensRepository struct {
	redis redis.Cmdable
}

func NewGitHubTokensRepository(redis redis.Cmdable) GitHubTokensRepository {
	return GitHubTokensRepository{redis: redis}
}

func (r GitHubTokensRepository) GetGitHubToken(ctx context.Context, owner auth.UserID) (string, error) {
	key := formatKey(githubTokenKeyPrefix, owner)
	val, err := r.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", preferences.ErrMissingToken
	}

	return val, err
}

func (r GitHubTokensRepository) SetGitHubToken(ctx context.Context, owner auth.UserID, token string) error {
	key := formatKey(githubTokenKeyPrefix, owner)
	return r.redis.Set(ctx, key, token, 0).Err()
}

func (r GitHubTokensRepository) RemoveGitHubToken(ctx context.Context, owner auth.UserID) error {
	key := formatKey(githubTokenKeyPrefix, owner)
	err := r.redis.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		return nil
	}

	return err
}
