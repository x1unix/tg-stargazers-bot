package chat

import (
	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"github.com/x1unix/tg-stargazers-bot/internal/services"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
)

func NewHandlers(
	cfg *config.Config,
	githubSvc *services.GitHubService,
) bot.Handlers {
	return bot.Handlers{
		Commands: map[string]bot.RoutedEventHandler{
			"start": NewStartCommandHandler(cfg.HTTP, githubSvc),
		},
		Default: DefaultCommandHandler{},
	}
}
