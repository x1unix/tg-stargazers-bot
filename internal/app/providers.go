package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"github.com/x1unix/tg-stargazers-bot/internal/handlers/chat"
	"github.com/x1unix/tg-stargazers-bot/internal/repository"
	"github.com/x1unix/tg-stargazers-bot/internal/services"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"go.uber.org/zap"
)

var dependenciesSet = wire.NewSet(
	config.ReadCommandFlags,
	config.FromEnv,
	provideLogger,
	provideRedis,
	provideBotConfig,
	provideBotEventRouter,
	provideGitHubService,
	repository.NewTokenRepository,
	bot.NewService,
	NewService,
)

func provideLogger(cfg *config.Config) (*zap.Logger, error) {
	l, err := cfg.Log.NewLogger()
	if err != nil {
		return nil, err
	}

	zap.ReplaceGlobals(l)
	return l, nil
}

func provideBotConfig(cfg *config.Config) config.BotConfig {
	return cfg.Bot
}

func provideBotEventRouter(cfg *config.Config, githubSvc *services.GitHubService) bot.EventHandler {
	return bot.NewRouter(chat.NewHandlers(cfg, githubSvc))
}

func provideGitHubService(cfg *config.Config) *services.GitHubService {
	return services.NewGitHubService(cfg.GitHub)
}

func provideRedis(cfg *config.Config) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.Redis.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis DSN: %w", err)
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	client := redis.NewClient(opts)
	if err := client.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return client, nil
}
