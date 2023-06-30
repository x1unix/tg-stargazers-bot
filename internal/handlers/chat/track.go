package chat

import (
	"context"
	"regexp"
	"strings"

	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
)

var repoRegEx = regexp.MustCompile(`(?i)^[a-z0-9\.\-_]{1,}\/[a-z0-9\.\-_]{1,}$`)

type TrackRepoMessageHandler struct {
	repoMgr RepositoryManager
}

func NewTrackRepoMessageHandler(repoMgr RepositoryManager) TrackRepoMessageHandler {
	return TrackRepoMessageHandler{repoMgr: repoMgr}
}

func (t TrackRepoMessageHandler) HandleBotEvent(ctx context.Context, e bot.RoutedEvent) (*bot.RouteEventResult, error) {
	// TODO: implement webhook installation
	if e.Message == nil {
		return nil, bot.ErrUnsupported
	}

	text := strings.TrimSpace(e.Message.Text)
	if !repoNameValid(text) {
		return e.NewResultWithMessage("This doesn't seem like a valid GitHub repository name"), nil
	}

	return e.NewResultWithMessage("✅ Done. I'll notify you when someone will star that repo."), nil
}

func repoNameValid(name string) bool {
	return repoRegEx.MatchString(name)
}
