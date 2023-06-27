package chat

import (
	"context"
	_ "embed"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
)

//go:embed templates/welcome.txt
var welcomeMessage []byte

type StartCommandHandler struct{}

func (s StartCommandHandler) HandleBotEvent(_ context.Context, e bot.RoutedEvent) (*bot.RouteEventResult, error) {
	return &bot.RouteEventResult{
		Message: tgbotapi.NewMessage(e.FromChat().ID, string(welcomeMessage)),
	}, nil
}
