package chat

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
)

type RepositoryManager interface {
	GetUntrackedRepositories(ctx context.Context, uid bot.ChatID) ([]string, error)
}

type AddRepoCommand struct {
	reposMgr RepositoryManager
}

func NewAddRepoCommand(reposMgr RepositoryManager) AddRepoCommand {
	return AddRepoCommand{reposMgr: reposMgr}
}

func (cmd AddRepoCommand) CommandDescription() string {
	return "Add a new repository"
}

func (cmd AddRepoCommand) HandleBotEvent(ctx context.Context, e bot.RoutedEvent) (*bot.RouteEventResult, error) {
	repos, err := cmd.reposMgr.GetUntrackedRepositories(ctx, e.ChatID)
	if err != nil {
		return nil, bot.NewErrorResponse("Failed to get a list of repositories", err)
	}

	if len(repos) == 0 {
		msg := tgbotapi.NewMessage(e.ChatID, "No available repositories to track.")
		return &bot.RouteEventResult{Message: msg}, nil
	}

	msg := tgbotapi.NewMessage(e.ChatID, "Please choose a repository you'd like to track")

	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		buildReposKeyboard(repos)...,
	)

	return &bot.RouteEventResult{Message: msg}, nil
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