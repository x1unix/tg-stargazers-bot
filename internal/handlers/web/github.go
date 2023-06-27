package web

import (
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
	return h.dumpRequest(c)
}

func (h GitHubHandler) HandleWebhook(c echo.Context) error {
	return h.dumpRequest(c)
}

func (h GitHubHandler) dumpRequest(c echo.Context) error {
	req := c.Request()
	code := c.QueryParam("code")
	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	return nil
}
