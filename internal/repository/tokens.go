package repository

import "github.com/redis/go-redis/v9"

type TokenRepository struct {
	redis *redis.Client
}

func NewTokenRepository(redis *redis.Client) *TokenRepository {
	return &TokenRepository{redis: redis}
}
