package bot

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"go.uber.org/zap"
)

// EventHandler handles messages from Telegram bot user
type EventHandler interface {
	HandleBotEvent(ctx context.Context, e *tgbotapi.Update) (tgbotapi.Chattable, error)
}

type Service struct {
	log     *zap.Logger
	cfg     config.BotConfig
	bot     *tgbotapi.BotAPI
	handler EventHandler

	messages chan *tgbotapi.Update
}

func NewService(
	log *zap.Logger,
	cfg config.BotConfig,
	handler EventHandler,
) (*Service, error) {
	logger := log.Named("bot")
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.BotToken)
	if err != nil {
		return nil, fmt.Errorf("failed to init Telegram bot client: %w", err)
	}

	bot.Debug = cfg.Env == config.DevEnvironment
	logger.Info("successfully authorized as Telegram bot",
		zap.String("user_name", bot.Self.UserName))
	return &Service{
		log:      logger,
		cfg:      cfg,
		bot:      bot,
		handler:  handler,
		messages: make(chan *tgbotapi.Update, cfg.MessageBufferSize),
	}, nil
}

func (svc Service) HandleUpdate(u *tgbotapi.Update) {
	svc.messages <- u
}

func (svc Service) UpdateWebHookURL(baseURL string) error {
	webhookUrl, err := url.JoinPath(baseURL, svc.cfg.WebHookURLPath)
	if err != nil {
		return fmt.Errorf("failed to format webhook URL: %w", err)
	}

	wh, err := tgbotapi.NewWebhook(webhookUrl)
	if err != nil {
		return err
	}

	_, err = svc.bot.Request(wh)
	if err != nil {
		return err
	}

	svc.log.Info("bot webhook URL was updated successfully",
		zap.String("webhook_url", webhookUrl))
	return nil
}

func (svc Service) StartConsumer(ctx context.Context) {
	defer close(svc.messages)

	wg := &sync.WaitGroup{}
	for i := 0; i < svc.cfg.WorkerPoolSize; i++ {
		workerID := i
		wg.Add(1)
		go func() {
			svc.log.Debug("message consumer started", zap.Int("worker_id", workerID))
			defer svc.log.Debug("message consumer stopped", zap.Int("worker_id", workerID))

			for {
				select {
				case <-ctx.Done():
					return
				case event, ok := <-svc.messages:
					if !ok {
						svc.log.Debug("messages channel closed, skip", zap.Int("worker_id", workerID))
						return
					}

					svc.handleMessage(ctx, event)
				}
			}
		}()
	}

	svc.log.Info("starting event consumers")
	wg.Wait()
	svc.log.Info("stopped event consumers")
}

func (svc Service) handleMessage(ctx context.Context, u *tgbotapi.Update) {
	svc.log.Debug("received bot event", zap.Any("update", u))
	response, err := svc.handler.HandleBotEvent(ctx, u)
	if err != nil {
		svc.log.Error("failed to handle bot event",
			zap.Any("update", u),
			zap.Error(err),
		)

		errRsp, ok := IsErrorResponse(err)
		if !ok {
			return
		}

		response = tgbotapi.NewMessage(u.FromChat().ID, errRsp.Error())
	}

	if response == nil {
		return
	}

	_, err = svc.bot.Send(response)
	if err != nil {
		svc.log.Error("failed to send event response", zap.Error(err))
	}
}
