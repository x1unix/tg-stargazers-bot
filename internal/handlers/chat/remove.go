package chat

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"go.uber.org/zap"
)

type RemoveRepoCommand struct {
	log      *zap.Logger
	reposMgr RepositoryManager
}

func NewRemoveRepoCommand(log *zap.Logger, reposMgr RepositoryManager) RemoveRepoCommand {
	return RemoveRepoCommand{log: log, reposMgr: reposMgr}
}

func (cmd RemoveRepoCommand) CommandDescription() string {
	return "Stop tracking a repository"
}

func (cmd RemoveRepoCommand) HandleBotEvent(ctx context.Context, e bot.RoutedEvent) (*bot.RouteEventResult, error) {
	repos, err := cmd.reposMgr.GetTrackedRepositories(ctx, e.ChatID)
	if err != nil {
		return nil, err
	}

	if len(repos) == 0 {
		return e.NewResultWithMessage(
			"Currently you don't track any repositories.\n\n" +
				"Use /add command to track a new repository.",
		), nil
	}

	markup := tgbotapi.NewOneTimeReplyKeyboard(
		buildReposKeyboard(repos)...,
	)

	result := e.NewResultWithMessage(
		"Please choose a repository you'd like to stop tracking",
		bot.WithReplyMarkup(markup),
	)

	result.NextMessageHandler = NewUntrackRepoMessageHandler(cmd.log, cmd.reposMgr)
	return result, nil
}
