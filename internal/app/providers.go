package app

import (
	"context"
	"fmt"
	"github.com/x1unix/tg-stargazers-bot/internal/services/feedback"
	"time"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"github.com/x1unix/tg-stargazers-bot/internal/handlers/chat"
	"github.com/x1unix/tg-stargazers-bot/internal/repository"
	"github.com/x1unix/tg-stargazers-bot/internal/services"
	"github.com/x1unix/tg-stargazers-bot/internal/services/auth"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"github.com/x1unix/tg-stargazers-bot/internal/services/preferences"
	"go.uber.org/zap"
)

var dependenciesSet = wire.NewSet(
	config.ReadCommandFlags,
	config.FromEnv,
	provideLogger,
	provideRedis,
	provideBotConfig,
	provideAuthConfig,
	provideGitHubService,
	repository.NewPreferencesRepository,
	repository.NewTokenRepository,
	auth.NewService,
	bot.NewService,
	services.NewEventRouter,
	feedback.NewNotificationsService,
	chat.NewHandlers,
	NewService,
	wire.Bind(new(chat.TokenProvider), new(*auth.Service)),
	wire.Bind(new(auth.TokenStorage), new(repository.TokenRepository)),
	wire.Bind(new(preferences.Store), new(repository.PreferencesRepository)),
	wire.Bind(new(bot.MessageSender), new(*bot.Service)),
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

func provideAuthConfig(cfg *config.Config) (config.ResolvedAuthConfig, error) {
	authCfg, err := cfg.Auth.ResolvedAuthConfig()
	if err != nil {
		return config.ResolvedAuthConfig{}, nil
	}

	return *authCfg, nil
}

func provideGitHubService(cfg *config.Config, store preferences.Store) *preferences.GitHubService {
	return preferences.NewGitHubService(cfg.GitHub, store)
}

func provideRedis(cfg *config.Config) (redis.Cmdable, error) {
	opts, err := redis.ParseURL(cfg.Redis.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis DSN: %w", err)
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	client := redis.NewClient(opts)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return client, nil
}
