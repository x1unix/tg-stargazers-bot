package feedback

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v53/github"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
)

var reservedChars = []string{"_", "-"}

// NotificationsService provides feedback to async user actions or new events.
type NotificationsService struct {
	msgSender bot.MessageSender
}

func NewNotificationsService(msgSender bot.MessageSender) *NotificationsService {
	return &NotificationsService{msgSender: msgSender}
}

func (svc NotificationsService) NotifyAuthSuccessful(id bot.ChatID) {
	svc.msgSender.SendMessage(id,
		"🎉 Authorized successfully\nUse /add command to track new repository.",
	)
}

func (svc NotificationsService) NotifyError(id bot.ChatID, msg string, code ErrorCode) {
	svc.msgSender.SendMessage(
		id, NewErrorMessage(msg, code), bot.WithParseMode(bot.ParseModeMarkdown),
	)
}

func (svc NotificationsService) NotifyAuthFailure(id bot.ChatID, code ErrorCode) {
	svc.NotifyError(id, "GitHub auth failed", code)
}

func (svc NotificationsService) NotifyStarEvent(id bot.ChatID, event *github.StarEvent) {
	if event.Action == nil {
		return
	}

	var msg string
	switch *event.Action {
	case "created":
		msg = fmt.Sprintf(
			`🎉 User [%s](%s) starred the repo [%s](%s) \(⭐ %d\)`,
			sanitizeTelegramString(*event.Sender.Login),
			sanitizeTelegramString(*event.Sender.HTMLURL),
			sanitizeTelegramString(*event.Repo.FullName),
			sanitizeTelegramString(*event.Repo.HTMLURL),
			*event.Repo.StargazersCount,
		)
	case "deleted":
		// Deleted event doesn't contain stargazers count
		msg = fmt.Sprintf(
			`👎 User [%s](%s) un\-starred the repo [%s](%s)`,
			sanitizeTelegramString(*event.Sender.Login),
			sanitizeTelegramString(*event.Sender.HTMLURL),
			sanitizeTelegramString(*event.Repo.FullName),
			sanitizeTelegramString(*event.Repo.HTMLURL),
		)
	default:
		return
	}

	svc.msgSender.SendMessage(id, msg, bot.WithParseMode(bot.ParseModeMarkdown))
}

func sanitizeTelegramString(str string) string {
	newStr := str
	for _, reservedChar := range reservedChars {
		newStr = strings.ReplaceAll(newStr, reservedChar, `\`+reservedChar)
	}
	return newStr
}
