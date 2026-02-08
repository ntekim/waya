package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func APIKeyAuth(next echo.HandlerFunc, apiKey string) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Request().Header.Get("x-api-key")

		if key == "" || key != apiKey {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: Invalid or missing x-api-key"})
		}
		return next(c)
	}
}