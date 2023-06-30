package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// ParseMode is message parse mode.
//
// See: https://core.telegram.org/bots/api#formatting-options
type ParseMode = string

const (
	ParseModeMarkdownLegacy ParseMode = "Markdown"
	ParseModeMarkdown       ParseMode = "MarkdownV2"
	ParseModeHTML           ParseMode = "HTML"
)

type MessageSender interface {
	SendMessage(chatID ChatID, msg string, opts ...MessageOption)
}

type MessageOption func(cfg *tgbotapi.MessageConfig)

func applyMessageOpts(msg *tgbotapi.MessageConfig, opts []MessageOption) {
	if len(opts) == 0 {
		return
	}

	for _, opt := range opts {
		opt(msg)
	}
}

// WithParseMode sets message parse mode.
func WithParseMode(mode ParseMode) MessageOption {
	return func(cfg *tgbotapi.MessageConfig) {
		cfg.ParseMode = mode
	}
}
