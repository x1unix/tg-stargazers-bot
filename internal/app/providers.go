package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"github.com/x1unix/tg-stargazers-bot/internal/handlers/chat"
	"github.com/x1unix/tg-stargazers-bot/internal/handlers/web"
	"github.com/x1unix/tg-stargazers-bot/internal/repository"
	"github.com/x1unix/tg-stargazers-bot/internal/services"
	"github.com/x1unix/tg-stargazers-bot/internal/services/auth"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"github.com/x1unix/tg-stargazers-bot/internal/services/feedback"
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
	provideURLBuilder,
	repository.NewGitHubTokensRepository,
	repository.NewTokenRepository,
	repository.NewHookRepository,
	auth.NewService,
	bot.NewService,
	services.NewEventRouter,
	feedback.NewNotificationsService,
	chat.NewHandlers,
	NewService,
	wire.Bind(new(bot.MessageSender), new(*bot.Service)),
	wire.Bind(new(chat.TokenManager), new(*auth.Service)),
	wire.Bind(new(chat.RepositoryManager), new(*preferences.GitHubService)),
	wire.Bind(new(auth.TokenStorage), new(repository.TokenRepository)),
	wire.Bind(new(preferences.GitHubTokenStore), new(repository.GitHubTokensRepository)),
	wire.Bind(new(preferences.HookStore), new(repository.HookRepository)),
	wire.Bind(new(preferences.TokenProvider), new(*auth.Service)),
	wire.Bind(new(preferences.WebhookURLBuilder), new(web.URLBuilder)),
	wire.Bind(new(chat.CallbackURLBuilder), new(web.URLBuilder)),
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
		return config.ResolvedAuthConfig{}, err
	}

	return *authCfg, nil
}

func provideGitHubService(
	cfg *config.Config,
	authSvc *auth.Service,
	tokenStore preferences.GitHubTokenStore,
	hookStore preferences.HookStore,
	urlBuilder web.URLBuilder,
) *preferences.GitHubService {
	return preferences.NewGitHubService(cfg.GitHub, urlBuilder, authSvc, hookStore, tokenStore)
}

func provideURLBuilder(cfg *config.Config) web.URLBuilder {
	return web.NewURLBuilder(cfg.HTTP.BaseURL)
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
