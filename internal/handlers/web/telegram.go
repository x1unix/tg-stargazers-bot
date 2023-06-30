package web

import (
	"encoding/json"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/labstack/echo/v4"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"go.uber.org/zap"
)

const webhookSecretParam = "s"

type TelegramHandler struct {
	log     *zap.Logger
	secrets WebhookSecrets
	botSvc  *bot.Service
}

func NewTelegramHandler(log *zap.Logger, secrets WebhookSecrets, botSvc *bot.Service) TelegramHandler {
	return TelegramHandler{log: log, botSvc: botSvc, secrets: secrets}
}

// HandleTelegramWebhook handles bot webhook calls from Telegram.
func (h TelegramHandler) HandleTelegramWebhook(c echo.Context) (err error) {
	webhookSecret := c.QueryParam(webhookSecretParam)
	if webhookSecret != h.secrets.Telegram {
		h.log.Warn("invalid Telegram webhook secret in request",
			zap.Any("headers", c.Request().Header),
			zap.String("request_url", c.Request().RequestURI),
		)
		return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
	}

	h.log.Debug("received Telegram webhook call",
		zap.Any("headers", c.Request().Header),
		zap.String("request_url", c.Request().RequestURI),
	)

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
		return WrapHTTPError(http.StatusBadRequest, err)
	}

	h.botSvc.HandleUpdate(update)
	return nil
}
