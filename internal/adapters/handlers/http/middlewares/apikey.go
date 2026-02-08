package middlewares

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func APIKeyAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Request().Header.Get("x-api-key")
		expectedAPIKey := os.Getenv("WAYA_API_KEY")

		if key == "" || key != expectedAPIKey {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: Invalid or missing x-api-key"})
		}
		return next(c)
	}
}