package app

import (
	"context"
	"errors"
	"net/http"

	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"github.com/x1unix/tg-stargazers-bot/internal/handlers/web"
	"github.com/x1unix/tg-stargazers-bot/internal/services/auth"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"github.com/x1unix/tg-stargazers-bot/internal/services/feedback"
	"github.com/x1unix/tg-stargazers-bot/internal/services/preferences"
	"go.uber.org/zap"
)

// Version is application version.
//
// Value will be supplied by linker during build.
var Version = "1.0.0-snapshot"

type Service struct {
	log           *zap.Logger
	cfg           *config.Config
	botSvc        *bot.Service
	authSvc       *auth.Service
	githubSvc     *preferences.GitHubService
	notifications *feedback.NotificationsService
}

func NewService(
	log *zap.Logger,
	cfg *config.Config,
	botSvc *bot.Service,
	authSvc *auth.Service,
	githubSvc *preferences.GitHubService,
	notifications *feedback.NotificationsService,
) *Service {
	return &Service{
		log:           log,
		cfg:           cfg,
		botSvc:        botSvc,
		authSvc:       authSvc,
		githubSvc:     githubSvc,
		notifications: notifications,
	}
}

func (svc Service) Start(ctx context.Context) error {
	svc.log.Info("starting the service...")
	defer svc.log.Info("service stopped, goodbye")

	srvCfg := web.ServerConfig{
		HTTPConfig: svc.cfg.HTTP,
		Version:    Version,
		Env:        svc.cfg.Log.Environment,
		WebhookSecrets: web.WebhookSecrets{
			Telegram: svc.cfg.Bot.WebHookSecret,
		},
	}

	srv := web.NewServer(
		srvCfg,
		svc.botSvc,
		svc.authSvc,
		svc.githubSvc,
		svc.notifications,
	)
	go svc.botSvc.StartConsumer(ctx)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}

			svc.log.Fatal("failed to start http server", zap.Error(err))
		}
	}()

	<-ctx.Done()
	_ = srv.Shutdown(ctx)
	return nil
}
