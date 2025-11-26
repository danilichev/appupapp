package handlers

import (
	"apps/api/internal/services"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type PingHandler struct{}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) GetPing(c echo.Context) error {
	logger := c.Logger()

	if u := c.Get("user"); u != nil {
		if token, ok := u.(*jwt.Token); ok {
			if claims, ok := token.Claims.(*services.JwtClaims); ok {
				logger.Infof("User ID: %s", claims.UserId)
			}
		}
	} else {
		logger.Info("Unauthenticated ping")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
}
