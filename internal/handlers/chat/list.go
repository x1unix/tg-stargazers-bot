package chat

import (
	"context"
	"strings"

	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
)

type ListRepoCommand struct {
	reposMgr RepositoryManager
}

func NewListRepoCommand(reposMgr RepositoryManager) ListRepoCommand {
	return ListRepoCommand{reposMgr: reposMgr}
}

func (cmd ListRepoCommand) HandleBotEvent(ctx context.Context, e bot.RoutedEvent) (*bot.RouteEventResult, error) {
	repos, err := cmd.reposMgr.GetTrackedRepositories(ctx, e.ChatID)
	if err != nil {
		return nil, err
	}

	if len(repos) == 0 {
		return e.NewResultWithMessage(
			"Right now, I don't track any repositories.\n" +
				"Feel free to add new repository using the /add command.",
		), nil
	}

	sb := strings.Builder{}
	sb.WriteString("Here is a list of repositories that I track:\n\n")
	for _, repo := range repos {
		sb.WriteString("• ")
		sb.WriteString(repo)
		sb.WriteRune('\n')
	}

	return e.NewResultWithMessage(sb.String()), nil
}

func (cmd ListRepoCommand) CommandDescription() string {
	return "List tracked repositories"
}
