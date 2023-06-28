package web

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/x1unix/tg-stargazers-bot/internal/services/auth"
	"go.uber.org/zap"
)

type DebugHandler struct {
	log     *zap.Logger
	cfg     ServerConfig
	authSvc *auth.Service
}

func NewDebugHandler(log *zap.Logger, cfg ServerConfig, authSvc *auth.Service) *DebugHandler {
	return &DebugHandler{log: log, cfg: cfg, authSvc: authSvc}
}

func (h DebugHandler) HandleNewToken(c echo.Context) error {
	userParam := c.QueryParam("uid")
	userID, err := auth.ParseUserID(userParam)
	if err != nil {
		return WrapHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	token, err := h.authSvc.CreateUserToken(ctx, userID)
	if err != nil {
		return err
	}

	queryParams := url.Values{
		"t": []string{token},
	}

	authUrl := h.cfg.BaseURL.JoinPath(githubAuthPath)
	authUrl.RawQuery = queryParams.Encode()

	return c.JSON(http.StatusOK, map[string]any{
		"token":    token,
		"auth_url": authUrl.String(),
	})
}
