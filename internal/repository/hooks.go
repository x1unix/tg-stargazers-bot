package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/x1unix/tg-stargazers-bot/internal/services/auth"
	"github.com/x1unix/tg-stargazers-bot/internal/services/preferences"
)

const (
	hooksKeyPrefix = "github:hooks:"
)

var _ preferences.HookStore = (*HookRepository)(nil)

type HookRepository struct {
	redis redis.Cmdable
}

func NewHookRepository(redis redis.Cmdable) HookRepository {
	return HookRepository{redis: redis}
}

func (r HookRepository) AddHook(ctx context.Context, uid auth.UserID, repo string, hookID int64) error {
	key := formatKey(hooksKeyPrefix, uid)
	return r.redis.HSet(ctx, key, repo, strconv.FormatInt(hookID, 10)).Err()
}

func (r HookRepository) GetHook(ctx context.Context, uid auth.UserID, repo string) (int64, error) {
	key := formatKey(hooksKeyPrefix, uid)
	hook, err := r.redis.HGet(ctx, key, repo).Result()
	if errors.Is(err, redis.Nil) {
		return 0, preferences.ErrHookNotExists
	}

	hookNum, err := parseHook(hook)
	if err != nil {
		return 0, err
	}

	return hookNum, nil
}

func (r HookRepository) GetHooks(ctx context.Context, uid auth.UserID) (preferences.Hooks, error) {
	key := formatKey(hooksKeyPrefix, uid)
	res, err := r.redis.HGetAll(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	result := make(preferences.Hooks, len(res))
	for repo, hookStr := range res {
		hook, err := parseHook(hookStr)
		if err != nil {
			return nil, err
		}

		result[repo] = hook
	}

	return result, nil
}

func (r HookRepository) GetHookRepositories(ctx context.Context, uid auth.UserID) ([]string, error) {
	key := formatKey(hooksKeyPrefix, uid)
	res, err := r.redis.HKeys(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	return res, err
}

func (r HookRepository) RemoveHook(ctx context.Context, uid auth.UserID, repo string) error {
	key := formatKey(hooksKeyPrefix, uid)
	err := r.redis.HDel(ctx, key, repo).Err()
	if errors.Is(err, redis.Nil) {
		return nil
	}

	return err
}

func (r HookRepository) TruncateHooks(ctx context.Context, uid auth.UserID) error {
	key := formatKey(hooksKeyPrefix, uid)
	err := r.redis.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		return nil
	}

	return err
}

func parseHook(str string) (int64, error) {
	h, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("malformed stored hook ID: %w", err)
	}

	return h, nil
}
