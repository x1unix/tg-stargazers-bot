package chat

import (
	"context"
	_ "embed"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"github.com/x1unix/tg-stargazers-bot/internal/services/preferences"
	"go.uber.org/zap"
)

//go:embed templates/auth.txt
var authMessage []byte

type CallbackURLBuilder interface {
	BuildAuthCallbackURL(token string) *url.URL
}

type AuthCommandHandler struct {
	log           *zap.Logger
	urlBuilder    CallbackURLBuilder
	githubSvc     *preferences.GitHubService
	tokenProvider TokenProvider
}

func (s AuthCommandHandler) CommandDescription() string {
	return "Authorize the bot on GitHub"
}

func NewAuthCommandHandler(
	log *zap.Logger,
	urlBuilder CallbackURLBuilder,
	githubSvc *preferences.GitHubService,
	tokenProvider TokenProvider,
) *AuthCommandHandler {
	return &AuthCommandHandler{
		log:           log,
		urlBuilder:    urlBuilder,
		githubSvc:     githubSvc,
		tokenProvider: tokenProvider,
	}
}

func (s AuthCommandHandler) HandleBotEvent(ctx context.Context, e bot.RoutedEvent) (*bot.RouteEventResult, error) {
	actorID := e.FromChat().ID

	token, err := s.tokenProvider.ProvideUserToken(ctx, actorID)
	if err != nil {
		s.log.Error("failed to generate jwt token",
			zap.Int64("uid", actorID),
			zap.Error(err),
		)

		return nil, bot.NewErrorResponse("error occurred, please try again later", nil)
	}

	callbackUrl := s.urlBuilder.BuildAuthCallbackURL(token)
	redirectUrl := s.githubSvc.BuildAuthURL(callbackUrl)
	msg := tgbotapi.NewMessage(e.FromChat().ID, string(authMessage))
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
