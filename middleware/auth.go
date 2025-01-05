package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Auth ...
func Auth(apiKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request().Header.Get("x-api-key")
			if req != apiKey {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			return next(c)
		}
	}
}
