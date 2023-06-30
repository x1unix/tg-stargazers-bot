package web

import (
	"github.com/x1unix/tg-stargazers-bot/internal/services/feedback"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"github.com/x1unix/tg-stargazers-bot/internal/services/auth"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"github.com/x1unix/tg-stargazers-bot/internal/services/preferences"
	"go.uber.org/zap"
)

const (
	githubAuthPath      = "/auth/github"
	githubWebHookPath   = "/webhook/github"
	telegramWebHookPath = "/webhook/telegram"

	tokenQueryParam = "t"
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
	githubSvc *preferences.GitHubService,
	notificationSvc *feedback.NotificationsService,

) *http.Server {
	logger := zap.L().Named("web")
	telegramHandler := NewTelegramHandler(logger, cfg.WebhookSecrets, botSvc)
	githubHandler := NewGitHubHandler(logger, cfg.HTTPConfig, githubSvc, notificationSvc)

	signParams := authSvc.JWTSignParams()
	authMiddleware := echojwt.WithConfig(echojwt.Config{
		TokenLookup:   "query:" + tokenQueryParam,
		SigningMethod: signParams.Method,
		SigningKey:    signParams.SigningKey,
		ErrorHandler: func(c echo.Context, err error) error {
			logWithContext(logger, c).
				Warn("unauthorized error", zap.Error(err))

			return echo.NewHTTPError(http.StatusUnauthorized)
		},
		NewClaimsFunc: func(_ echo.Context) jwt.Claims {
			return &auth.Claims{}
		},
	})

	e := echo.New()
	e.Use(middleware.Recover())

	secretMw := requireSecretMiddleware(logger, cfg.WebhookSecrets.Telegram)
	userTokenMw := userTokenMiddleware(logger, authSvc)

	e.POST(telegramWebHookPath, telegramHandler.HandleTelegramWebhook, secretMw)
	e.GET(githubAuthPath, githubHandler.HandleLogin, authMiddleware, userTokenMw)
	e.POST(githubWebHookPath, githubHandler.HandleWebhook, authMiddleware, userTokenMw)

	if !cfg.Env.IsProduction() {
		debugHandler := NewDebugHandler(logger, cfg, authSvc)
		e.GET("/auth/debug/login", debugHandler.HandleNewToken)
		e.GET("/auth/debug/info", debugHandler.HandleTestToken, authMiddleware, userTokenMw)
	}

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if _, ok := err.(*echo.HTTPError); !ok {
			logWithContext(logger, c).
				Error("internal server error", zap.Error(err))
		}

		e.DefaultHTTPErrorHandler(err, c)
	}

	srv := &http.Server{
		Addr:    cfg.ListenAddress,
		Handler: e,
	}

	return srv
}
