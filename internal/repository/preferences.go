package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
)

const (
	reposKeyPrefix = "prefs:repos:"
)

type PreferencesRepository struct {
	redis redis.Cmdable
}

func (r PreferencesRepository) GetTrackedRepos(ctx context.Context, chatId bot.ChatID) []string {
	//r.redis.SMembers(ctx)
	return nil
}
