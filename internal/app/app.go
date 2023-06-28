package app

import (
	"context"
	"errors"
	"net/http"

	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"github.com/x1unix/tg-stargazers-bot/internal/handlers/web"
	"github.com/x1unix/tg-stargazers-bot/internal/services/auth"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"go.uber.org/zap"
)

type Service struct {
	log     *zap.Logger
	cfg     *config.Config
	botSvc  *bot.Service
	authSvc *auth.Service
}

func NewService(log *zap.Logger, cfg *config.Config, botSvc *bot.Service, authSvc *auth.Service) *Service {
	return &Service{
		log:     log,
		cfg:     cfg,
		botSvc:  botSvc,
		authSvc: authSvc,
	}
}

func (svc Service) Start(ctx context.Context) error {
	svc.log.Info("starting the service...")
	defer svc.log.Info("service stopped, goodbye")

	srvCfg := web.ServerConfig{
		HTTPConfig: svc.cfg.HTTP,
		Env:        svc.cfg.Log.Environment,
		WebhookSecrets: web.WebhookSecrets{
			Telegram: svc.cfg.Bot.WebHookSecret,
		},
	}

	srv := web.NewServer(srvCfg, svc.botSvc, svc.authSvc)
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
