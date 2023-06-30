package bot

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type defaultCommandHandler struct{}

func (h defaultCommandHandler) HandleBotEvent(_ context.Context, e RoutedEvent) (*RouteEventResult, error) {
	msg := tgbotapi.NewMessage(
		e.ChatID,
		"😕 Sorry, I didn't get that.\n\n"+
			"Use /help command to get a list of available commands",
	)
	return &RouteEventResult{
		Message: msg,
	}, nil
}

type helpCommandHandler struct {
	message string
}

func newHelpCommandHandler(commands map[string]CommandHandler) helpCommandHandler {
	sb := strings.Builder{}
	sb.WriteString("ℹ️ Here is a list of available commands:\n\n")
	for cmdName, cmd := range commands {
		if cmdName == "start" {
			continue
		}

		sb.WriteRune('/')
		sb.WriteString(cmdName)
		sb.WriteString(" - ")
		sb.WriteString(cmd.CommandDescription())
		sb.WriteRune('\n')
		sb.WriteString("/help - Shows this help")
	}

	return helpCommandHandler{
		message: sb.String(),
	}
}

func (h helpCommandHandler) HandleBotEvent(_ context.Context, e RoutedEvent) (*RouteEventResult, error) {
	return &RouteEventResult{
		Message: tgbotapi.NewMessage(e.ChatID, h.message),
	}, nil
}

func applyDefaults(h Handlers) Handlers {
	if h.Default == nil {
		h.Default = defaultCommandHandler{}
	}

	if h.Help == nil {
		h.Help = newHelpCommandHandler(h.Commands)
	}

	return h
}
