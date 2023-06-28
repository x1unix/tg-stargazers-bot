package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"github.com/x1unix/tg-stargazers-bot/internal/services/auth"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"go.uber.org/zap"
)

const (
	githubAuthPath      = "/auth/github"
	githubWebHookPath   = "/webhook/github"
	telegramWebHookPath = "/webhook/telegram"
)

type WebhookSecrets struct {
	Telegram string
}

type ServerConfig struct {
	config.HTTPConfig
	Env            config.Environment
	WebhookSecrets WebhookSecrets
}

// NewServer constructs a new HTTP server.
func NewServer(
	cfg ServerConfig,
	botSvc *bot.Service,
	authSvc *auth.Service,
) *http.Server {
	logger := zap.L().Named("web")
	telegramHandler := NewTelegramHandler(logger, cfg.WebhookSecrets, botSvc)
	githubHandler := NewGitHubHandler(logger)

	e := echo.New()
	e.Use(middleware.Recover())
	e.POST(telegramWebHookPath, telegramHandler.HandleTelegramWebhook)
	e.GET(githubAuthPath, githubHandler.HandleLogin)
	e.POST(githubWebHookPath, githubHandler.HandleWebhook)

	if !cfg.Env.IsProduction() {
		debugHandler := NewDebugHandler(logger, cfg, authSvc)
		e.GET("/auth/debug", debugHandler.HandleNewToken)
	}

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if _, ok := err.(*echo.HTTPError); !ok {
			logger.Error("internal server error",
				zap.String("method", c.Request().Method),
				zap.String("url", c.Request().RequestURI),
				zap.Error(err),
			)
		}

		e.DefaultHTTPErrorHandler(err, c)
	}

	srv := &http.Server{
		Addr:    cfg.ListenAddress,
		Handler: e,
	}

	return srv
}
