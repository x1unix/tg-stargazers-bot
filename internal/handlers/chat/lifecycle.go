package chat

import (
	"context"

	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"go.uber.org/zap"
)

type TokenRemover interface {
	RemoveUserToken(ctx context.Context, subjectID bot.ChatID) error
}

type LifecycleHandler struct {
	log          *zap.Logger
	tokenRemover TokenRemover
	repoManager  RepositoryManager
}

func NewLifecycleHandler(
	log *zap.Logger,
	tokenRemover TokenRemover,
	repoManager RepositoryManager,
) *LifecycleHandler {
	return &LifecycleHandler{
		log:          log,
		tokenRemover: tokenRemover,
		repoManager:  repoManager,
	}
}

func (l LifecycleHandler) HandleUserJoin(_ context.Context, e bot.RoutedEvent) (*bot.RouteEventResult, error) {
	l.log.Info("new chat created",
		zap.String("username", e.FromChat().UserName),
		zap.Int64("chat_id", e.ChatID),
	)

	return nil, nil
}

func (l LifecycleHandler) HandleUserLeave(ctx context.Context, e bot.RoutedEvent) error {
	l.log.Info("chat deleted",
		zap.String("username", e.FromChat().UserName),
		zap.Int64("chat_id", e.ChatID),
	)
	if err := l.tokenRemover.RemoveUserToken(ctx, e.ChatID); err != nil {
		l.log.Error("failed to remove saved auth token",
			zap.String("username", e.FromChat().UserName),
			zap.Int64("chat_id", e.ChatID),
			zap.Error(err),
		)
	}

	if err := l.repoManager.TruncateUserData(ctx, e.ChatID); err != nil {
		l.log.Error("failed to remove github data",
			zap.String("username", e.FromChat().UserName),
			zap.Int64("chat_id", e.ChatID),
			zap.Error(err),
		)
	}

	return nil
}
