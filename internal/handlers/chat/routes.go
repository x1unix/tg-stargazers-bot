package chat

import (
	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"github.com/x1unix/tg-stargazers-bot/internal/services/preferences"
	"go.uber.org/zap"
)

func NewHandlers(
	log *zap.Logger,
	cfg *config.Config,
	githubSvc *preferences.GitHubService,
	tokenProvider TokenProvider,
) bot.Handlers {
	logger := log.With(zap.String("tag", "bot"))
	return bot.Handlers{
		Start:            NewStartCommandHandler(),
		LifecycleHandler: NewLifecycleHandler(log),
		Commands: bot.CommandHandlers{
			"auth": NewAuthCommandHandler(logger, cfg.HTTP, githubSvc, tokenProvider),
			"add":  NewAddRepoCommand(githubSvc),
		},
	}
}
