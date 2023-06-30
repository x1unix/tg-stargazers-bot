package chat

import (
	"context"
	"go.uber.org/zap"
	"regexp"
	"strings"

	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
)

var repoRegEx = regexp.MustCompile(`(?i)^[a-z0-9\.\-_]{1,}\/[a-z0-9\.\-_]{1,}$`)

type TrackRepoMessageHandler struct {
	logger  *zap.Logger
	repoMgr RepositoryManager
}

func NewTrackRepoMessageHandler(logger *zap.Logger, repoMgr RepositoryManager) TrackRepoMessageHandler {
	return TrackRepoMessageHandler{logger: logger, repoMgr: repoMgr}
}

func (t TrackRepoMessageHandler) HandleBotEvent(ctx context.Context, e bot.RoutedEvent) (*bot.RouteEventResult, error) {
	if e.Message == nil {
		return nil, bot.ErrUnsupported
	}

	repo := strings.TrimSpace(e.Message.Text)
	if !repoNameValid(repo) {
		return e.NewResultWithMessage("This doesn't seem like a valid GitHub repository name"), nil
	}

	if err := t.repoMgr.TrackRepository(ctx, e.ChatID, repo); err != nil {
		return nil, bot.NewErrorResponse("Error occurred", err)
	}

	t.logger.Info("added new repository for tracking",
		zap.String("repo", repo),
		zap.Int64("uid", e.ChatID),
	)
	return e.NewResultWithMessage("✅ Done. I'll notify you when someone will star that repo."), nil
}

func repoNameValid(name string) bool {
	return repoRegEx.MatchString(name)
}
