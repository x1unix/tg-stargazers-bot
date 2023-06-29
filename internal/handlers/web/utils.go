package web

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/x1unix/tg-stargazers-bot/internal/services/auth"
	"go.uber.org/zap"
)

const secretQueryParam = "s"

type tokenValidator interface {
	ValidateToken(ctx context.Context, token *auth.Claims) (auth.UserID, error)
}

func WrapHTTPError(code int, err error) *echo.HTTPError {
	return echo.NewHTTPError(code, err.Error())
}

func logWithContext(l *zap.Logger, c echo.Context) *zap.Logger {
	return l.With(
		zap.String("method", c.Request().Method),
		zap.String("url", c.Request().RequestURI),
		zap.String("ua", c.Request().UserAgent()),
	)
}
