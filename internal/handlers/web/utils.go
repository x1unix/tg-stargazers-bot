package web

import "github.com/labstack/echo/v4"

func WrapHTTPError(code int, err error) *echo.HTTPError {
	return echo.NewHTTPError(code, err.Error())
}
