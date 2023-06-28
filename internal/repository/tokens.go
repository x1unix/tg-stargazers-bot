package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/x1unix/tg-stargazers-bot/internal/services/auth"
)

const (
	tokenKeyPrefix = "tokens:"
)

var _ auth.TokenStorage = (*TokenRepository)(nil)

type TokenRepository struct {
	redis *redis.Client
}

func NewTokenRepository(redis *redis.Client) TokenRepository {
	return TokenRepository{redis: redis}
}

func (t TokenRepository) TokenExists(ctx context.Context, tokenID string, subjectID auth.UserID) (bool, error) {
	key := formatKey(tokenKeyPrefix, subjectID)
	return t.redis.SIsMember(ctx, key, tokenID).Result()
}

func (t TokenRepository) AddToken(ctx context.Context, tokenID string, subjectID auth.UserID) error {
	key := formatKey(tokenKeyPrefix, subjectID)
	return t.redis.SAdd(ctx, key, tokenID).Err()
}
