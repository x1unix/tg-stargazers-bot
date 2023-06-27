package chat

import "github.com/x1unix/tg-stargazers-bot/internal/services/bot"

func NewHandlers() bot.Handlers {
	return bot.Handlers{
		Commands: map[string]bot.RoutedEventHandler{
			"start": StartCommandHandler{},
		},
		Default: DefaultCommandHandler{},
	}
}
