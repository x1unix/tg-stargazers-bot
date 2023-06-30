package repository

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/x1unix/tg-stargazers-bot/internal/services/auth"
)

const (
	tokenKeyPrefix = "tokens:"
)

var _ auth.TokenStorage = (*TokenRepository)(nil)

type TokenRepository struct {
	redis redis.Cmdable
}

func NewTokenRepository(redis redis.Cmdable) TokenRepository {
	return TokenRepository{redis: redis}
}

func (t TokenRepository) TokenExists(ctx context.Context, token string, subjectID auth.UserID) (bool, error) {
	key := formatKey(tokenKeyPrefix, subjectID)
	gotToken, err := t.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	ok := token == gotToken
	return ok, nil
}

func (t TokenRepository) AddToken(ctx context.Context, token string, subjectID auth.UserID) error {
	key := formatKey(tokenKeyPrefix, subjectID)
	return t.redis.Set(ctx, key, token, 0).Err()
}
