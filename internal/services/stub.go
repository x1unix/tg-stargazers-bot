package services

import (
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
)

var (
	_ bot.EventHandler = (*bot.Router)(nil)
)

// NewEventRouter is dirty hack around google wire bug when it doesn't see "bot" package.
func NewEventRouter(handlers bot.Handlers) bot.EventHandler {
	return bot.NewRouter(handlers)
}
