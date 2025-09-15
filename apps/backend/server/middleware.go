package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (s *Server) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if auth == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing auth"})
		}
		_, err := ParseJWT(s.Config.JWTSecret, auth)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
		}
		// token valid - continue
		return next(c)
	}
}

func ExtractTokenFromHeader(h string) string {
	h = strings.TrimSpace(h)
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return h
}
