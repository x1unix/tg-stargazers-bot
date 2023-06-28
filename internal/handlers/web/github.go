package web

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type GitHubHandler struct {
	log *zap.Logger
}

func NewGitHubHandler(log *zap.Logger) *GitHubHandler {
	return &GitHubHandler{log: log}
}

func (h GitHubHandler) HandleLogin(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	h.log.Info("login!", zap.String("code", code))
	return h.dumpRequest(c)
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
