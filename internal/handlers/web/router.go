package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"go.uber.org/zap"
)

const (
	githubAuthPath    = "/github/auth"
	githubWebHookPath = "/github/webhook"
)

// NewServer constructs a new HTTP server.
func NewServer(
	cfg *config.Config,
	botSvc *bot.Service,
) *http.Server {
	telegramHandler := NewTelegramHandler(zap.L(), botSvc)
	githubHandler := NewGitHubHandler(zap.L())

	e := echo.New()
	e.Use(middleware.Recover())
	e.POST(cfg.Bot.WebHookURLPath, telegramHandler.HandleTelegramWebhook)
	e.GET(githubAuthPath, githubHandler.HandleLogin)
	e.POST(githubWebHookPath, githubHandler.HandleWebhook)

	srv := &http.Server{
		Addr:    cfg.HTTP.ListenAddress,
		Handler: e,
	}

	return srv
}
