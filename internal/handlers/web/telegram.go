package web

import (
	"encoding/json"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type TelegramHandler struct {
	log    *zap.Logger
	botSvc *bot.Service
}

func NewTelegramHandler(log *zap.Logger, botSvc *bot.Service) TelegramHandler {
	return TelegramHandler{log: log, botSvc: botSvc}
}

// HandleTelegramWebhook handles bot webhook calls from Telegram.
func (h TelegramHandler) HandleTelegramWebhook(c echo.Context) (err error) {
	h.log.Debug("received Telegram webhook call")

	update := new(tgbotapi.Update)
	body := c.Request().Body
	defer body.Close()

	if err := json.NewDecoder(body).Decode(update); err != nil {
		headers := c.Request().Header
		h.log.Error(
			"failed to decode webhook request body",
			zap.Error(err),
			zap.Any("headers", headers),
		)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	h.botSvc.HandleUpdate(update)
	return nil
}
