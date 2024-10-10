package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/taufiksty/be-socket/pkg/shared/util"
)

func JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Missing or invalid JWT token",
			})
		}

		// Parse the token
		claims, err := util.ParseToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Invalid JWT token",
			})
		}

		// Store the email in the context for later use
		c.Set("email", claims.Email)
		c.Set("user_id", claims.UserId)
		return next(c)
	}
}
