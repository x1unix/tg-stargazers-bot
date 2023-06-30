package web

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"github.com/x1unix/tg-stargazers-bot/internal/services/feedback"
	"github.com/x1unix/tg-stargazers-bot/internal/services/preferences"
	"go.uber.org/zap"
)

type GitHubHandler struct {
	log             *zap.Logger
	httpCfg         config.HTTPConfig
	githubSvc       *preferences.GitHubService
	notificationSvc *feedback.NotificationsService
}

func NewGitHubHandler(
	log *zap.Logger,
	httpCfg config.HTTPConfig,
	githubSvc *preferences.GitHubService,
	notificationSvc *feedback.NotificationsService,
) *GitHubHandler {
	return &GitHubHandler{
		log:             log,
		httpCfg:         httpCfg,
		githubSvc:       githubSvc,
		notificationSvc: notificationSvc,
	}
}

func (h GitHubHandler) HandleLogin(c echo.Context) error {
	user, err := getUserInfo(c)
	if err != nil {
		h.notificationSvc.NotifyAuthFailure(user.UserID, feedback.ErrInvalidToken)
		return err
	}

	code := c.QueryParam("code")
	if code == "" {
		h.notificationSvc.NotifyAuthFailure(user.UserID, feedback.ErrBadAuthCallbackCall)
		return echo.NewHTTPError(http.StatusBadRequest, "missing code")
	}

	h.log.Debug("github login completed",
		zap.String("code", code),
		zap.Int64("uid", user.UserID),
	)

	ctx := c.Request().Context()
	if err := h.githubSvc.SetGitHubToken(ctx, user.UserID, code); err != nil {
		h.notificationSvc.NotifyAuthFailure(user.UserID, feedback.ErrTokenSaveError)
		return fmt.Errorf("failed to save GitHub token: %w", err)
	}

	h.notificationSvc.NotifyAuthSuccessful(user.UserID)
	return c.File(h.httpCfg.StaticFilePath("post-login.html"))
}

func (h GitHubHandler) HandleWebhook(c echo.Context) error {
	return h.dumpRequest(c)
}

func (h GitHubHandler) dumpRequest(c echo.Context) error {
	req := c.Request()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	defer req.Body.Close()
	h.log.Debug("request",
		zap.String("path", req.URL.RawPath),
		zap.ByteString("body", body),
		zap.Any("headers", req.Header))
	return nil
}
