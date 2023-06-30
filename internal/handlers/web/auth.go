package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/x1unix/tg-stargazers-bot/internal/services/auth"
	"go.uber.org/zap"
)

const userInfoKey = "userInfo"

type userInfo struct {
	UserID  auth.UserID
	TokenID string
}

func userTokenMiddleware(l *zap.Logger, validator tokenValidator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				logWithContext(l, c).Info("missing required token")
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			claims, ok := token.Claims.(*auth.Claims)
			if !ok {
				logWithContext(l, c).Error(
					"unexpected JWT token clains type",
					zap.String("got_type", fmt.Sprintf("%T", claims)),
				)
				return echo.NewHTTPError(http.StatusForbidden)
			}

			userID, err := validator.ValidateToken(c.Request().Context(), claims)
			if errors.Is(err, auth.ErrInvalidToken) {
				logWithContext(l, c).Info("token ID is expired",
					zap.String("token_id", claims.ID),
					zap.String("sub_id", claims.Subject),
				)
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			if err != nil {
				logWithContext(l, c).Error("failed to validate token",
					zap.String("token_id", claims.ID),
					zap.String("sub_id", claims.Subject),
					zap.Error(err),
				)
				return echo.NewHTTPError(http.StatusInternalServerError)
			}

			c.Set(userInfoKey, userInfo{
				UserID:  userID,
				TokenID: claims.ID,
			})

			return next(c)
		}
	}
}

func requireSecretMiddleware(l *zap.Logger, secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			gotSecret := c.QueryParam(secretQueryParam)
			if gotSecret == "" {
				logWithContext(l, c).Warn("missing secret param")
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			if gotSecret != secret {
				logWithContext(l, c).
					Warn("invalid secret param", zap.String("got", gotSecret))
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}

func getUserInfo(c echo.Context) (userInfo, error) {
	info, ok := c.Get(userInfoKey).(userInfo)
	if !ok {
		return userInfo{}, echo.NewHTTPError(http.StatusInternalServerError, "Missing auth middleware")
	}

	return info, nil
}
