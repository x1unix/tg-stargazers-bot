package chat

import (
	"context"
	"errors"
	"github.com/x1unix/tg-stargazers-bot/internal/services/preferences"
	"go.uber.org/zap"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
)

type RepositoryManager interface {
	TrackRepository(ctx context.Context, uid bot.ChatID, repo string) error
	GetUntrackedRepositories(ctx context.Context, uid bot.ChatID) ([]string, error)
}

type AddRepoCommand struct {
	log      *zap.Logger
	reposMgr RepositoryManager
}

func NewAddRepoCommand(log *zap.Logger, reposMgr RepositoryManager) AddRepoCommand {
	return AddRepoCommand{log: log, reposMgr: reposMgr}
}

func (cmd AddRepoCommand) CommandDescription() string {
	return "Add a new repository"
}

func (cmd AddRepoCommand) HandleBotEvent(ctx context.Context, e bot.RoutedEvent) (*bot.RouteEventResult, error) {
	repos, err := cmd.reposMgr.GetUntrackedRepositories(ctx, e.ChatID)
	if errors.Is(err, preferences.ErrMissingToken) {
		return e.NewResultWithMessage("Please authorize on GitHub using /auth command."), nil
	}
	if err != nil {
		return nil, bot.NewErrorResponse("Failed to get a list of repositories", err)
	}

	if len(repos) == 0 {
		return e.NewResultWithMessage("No available repositories to track."), nil
	}

	markup := tgbotapi.NewOneTimeReplyKeyboard(
		buildReposKeyboard(repos)...,
	)

	result := e.NewResultWithMessage(
		"Please choose a repository you'd like to track",
		bot.WithReplyMarkup(markup),
	)
	result.NextMessageHandler = NewTrackRepoMessageHandler(cmd.log, cmd.reposMgr)
	return result, nil
}

func buildReposKeyboard(repos []string) [][]tgbotapi.KeyboardButton {
	if len(repos) == 1 {
		return [][]tgbotapi.KeyboardButton{
			{
				{
					Text: repos[0],
				},
			},
		}
	}

	tail := len(repos) % 2
	count := len(repos) - tail
	rows := make([][]tgbotapi.KeyboardButton, 0, len(repos)/2+tail)
	for i := 0; i < count; i += 2 {
		cols := []tgbotapi.KeyboardButton{
			{Text: repos[i]},
			{Text: repos[i+1]},
		}
		rows = append(rows, cols)
	}

	if tail > 0 {
		rows = append(rows, []tgbotapi.KeyboardButton{
			{
				Text: repos[len(repos)-tail],
			},
		})
	}

	return rows
}
