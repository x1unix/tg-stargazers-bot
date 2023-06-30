package repository

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"github.com/x1unix/tg-stargazers-bot/internal/services/preferences"
)

const (
	reposKeyPrefix       = "prefs:repos:"
	githubTokenKeyPrefix = "prefs:token:"
)

var _ preferences.Store = (*PreferencesRepository)(nil)

type PreferencesRepository struct {
	redis redis.Cmdable
}

func NewPreferencesRepository(redis redis.Cmdable) PreferencesRepository {
	return PreferencesRepository{redis: redis}
}

func (r PreferencesRepository) GetGitHubToken(ctx context.Context, owner bot.ChatID) (string, error) {
	key := formatKey(githubTokenKeyPrefix, owner)
	val, err := r.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", preferences.ErrMissingToken
	}

	return val, nil
}

func (r PreferencesRepository) SetGitHubToken(ctx context.Context, owner bot.ChatID, token string) error {
	key := formatKey(githubTokenKeyPrefix, owner)
	return r.redis.Set(ctx, key, token, 0).Err()
}

func (r PreferencesRepository) GetRepositories(ctx context.Context, owner bot.ChatID) ([]string, error) {
	key := formatKey(reposKeyPrefix, owner)
	return r.redis.SMembers(ctx, key).Result()
}

func (r PreferencesRepository) AddRepository(ctx context.Context, owner bot.ChatID, repo string) error {
	key := formatKey(reposKeyPrefix, owner)
	return r.redis.SAdd(ctx, key, repo).Err()
}

func (r PreferencesRepository) RemoveRepository(ctx context.Context, owner bot.ChatID, repo string) error {
	key := formatKey(reposKeyPrefix, owner)
	return r.redis.SAdd(ctx, key, repo).Err()
}
