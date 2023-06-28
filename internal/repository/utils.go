package repository

import (
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"strconv"
)

func formatKey(keyPrefix string, chatId bot.ChatID) string {
	return keyPrefix + strconv.FormatInt(chatId, 10)
}
