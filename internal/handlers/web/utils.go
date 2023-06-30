package web

import (
	"context"
	"net/url"

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

// BuildAuthCallbackURL builds GitHub post-auth callback URL.
func BuildAuthCallbackURL(baseUrl *url.URL, token string) *url.URL {
	params := url.Values{
		tokenQueryParam: []string{token},
	}

	newUrl := baseUrl.JoinPath(githubAuthPath)
	newUrl.RawQuery = params.Encode()
	return newUrl
}

func BuildWebhookUrl(baseUrl *url.URL, token string) *url.URL {
	params := url.Values{
		tokenQueryParam: []string{token},
	}

	newUrl := baseUrl.JoinPath(githubWebHookPath)
	newUrl.RawQuery = params.Encode()
	return newUrl
}
