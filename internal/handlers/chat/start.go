package chat

import (
	"context"
	_ "embed"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"github.com/x1unix/tg-stargazers-bot/internal/services"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
)

//go:embed templates/welcome.txt
var welcomeMessage []byte

type StartCommandHandler struct {
	cfg       config.HTTPConfig
	githubSvc *services.GitHubService
}

func NewStartCommandHandler(cfg config.HTTPConfig, githubSvc *services.GitHubService) *StartCommandHandler {
	return &StartCommandHandler{cfg: cfg, githubSvc: githubSvc}
}

func (s StartCommandHandler) HandleBotEvent(_ context.Context, e bot.RoutedEvent) (*bot.RouteEventResult, error) {
	redirectUrl := s.githubSvc.GetAuthURL()

	msg := tgbotapi.NewMessage(e.FromChat().ID, string(welcomeMessage))
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				{
					Text: "Authorize Starbot",
					URL:  &redirectUrl,
				},
			},
		},
	}
	return &bot.RouteEventResult{
		Message: msg,
	}, nil
}
