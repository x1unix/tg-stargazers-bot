package chat

import (
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"github.com/x1unix/tg-stargazers-bot/internal/services/preferences"
	"go.uber.org/zap"
)

type TokenManager interface {
	TokenProvider
	TokenRemover
}

func NewHandlers(
	log *zap.Logger,
	urlBuilder CallbackURLBuilder,
	githubSvc *preferences.GitHubService,
	tokenMgr TokenManager,
) bot.Handlers {
	logger := log.With(zap.String("tag", "bot"))
	return bot.Handlers{
		Start:            NewStartCommandHandler(),
		LifecycleHandler: NewLifecycleHandler(log, tokenMgr, githubSvc),
		Commands: bot.CommandHandlers{
			"auth":   NewAuthCommandHandler(logger, urlBuilder, githubSvc, tokenMgr),
			"add":    NewAddRepoCommand(logger, githubSvc),
			"remove": NewRemoveRepoCommand(log, githubSvc),
			"list":   NewListRepoCommand(githubSvc),
		},
	}
}
