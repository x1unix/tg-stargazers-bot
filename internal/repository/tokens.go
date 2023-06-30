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

func (t TokenRepository) GetToken(ctx context.Context, subjectID auth.UserID) (string, error) {
	key := formatKey(tokenKeyPrefix, subjectID)
	token, err := t.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", auth.ErrTokenNotExists
	}

	return token, err
}

func (t TokenRepository) AddToken(ctx context.Context, token string, subjectID auth.UserID) error {
	key := formatKey(tokenKeyPrefix, subjectID)
	return t.redis.Set(ctx, key, token, 0).Err()
}
