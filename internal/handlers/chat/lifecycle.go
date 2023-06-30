package chat

import (
	"context"

	"github.com/labstack/gommon/log"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"go.uber.org/zap"
)

type LifecycleHandler struct {
	log *zap.Logger
}

func (l LifecycleHandler) HandleUserJoin(_ context.Context, e bot.RoutedEvent) (*bot.RouteEventResult, error) {
	// TODO: handle user join
	log.Info("new chat created",
		zap.String("username", e.FromChat().UserName),
		zap.Int64("chat_id", e.ChatID),
	)

	return nil, nil
}

func (l LifecycleHandler) HandleUserLeave(_ context.Context, e bot.RoutedEvent) error {
	// TODO: remove userdata
	log.Info("chat deleted",
		zap.String("username", e.FromChat().UserName),
		zap.Int64("chat_id", e.ChatID),
	)

	return nil
}

func NewLifecycleHandler(log *zap.Logger) *LifecycleHandler {
	return &LifecycleHandler{log: log}
}
