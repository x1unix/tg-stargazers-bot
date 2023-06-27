package web

import (
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"go.uber.org/zap"
)

// NewServer constructs a new HTTP server.
func NewServer(cfg *config.Config, botSvc *bot.Service) *http.Server {
	telegramHandler := NewTelegramHandler(zap.L(), botSvc)

	e := echo.New()
	e.Use(middleware.Recover())
	e.POST(cfg.Bot.WebHookURLPath, telegramHandler.HandleTelegramWebhook)

	srv := &http.Server{
		Addr:    cfg.HTTP.ListenAddress,
		Handler: e,
	}

	return srv
}
