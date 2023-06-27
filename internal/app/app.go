package app

import (
	"context"
	"errors"
	"github.com/x1unix/tg-stargazers-bot/internal/handlers/web"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"net/http"

	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"go.uber.org/zap"
)

type Service struct {
	log    *zap.Logger
	cfg    *config.Config
	botSvc *bot.Service
}

func NewService(log *zap.Logger, cfg *config.Config, botSvc *bot.Service) *Service {
	return &Service{log: log, cfg: cfg, botSvc: botSvc}
}

func (svc Service) Start(ctx context.Context) error {
	svc.log.Info("starting the service...")
	defer svc.log.Info("service stopped, goodbye")

	srv := web.NewServer(svc.cfg, svc.botSvc)
	go svc.botSvc.StartConsumer(ctx)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}

			svc.log.Fatal("failed to start http server", zap.Error(err))
		}
	}()

	if svc.cfg.Bot.UpdateWebhookOnBoot {
		if err := svc.botSvc.UpdateWebHookURL(svc.cfg.HTTP.BaseURL); err != nil {
			svc.log.Error("failed to update bot webhook url", zap.Error(err))
		}
	}

	<-ctx.Done()
	_ = srv.Shutdown(ctx)
	return nil
}
