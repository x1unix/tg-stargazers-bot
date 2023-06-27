package chat

import (
	"context"
	_ "embed"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
)

//go:embed templates/default.txt
var defaultMessage []byte

type DefaultCommandHandler struct{}

func (d DefaultCommandHandler) HandleBotEvent(_ context.Context, e bot.RoutedEvent) (*bot.RouteEventResult, error) {
	return &bot.RouteEventResult{
		Message: tgbotapi.NewMessage(e.ChatID, string(defaultMessage)),
	}, nil
}
