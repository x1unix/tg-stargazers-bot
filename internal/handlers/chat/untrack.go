package chat

import (
	"context"
	"strings"

	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"go.uber.org/zap"
)

type UntrackRepoMessageHandler struct {
	logger  *zap.Logger
	repoMgr RepositoryManager
}

func NewUntrackRepoMessageHandler(logger *zap.Logger, repoMgr RepositoryManager) UntrackRepoMessageHandler {
	return UntrackRepoMessageHandler{logger: logger, repoMgr: repoMgr}
}

func (t UntrackRepoMessageHandler) HandleBotEvent(ctx context.Context, e bot.RoutedEvent) (*bot.RouteEventResult, error) {
	if e.Message == nil {
		return nil, bot.ErrUnsupported
	}

	repo := strings.TrimSpace(e.Message.Text)
	if !repoNameValid(repo) {
		return e.NewResultWithMessage("This doesn't seem like a valid GitHub repository name"), nil
	}

	if err := t.repoMgr.UntrackRepository(ctx, e.ChatID, repo); err != nil {
		return nil, bot.NewErrorResponse("Error occurred", err)
	}

	return e.NewResultWithMessage(
		"✅ Done. I uninstalled the webhook and removed the repo from tracking list.",
	), nil
}
