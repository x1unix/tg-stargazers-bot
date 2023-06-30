package feedback

import (
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
)

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
